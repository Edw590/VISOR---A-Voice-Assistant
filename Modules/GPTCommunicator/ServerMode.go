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
	"GPTComm"
	"OnlineInfoChk"
	"Utils"
	"Utils/ModsFileInfo"
	"os"
	"strings"
	"time"
)

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const STOP_CMD string = "/stop"

const _INACTIVE_SESSION_TIME_S int64 = 30*60

const _TIME_SLEEP_S int = 1

func serverMode() {
	getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STARTING

	if getModUserInfo().Models_to_use == "" || getModUserInfo().Context_size == 0 {
		time.Sleep(2 * time.Second)

		getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STOPPED

		return
	}

	if getModGenSettings().N_mems_when_last_memorized == 0 {
		getModGenSettings().N_mems_when_last_memorized = 25 // So that the double is 50 for the first time
	}

	// Prepare the session for the temp and dumb sessions
	// Very low timestamp to avoid being selected as latest sessions
	addSessionEntry("temp", -1000000, "")
	addSessionEntry("dumb", -1000000, "")

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	// In case Ollama was started (as opposed to already being running), send a test message for it to actually
	// start and be ready.
	var chatWithGPT_params _ChatWithGPTParams = _ChatWithGPTParams{
		Device_id:    Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id,
		User_message: "test",
		Session_id:   "temp",
		Role:         GPTComm.ROLE_USER,
		More_coming:  false,
		Model_type:   GPTComm.MODEL_TYPE_TEXT,
	}
	chatWithGPT(chatWithGPT_params)

	//go autoMemorize() TODO

	// Process the text to input to the LLM model
	for {
		var to_process_dir Utils.GPath = modDirsInfo_GL.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
		var file_list []Utils.FileInfo = to_process_dir.GetFileList()
		for len(file_list) > 0 {
			file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
			var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

			var to_process string = *file_path.ReadTextFile()
			if to_process != "" {
				// It comes like: "[device ID|session type|role|true/false|model type]text"
				var params_split []string = strings.Split(to_process[1:strings.Index(to_process, "]")], "|")
				var device_id string = params_split[0]
				// Session type to use
				var session_type string = params_split[1]

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

				var text string = to_process[strings.Index(to_process, "]") + 1:]

				chatWithGPT_params = _ChatWithGPTParams{
					Device_id:    params_split[0],
					User_message: to_process[strings.Index(to_process, "]") + 1:],
					Session_id:   session_id,
					Role:         params_split[2],
					More_coming:  params_split[3] == "true",
					Model_type:   params_split[4],
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
						chatWithGPT_params.User_message = "I've got this from WolframAlpha. Summarize it for me: " +
							result + "]"
						chatWithGPT(chatWithGPT_params)
					}
				} else if strings.Contains(text, SEARCH_WIKIPEDIA) {
					// Search for the Wikipedia page title
					var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

					_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + OnlineInfoChk.RetrieveWikipedia(query) +
						getEndString(), true)
				} else if !strings.Contains(text, STOP_CMD) {
					chatWithGPT(chatWithGPT_params)
				}
			}

			Utils.DelElemSLICES(&file_list, idx_to_remove)
			_ = os.Remove(file_path.GPathToStringConversion())
		}

		if Utils.WaitWithStopDATETIME(module_stop_GL, _TIME_SLEEP_S) {
			getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STOPPED

			return
		}
	}
}

func addSessionEntry(session_id string, last_interaction_s int64, user_message string) bool {
	var session_exists bool = false
	for _, session := range getModGenSettings().Sessions {
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
				"inside \"double quotation marks\", please. Don't include the date and time. Text --> " +
				message_without_add_info + " <--"
			var chatWithGPT_params _ChatWithGPTParams = _ChatWithGPTParams{
				Device_id:    Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id,
				User_message: prompt,
				Session_id:   "temp",
				Role:         GPTComm.ROLE_USER,
				More_coming:  false,
				Model_type:   GPTComm.MODEL_TYPE_TEXT,
			}
			session_name = chatWithGPT(chatWithGPT_params)
			if strings.Contains(session_name, "\"") {
				session_name = strings.Split(session_name, "\"")[1]
				// Sometimes the name may come like "[name here]", so remove the brackets.
				session_name = strings.Replace(session_name, "[", "", -1)
				session_name = strings.Replace(session_name, "]", "", -1)
			} else {
				session_name = "[Error naming the session]"
			}
		}

		getModGenSettings().Sessions = append(getModGenSettings().Sessions, ModsFileInfo.Session{
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

func getActiveSessionId() string {
	// The latest session with less than 30 minutes of inactivity is considered the active one
	var latest_interaction int64 = 0
	var active_session_id string = ""
	for _, session := range getModGenSettings().Sessions {
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

func setBusyState() {
	if getModGenSettings().State != ModsFileInfo.MOD_7_STATE_STARTING {
		// If the module is starting, keep it on the starting state until it becomes ready. Else, set it to busy.
		getModGenSettings().State = ModsFileInfo.MOD_7_STATE_BUSY
	}
}

func setReadyState() {
	getModGenSettings().State = ModsFileInfo.MOD_7_STATE_READY
}

func getModGenSettings() *ModsFileInfo.Mod7GenInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7
}
