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
	"bytes"
	"strconv"
	"strings"
	"time"
)

const _TO_PROCESS_REL_FOLDER string = "to_process"

const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "

const STOP_CMD string = "/stop"

const _INACTIVE_SESSION_TIME_S int64 = 30*60

func serverMode() {
	getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STARTING

	if len(getModUserInfo().Models) == 0 {
		// This is here to signal users that they need to do something to make the module work (add models).
		time.Sleep(2 * time.Second)

		getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STOPPED

		return
	}

	if getModGenSettings().N_mems_when_last_memorized == 0 {
		getModGenSettings().N_mems_when_last_memorized = 25 // So that the double is 50 for the first time
	}

	// Prepare the session for the temp and dumb sessions
	// Very low timestamp to avoid being selected as latest sessions
	addSessionEntry("temp", -1000000, "", nil)
	addSessionEntry("dumb", -1000000, "", nil)

	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")

	// In case Ollama was started (as opposed to already being running), send a test message for it to actually start
	// and be ready.
	var chatWithGPT_params _ChatWithGPTParams = _ChatWithGPTParams{
		Device_id:    Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id,
		User_message: "test",
		Session_id:   "temp",
		Role:         GPTComm.ROLE_USER,
		More_coming:  false,
	}
	chatWithGPT(chatWithGPT_params)

	//go autoMemorize() TODO

	// Process the text to input to the LLM model
	for {
		var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 2, -1)
		if comms_map == nil {
			break
		}

		var to_process []byte = comms_map["Process"].([]byte)
		var to_process_str string = string(to_process)

		chatWithGPT_params = _ChatWithGPTParams{}

		// It comes like: "[number of images|number of audios|size of first file|size of second file|...|device ID|
		// session type|role|true/false|model type]text + '\x00' [+ each file bytes one after the other - optional]"
		var idx_first_closing_bracket int = strings.Index(to_process_str, "]")
		var params_split []string = strings.Split(to_process_str[1:idx_first_closing_bracket], "|")

		// Get number of images and audios and their sizes
		n_images, _ := strconv.Atoi(params_split[0])
		n_audios, _ := strconv.Atoi(params_split[1])
		var size_images []int = nil
		var size_audios []int = nil
		for i := 0; i < n_images; i++ {
			size, _ := strconv.Atoi(params_split[2+i])
			size_images = append(size_images, size)
		}
		for i := 0; i < n_audios; i++ {
			size, _ := strconv.Atoi(params_split[2+n_images+i])
			size_audios = append(size_audios, size)
		}

		var idx_begin_other_params int = 2 + n_images + n_audios

		var device_id string = params_split[idx_begin_other_params+0]
		chatWithGPT_params.Device_id = device_id

		var session_type string = params_split[idx_begin_other_params+1]
		var session_id string = ""
		switch session_type {
		case GPTComm.SESSION_TYPE_NEW:
			session_id = Utils.RandStringGENERAL(50)
		case GPTComm.SESSION_TYPE_TEMP:
			session_id = "temp"
		case GPTComm.SESSION_TYPE_ACTIVE:
			session_id = getActiveSessionId()
			if session_id == "" {
				session_id = Utils.RandStringGENERAL(50)
			}
		default:
			session_id = session_type
		}
		chatWithGPT_params.Session_id = session_id

		chatWithGPT_params.Role = params_split[idx_begin_other_params+2]
		chatWithGPT_params.More_coming = params_split[idx_begin_other_params+3] == "true"

		var idx_null int = bytes.Index(to_process, []byte{0})
		var text string = to_process_str[idx_first_closing_bracket+1 : idx_null]
		chatWithGPT_params.User_message = text

		// Get the images and audios
		var files []byte = to_process[idx_null+1:]
		var files_param []GPTComm.File = nil
		for i := 0; i < n_images; i++ {
			files_param = append(files_param, GPTComm.File{
				Is_image: true,
				Size:     size_images[i],
				Contents: files[:size_images[i]],
			})
			files = files[size_images[i]:]
		}
		for i := 0; i < n_audios; i++ {
			files_param = append(files_param, GPTComm.File{
				Is_image: false,
				Size:     size_audios[i],
				Contents: files[:size_audios[i]],
			})
			files = files[size_audios[i]:]
		}

		chatWithGPT_params.Files = files_param

		// Control commands begin with a slash
		if strings.HasSuffix(text, ASK_WOLFRAM_ALPHA) {
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
		} else if strings.HasSuffix(text, SEARCH_WIKIPEDIA) {
			// Search for the Wikipedia page title
			var query string = text[strings.Index(text, SEARCH_WIKIPEDIA)+len(SEARCH_WIKIPEDIA):]

			_ = gpt_text_txt.WriteTextFile(getStartString(device_id)+OnlineInfoChk.RetrieveWikipedia(query)+
				getEndString(), true)
		} else if !strings.HasSuffix(text, STOP_CMD) {
			chatWithGPT(chatWithGPT_params)
		}
	}

	getModGenSettings().State = ModsFileInfo.MOD_7_STATE_STOPPED
}

func addSessionEntry(session_id string, last_interaction_s int64, user_message string, files []GPTComm.File) {
	var session_exists bool = false
	for id := range getModGenSettings().Sessions {
		if id == session_id {
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
			var prompt string = "Create a title for the following text and possible files (beginning of a " +
				"conversation) and put it inside \"double quotation marks\", please. Don't include the date and time. " +
				"Text --> " + message_without_add_info + " <--"
			var chatWithGPT_params _ChatWithGPTParams = _ChatWithGPTParams{
				Device_id:    Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id,
				User_message: prompt,
				Session_id:   "temp",
				Role:         GPTComm.ROLE_USER,
				More_coming:  false,
				Files:        files,
			}
			session_name = chatWithGPT(chatWithGPT_params)
			var delimiter_start string = ""
			var delimiter_end string = ""
			if strings.Contains(session_name, "\"") {
				delimiter_start = "\""
				delimiter_end = "\""
			} else if strings.Contains(session_name, "“") && strings.Contains(session_name, "”") {
				delimiter_start = "“"
				delimiter_end = "”"
			}
			if delimiter_start == "" {
				session_name = "[Error naming the session]"
			} else {
				var idx_start int = strings.Index(session_name, delimiter_start)
				var idx_end int = strings.Index(session_name[idx_start+len(delimiter_start):], delimiter_end)
				if idx_start != -1 && idx_end != -1 {
					session_name = session_name[idx_start+len(delimiter_start) : idx_end+1]
				} else {
					session_name = "[Error naming the session]"
				}
				// Sometimes the name may come like "[name here]", so remove the brackets.
				session_name = strings.Replace(session_name, "[", "", -1)
				session_name = strings.Replace(session_name, "]", "", -1)
				session_name = strings.TrimSpace(session_name)
			}
		}

		getModGenSettings().Sessions[session_id] = &ModsFileInfo.Session{
			Name:               session_name,
			Created_time_s:     time.Now().Unix(),
			History:            nil,
			Last_interaction_s: last_interaction_s,
			Memorized:          false,
		}
	}
}

/*
checkStopSpeech checks if the text to process contains the /stop command.

-----------------------------------------------------------

– Returns:
  - true if the /stop command was found, false otherwise
*/
func checkStopSpeechServer() bool {
	if !Utils.VISOR_server_GL {
		return false
	}

	var stop bool = false

	// Retrieve all values from the channel
	var temp_values []any
	for {
		var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 2, 0)
		if comms_map == nil {
			break
		}

		var map_value []byte = comms_map["Process"].([]byte)
		var to_process_str string = string(map_value)
		if strings.HasSuffix(to_process_str[:len(to_process_str)-1], STOP_CMD) {
			stop = true

			continue
		}

		temp_values = append(temp_values, map_value)
	}

	// Put all values back on the channel except the stop command
	for _, value := range temp_values {
		Utils.SendToModChannel(Utils.NUM_MOD_GPTCommunicator, 2, "Process", value)
	}

	return stop
}

func getActiveSessionId() string {
	// The latest session with less than 30 minutes of inactivity is considered the active one
	var latest_interaction int64 = 0
	var active_session_id string = ""
	for session_id, session := range getModGenSettings().Sessions {
		if session.Last_interaction_s > latest_interaction &&
			time.Now().Unix() - session.Last_interaction_s < _INACTIVE_SESSION_TIME_S {
			active_session_id = session_id
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
func reduceGptTextTxt() {
	var gpt_text_txt Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "gpt_text.txt")
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
