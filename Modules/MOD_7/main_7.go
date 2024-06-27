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
	"bufio"
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

var is_speaking_GL bool = false

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		var modUserInfo _ModUserInfo
		if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
			panic(err)
		}

		// Force stop Llama to start fresh, in case for any reason it's running without the module being running too,
		// like a force-stop on the module which doesn't call forceStopLlama().
		forceStopLlama()

		cmd := exec.Command(Utils.GetShell("", ""))
		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		_ = cmd.Start()

		// Begin with the server ID (to say the first hello)
		var device_id string = Utils.PersonalConsts_GL.DEVICE_ID

		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")
		// Start a goroutine to print to the screen and write to a file the output of the LLM model
		go func() {

			buf := bufio.NewReader(stdout)
			var last_answer string = ""
			var last_word string = ""
			for {
				var one_byte []byte = make([]byte, 1)
				n, _ := buf.Read(one_byte)
				if n == 0 {
					// End of the stream (pipe closed by the main module thread)

					return
				}

				var one_byte_str string = string(one_byte)
				last_answer += one_byte_str
				//fmt.Print(one_byte_str)

				if is_speaking_GL {
					if one_byte_str == " " || one_byte_str == "\n" {
						if last_word != "[3234_START]" && last_word != "[3234_END]" {
							_ = gpt_text_txt.WriteTextFile(last_word + one_byte_str, true)
						}

						last_word = ""
					} else {
						last_word += one_byte_str
					}
				}

				if strings.Contains(last_answer, "[3234_START]") {
					is_speaking_GL = true
					last_answer = strings.Replace(last_answer, "[3234_START]", "", -1)

					_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)
				} else if strings.Contains(last_answer, "[3234_END]") {
					is_speaking_GL = false

					_ = gpt_text_txt.WriteTextFile(getEndString(), true)

					last_word = ""
					last_answer = ""
				}
			}
		}()

		// Configure the LLM model
		var config_str string = *moduleInfo_GL.ModDirsInfo.UserData.Add2(false, "config_string.txt").ReadTextFile()
		writer := bufio.NewWriter(stdin)
		_, _ = writer.WriteString("llama-cli -m " + modUserInfo.Model_loc + " " +
			"--in-suffix [3234_START] --interactive-first --ctx-size 0 --threads 4 --temp 0.2 --mlock " +
			"--prompt \"" + config_str + "\"\n")
		_ = writer.Flush()

		sendToGPT := func(to_send string) {
			_, _ = writer.WriteString(Utils.RemoveNonGraphicChars(to_send) + "\n")
			_ = writer.Flush()
		}

		sendToGPT("hello")

		// Process the files to input to the LLM model
		for {
			var to_process_dir Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
			var file_list []Utils.FileInfo = to_process_dir.GetFileList()
			for len(file_list) > 0 {
				var shut_down bool = false
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
							var query string = text[len(ASK_WOLFRAM_ALPHA):]

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
					forceStopLlama()
					_ = stdout.Close()

					return
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				forceStopLlama()
				_ = stdout.Close()

				return
			}
		}
	}
}

/*
SpeakOnDevice sends a text to be spoken on a device.

-----------------------------------------------------------

– Params:
  - device_id – the device ID
  - text – the text to be spoken

– Returns:
  - true if the text was sent to be spoken, false if the device is already speaking
 */
func SpeakOnDevice(device_id string, text string) bool {
	if is_speaking_GL {
		return false
	}

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + text + getEndString(), true)

	return true
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
