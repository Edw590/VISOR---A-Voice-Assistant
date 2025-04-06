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
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func chatWithGPT(device_id string, user_message string, session_id string, role string, more_coming bool) string {
	setBusyState()
	defer setReadyState()

	if session_id == "" {
		// Get latest session ID if none is provided
		var latest_interaction int64 = -1
		for _, session := range getModGenSettings().Sessions {
			if session.Last_interaction_s > latest_interaction {
				session_id = session.Id
				latest_interaction = session.Last_interaction_s
			}
		}
	}

	addSessionEntry(session_id, time.Now().Unix(), user_message)
	setBusyState() // Again because chatWithGPT() is called inside addSessionEntry() and will set it to Ready on ending

	var curr_session ModsFileInfo.Session = *getSession(session_id)
	curr_session.Memorized = false

	var actual_role string = ""
	switch role {
		case GPTComm.ROLE_USER:
			actual_role = "user"
		case GPTComm.ROLE_TOOL:
			user_message = "As per user request, inform them that: \"" + user_message + "\"."
			if getModUserInfo().Model_has_tool_role {
				actual_role = "tool"
			} else {
				actual_role = "user"

				// Keep the last part. He'll say less random stuff this way.
				user_message = "[SYSTEM TASK - " + user_message + " NO SAYING YOU'RE REWORDING IT]"
			}
		default:
			actual_role = role
	}

	// Append user message to history
	curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
		Role:        actual_role,
		Content:     UtilsSWA.RemoveNonGraphicCharsGENERAL(user_message),
		Timestamp_s: time.Now().Unix(),
	})

	var response string = ""
	if !more_coming {
		var history_with_system_prompt []ModsFileInfo.OllamaMessage = Utils.CopyOuterSLICES(curr_session.History)
		var system_prompt string = ""
		// Add the system prompt every time *temporarily*, so that if it's updated, it's updated in all sessions when
		// they're used - because it's not stored in any session.
		if session_id == "dumb" {
			system_prompt = getModUserInfo().System_info + "\n\n" + "You're a voice assistant"
		} else {
			var visor_intro, visor_memories string = getVisorIntroAndMemories()
			system_prompt = getModUserInfo().System_info + "\n\n" + "Long-term memories stored about the user: " +
				visor_memories + "\n\n" + "About you: " + visor_intro
		}
		Utils.AddElemSLICES(&history_with_system_prompt, ModsFileInfo.OllamaMessage{
			Role:        "system",
			Content:     system_prompt,
			Images:      nil,
			Timestamp_s: 0,
		}, 0)

		// Ready to function, but when creating a title for "What's the battery percentage", he called the function to
		// get the battery percentage - not too useful. So the code is disabled (for now?).
		//var tools_json []byte = modDirsInfo_GL.ProgramData.Add2(false, "tools.json").ReadFile()
		//var tools ModsFileInfo.OllamaTools = nil
		//if tools_json != nil {
		//	err := Utils.FromJsonGENERAL(tools_json, &tools)
		//	if err != nil {
		//		log.Println("Error unmarshalling tools JSON: ", err)
		//
		//		return ""
		//	}
		//}

		// Create payload
		var ollama_request ModsFileInfo.OllamaChatRequest = ModsFileInfo.OllamaChatRequest{
			Model:    getModUserInfo().Model_name,
			Messages: history_with_system_prompt,
			Options: ModsFileInfo.OllamaOptions{
				Num_keep:    99999999,
				Num_ctx:     getModUserInfo().Context_size,
				Temperature: getModUserInfo().Temperature,
			},
			Stream:     true,
			Keep_alive: "9999m",
			//Tools: tools,
		}

		jsonData, err := json.Marshal(ollama_request)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)

			return ""
		}

		log.Println("Posting to Ollama: ", string(jsonData))

		resp, err := http.Post("http://" + getModUserInfo().Server_url + "/api/chat", "application/json; charset=utf-8",
			bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error posting to Ollama: ", err)

			// Wait 2 seconds before stopping the module for the clients to receive the STARTING state before the
			// STOPPING one (they check every second).
			time.Sleep(2 * time.Second)

			// Ollama stopped running, so stop the module
			*module_stop_GL = true

			return ""
		}
		defer resp.Body.Close()

		response_str, timestamp := readGPT(device_id, resp, true)
		response = response_str
		if response != "" {
			response = response[:len(response)-1] // Remove the last character, which is a null character
		}

		if session_id != "temp" && session_id != "dumb" {
			curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
				Role:        "assistant",
				Content:     response,
				Timestamp_s: timestamp,
			})
			curr_session.Last_interaction_s = time.Now().Unix()
		}
	}

	if session_id != "temp" && session_id != "dumb" {
		// Save the session unless it's to use the temp or dumb sessions
		for i, session := range getModGenSettings().Sessions {
			if session.Id == session_id {
				getModGenSettings().Sessions[i] = curr_session

				break
			}
		}
	}

	return response
}

func readGPT(device_id string, http_response *http.Response, print bool) (string, int64) {
	var writing_to_self bool = device_id == Utils.GetGenSettings().Device_settings.Id

	var timestamp_s int64 = -1

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	if !writing_to_self {
		reduceGptTextTxt(gpt_text_txt)
		_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)
	}

	// Send a message to LIB_2 saying the GPT just started writing
	Utils.SendToModChannel(Utils.NUM_MOD_WebsiteBackend, 0, "Message", []byte(device_id + "|L_2_0|start"))

	// Use a JSON decoder to handle the streamed response
	var decoder *json.Decoder = json.NewDecoder(http_response.Body)

	var message string = ""
	var last_word string = ""
	var save_words bool = true
	var curr_idx int = 0
	for {
		if *module_stop_GL || checkStopSpeech() {
			// Closing the connection makes Ollama stop generating the response
			http_response.Body.Close()

			break
		}

		if strings.Contains(message, "\000") {
			break
		}

		var response ModsFileInfo.OllamaChatResponse
		if err := decoder.Decode(&response); err == nil {
			message += response.Message.Content
		} else {
			message += "\000"
		}

		if timestamp_s == -1 {
			timestamp_s = time.Now().Unix()
		}

		for {
			if curr_idx >= len(message) {
				break
			}

			var one_byte_str string = string(message[curr_idx])

			if print {
				fmt.Print(one_byte_str)
			}

			if one_byte_str == " " || one_byte_str == "\n" || one_byte_str == "\000" {
				if !writing_to_self {
					// Meaning: new word written
					if one_byte_str == "\000" {
						_ = gpt_text_txt.WriteTextFile(last_word, true)
					} else {
						_ = gpt_text_txt.WriteTextFile(last_word + one_byte_str, true)
					}
				}

				last_word = ""
			} else {
				// VISOR may start by writing the current date and time like "[date and time here]" - this below cuts
				// that out of the answer.
				//if last_word == "" {
				//	if one_byte_str == "[" {
				//		save_words = false
				//	} else if one_byte_str == "]" {
				//		save_words = true
				//
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

	if !writing_to_self {
		_ = gpt_text_txt.WriteTextFile("\n" + getEndString(), true)
	}

	return message, timestamp_s
}

func getVisorIntroAndMemories() (string, string) {
	// Load visor introduction text
	var visor_intro string = *modDirsInfo_GL.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
	visor_intro = strings.Replace(visor_intro, "3234_NICK", getModUserInfo().User_nickname, -1)
	if !getModUserInfo().Model_has_tool_role {
		// If the model tool role is not set, the user one will be used instead - but in that case VISOR has to
		// differentiate from the actual user input. So "SYSTEM TASK"s are used.
		visor_intro = strings.Replace(visor_intro, "3234_SYS_TASKS", "Sometimes there will be \"SYSTEM TASK\"s. "+
			"These are tasks that the system has set for you to do. You must do as written.", -1)
	} else {
		visor_intro = strings.Replace(visor_intro, "3234_SYS_TASKS\n", "", -1)
	}

	// Initialize memory string
	var visor_memories string = strings.Join(getModGenSettings().Memories, "\n")

	return visor_intro, visor_memories
}

func getSession(session_id string) *ModsFileInfo.Session {
	for i, session := range getModGenSettings().Sessions {
		if session.Id == session_id {
			return &getModGenSettings().Sessions[i]
		}
	}

	getModGenSettings().Sessions = append(getModGenSettings().Sessions, ModsFileInfo.Session{
		Id: session_id,
	})

	return &getModGenSettings().Sessions[len(getModGenSettings().Sessions)-1]
}
