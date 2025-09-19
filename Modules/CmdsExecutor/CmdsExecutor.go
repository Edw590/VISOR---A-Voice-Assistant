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
	"DialogMan"
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
	"strconv"
	"strings"
)

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	ACD.ReloadCmdsArray(prepareCommandsString())
	DialogMan.ReloadIntentList(getIntentList())

	var handle_input_result *DialogMan.HandleInputResult = nil
	for {
		var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_CmdsExecutor, 0, -1)
		if comms_map == nil {
			return
		}

		var speech_priority int32 = SpeechQueue.PRIORITY_MEDIUM
		var sentence string = ""
		if map_value, ok := comms_map["SentenceInternal"]; ok {
			sentence = map_value.(string)
		} else if map_value, ok = comms_map["Sentence"]; ok {
			speech_priority = SpeechQueue.PRIORITY_USER_ACTION
			sentence = map_value.(string)
		} else {
			continue
		}

		handle_input_result = DialogMan.HandleInput(sentence, handle_input_result)
		if handle_input_result == nil {
			sendToGPT(sentence)

			continue
		}

		if handle_input_result.Response != "" {
			speakInternal(handle_input_result.Response, speech_priority, SpeechQueue.MODE_DEFAULT, _SESSION_TYPE_NONE,
				false, true)
		}

		var any_intent_detected bool = false
		for _, intent := range handle_input_result.Intents {
			if intent == nil {
				break
			}
			any_intent_detected = true

			var speech_mode2 = SpeechQueue.MODE_DEFAULT
			if cmdi_info[intent.Acd_cmd_id] == CMDi_INF1_ONLY_SPEAK {
				speech_mode2 = SpeechQueue.MODE2_BYPASS_NO_SND
			}

			switch intent.Acd_cmd_id {
				case CMD_ASK_TIME:
					var speak string = "It's " + Utils.GetTimeStrTIMEDATE(-1)
					speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false, true)

				case CMD_ASK_DATE:
					var speak string = "Today's " + Utils.GetDateStrDATETIME(-1)
					speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false, true)

				case CMD_TOGGLE_WIFI:
					if Utils.ToggleWifiCONNECTIVITY(intent.Value == RET_ON) {
						var speak string
						if intent.Value == RET_ON {
							speak = "Wi-Fi turned on."
						} else {
							speak = "Wi-Fi turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					} else {
						var on_off string = "off"
						if intent.Value == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the Wi-Fi " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_ASK_BATTERY_PERCENT:
					var battery_percentage int = int(UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).
						GetInt(true))
					var speak string = "Battery percentage: " + strconv.Itoa(battery_percentage) + "%"
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

				case CMD_TELL_WEATHER:
					var speak string = "Obtaining the weather..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var weather_locs []string = strings.Split(OICComm.GetWeatherLocationsList(), "|")

						for _, weather_loc := range weather_locs {
							var weather *ModsFileInfo.Weather = OICComm.GetWeather(weather_loc)
							if weather == nil {
								speak = "There is no weather data associated with the location " + weather_loc + "."
								speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)

								continue
							}

							if weather.Temperature == "" {
								// One being empty means the whole weather is empty
								speak = "There was a problem obtaining the weather for " + weather.Location + "."
								speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)

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
							speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, false, true)
						}
					} else {
						speak = "Not connected to the server to get the weather."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_TELL_NEWS:
					var speak string = "Obtaining the latest news..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var news_locs []string = strings.Split(OICComm.GetNewsLocationsList(), "|")

						for _, news_loc := range news_locs {
							var news *ModsFileInfo.News = OICComm.GetNews(news_loc)

							speak = "News in " + news.Location + ". "

							for _, n := range news.News {
								speak += n + ". "
							}
							speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)
						}
					} else {
						speak = "Not connected to the server to get the news."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_TOGGLE_ETHERNET:
					if Utils.ToggleEthernetCONNECTIVITY(intent.Value == RET_ON) {
						var speak string
						if intent.Value == RET_ON {
							speak = "Ethernet turned on."
						} else {
							speak = "Ethernet turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)
					} else {
						var on_off string = "off"
						if intent.Value == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the Ethernet " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_TOGGLE_NETWORKING:
					if Utils.ToggleNetworkingCONNECTIVITY(intent.Value == RET_ON) {
						var speak string
						if intent.Value == RET_ON {
							speak = "Networking turned on."
						} else {
							speak = "Networking turned off."
						}
						speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)
					} else {
						var on_off string = "off"
						if intent.Value == RET_ON {
							on_off = "on"
						}
						var speak string = "Sorry, I couldn't turn the networking " + on_off + "."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_ASK_EVENTS:
					var speak string = "Obtaining the tasks and events..."
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

					// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

					if UtilsSWA.WaitForNetwork(0) {
						var events_ids []string = strings.Split(GMan.GetEventsIdsList(true), "|")
						var tasks_ids []string = strings.Split(GMan.GetTasksIdsList(), "|")

						speak = getEventsList(events_ids, intent.Value)

						if intent.Value == RET_31_TODAY || intent.Value == RET_31_TOMORROW {
							speak += " " + getTasksList(tasks_ids, intent.Value)
						}

						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, true, true)
					} else {
						speak = "Not connected to the server to get the tasks and events."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)
					}

				case CMD_HELP_PICTURE:
					var clipboard []byte = Utils.GetClipboardGENERAL()
					var png []byte = isPng(clipboard)
					if png == nil {
						var speak string = "There is no PNG image in the clipboard. Remember, it has to be a PNG image."
						speakInternal(speak, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_TEMP, false, true)

						continue
					}

					GPTComm.AddFileToSend(true, png)

					speakInternal(sentence, speech_priority, speech_mode2, GPTComm.SESSION_TYPE_ACTIVE, true, false)

				case CMD_CREATE_EVENT:
					var about_str string = ""
					var when_s int64 = 0
					var how_long_min int64 = 0
					for slot_idx, slot := range intent.Slots {
						switch slot_idx {
							case 0: // About?
								about_str = slot.Value
							case 1: // When?
								when_s = UtilsSWA.TimeDateToTimestampDATETIME(slot.Value)
								if when_s == -1 {
									var speak string = "Sorry, I couldn't understand the date you mentioned. Set the " +
										"event again please with another format for the date."
									speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

									break
								}
							case 2: // For how long?
								how_long_min = UtilsSWA.ParseDurationDATETIME(slot.Value)
								if how_long_min == -1 {
									var speak string = "Sorry, I couldn't understand the duration you mentioned. Set " +
										"the event again please with another format for the duration."
									speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

									break
								}

								how_long_min /= 60
						}
					}
					if about_str == "" || when_s == 0 || how_long_min == 0 {
						break
					}

					var speak string
					if UtilsSWA.WaitForNetwork(0) {
						GMan.AddEvent(&ModsFileInfo.GEvent{
							Summary:      about_str,
							Start_time_s: when_s,
							Duration_min: how_long_min,
						})

						if GMan.IsTokenValid() {
							speak = "The event will be created now."
						} else {
							speak = "Apologies Sir, but the event will not be created: the Google Manager " +
								"token is not valid. Please set it up again."
						}
					} else {
						speak = "Not connected to the server to add event."
					}
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

				case CMD_CREATE_TASK:
					var about_str string = ""
					var when_s int64 = 0
					for slot_idx, slot := range intent.Slots {
						switch slot_idx {
							case 0:
								about_str = slot.Value
							case 1:
								when_s = UtilsSWA.TimeDateToTimestampDATETIME(slot.Value)
								if when_s == -1 {
									var speak string = "Sorry, I couldn't understand the date you mentioned. Set the " +
										"task again please with another format for the date."
									speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)

									break
								}
						}
					}
					if about_str == "" || when_s == 0 {
						break
					}

					var speak string
					if UtilsSWA.WaitForNetwork(0) {
						GMan.AddTask(&ModsFileInfo.GTask{
							Title:   about_str,
							Date_s:  when_s,
						})

						if GMan.IsTokenValid() {
							speak = "The task will be created now."
						} else {
							speak = "Apologies Sir, but the task will not be created: the Google Manager " +
								"token is not valid. Please set it up again."
						}
					} else {
						speak = "Not connected to the server to add task."
					}
					speakInternal(speak, speech_priority, speech_mode2, _SESSION_TYPE_NONE, false, true)
			}
		}

		if !any_intent_detected {
			sendToGPT(sentence)

			continue
		}


		if Utils.WaitWithStopDATETIME(module_stop, 1) {
			TEHelper.StopChecker()

			return
		}
	}
}

const _SESSION_TYPE_NONE string = "NONE"
func speakInternal(txt_to_speak string, speech_priority int32, mode int32, session_type string, wait_for_gpt bool,
				   use_tool_role bool) {
	if session_type != _SESSION_TYPE_NONE && speech_priority <= SpeechQueue.PRIORITY_USER_ACTION &&
				Utils.IsCommunicatorConnectedSERVER() && (wait_for_gpt ||
				GPTComm.SendText("", "", "", false) == ModsFileInfo.MOD_7_STATE_READY) {
		var role_to_use string = GPTComm.ROLE_USER
		if use_tool_role {
			role_to_use = GPTComm.ROLE_TOOL
		}
		var speak string = ""
		switch GPTComm.SendText(txt_to_speak, session_type, role_to_use, false) {
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

/*
isPng checks if the given byte slice is a PNG image (path or PNG bytes).

-----------------------------------------------------------

– Params:
  - clipboard – the byte slice to check

– Returns:
  - the PNG image bytes if it's a PNG image, nil otherwise
 */
func isPng(clipboard []byte) []byte {
	if clipboard == nil {
		return nil
	}

	var png []byte = clipboard

	var png_path string = string(png)
	if strings.HasPrefix(png_path, "\"") && strings.HasSuffix(png_path, "\"") {
		png_path = png_path[1 : len(png_path)-1]
	}
	var image_path Utils.GPath = Utils.PathFILESDIRS(false, "", png_path)
	if image_path.Exists() {
		// If it's a file path, use its contents. If it's not, check if it's a PNG already.
		png = image_path.ReadFile()
	}

	var png_header []byte = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if len(png) < len(png_header) {
		return nil
	}

	var is_png bool = true
	for i := 0; i < len(png_header); i++ {
		if png_header[i] != png[i] {
			is_png = false

			break
		}
	}
	if !is_png {
		return nil
	}

	return png
}
