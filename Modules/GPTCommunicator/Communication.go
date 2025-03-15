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

func chatWithGPT(device_id string, user_message string, session_id string) string {
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

	addSessionEntry(session_id, nil, 0, user_message)

	var curr_session ModsFileInfo.Session = *modGenInfo_GL.Sessions[session_id]
	curr_session.Memorized = false

	// Append user message to history
	curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
		Role:    "user",
		Content: UtilsSWA.RemoveNonGraphicCharsGENERAL(user_message),
		Timestamp_s: time.Now().Unix(),
	})

	var history_with_system_prompt []ModsFileInfo.OllamaMessage = Utils.CopyOuterSLICES(curr_session.History)
	var system_prompt string = ""
	// Add the system prompt every time *temporarily*, so that if it's updated, it's updated in all sessions when
	// they're used - because it's not stored in any session.
	if session_id == "dumb" {
		system_prompt = modUserInfo_GL.System_info + "\n\n" + "You're a voice assistant"
	} else {
		var visor_intro, visor_memories string = getVisorIntroAndMemories()
		system_prompt = modUserInfo_GL.System_info + "\n\n" + "Long-term memories stored about the user: " +
			visor_memories + "\n\n" + "About you: " + visor_intro
	}
	Utils.AddElemSLICES(&history_with_system_prompt, ModsFileInfo.OllamaMessage{
		Role:        "system",
		Content:     system_prompt,
		Images:      nil,
		Timestamp_s: 0,
	}, 0)

	// Create payload
	var ollama_request _OllamaRequest = _OllamaRequest{
		Model:    modUserInfo_GL.Model_name,
		Messages: history_with_system_prompt,
		Options: _OllamaOptions{
			Num_keep:    9999999,
			Num_ctx:     4096,
			Temperature: 0.8,
		},
		Stream: true,
		Keep_alive: "99999999m",
	}

	jsonData, err := json.Marshal(ollama_request)
	if err != nil {
		log.Println("Error marshalling JSON: ", err)

		return ""
	}

	log.Println("Posting to Ollama: ", string(jsonData))

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json; charset=utf-8", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error posting to Ollama: ", err)

		// Ollama stopped running, so stop the module
		*module_stop_GL = true

		return ""
	}
	defer resp.Body.Close()

	var response, timestamp = readGPT(device_id, resp, true)
	response = response[:len(response)-1] // Remove the last character, which is a null character

	if session_id != "temp" && session_id != "dumb" {
		curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
			Role:    "assistant",
			Content: response,
			Timestamp_s: timestamp,
		})
		curr_session.Last_interaction_s = time.Now().Unix()

		// Save the session unless it's to use the temp or dumb sessions
		modGenInfo_GL.Sessions[session_id] = &curr_session
	}

	return response
}

func readGPT(device_id string, http_response *http.Response, print bool) (string, int64) {
	var message string = ""
	var last_word string = ""
	var save_words bool = true
	var curr_idx int = 0

	var writing_to_self bool = device_id == Utils.Gen_settings_GL.Device_settings.Id

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

	for {
		if *module_stop_GL || checkStopSpeech() {
			// Write the end string before exiting
			if !writing_to_self {
				_ = gpt_text_txt.WriteTextFile(getEndString(), true)
			}

			// Closing the connection makes Ollama stop generating the response
			http_response.Body.Close()

			break
		}

		if strings.Contains(message, "\000") {
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

			if print {
				fmt.Print(one_byte_str)
			}

			if one_byte_str == " " || one_byte_str == "\n" || one_byte_str == "\000" {
				if !writing_to_self {
					// Meaning: new word written
					if one_byte_str == "\000" {
						_ = gpt_text_txt.WriteTextFile(last_word, true)
					} else {
						_ = gpt_text_txt.WriteTextFile(last_word+one_byte_str, true)
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

	modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_READY

	return message, timestamp_s
}
