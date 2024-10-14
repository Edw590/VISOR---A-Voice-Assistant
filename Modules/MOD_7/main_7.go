/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

package MOD_7

import (
	MOD_6 "OnlineInfoChk"
	"Utils"
	"Utils/ModsFileInfo"
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

// GPT Communicator //

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const _START_TOKENS string = "<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
const _END_TOKENS string = "<|start_header_id|>user<|end_header_id|>"

const _TIME_SLEEP_S int = 1

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
		modUserInfo_GL = &Utils.User_settings_GL.MOD_7

		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

		// Force stop Llama to start fresh, in case for any reason it's running without the module being running too,
		// like a force-stop on VISOR which then won't call forceStopLlama().
		forceStopLlama()

		var visor_intro string = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
		visor_intro = strings.Replace(visor_intro, "\n", " ", -1)
		visor_intro = strings.Replace(visor_intro, "\"", "\\\"", -1)

		writer_smart, stdout_smart, stderr_smart := startLlama(12288, 4, 0.8, modUserInfo_GL.Model_smart_loc,
			modUserInfo_GL.User_intro, visor_intro)
		if writer_smart == nil {
			log.Println("Error starting the Llama model (smart)")

			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
			forceStopLlama()

			return
		}
		reader_smart := bufio.NewReader(stdout_smart)

		writer_dumb, stdout_dumb, stderr_dumb := startLlama(4096, 4, 1.5, modUserInfo_GL.Model_dumb_loc, "", visor_intro)
		if writer_dumb == nil {
			log.Println("Error starting the Llama model (dumb)")

			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
			forceStopLlama()
			_ = stdout_smart.Close()

			return
		}
		reader_dumb := bufio.NewReader(stdout_dumb)

		// Leave these 2 threads running. Without them the processes stop for some reason.
		go func() {
			buf := bufio.NewReader(stderr_smart)
			for {
				var one_byte []byte = make([]byte, 1)
				n, _ := buf.Read(one_byte)
				if n == 0 {
					// End of the stream (pipe closed by the main module thread)

					return
				}

				//var one_byte_str string = string(one_byte)
				//fmt.Print(one_byte_str)
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

				//var one_byte_str string = string(one_byte)
				//fmt.Print(one_byte_str)
			}
		}()

		var device_id string = ""
		var shut_down bool = false

		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

		readGPT := func(reader *bufio.Reader, print bool) {
			var last_answer string = ""
			var last_word string = ""
			var is_writing bool = false
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
						last_word += one_byte_str
					}
				}

				if strings.Contains(last_answer, _START_TOKENS) {
					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
					is_writing = true
					last_answer = strings.Replace(last_answer, _START_TOKENS, "", -1)

					reduceGptTextTxt(gpt_text_txt)
					_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)

					Utils.ModsCommsChannels_GL[Utils.NUM_MOD_WebsiteBackend] <- map[string]any{
						// Send a message to LIB_2 saying the GPT just started writing
						"Message": []byte(device_id + "|L_2|start"),
					}
				} else if strings.Contains(last_answer, _END_TOKENS) {
					is_writing = false

					_ = gpt_text_txt.WriteTextFile(getEndString(), true)

					last_word = ""
					last_answer = ""

					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_READY
				}
			}
		}

		// Read from both GPTs --> but make them *never* work both at the same time. Only one at a time answering.
		go func() {
			readGPT(reader_smart, true)
		}()
		go func() {
			readGPT(reader_dumb, true)
		}()

		// Wait for the LLM model to start
		Utils.WaitWithStopTIMEDATE(module_stop, 30)

		sendToGPT := func(to_send string, use_smart bool) {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
			var to_write string = Utils.RemoveNonGraphicCharsGENERAL(to_send) + "\n"
			if use_smart {
				_, _ = writer_smart.WriteString(to_write)
				_ = writer_smart.Flush()
			} else {
				_, _ = writer_dumb.WriteString(to_write)
				_ = writer_dumb.Flush()
			}

			time.Sleep(5 * time.Second)

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

		// Process the files to input to the LLM model
		for {
			if *module_stop {
				modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
				forceStopLlama()
				_ = stdout_smart.Close()
				_ = stdout_dumb.Close()

				return
			}

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
					var text string = to_process[strings.Index(to_process, "]") + 1:]

					if use_smart && strings.HasPrefix(text, "/") {
						// Control commands begin with a slash
						if text == "/clear" {
							// Clear the context of the LLM model by stopping the module (the Manager will restart it)
							shut_down = true
						} else if strings.HasPrefix(text, ASK_WOLFRAM_ALPHA) {
							// Ask Wolfram Alpha the question
							var question string = text[len(ASK_WOLFRAM_ALPHA):]
							result, direct_result := MOD_6.RetrieveWolframAlpha(question)

							if direct_result {
								_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + "The answer is: " + result +
									". " + getEndString(), true)
							} else {
								sendToGPT("Summarize in sentences the following: " + result, true)
							}
						} else if strings.HasPrefix(text, SEARCH_WIKIPEDIA) {
							// Search for the Wikipedia page title
							var query string = text[len(SEARCH_WIKIPEDIA):]

							_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + MOD_6.RetrieveWikipedia(query) +
								getEndString(), true)
						}
					} else {
						sendToGPT(text, use_smart)
					}
				}

				Utils.DelElemSLICES(&file_list, idx_to_remove)
				_ = os.Remove(file_path.GPathToStringConversion())

				if shut_down {
					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
					forceStopLlama()
					_ = stdout_smart.Close()
					_ = stdout_dumb.Close()

					return
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
				forceStopLlama()
				_ = stdout_smart.Close()
				_ = stdout_dumb.Close()

				return
			}
		}
	}
}

func startLlama(ctx_size int, threads int, temp float32, model_loc string, user_intro string, visor_intro string) (*bufio.Writer, io.ReadCloser, io.ReadCloser) {
	cmd := exec.Command(Utils.GetShell("", ""))
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Println("Error starting a command shell:", err)

		return nil, nil, nil
	}

	// Configure the LLM model (Llama3/3.1/3.2's prompt)
	writer := bufio.NewWriter(stdin)
	_, _ = writer.WriteString("llama-cli " +
		"--model \"" + model_loc + "\" " +
		"--interactive-first " +
		"--ctx-size " + strconv.Itoa(ctx_size) + " " +
		"--threads " + strconv.Itoa(threads) + " " +
		"--temp " + strconv.FormatFloat(float64(temp), 'f', -1, 32) + " " +
		"--keep -1 " +
		"--mlock " +
		"--prompt \"<|begin_of_text|><|start_header_id|>system<|end_header_id|>" +
		strings.Replace(modUserInfo_GL.System_info, "3234_YEAR", strconv.Itoa(time.Now().Year()), -1) + " " +
		user_intro + ". " + visor_intro + "<|eot_id|>\" " +
		"--reverse-prompt \"<|eot_id|>\" " +
		"--in-prefix \"" + _END_TOKENS + "\" " +
		"--in-suffix \"" + _START_TOKENS + "\" " +
		"\n")
	_ = writer.Flush()

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
