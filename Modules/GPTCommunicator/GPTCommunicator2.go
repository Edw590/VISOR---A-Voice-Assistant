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
	"GPTComm/GPTComm"
	"OnlineInfoChk"
	"Utils"
	"Utils/ModsFileInfo"
	"os"
	"strconv"
	"strings"
	"time"
)

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const STOP_CMD string = "/stop"

const _INACTIVE_SESSION_TIME_S int64 = 30*60

const _TIME_SLEEP_S int = 1

var module_stop_GL *bool = nil

var (
	modDirsInfo_GL  Utils.ModDirsInfo
	modGenInfo_GL  *ModsFileInfo.Mod7GenInfo
	modUserInfo_GL *ModsFileInfo.Mod7UserInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)
	modGenInfo_GL = &Utils.Gen_settings_GL.MOD_7
	modUserInfo_GL = &Utils.User_settings_GL.GPTCommunicator

	module_stop_GL = module_stop

	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

	if modUserInfo_GL.Model_name == "" || modUserInfo_GL.Context_size == 0 {
		time.Sleep(2 * time.Second)

		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPED

		return
	}

	if modGenInfo_GL.N_mems_when_last_memorized == 0 {
		modGenInfo_GL.N_mems_when_last_memorized = 25 // So that the double is 50 for the first time
	}

	// Prepare the session for the temp and dumb sessions
	// Very low timestamp to avoid being selected as latest sessions
	addSessionEntry("temp", -1000000, "")
	addSessionEntry("dumb", -1000000, "")

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	// In case Ollama was started (as opposed to already being running), send a test message for it to actually
	// start and be ready.
	chatWithGPT(Utils.Gen_settings_GL.Device_settings.Id, "test", "temp", GPTComm.ROLE_USER, false)

	go autoMemorize()

	// Process the text to input to the LLM model
	for {
		var to_process_dir Utils.GPath = modDirsInfo_GL.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
		var file_list []Utils.FileInfo = to_process_dir.GetFileList()
		for len(file_list) > 0 {
			file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
			var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

			var to_process string = *file_path.ReadTextFile()
			if to_process != "" {
				// It comes like: "[device ID|session type|role|true/false]text"
				var params_split []string = strings.Split(to_process[1:strings.Index(to_process, "]")], "|")
				// ID of the device sending the text
				var device_id string = params_split[0]
				// Session type to use
				var session_type string = params_split[1]
				// Role to use
				var role string = params_split[2]
				// If more text is coming and VISOR should wait before calling the LLM
				var more_coming bool = params_split[3] == "true"

				var text string = to_process[strings.Index(to_process, "]") + 1:]

				var session_id string = ""
				switch session_type {
					case GPTComm.SESSION_TYPE_NEW:
						session_id = Utils.RandStringGENERAL(10)
					case GPTComm.SESSION_TYPE_TEMP:
						session_id = "temp"
					case GPTComm.SESSION_TYPE_ACTIVE:
						session_id = getActiveSessionId()
						if session_id == "" {
							session_id = Utils.RandStringGENERAL(10)
						}
					default:
						session_id = session_type
				}

				// Control commands begin with a slash
				if strings.Contains(text, ASK_WOLFRAM_ALPHA) {
					// Ask Wolfram Alpha the question
					var question string = text[strings.Index(text, ASK_WOLFRAM_ALPHA)+len(ASK_WOLFRAM_ALPHA):]
					result, direct_result := OnlineInfoChk.RetrieveWolframAlpha(question)

					if direct_result {
						_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + "The answer is: " + result +
							". " + getEndString(), true)
					} else {
						chatWithGPT(device_id, "I've got this from WolframAlpha. Summarize it for me: " + result + "]",
							session_id, role, more_coming)
					}
				} else if strings.Contains(text, SEARCH_WIKIPEDIA) {
					// Search for the Wikipedia page title
					var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

					_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + OnlineInfoChk.RetrieveWikipedia(query) +
						getEndString(), true)
				} else if !strings.Contains(text, STOP_CMD) {
					chatWithGPT(device_id, text, session_id, role, more_coming)
				}
			}

			Utils.DelElemSLICES(&file_list, idx_to_remove)
			_ = os.Remove(file_path.GPathToStringConversion())
		}

		if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPED

			return
		}
	}
}

func addSessionEntry(session_id string, last_interaction_s int64, user_message string) bool {
	var session_exists bool = false
	for _, session := range modGenInfo_GL.Sessions {
		if session.Id == session_id {
			session_exists = true

			break
		}
	}
	if !session_exists {
		// If the session doesn't exist, create it

		var session_name = ""
		if session_id == "temp" {
			session_name = "Temporary session"
		} else if session_id == "dumb" {
			session_name = "Dumb session"
		} else {
			var message_without_add_info string = user_message[strings.Index(user_message, "]") + 1:]
			// I've titled the text for you, Sir: "App Notification Settings on OnePlus Watch".
			// Get the text inside the quotation marks.
			var prompt string = "Create a title for the following text (beginning of a conversation) and put it " +
				"inside \"double quotation marks\", please. Don't include the date and time. Text: " +
				message_without_add_info
			session_name = chatWithGPT(Utils.Gen_settings_GL.Device_settings.Id, prompt, "temp", GPTComm.ROLE_USER, false)
			if strings.Contains(session_name, "\"") {
				session_name = strings.Split(session_name, "\"")[1]
				// Sometimes the name may come like "[name here]", so remove the brackets.
				session_name = strings.Replace(session_name, "[", "", -1)
				session_name = strings.Replace(session_name, "]", "", -1)
			} else {
				session_name = "[Error naming the session]"
			}
		}

		modGenInfo_GL.Sessions = append(modGenInfo_GL.Sessions, ModsFileInfo.Session{
			Id:                 session_id,
			Name:               session_name,
			Created_time_s:     time.Now().Unix(),
			History:            nil,
			Last_interaction_s: last_interaction_s,
			Memorized:          false,
		})

		return true
	}

	return false
}

func getActiveSessionId() string {
	// The latest session with less than 30 minutes of inactivity is considered the active one
	var latest_interaction int64 = 0
	var active_session_id string = ""
	for _, session := range modGenInfo_GL.Sessions {
		if session.Last_interaction_s > latest_interaction &&
				time.Now().Unix() - session.Last_interaction_s < _INACTIVE_SESSION_TIME_S {
			active_session_id = session.Id
			latest_interaction = session.Last_interaction_s
		}
	}

	return active_session_id
}

/*
reduceGptTextTxt reduces the GPT text file to the last 5 entries.

-----------------------------------------------------------

– Params:
  - gpt_text_txt – the GPT text file
*/
func reduceGptTextTxt(gpt_text_txt Utils.GPath) {
	var p_text *string = gpt_text_txt.ReadTextFile()
	if p_text == nil {
		// The file doesn't yet exist
		return
	}

	var entries []string = strings.Split(*p_text, "[3234_START:")
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
	var to_process_dir Utils.GPath = modDirsInfo_GL.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
	var file_list []Utils.FileInfo = to_process_dir.GetFileList()
	for len(file_list) > 0 {
		file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
		var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

		var to_process string = *file_path.ReadTextFile()
		if to_process != "" {
			var text string = to_process[strings.Index(to_process, "]") + 1:]

			if strings.HasSuffix(text, STOP_CMD) {
				_ = os.Remove(file_path.GPathToStringConversion())

				return true
			}
		}

		Utils.DelElemSLICES(&file_list, idx_to_remove)
	}

	return false
}

func getStartString(device_id string) string {
	return "[3234_START:" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + device_id + "|]"
}

func getEndString() string {
	return "[3234_END]\n"
}
