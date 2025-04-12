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

// _ChatWithGPTParams is a struct containing the parameters for the chatWithGPT function.
type _ChatWithGPTParams struct {
	// Device_id is the ID of the device sending the text
	Device_id string
	// Session_id is the session type to use
	Session_id string
	// User_message is the message to send to the model
	User_message string
	// Role is the role of the message
	Role string
	// More_coming is true if more messages are coming and VISOR should wait before calling the LLM
	More_coming bool
	// Model_type is the type of model to use
	Model_type string
}

func chatWithGPT(params _ChatWithGPTParams) string {
	setBusyState()
	defer setReadyState()

	if params.Session_id == "" {
		// Get latest session ID if none is provided
		var latest_interaction int64 = -1
		for _, session := range getModGenSettings().Sessions {
			if session.Last_interaction_s > latest_interaction {
				params.Session_id = session.Id
				latest_interaction = session.Last_interaction_s
			}
		}
	}

	addSessionEntry(params.Session_id, time.Now().Unix(), params.User_message)
	setBusyState() // Again because chatWithGPT() is called inside addSessionEntry() and will set it to Ready on ending

	var curr_session ModsFileInfo.Session = *getSession(params.Session_id)
	curr_session.Memorized = false

	switch params.Role {
		case GPTComm.ROLE_USER:
			params.Role = "user"
		case GPTComm.ROLE_TOOL:
			params.User_message = "As per system request, inform the user that: \"" + params.User_message + "\"."
			if getModUserInfo().Model_has_tool_role {
				params.Role = "tool"
			} else {
				params.Role = "user"

				// Keep the last part. He'll say less random stuff this way it seems.
				params.User_message = "[SYSTEM TASK - " + params.User_message + " NO SAYING YOU'RE REWORDING IT]"
			}
		default:
			// Keep the original
	}

	// Append user message to history
	curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
		Role:        params.Role,
		Content:     UtilsSWA.RemoveNonGraphicCharsGENERAL(params.User_message),
		Timestamp_s: time.Now().Unix(),
	})

	var response string = ""
	if !params.More_coming {
		var history_with_system_prompt []ModsFileInfo.OllamaMessage = Utils.CopyOuterSLICES(curr_session.History)
		var system_prompt string = ""
		// Add the system prompt every time *temporarily*, so that if it's updated, it's updated in all sessions when
		// they're used - because it's not stored in any session.
		if params.Session_id == "dumb" {
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

		model_name, device_id_with_model := getModelName(params.Model_type)
		if model_name == "" {
			return ""
		}

		// Create payload
		var ollama_request ModsFileInfo.OllamaChatRequest = ModsFileInfo.OllamaChatRequest{
			Model:    model_name,
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

		request_json, err := json.Marshal(ollama_request)
		if err != nil {
			log.Println("Error marshalling JSON: ", err)

			return ""
		}

		response_str, timestamp := sendReceiveOllamaRequest(params.Device_id, request_json, device_id_with_model)
		if timestamp == -1 {
			return ""
		}

		response = response_str
		if response != "" {
			response = response[:len(response)-1] // Remove the last character, which is a null character
		}

		if params.Session_id != "temp" && params.Session_id != "dumb" {
			curr_session.History = append(curr_session.History, ModsFileInfo.OllamaMessage{
				Role:        "assistant",
				Content:     response,
				Timestamp_s: timestamp,
			})
			curr_session.Last_interaction_s = time.Now().Unix()
		}
	}

	if params.Session_id != "temp" && params.Session_id != "dumb" {
		// Save the session unless it's to use the temp or dumb sessions
		for i, session := range getModGenSettings().Sessions {
			if session.Id == params.Session_id {
				getModGenSettings().Sessions[i] = curr_session

				break
			}
		}
	}

	return response
}

func sendReceiveOllamaRequest(device_id string, request_json []byte, device_id_with_model string) (string, int64) {
	if device_id_with_model == Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id {
		log.Println("Posting to Ollama locally: ", string(request_json))

		resp, err := http.Post("http://localhost:11434/api/chat", "application/json; charset=utf-8",
			bytes.NewBuffer(request_json))
		if err != nil {
			log.Println("Error posting to Ollama: ", err)

			// Wait 2 seconds before stopping the module for the clients to receive the STARTING state before the
			// STOPPING one (they check every second).
			time.Sleep(2 * time.Second)

			// Ollama stopped running, so stop the module
			*module_stop_GL = true

			return "", -1
		}
		defer resp.Body.Close()

		return readGPT(device_id, resp, true)
	} else {
		log.Println("Posting to Ollama on \"" + device_id_with_model + "\": ", string(request_json))

		Utils.QueueMessageBACKEND(true, Utils.NUM_MOD_GPTCommunicator, 0, device_id_with_model,
			[]byte(device_id + "|" + string(request_json)))
	}

	return "", -1
}

func readGPT(device_id string, http_response *http.Response, print bool) (string, int64) {
	var timestamp_s int64 = -1

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	// Default is false because the client can generate text to itself (if it's true it will ignore the text - not
	// good), but it goes to the server before returning to the client.
	var writing_to_self bool = false
	if Utils.VISOR_server_GL {
		writing_to_self = device_id == Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id
		if !writing_to_self {
			reduceGptTextTxt(gpt_text_txt)
			_ = gpt_text_txt.WriteTextFile(getStartString(device_id), true)
		}

		// Send a message to LIB_2 saying the GPT just started writing
		Utils.QueueMessageBACKEND(false, Utils.NUM_LIB_GPTComm, 0, device_id, []byte("start"))
	} else {
		sendToServer(getStartString(device_id))
	}

	// Use a JSON decoder to handle the streamed response
	var decoder *json.Decoder = json.NewDecoder(http_response.Body)

	var message string = ""
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
			var content string = response.Message.Content
			message += content
			if !writing_to_self {
				sendWriteText(content)
			}
			if print {
				fmt.Print(content)
			}
		} else {
			break
		}

		if timestamp_s == -1 {
			timestamp_s = time.Now().Unix()
		}
	}

	if !writing_to_self {
		sendWriteText(getEndString())
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

func sendWriteText(text string) {
	if Utils.VISOR_server_GL {
		var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")
		_ = gpt_text_txt.WriteTextFile(text, true)
	} else {
		sendToServer(text)
	}
}

func sendToServer(text string) {
	var message []byte = []byte("GPT|[redirect]")
	message = append(message, text...)
	if !Utils.QueueNoResponseMessageSERVER(message) {
		log.Println("Error sending message to server: ", text)
	}
}
