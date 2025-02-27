/*******************************************************************************
 * Copyright 2023-2025 The V.I.S.O.R. authors
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package GPTCommunicator

import (
	"OnlineInfoChk"
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const _START_TOKENS string = "<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
const _END_TOKENS string = "<|start_header_id|>user<|end_header_id|>"

const _TIME_SLEEP_S int = 1

// Module information and state variables
var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
	modGenInfo_GL  *ModsFileInfo.Mod7GenInfo
	modUserInfo_GL *ModsFileInfo.Mod7UserInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_7
		modUserInfo_GL = &Utils.User_settings_GL.GPTCommunicator

		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

		forceStopLlama()

		var visor_intro string = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
		//var visor_functions = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "functions.json").ReadTextFile()
		//visor_intro = strings.Replace(visor_intro, "3234_FUNCTIONS", visor_functions, -1)
		visor_intro = strings.Replace(visor_intro, "3234_NICK", modUserInfo_GL.User_nickname, -1)
		visor_intro = strings.Replace(visor_intro, "\n", "\\n", -1)
		visor_intro = strings.Replace(visor_intro, "\"", "\\\"", -1)

		var memories string = strings.Join(modGenInfo_GL.Memories, ". ")
		memories = strings.Replace(memories, "\"", "\\\"", -1)

		var smart_ctx_size int = 12288
		var dumb_ctx_size int = 4096

		// Log the configuration of both LLaMa instances
		//log.Printf("Starting 2 LLaMa instances with the following configurations:\n" +
		//	"	- Smart: ctxSize=%d, threads=%d, temp=%f\n" +
		//	"	- Dumb: ctxSize=%d, threads=%d, temp=%f\n",
		//	smart_ctx_size, 4, 0.8, dumb_ctx_size, 4, 1.5,
		//)

		writer_smart, stdout_smart, stderr_smart := startLlama("Smart", smart_ctx_size, 4, 0.8,
			modUserInfo_GL.Model_smart_loc, memories, visor_intro)
		if writer_smart == nil {
			log.Println("Error starting the Llama model (smart)")

			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
			forceStopLlama()

			return
		}
		reader_smart := bufio.NewReader(stdout_smart)

		writer_dumb, stdout_dumb, stderr_dumb := startLlama("Dumb", dumb_ctx_size, 4, 1.5,
			modUserInfo_GL.Model_dumb_loc, "", "You're a voice assistant")
		if writer_dumb == nil {
			log.Println("Error starting the Llama model (dumb)")

			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
			forceStopLlama()
			_ = stdout_smart.Close()

			return
		}
		reader_dumb := bufio.NewReader(stdout_dumb)

		// Run background threads to read stderr output (prevent process termination)
		go func() {
			buf := bufio.NewReader(stderr_smart)
			for {
				var one_byte []byte = make([]byte, 1)
				n, _ := buf.Read(one_byte)
				if n == 0 {
					// End of the stream (pipe closed by the main module thread)

					return
				}

				//fmt.Print(string(one_byte))
			}
		}()
		go func() {
			buf := bufio.NewReader(stderr_dumb)
			for {
				var one_byte []byte = make([]byte, 1)
				n, _ := buf.Read(one_byte)
				if n == 0 {
					// End of the stream (pipe closed by the main module thread)

					return
				}

				//fmt.Print(string(one_byte))
			}
		}()

		var device_id string = ""
		var shut_down bool = false
		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

		if !gpt_text_txt.Exists() {
			// The file must start with the start string, even if it's empty.
			_ = gpt_text_txt.WriteTextFile("[3234_START:-1|5678iuytrrtyuj8uy|]", false)
		}

		var first_3234_end bool = true
		var memorizing bool = false
		var to_memorize string = ""

		readGPT := func(reader *bufio.Reader, print bool) {
			var last_answer string = ""
			var last_word string = ""
			var is_writing bool = false
			var save_words bool = true
			for {
				var one_byte []byte = make([]byte, 1)
				n, err := reader.Read(one_byte)
				if n == 0 || err != nil {
					// End of the stream (pipe closed by the main module thread or some error occurred - so shut down)
					shut_down = true

					return
				}

				var one_byte_str string = string(one_byte)
				last_answer += one_byte_str
				if memorizing {
					to_memorize += one_byte_str
				}
				if print {
					fmt.Print(one_byte_str)
				}

				if is_writing {
					if one_byte_str == " " || one_byte_str == "\n" {
						if last_word != _START_TOKENS && last_word != _END_TOKENS {
							// Meaning: new word written
							_ = gpt_text_txt.WriteTextFile(last_word + one_byte_str, true)
						}

						last_word = ""
					} else {
						// VISOR may start by writing the current date and time like "[date and time here]" - this
						// below cuts that out of the answer.
						//if last_word == "" {
						//	if one_byte_str == "[" {
						//		save_words = false
						//	} else if one_byte_str == "]" {
						//		save_words = true
						//		continue
						//	}
						//}

						if save_words {
							last_word += one_byte_str
						}
					}
				}

				if strings.Contains(last_answer, _START_TOKENS) {
					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
					is_writing = true
					last_answer = strings.Replace(last_answer, _START_TOKENS, "", -1)

					reduceGptTextTxt(gpt_text_txt)
					_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)

					// Send a message to LIB_2 saying the GPT just started writing
					Utils.SendToModChannel(Utils.NUM_MOD_WebsiteBackend, 0, "Message", []byte(device_id+"|L_2_0|start"))
				} else if strings.Contains(last_answer, _END_TOKENS) {
					is_writing = false

					_ = gpt_text_txt.WriteTextFile(getEndString(), true)

					last_word = ""
					last_answer = ""

					// The first time is the "dumb" LLM being ready. The 2nd time is the "smart" one.
					if first_3234_end {
						first_3234_end = false
					} else {
						modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_READY
					}
				}
			}
		}

		// Read from both GPTs --> but make them *never* work both at the same time. Only one at a time answering.
		go func() {
			readGPT(reader_smart, true)
		}()
		go func() {
			readGPT(reader_dumb, false)
		}()

		// Wait for the LLM models to start
		for modGenInfo_GL.State != ModsFileInfo.MOD_7_STATE_READY {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				return
			}
		}

		var user_text string = ""
		sendToGPT := func(to_send string, use_smart bool) {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
			var to_write string = UtilsSWA.RemoveNonGraphicCharsGENERAL(to_send)
			if use_smart {
				user_text += to_write + ". "
				_, _ = writer_smart.WriteString(to_write + "\n")
				_ = writer_smart.Flush()
			} else {
				_, _ = writer_dumb.WriteString(to_write + "\n")
				_ = writer_dumb.Flush()
			}

			for modGenInfo_GL.State != ModsFileInfo.MOD_7_STATE_READY && !*module_stop {
				if checkStopSpeech() {
					// Write the end string before exiting
					_ = gpt_text_txt.WriteTextFile(getEndString(), true)

					// Not sure how to send a Ctrl+C signal to the process in a way that works (Linux, at least). So
					// plan B and the process is killed, also clearing the context, unfortunately.
					shut_down = true

					break
				}

				time.Sleep(1 * time.Second)
			}
		}

		memorizeThings := func(input_text string, erase_mem bool) {
			device_id = Utils.Gen_settings_GL.Device_settings.Id // To not speak in any device the memorization strings
			memorizing = true
			var text string = "User input: " + input_text + ". PROFILE the USER based on their behavior, preferences, " +
				"personality traits, or habits revealed in their input. IGNORE specific, temporary events, schedules, " +
				"or day-to-day plans. Summarize as KEY GENERAL user information in BULLET points (no + or - or " +
				"anything. ONLY *). Format the output as \"* The user [detail]\". Example: \"* The user is interested " +
				"in technology\"."
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
			_, _ = writer_dumb.WriteString(text + "\n")
			_ = writer_dumb.Flush()

			for modGenInfo_GL.State != ModsFileInfo.MOD_7_STATE_READY {
				time.Sleep(1 * time.Second)
			}
			memorizing = false

			if erase_mem {
				modGenInfo_GL.Memories = nil
			}
			var memories_split []string = strings.Split(to_memorize, "\n")
			for _, memory := range memories_split {
				if UtilsSWA.StringHasLettersGENERAL(memory) && strings.Contains(memory, "* ") &&
						!strings.Contains(strings.ToLower(memory), "none") {
					memory = strings.Replace(memory, "* The user's ", "* The user ", -1)
					var the_user_idx int = strings.LastIndex(memory, "* The user ")
					if the_user_idx == -1 {
						continue
					}

					modGenInfo_GL.Memories = append(modGenInfo_GL.Memories, memory[the_user_idx + len("* The user "):])
				}
			}

			// Give time to write everything down
			time.Sleep(6 * time.Second)
		}

		shutDown := func() {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
			forceStopLlama()
			_ = stdout_smart.Close()
			_ = stderr_smart.Close()
			_ = stdout_dumb.Close()
			_ = stderr_dumb.Close()
		}

		// Process the text to input to the LLM model
		for {
			var to_process_dir Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
			var file_list []Utils.FileInfo = to_process_dir.GetFileList()
			for len(file_list) > 0 {
				file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
				var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

				var to_process string = *file_path.ReadTextFile()
				if to_process != "" {
					// It comes like: "[device_id|[true or false]]text"
					var params_split []string = strings.Split(to_process[1:strings.Index(to_process, "]")], "|")
					device_id = params_split[0]
					var use_smart bool = params_split[1] == "true"
					var text string = to_process[strings.Index(to_process, "]")+1:]
					text = strings.Replace(text, "\n", "\\n", -1)

					if use_smart {
						// Control commands begin with a slash
						if strings.HasSuffix(text, "/clear") {
							// Clear the context of the LLM model by stopping the module (the Manager will restart it)
							shut_down = true
						} else if strings.HasSuffix(text, "/mem") {
							// Memorize and clear the context
							if user_text != "" {
								memorizeThings(user_text, false)
							}

							shut_down = true
						} else if strings.HasSuffix(text, "/memmem") {
							// Summarize the list of memories (sometimes VISOR may memorize useless sentences, so this will
							// cut them out - will cut out other things too though, so use with caution).
							if len(modGenInfo_GL.Memories) > 0 {
								var memories_str string = strings.Join(modGenInfo_GL.Memories, ". ")
								memorizeThings(memories_str, true)
							}
						} else if strings.Contains(text, ASK_WOLFRAM_ALPHA) {
							// Ask Wolfram Alpha the question
							var question string = text[strings.Index(text, ASK_WOLFRAM_ALPHA)+len(ASK_WOLFRAM_ALPHA):]
							result, direct_result := OnlineInfoChk.RetrieveWolframAlpha(question)

							if direct_result {
								_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + "The answer is: " + result +
									". " + getEndString(), true)
							} else {
								sendToGPT("Summarize in sentences the following: " + result, false)
							}
						} else if strings.Contains(text, SEARCH_WIKIPEDIA) {
							// Search for the Wikipedia page title
							var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

							_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + OnlineInfoChk.RetrieveWikipedia(query) +
								getEndString(), true)
						} else {
							sendToGPT(text, true)
						}
					} else {
						sendToGPT(text, false)
					}
				}

				Utils.DelElemSLICES(&file_list, idx_to_remove)
				_ = os.Remove(file_path.GPathToStringConversion())

				if shut_down {
					shutDown()

					return
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				shutDown()

				return
			}
		}
	}
}

func startLlama(instance_type string, ctx_size int, threads int, temp float32, model_loc string, memories string,
				visor_intro string) (*bufio.Writer, io.ReadCloser, io.ReadCloser) {
	cmd := exec.Command(Utils.GetShellSHELL("", ""))
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting %s LLaMA instance: %v", instance_type, err)

		return nil, nil, nil
	}

	var system_info string = modUserInfo_GL.System_info
	system_info = strings.Replace(system_info, "3234_WEEKDAY", time.Now().Weekday().String(), -1)
	system_info = strings.Replace(system_info, "3234_DAY", strconv.Itoa(time.Now().Day()), -1)
	system_info = strings.Replace(system_info, "3234_MONTH", time.Now().Month().String(), -1)
	system_info = strings.Replace(system_info, "3234_YEAR", strconv.Itoa(time.Now().Year()), -1)

	// Configure the LLM model (Llama3/3.1/3.2's prompt)
	writer := bufio.NewWriter(stdin)
	_, err = writer.WriteString("llama-cli " +
		"--model \"" + model_loc + "\" " +
		"--interactive-first " +
		"--ctx-size " + strconv.Itoa(ctx_size) + " " +
		"--threads " + strconv.Itoa(threads) + " " +
		"--temp " + strconv.FormatFloat(float64(temp), 'f', -1, 32) + " " +
		"--keep -1 " +
		"--mlock " +
		"--prompt \"<|begin_of_text|><|start_header_id|>system<|end_header_id|>" + system_info + "\\n\\n\\n" +
		"Memories stored about the user: " + memories + ".\\n\\n\\nAbout you:\\n" + visor_intro + "<|eot_id|>\" " +
		"--reverse-prompt \"<|eot_id|>\" " +
		"--in-prefix \"" + _END_TOKENS + "\" " +
		"--in-suffix \"" + _START_TOKENS + "\" " +
		"\n")
	if err != nil {
		log.Printf("Error writing to %s LLaMA instance stdin: %v", instance_type, err)
		return nil, nil, nil
	}
	err = writer.Flush()
	if err != nil {
		log.Printf("Error flushing %s LLaMA instance stdin: %v", instance_type, err)
		return nil, nil, nil
	}

	return writer, stdout, stderr
}

/*
reduceGptTextTxt reduces the GPT text file to the last 5 entries.

-----------------------------------------------------------

– Params:
  - gpt_text_txt – the GPT text file
*/
func reduceGptTextTxt(gpt_text_txt Utils.GPath) {
	var text string = *gpt_text_txt.ReadTextFile()
	var entries []string = strings.Split(text, "[3234_START:")
	if len(entries) > 5 {
		_ = gpt_text_txt.WriteTextFile("[3234_START:" + entries[len(entries)-5], false)

		for i := len(entries) - 4; i < len(entries); i++ {
			_ = gpt_text_txt.WriteTextFile("[3234_START:" + entries[i], true)
		}
	}
}

/*
checkStopSpeech checks if the text to process contains the /stop command.

-----------------------------------------------------------

– Returns:
  - true if the /stop command was found, false otherwise
*/
func checkStopSpeech() bool {
	var to_process_dir Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
	var file_list []Utils.FileInfo = to_process_dir.GetFileList()
	for len(file_list) > 0 {
		file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
		var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

		var to_process string = *file_path.ReadTextFile()
		if to_process != "" {
			var text string = to_process[strings.Index(to_process, "]") + 1:]

			if text == "/stop" {
				_ = os.Remove(file_path.GPathToStringConversion())

				return true
			}
		}

		Utils.DelElemSLICES(&file_list, idx_to_remove)
	}

	return false
}

/*
forceStopLlama stops the LLM model by killing its processes.
*/
func forceStopLlama() {
	Utils.KillAllPROCESSES("llama-cli")
}

func getStartString(device_id string) string {
	return "[3234_START:" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + device_id + "|]"
}

func getEndString() string {
	return "[3234_END]\n"
}
