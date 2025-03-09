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

// Prepared for Ollama 0.5.11

type _OllamaRequest struct {
	Model string `json:"model"`
	Messages []ModsFileInfo.OllamaMessage `json:"messages"`

	Format  string `json:"format"`
	Options _OllamaOptions `json:"options"`
	Stream bool `json:"stream"`
	Keep_alive string `json:"keep_alive"`
}

type _OllamaOptions struct {
	Num_ctx int `json:"num_ctx"`
	Temperature float32 `json:"temperature"`
}

type _OllamaResponse struct {
	Model string `json:"model"`
	Created_at string `json:"created_at"`
	Message ModsFileInfo.OllamaMessage `json:"message"`
	Done bool `json:"done"`
	Total_duration int `json:"total_duration"`
	Load_duration int `json:"load_duration"`
	Prompt_eval_count int `json:"prompt_eval_count"`
	Prompt_eval_duration int `json:"prompt_eval_duration"`
	Eval_count int `json:"eval_count"`
	Eval_duration int `json:"eval_duration"`
}

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const _INACTIVE_SESSION_TIME_S int64 = 30*60

const _TIME_SLEEP_S int = 1

var module_stop_GL *bool = nil

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

		module_stop_GL = module_stop

		if modUserInfo_GL.Model_name == "" {
			return
		}

		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

		if modGenInfo_GL.N_mems_when_last_memorized == 0 {
			modGenInfo_GL.N_mems_when_last_memorized = 25 // So that the double is 50 for the first time
		}

		// Start Ollama in case it's not running
		startOllama()
		time.Sleep(2 * time.Second)

		if modGenInfo_GL.Sessions == nil {
			modGenInfo_GL.Sessions = make(map[string]*ModsFileInfo.Session)
		}

		// Prepare the session for the temp and dumb sessions
		addSessionEntry("temp", nil, -1000000, "") // A very low timestamp to avoid being selected as latest session
		addSessionEntry("dumb", nil, -1000000, "") // A very low timestamp to avoid being selected as latest session

		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

		// In case Ollama was started (as opposed to already being running), send a test message for it to actually
		// start and be ready.
		chatWithGPT(Utils.Gen_settings_GL.Device_settings.Id, "test", "dumb")

		go autoMemorize()

		// Process the text to input to the LLM model
		for {
			var to_process_dir Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _TO_PROCESS_REL_FOLDER)
			var file_list []Utils.FileInfo = to_process_dir.GetFileList()
			for len(file_list) > 0 {
				file_to_process, idx_to_remove := Utils.GetOldestFileFILESDIRS(file_list)
				var file_path Utils.GPath = to_process_dir.Add2(false, file_to_process.Name)

				var to_process string = *file_path.ReadTextFile()
				if to_process != "" {
					// It comes like: "[device ID|session type]text"
					var params_split []string = strings.Split(to_process[1:strings.Index(to_process, "]")], "|")
					var device_id = params_split[0]
					var session_type = params_split[1]

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
					if strings.HasSuffix(text, "/restartGPT") {
						restartOllama()
						time.Sleep(2 * time.Second)
						chatWithGPT(Utils.Gen_settings_GL.Device_settings.Id, "test", "dumb")
					} else if strings.Contains(text, ASK_WOLFRAM_ALPHA) {
						// Ask Wolfram Alpha the question
						var question string = text[strings.Index(text, ASK_WOLFRAM_ALPHA)+len(ASK_WOLFRAM_ALPHA):]
						result, direct_result := OnlineInfoChk.RetrieveWolframAlpha(question)

						if direct_result {
							_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + "The answer is: " + result +
								". " + getEndString(), true)
						} else {
							sendToGPT(device_id, "I've got this from WolframAlpha. Summarize it for me: " + result + "]",
								session_id)
						}
					} else if strings.Contains(text, SEARCH_WIKIPEDIA) {
						// Search for the Wikipedia page title
						var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

						_ = gpt_text_txt.WriteTextFile(getStartString(device_id) + OnlineInfoChk.RetrieveWikipedia(query) +
							getEndString(), true)
					} else {
						sendToGPT(device_id, text, session_id)
					}
				}

				Utils.DelElemSLICES(&file_list, idx_to_remove)
				_ = os.Remove(file_path.GPathToStringConversion())
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING

				return
			}
		}
	}
}

func sendToGPT(device_id string, user_message string, session_id string) string {
	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY

	return chatWithGPT(device_id, user_message, session_id)
}

func addSessionEntry(session_id string, history []ModsFileInfo.OllamaMessage, last_interaction_s int64, user_message string) bool {
	if _, ok := modGenInfo_GL.Sessions[session_id]; !ok {
		// If the session doesn't exist, create it

		var session_name = ""
		if session_id == "temp" {
			session_name = "Temporary session"
		} else if session_id == "dumb" {
			session_name = "Dumb session"
		} else {
			// I've titled the text for you, Sir: "App Notification Settings on OnePlus Watch".
			// Get the text inside the quotation marks.
			var prompt string = "Create a title for the following text and put it inside \"quotation marks\", please. " +
				"Text: " + user_message
			session_name = chatWithGPT(Utils.Gen_settings_GL.Device_settings.Id, prompt, "temp")
			if strings.Contains(session_name, "\"") {
				session_name = strings.Split(session_name, "\"")[1]
				// Sometimes the name may come like "[name here]", so remove the brackets.
				session_name = strings.Replace(session_name, "[", "", -1)
				session_name = strings.Replace(session_name, "]", "", -1)
			} else {
				session_name = "[Error naming the session]"
			}
		}

		modGenInfo_GL.Sessions[session_id] = &ModsFileInfo.Session{
			Name:               session_name,
			Created_time_s:     time.Now().Unix(),
			History:            history,
			Last_interaction_s: last_interaction_s,
			Memorized:          false,
		}

		return true
	}

	return false
}

func getVisorIntroAndMemories() (string, string) {
	// Load visor introduction text
	var visor_intro string = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
	//var visor_functions = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "functions.json").ReadTextFile()
	//visor_intro = strings.Replace(visor_intro, "3234_FUNCTIONS", visor_functions, -1)
	visor_intro = strings.Replace(visor_intro, "\n", " ", -1)
	visor_intro = strings.Replace(visor_intro, "\"", "\\\"", -1)
	visor_intro = strings.Replace(visor_intro, "3234_NICK", modUserInfo_GL.User_nickname, -1)

	// Initialize memory string
	var visor_memories string = strings.Join(modGenInfo_GL.Memories, ". ")
	visor_memories = strings.Replace(visor_memories, "\"", "\\\"", -1)

	return visor_intro, visor_memories
}

func removeSessionEntry(session_id string) {
	delete(modGenInfo_GL.Sessions, session_id)
}

func getActiveSessionId() string {
	// The latest session with less than 30 minutes of inactivity is considered the active one
	var latest_interaction int64 = 0
	var active_session_id string = ""
	for id, session := range modGenInfo_GL.Sessions {
		if session.Last_interaction_s > latest_interaction &&
				time.Now().Unix() - session.Last_interaction_s < _INACTIVE_SESSION_TIME_S {
			active_session_id = id
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

			if strings.HasSuffix(text, "/stop") {
				_ = os.Remove(file_path.GPathToStringConversion())

				return true
			}
		}

		Utils.DelElemSLICES(&file_list, idx_to_remove)
	}

	return false
}

func startOllama() {
	_, _ = Utils.ExecCmdSHELL([]string{"sudo systemctl start ollama.service"})

}

/*
stopOllama stop the Ollama service.
*/
func stopOllama() {
	_, _ = Utils.ExecCmdSHELL([]string{"sudo systemctl stop ollama.service"})
}

/*
restartOllama restarts the Ollama service.
*/
func restartOllama() {
	_, _ = Utils.ExecCmdSHELL([]string{"sudo systemctl restart ollama.service"})
}

func getStartString(device_id string) string {
	return "[3234_START:" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + device_id + "|]"
}

func getEndString() string {
	return "[3234_END]\n"
}
