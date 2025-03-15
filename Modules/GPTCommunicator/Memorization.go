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
	"encoding/json"
	"log"
	"strings"
	"time"
)

func autoMemorize() {
	for {
		if modGenInfo_GL.State == ModsFileInfo.MOD_7_STATE_READY {
			for session_id, session := range modGenInfo_GL.Sessions {
				if session_id == getActiveSessionId() || session.Memorized || session_id == "temp" || session_id == "dumb" {
					continue
				}

				// If the session is no longer the active one, memorize it
				if memorizeSession(session_id) {
					session.Memorized = true
				}
			}

			//if len(modGenInfo_GL.Memories) >= modGenInfo_GL.N_mems_when_last_memorized * 2 {
			//	for !summarizeMemories() {
			//		// VISOR may not memorize because of for example romantic stuff being on the memories, or just
			//		// because they're of a user. In that case, just try again.
			//	}
			//
			//	modGenInfo_GL.N_mems_when_last_memorized = len(modGenInfo_GL.Memories)
			//}
		}

		if Utils.WaitWithStopTIMEDATE(module_stop_GL, 1*60) {
			return
		}
	}
}

func memorizeSession(session_id string) bool {
	var session_history []ModsFileInfo.OllamaMessage = Utils.CopyOuterSLICES(modGenInfo_GL.Sessions[session_id].History)
	for i, message := range session_history {
		if message.Role == "user" && !strings.Contains(message.Content, "[SYSTEM TASK - ") {
			// Remove the first part of the user message (like time and date and location, all inside square brackets)
			session_history[i].Content = message.Content[strings.Index(message.Content, "]") + 1:]
		} else {
			// Remove the system and assistant prompts, and the "system messages" from the user prompts
			session_history = append(session_history[:i], session_history[i+1:]...)
		}
	}
	session_history_json, err := json.Marshal(session_history)
	if err != nil || len(session_history) == 0 {
		log.Println("Error memorizing session " + session_id)

		return false
	}

	var prompt string = "User messages (in JSON): " + string(session_history_json) + ". Write NEW things you've " +
		"learned from this specific conversation (EXCLUDING YOUR MEMORIES) in BULLET points (no + or - or anything. " +
		"ONLY *). Format the output as \"* [detail]\". IGNORE specific, temporary events, scheedules, or day-to-day " +
		"plans. Summarize as KEY GENERAL information. If there is nothing, write \"* 3234_NONE\"."

	var response string = sendToGPT(Utils.Gen_settings_GL.Device_settings.Id, prompt, "temp")

	var lines []string = strings.Split(response, "\n")
	for _, line := range lines {
		if UtilsSWA.StringHasLettersGENERAL(line) && strings.Contains(line, "* ") && !strings.Contains(line, "3234_NONE") {
			line = strings.Replace(line, "\r ", "", -1)
			line = strings.Replace(line, "You ", "The user ", -1)
			line = strings.Replace(line, "He ", "The user ", -1)
			line = strings.Replace(line, "She ", "The user ", -1)
			line = strings.Replace(line, "They ", "The user ", -1)
			var the_user_idx int = strings.LastIndex(line, "* ")
			if the_user_idx == -1 {
				continue
			}

			modGenInfo_GL.Memories = append(modGenInfo_GL.Memories, line[the_user_idx + len("* "):])
		}
	}

	// Give time to write everything down
	time.Sleep(6 * time.Second)

	return true
}

func summarizeMemories() bool {
	var prompt string = "Summarize your memories about the user in BULLET points (no + or - or anything. ONLY *). " +
		"Format the output as \"* [detail]\". Write as much as you need. If newer memories contradict old " +
		"ones, update them. ALL MEMORIES ARE IMPORTANT, EVEN MINOR ONES!!! But again, SUMMARIZE them."

	var response string = sendToGPT(Utils.Gen_settings_GL.Device_settings.Id, prompt, "temp")

	var new_memories []string = nil
	var lines []string = strings.Split(response, "\n")
	for _, line := range lines {
		if UtilsSWA.StringHasLettersGENERAL(line) && strings.Contains(line, "* ") {
			line = strings.Replace(line, "\r", "", -1)
			line = strings.Replace(line, "You ", "The user", -1)
			line = strings.Replace(line, "He ", "The user", -1)
			line = strings.Replace(line, "She ", "The user", -1)
			line = strings.Replace(line, "They ", "The user", -1)
			var the_user_idx int = strings.LastIndex(line, "* ")
			if the_user_idx == -1 {
				continue
			}

			new_memories = append(new_memories, line[the_user_idx + len("* "):])
		}
	}

	if len(new_memories) == 0 {
		return false
	}

	modGenInfo_GL.Memories = new_memories

	// Give time to write everything down
	time.Sleep(6 * time.Second)

	return true
}
