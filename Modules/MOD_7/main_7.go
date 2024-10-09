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

const _TIME_SLEEP_S int = 1

var is_writing_GL bool = false

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
		// like a force-stop on the module which doesn't call forceStopLlama().
		forceStopLlama()

		cmd := exec.Command(Utils.GetShell("", ""))
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		//stderr, _ := cmd.StderrPipe()
		err := cmd.Start()
		if err != nil {
			log.Println("Error starting GPT:", err)

			return
		}

		// Begin with the server ID (to say the first hello)
		var device_id string = Utils.Device_settings_GL.Device_ID

		var shut_down bool = false

		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")
		// Start a goroutine to print to the screen and write to a file the output of the LLM model
		go func() {
			buf := bufio.NewReader(stdout)
			var last_answer string = ""
			var last_word string = ""
			for {
				var one_byte []byte = make([]byte, 1)
				n, err := buf.Read(one_byte)
				if n == 0 || err != nil {
					// End of the stream (pipe closed by the main module thread or some error occurred - so shut down)
					shut_down = true

					return
				}

				var one_byte_str string = string(one_byte)
				last_answer += one_byte_str
				//fmt.Print(one_byte_str)

				if is_writing_GL {
					if one_byte_str == " " || one_byte_str == "\n" {
						if last_word != "[3234_START]" && last_word != "[3234_END]" {
							// Meaning: new word written
							_ = gpt_text_txt.WriteTextFile(last_word + one_byte_str, true)
						}

						last_word = ""
					} else {
						last_word += one_byte_str
					}
				}

				if strings.Contains(last_answer, "[3234_START]") {
					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
					is_writing_GL = true
					last_answer = strings.Replace(last_answer, "[3234_START]", "", -1)

					reduceGptTextTxt(gpt_text_txt)
					_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)

					Utils.ModsCommsChannels_GL[Utils.NUM_MOD_WebsiteBackend] <- map[string]any{
						// Send a message to LIB_2 saying the GPT just started writing
						"Message": []byte(device_id + "|L_2|start"),
					}
				} else if strings.Contains(last_answer, "[3234_END]") {
					is_writing_GL = false

					_ = gpt_text_txt.WriteTextFile(getEndString(), true)

					last_word = ""
					last_answer = ""


					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_READY
				}
			}
		}()
		/*go func() {
			buf := bufio.NewReader(stderr)
			for {
				var one_byte []byte = make([]byte, 1)
				n, _ := buf.Read(one_byte)
				if n == 0 {
					// End of the stream (pipe closed by the main module thread)

					return
				}

				var one_byte_str string = string(one_byte)
				fmt.Print(one_byte_str)
			}
		}()*/

		// Configure the LLM model
		writer := bufio.NewWriter(stdin)
		_, _ = writer.WriteString("llama-cli -m " + modUserInfo_GL.Model_loc + " " +
			"--in-suffix [3234_START] --interactive-first --ctx-size 8192 --threads 4 --temp 1.0 --keep -1 --mlock " +
			"--prompt \"" + modUserInfo_GL.Config_str + "\"\n")
		_ = writer.Flush()

		// Wait for the LLM model to start
		Utils.WaitWithStopTIMEDATE(module_stop, 30)

		sendToGPT := func(to_send string) {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY
			_, _ = writer.WriteString(Utils.RemoveNonGraphicChars(to_send) + "\n")
			_ = writer.Flush()

			time.Sleep(5 * time.Second)

			for modGenInfo_GL.State != ModsFileInfo.MOD_7_STATE_READY && !*module_stop {
				// TODO: So if this now waits, how do we /stop him...?
				time.Sleep(1 * time.Second)
			}
		}

		// Keep this here. Seems sometimes it's necessary to say the first hello to Llama3 or it will say it even if we
		// ask something else (or Llama3 might start writing random things).
		sendToGPT("hello")

		// Process the files to input to the LLM model
		for {
			if *module_stop {
				modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
				forceStopLlama()
				_ = stdout.Close()

				return
			}

			var to_process_dir Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
			var file_list []Utils.FileInfo = to_process_dir.GetFileList()
			for len(file_list) > 0 {
				file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
				var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

				var to_process string = *file_path.ReadTextFile()
				if to_process != "" {
					// It comes like: "[device_id]text"
					device_id = to_process[1:strings.Index(to_process, "]")]
					var text string = to_process[strings.Index(to_process, "]") + 1:]

					if strings.HasPrefix(text, "/") {
						// Control commands begin with a slash
						if text == "/clear" || text == "/stop" {
							// Clear the context of the LLM model or stop while its writing by stopping the module (the
							// Manager will restart it)
							shut_down = true
						} else if strings.HasPrefix(text, ASK_WOLFRAM_ALPHA) {
							// Ask Wolfram Alpha the question
							var question string = text[len(ASK_WOLFRAM_ALPHA):]
							result, direct_result := MOD_6.RetrieveWolframAlpha(question)

							if direct_result {
								_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + "The answer is: " + result +
									". " + getEndString(), true)
							} else {
								sendToGPT("Summarize in sentences the following: " + result)
							}
						} else if strings.HasPrefix(text, SEARCH_WIKIPEDIA) {
							// Search for the Wikipedia page title
							var query string = text[len(SEARCH_WIKIPEDIA):]

							_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + MOD_6.RetrieveWikipedia(query) +
								getEndString(), true)
						}
					} else {
						sendToGPT(text)
					}
				}

				Utils.DelElemSLICES(&file_list, idx_to_remove)
				_ = os.Remove(file_path.GPathToStringConversion())

				if shut_down {
					modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
					forceStopLlama()
					_ = stdout.Close()

					return
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
				forceStopLlama()
				_ = stdout.Close()

				return
			}
		}
	}
}

const NO_ERRORS int = 0
const ALREADY_WRITING int = 1
const DEVICE_NOT_ACTIVE int = 2
/*
SpeakOnDevice sends a text to be spoken on a device.

-----CONSTANTS-----

  - NO_ERRORS – no errors
  - ALREADY_WRITING – VISOR is already generating text to some device
  - DEVICE_NOT_ACTIVE – the device is not active

-----CONSTANTS-----

-----------------------------------------------------------

– Params:
  - device_id – the device ID
  - text – the text to be spoken

– Returns:
  - true if the text was sent to be spoken, false if the device is already speaking or the device is not active
 */
func SpeakOnDevice(device_id string, text string) int {
	if is_writing_GL {
		return ALREADY_WRITING
	}
	// TODO: Disabled because the function is now gone - implement it again if this function is needed
	//if !MOD_12.IsDeviceActive(device_id) {
	//	return DEVICE_NOT_ACTIVE
	//}

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")
	reduceGptTextTxt(gpt_text_txt)

	_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + text + getEndString(), true)

	return NO_ERRORS
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
forceStopLlama stops the LLM model by killing its processes.
 */
func forceStopLlama() {
	_, _ = Utils.ExecCmdSHELL([]string{"killall llama-cli"})
}

func getStartString(device_id string) string {
	return "[3234_START:" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + device_id + "|]"
}

func getEndString() string {
	return "[3234_END]\n"
}
