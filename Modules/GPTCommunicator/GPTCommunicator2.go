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
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

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

// Prepared for Ollama 0.5.11

var memorizing_GL bool = false
var device_id_GL string = ""

var module_stop_GL *bool = nil

var visor_intro_GL string = ""
var visor_memories_GL string = ""

func readGPT(http_response *http.Response, print bool) (string, int64) {
	var message string = ""
	var last_word string = ""
	var save_words bool = true
	var curr_idx int = 0

	var timestamp_s int64 = -1

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY

	reduceGptTextTxt(gpt_text_txt)
	_ = gpt_text_txt.WriteTextFile(getStartString(device_id_GL), true)

	// Send a message to LIB_2 saying the GPT just started writing
	Utils.SendToModChannel(Utils.NUM_MOD_WebsiteBackend, 0, "Message", []byte(device_id_GL + "|L_2_0|start"))

	// Use a JSON decoder to handle the streamed response
	var decoder *json.Decoder = json.NewDecoder(http_response.Body)

	for {
		if *module_stop_GL || checkStopSpeech() {
			// Write the end string before exiting
			_ = gpt_text_txt.WriteTextFile(getEndString(), true)

			// Closing the connection makes Ollama stop generating the response
			http_response.Body.Close()

			break
		}

		if strings.Contains(message, "\000") {
			break
		}

		if checkStopSpeech() {

			break
		}

		var response _OllamaResponse
		if err := decoder.Decode(&response); err != nil {
			message += "\000"
		} else {
			message += response.Message.Content
		}

		if timestamp_s == -1 {
			timestamp_s = time.Now().Unix()
		}

		for {
			if curr_idx >= len(message) {
				break
			}

			var one_byte_str string = string(message[curr_idx])

			if memorizing_GL {
				//to_memorize += one_byte_str todo
			}
			if print {
				fmt.Print(one_byte_str)
			}

			if one_byte_str == " " || one_byte_str == "\n" || one_byte_str == "\000" {
				// Meaning: new word written
				if one_byte_str == "\000" {
					_ = gpt_text_txt.WriteTextFile(last_word, true)
				} else {
					_ = gpt_text_txt.WriteTextFile(last_word+one_byte_str, true)
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

			curr_idx++
		}
	}

	_ = gpt_text_txt.WriteTextFile(getEndString(), true)

	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_READY

	return message, timestamp_s
}

func chatWithGPT(user_message string, session_id string) string {
	if session_id == "" {
		// Get latest session ID if none is provided
		var latest_interaction int64 = -1
		for id, session := range modGenInfo_GL.Sessions {
			if session.Last_interaction_s > latest_interaction {
				session_id = id
				latest_interaction = session.Last_interaction_s
			}
		}
	}

	var curr_session *ModsFileInfo.Session = &ModsFileInfo.Session{}
	if session_id != "" {
		curr_session = modGenInfo_GL.Sessions[session_id]
	}

	// Append user message to history
	curr_session.History = append((*curr_session).History, ModsFileInfo.OllamaMessage{
		Role:    "user",
		Content: user_message,
		Timestamp_s: time.Now().Unix(),
	})

	var temperature float32 = 1.5
	if session_id != "dumb" {
		temperature = 0.8
	}

	// Create payload
	var ollama_request _OllamaRequest = _OllamaRequest{
		Model: "llama3.2:latest",
		Messages: curr_session.History,
		Options: _OllamaOptions{
			Num_ctx:     4096,
			Temperature: temperature,
		},
		Stream: true,
		Keep_alive: "99999999m",
	}

	jsonData, err := json.Marshal(ollama_request)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)

		return ""
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error posting to Ollama: ", err)

		return ""
	}
	defer resp.Body.Close()

	var response, timestamp = readGPT(resp, true)

	curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
		Role:    "assistant",
		Content: response[:len(response)-1],
		Timestamp_s: timestamp,
	})
	curr_session.Last_interaction_s = time.Now().Unix()

	if session_id != "dumb" {
		// Save the session unless it's to use the dumb LLM
		modGenInfo_GL.Sessions[session_id] = curr_session
	}

	return session_id
}

func sendToGPT(to_send string, use_smart bool) {
	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_BUSY

	var to_write string = UtilsSWA.RemoveNonGraphicCharsGENERAL(to_send)
	if use_smart {
		var system_prompt_smart string = modUserInfo_GL.System_info + "\n\n" + "Memories stored about the user: " +
			visor_memories_GL + "\n\n" + "About you: " + visor_intro_GL

		var session_id string = "test"
		if !addSessionEntry(session_id, []ModsFileInfo.OllamaMessage{
			{
				Role:    "system",
				Content: system_prompt_smart,
				Images:  nil,
				Timestamp_s: 0,
			},
		}, 0) {
			// If the session already exists, update the system message
			modGenInfo_GL.Sessions[session_id].History[0].Content = system_prompt_smart
		}

		chatWithGPT(to_write, session_id)
	} else {
		chatWithGPT(to_write, "dumb")
	}
}


// todo Edit each session to have an updated SYSTEM message having all user memories, each time it's reused


func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_7
		modUserInfo_GL = &Utils.User_settings_GL.GPTCommunicator

		module_stop_GL = module_stop


		// Set initial module state
		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

		// Start Ollama in case it's not running
		startOllama()

		// Load visor introduction text
		visor_intro_GL = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
		//var visor_functions = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "functions.json").ReadTextFile()
		//visor_intro = strings.Replace(visor_intro, "3234_FUNCTIONS", visor_functions, -1)
		visor_intro_GL = strings.Replace(visor_intro_GL, "\n", " ", -1)
		visor_intro_GL = strings.Replace(visor_intro_GL, "\"", "\\\"", -1)
		visor_intro_GL = strings.Replace(visor_intro_GL, "3234_NICK", modUserInfo_GL.User_nickname, -1)

		// Initialize memory string
		visor_memories_GL = strings.Join(modGenInfo_GL.Memories, ". ")
		visor_memories_GL = strings.Replace(visor_memories_GL, "\"", "\\\"", -1)

		var system_prompt_dumb string = modUserInfo_GL.System_info + "\n\n" + "You're a voice assistant"

		if modGenInfo_GL.Sessions == nil {
			modGenInfo_GL.Sessions = make(map[string]*ModsFileInfo.Session)
		}

		// Prepare the session for the dumb LLM
		removeSessionEntry("dumb")
		addSessionEntry("dumb", []ModsFileInfo.OllamaMessage{
			{
				Role:    "system",
				Content: system_prompt_dumb,
				Images:  nil,
				Timestamp_s: 0,
			},
		}, -1000000) // A very low timestamp to avoid being selected as latest session



		/*session_id := chatWithGPT("Who was the USA president in 2021?", Utils.RandStringGENERAL(10), true)

		_ = chatWithGPT("What's 2+2?", Utils.RandStringGENERAL(10), true)

		_ = chatWithGPT("And Portugal's?", session_id, true)

		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				return
			}
		}*/


		var shut_down bool = false
		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

		shutDown := func() {
			modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STOPPING
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
					device_id_GL = params_split[0]
					var use_smart bool = params_split[1] == "true"
					var text string = to_process[strings.Index(to_process, "]")+1:]
					text = strings.Replace(text, "\n", "\\n", -1)

					if use_smart {
						// Control commands begin with a slash
						if strings.HasSuffix(text, "/clear") {
							// Clear the context of the LLM model by stopping the module (the Manager will restart it)
							shut_down = true
							// todo No need to restart the module, just remove the session
						} else if strings.HasSuffix(text, "/restart") {
							restartOllama()
						} else if strings.Contains(text, ASK_WOLFRAM_ALPHA) {
							// Ask Wolfram Alpha the question
							var question string = text[strings.Index(text, ASK_WOLFRAM_ALPHA)+len(ASK_WOLFRAM_ALPHA):]
							result, direct_result := OnlineInfoChk.RetrieveWolframAlpha(question)

							if direct_result {
								_ = gpt_text_txt.WriteTextFile(getStartString(device_id_GL) + "The answer is: " + result +
									". " + getEndString(), true)
							} else {
								sendToGPT("Summarize in sentences the following: " + result, false)
							}
						} else if strings.Contains(text, SEARCH_WIKIPEDIA) {
							// Search for the Wikipedia page title
							var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

							_ = gpt_text_txt.WriteTextFile(getStartString(device_id_GL) + OnlineInfoChk.RetrieveWikipedia(query) +
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

func addSessionEntry(session_id string, history []ModsFileInfo.OllamaMessage, last_interaction_s int64) bool {
	if _, ok := modGenInfo_GL.Sessions[session_id]; !ok {
		// If the session doesn't exist, create it
		modGenInfo_GL.Sessions[session_id] = &ModsFileInfo.Session{
			Name:               session_id,
			Created_time_s:     time.Now().Unix(),
			History:            history,
			Last_interaction_s: last_interaction_s,
		}

		return true
	}

	return false
}

func removeSessionEntry(session_id string) {
	delete(modGenInfo_GL.Sessions, session_id)
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
