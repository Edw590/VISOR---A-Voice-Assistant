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

package CmdsExecutor

import (
	"ACD/ACD"
	"GMan"
	"GPTComm"
	"OICComm"
	"Speech"
	"SpeechQueue"
	"TEHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var last_it string = ""
var last_it_when int64 = 0
var last_and string = ""
var last_and_when int64 = 0

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	ACD.ReloadCmdsArray(prepareCommandsString())

	for {
		var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_CmdsExecutor, 0)
		if comms_map == nil {
			return
		}

		if time.Now().UnixMilli() > last_it_when + 60*1000 {
			last_it = ""
		}
		if time.Now().UnixMilli() > last_and_when + 60*1000 {
			last_and = ""
		}

		var speech_priority int32 = SpeechQueue.PRIORITY_MEDIUM
		var sentence_str string = ""
		if map_value, ok := comms_map["SentenceInternal"]; ok {
			sentence_str = map_value.(string)
		} else if map_value, ok = comms_map["Sentence"]; ok {
			speech_priority = SpeechQueue.PRIORITY_USER_ACTION
			sentence_str = map_value.(string)
		} else {
			continue
		}
		var cmds_info_str = ACD.Main(strings.ToLower(sentence_str), false, true, last_it + "|" + last_and)
		log.Println("*****************************")
		fmt.Println(sentence_str)
		log.Println(cmds_info_str)
		var cmds_info []string = strings.Split(cmds_info_str, ACD.INFO_CMDS_SEPARATOR)
		if len(cmds_info) < 2 {
			sendToGPT(sentence_str)

			time.Sleep(1 * time.Second)

			continue
		}
		var prev_cmd_info []string = strings.Split(cmds_info[0], ACD.PREV_CMD_INFO_SEPARATOR)
		var detected_cmds []string = strings.Split(cmds_info[1], ACD.CMDS_SEPARATOR)

		log.Println(last_it)
		log.Println(last_and)
		log.Println("***************")

		if prev_cmd_info[0] != "" {
			last_it = prev_cmd_info[0]
			last_it_when = time.Now().UnixMilli()
		}
		if prev_cmd_info[1] != "" {
			last_and = prev_cmd_info[1]
			last_and_when = time.Now().UnixMilli()
		}

		log.Println(last_it)
		log.Println(last_and)
		log.Println("*****************************")

		var send_to_GPT bool = false
		if strings.HasPrefix(cmds_info_str, ACD.ERR_CMD_DETECT) {
			var speak string = "WARNING! There was a problem processing the commands Sir. This needs a fix. The " +
				"error was the following: " + cmds_info_str + ". You said: " + sentence_str
			Speech.QueueSpeech(speak, speech_priority, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
			log.Println("EXECUTOR - ERR_PROC_CMDS")

			send_to_GPT = true
		}

		if len(detected_cmds) == 0 || detected_cmds[0] == "" {
			send_to_GPT = true
		} else {
			send_to_GPT = true
			for _, command := range detected_cmds {
				num, _ := strconv.ParseFloat(command, 32)
				if num >= 1 {
					// If there's any command detected, don't send to GPT
					send_to_GPT = false

					break
				}
			}
			// If there are only WARN_-started constants (negative numbers), send to GPT
		}
		if send_to_GPT {
			sendToGPT(sentence_str)

			return
		}

		for _, command := range detected_cmds {
			var dot_index int = strings.Index(command, ".")
			if dot_index == -1 {
				// No command.
				continue
			}

			var cmd_id string = command[:dot_index] // "14.3" --> "14"
			var cmd_variant string = command[dot_index:] // "14.3" --> ".3"

			var speech_mode2 = SpeechQueue.MODE_DEFAULT
			if cmdi_info[cmd_id] == CMDi_INF1_ONLY_SPEAK {
				speech_mode2 = SpeechQueue.MODE2_BYPASS_NO_SND
			}

			switch cmd_id {
				case CMD_ASK_TIME:
					var speak string = "It's " + Utils.GetTimeStrTIMEDATE(-1)
					speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false)

				case CMD_ASK_DATE:
					var speak string = "Today's " + Utils.GetDateStrDATETIME(-1)
					speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false)

				case CMD_TOGGLE_WIFI:
					if Utils.ToggleWifiCONNECTIVITY(cmd_variant == RET_ON) {
						var speak string
						if cmd_variant == RET_ON {
							speak = "Wi-Fi turned on."
						} else {
							speak = "Wi-Fi turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					} else {
						var on_off string = "off"
						if cmd_variant == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the Wi-Fi " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}

				case CMD_ASK_BATTERY_PERCENT:
					var battery_percentage int = int(UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).
						GetInt(true))
					var speak string = "Battery percentage: " + strconv.Itoa(battery_percentage) + "%"
					speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)

				case CMD_TELL_WEATHER:
					var speak string = "Obtaining the weather..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var weather_locs []string = strings.Split(OICComm.GetWeatherLocationsList(), "|")

						for _, weather_loc := range weather_locs {
							var weather *ModsFileInfo.Weather = OICComm.GetWeather(weather_loc)
							if weather == nil {
								speak = "There is no weather data associated with the location " + weather_loc + "."
								speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)

								continue
							}

							if weather.Temperature == "" {
								// One being empty means the whole weather is empty
								speak = "There was a problem obtaining the weather for " + weather.Location + "."
								speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)

								continue
							}

							var status_part string = " is "
							if weather.Status != "ERROR" {
								status_part += weather.Status + " with "
							}

							speak = "The weather in " + weather.Location + status_part + weather.Temperature +
								" degrees, a high of " + weather.Max_temp + " degrees and a low of " + weather.Min_temp +
								" degrees. The mean precipitation is of " + weather.Precipitation + ", mean " +
								"humidity of " + weather.Humidity + ", and mean wind of " + weather.Wind + "."
							speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false)
						}
					} else {
						speak = "Not connected to the server to get the weather."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}

				case CMD_TELL_NEWS:
					var speak string = "Obtaining the latest news..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var news_locs []string = strings.Split(OICComm.GetNewsLocationsList(), "|")

						for _, news_loc := range news_locs {
							var news *ModsFileInfo.News = OICComm.GetNews(news_loc)

							speak = "News in " + news.Location + ". "

							for _, n := range news.News {
								speak += n + ". "
							}
							speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)
						}
					} else {
						speak = "Not connected to the server to get the news."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}

				case CMD_TOGGLE_ETHERNET:
					if Utils.ToggleEthernetCONNECTIVITY(cmd_variant == RET_ON) {
						var speak string
						if cmd_variant == RET_ON {
							speak = "Ethernet turned on."
						} else {
							speak = "Ethernet turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)
					} else {
						var on_off string = "off"
						if cmd_variant == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the Ethernet " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}

				case CMD_TOGGLE_NETWORKING:
					if Utils.ToggleNetworkingCONNECTIVITY(cmd_variant == RET_ON) {
						var speak string
						if cmd_variant == RET_ON {
							speak = "Networking turned on."
						} else {
							speak = "Networking turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)
					} else {
						var on_off string = "off"
						if cmd_variant == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the networking " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}

				case CMD_ASK_EVENTS:
					var speak string = "Obtaining the tasks and events..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var events_ids []string = strings.Split(GMan.GetEventsIdsList(), "|")
						var tasks_ids []string = strings.Split(GMan.GetTasksIdsList(), "|")

						speak = getEventsList(events_ids, cmd_variant)

						if cmd_variant == RET_31_TODAY || cmd_variant == RET_31_TOMORROW {
							speak += " " + getTasksList(tasks_ids, cmd_variant)
						}

						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, true)
					} else {
						speak = "Not connected to the server to get the tasks and events."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false)
					}
			}
		}


		if Utils.WaitWithStopDATETIME(module_stop, 1) {
			TEHelper.StopChecker()

			return
		}
	}
}

const _SESSION_TYPE_NONE string = "NONE"
func speakInternal(txt_to_speak string, speech_priority int32, mode int32, session_type string, wait_for_gpt bool) {
	if session_type != _SESSION_TYPE_NONE && speech_priority <= SpeechQueue.PRIORITY_USER_ACTION &&
				Utils.IsCommunicatorConnectedSERVER() && (wait_for_gpt ||
				GPTComm.SendText("", "", "", false) == ModsFileInfo.MOD_7_STATE_READY) {
		var speak string = ""
		switch GPTComm.SendText(txt_to_speak, session_type, GPTComm.ROLE_TOOL, false) {
			case ModsFileInfo.MOD_7_STATE_STOPPED:
				speak = "The GPT is stopped. Text on hold."
			case ModsFileInfo.MOD_7_STATE_STARTING:
				speak = "The GPT is starting up. Text on hold."
			case ModsFileInfo.MOD_7_STATE_BUSY:
				speak = "The GPT is busy. Text on hold."
		}
		if speak != "" {
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}

		return
	}

	Speech.QueueSpeech(txt_to_speak, speech_priority, mode, "", 0)
}

func sendToGPT(txt_to_send string) {
	if !Utils.IsCommunicatorConnectedSERVER() {
		var speak string = "GPT unavailable. Not connected to the server."
		Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

		return
	}

	var speak string = ""
	switch GPTComm.SendText(txt_to_send, GPTComm.SESSION_TYPE_ACTIVE, GPTComm.ROLE_USER, false) {
		case ModsFileInfo.MOD_7_STATE_STOPPED:
			speak = "The GPT is stopped. Text on hold."
		case ModsFileInfo.MOD_7_STATE_STARTING:
			speak = "The GPT is starting up. Text on hold."
		case ModsFileInfo.MOD_7_STATE_BUSY:
			speak = "The GPT is busy. Text on hold."
	}
	if speak != "" && txt_to_send != "/stop" {
		Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
	}
}
