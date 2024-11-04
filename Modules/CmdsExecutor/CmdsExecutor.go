/*******************************************************************************
 * Copyright 2023-2024 The V.I.S.O.R. authors
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
	"GPTComm/GPTComm"
	"OICComm/OICComm"
	"Speech"
	"SpeechQueue/SpeechQueue"
	"TEHelper/TEHelper"
	"Utils"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Commands Executor //

var last_it string = ""
var last_it_when int64 = 0
var last_and string = ""
var last_and_when int64 = 0

var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		ACD.ReloadCmdsArray(prepareCommandsString())

		for {
			var comms_map map[string]any = <- Utils.ModsCommsChannels_GL[Utils.NUM_MOD_CmdsExecutor]
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
				// If there are only WARN_-started constants (nevative numbers), send to GPT
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
						speakInternal(speak, speech_priority, speech_mode2, true)

					case CMD_ASK_DATE:
						var speak string = "Today's " + Utils.GetDateStrTIMEDATE(-1)
						speakInternal(speak, speech_priority, speech_mode2, true)

					case CMD_TOGGLE_WIFI:
						if Utils.ToggleWifiCONNECTIVITY(cmd_variant == RET_ON) {
							var speak string
							if cmd_variant == RET_ON {
								speak = "Wi-Fi turned on."
							} else {
								speak = "Wi-Fi turned off."
							}
							speakInternal(speak, speech_priority, speech_mode2, false)
						} else {
							var on_off string = "off"
							if cmd_variant == RET_ON {
								on_off = "on"
							}
							var speak string = "Sorry, I couldn't turn the Wi-Fi " + on_off + "."
							speakInternal(speak, speech_priority, speech_mode2, true)
						}

					case CMD_ASK_BATTERY_PERCENT:
						var battery_percentage int = int(UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).
							GetInt(true))
						var speak string = "Battery percentage: " + strconv.Itoa(battery_percentage) + "%"
						speakInternal(speak, speech_priority, speech_mode2, true)

					case CMD_TELL_WEATHER:
						var speak string = "Obtaining the weather..."
						speakInternal(speak, speech_priority, speech_mode2, false)

						// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

						if UtilsSWA.WaitForNetwork(0) {
							var weather_str string = OICComm.GetWeather()
							if weather_str == "" {
								speak = "I'm sorry Sir, but I couldn't get the weather information."
								speakInternal(speak, speech_priority, speech_mode2, true)

								break
							}

							var weather_by_loc []string = strings.Split(weather_str, "\n")
							for _, weather := range weather_by_loc {
								if weather == "" {
									continue
								}

								var weather_data []string = strings.Split(weather, " ||| ")
								speak = "The weather in " + weather_data[0] + " is " + weather_data[5] +
									" with " + weather_data[1] + " degrees, a maximum of " + weather_data[6] +
									" degrees and a minimum of " + weather_data[7] + " degrees. The precipitation is of " +
									weather_data[2] + ", humidity of " + weather_data[3] + ", and wind of " +
									weather_data[4] + "."
								speakInternal(speak, speech_priority, speech_mode2, true)
							}
						} else {
							speak = "No network connection available to get the weather."
							speakInternal(speak, speech_priority, speech_mode2, true)
						}

					case CMD_TELL_NEWS:
						var speak string = "Obtaining the latest news..."
						speakInternal(speak, speech_priority, speech_mode2, true)

						// TODO: make him turn on Ethernet and Wi-Fi in case they're off and wait 10s instead of 0

						if UtilsSWA.WaitForNetwork(0) {
							var news_str string = OICComm.GetWeather()
							if news_str == "" {
								speak = "I'm sorry Sir, but I couldn't get the news information."
								speakInternal(speak, speech_priority, speech_mode2, true)

								break
							}

							var news_by_loc []string = strings.Split(news_str, "\n")
							for _, news_data := range news_by_loc {
								if news_data == "" {
									continue
								}

								var news []string = strings.Split(news_data, " ||| ")

								speak = "News in " + news[0] + ". "

								for _, n := range news[1:] {
									speak += n + ". "
								}
								speakInternal(speak, speech_priority, speech_mode2, false)
							}
						} else {
							speak = "No network connection available to get the weather."
							speakInternal(speak, speech_priority, speech_mode2, true)
						}

					case CMD_TOGGLE_ETHERNET:
						if Utils.ToggleEthernetCONNECTIVITY(cmd_variant == RET_ON) {
							var speak string
							if cmd_variant == RET_ON {
								speak = "Ethernet turned on."
							} else {
								speak = "Ethernet turned off."
							}
							speakInternal(speak, speech_priority, speech_mode2, false)
						} else {
							var on_off string = "off"
							if cmd_variant == RET_ON {
								on_off = "on"
							}
							var speak string = "Sorry, I couldn't turn the Ethernet " + on_off + "."
							speakInternal(speak, speech_priority, speech_mode2, true)
						}

					case CMD_TOGGLE_NETWORKING:
						if Utils.ToggleNetworkingCONNECTIVITY(cmd_variant == RET_ON) {
							var speak string
							if cmd_variant == RET_ON {
								speak = "Networking turned on."
							} else {
								speak = "Networking turned off."
							}
							speakInternal(speak, speech_priority, speech_mode2, false)
						} else {
							var on_off string = "off"
							if cmd_variant == RET_ON {
								on_off = "on"
							}
							var speak string = "Sorry, I couldn't turn the networking " + on_off + "."
							speakInternal(speak, speech_priority, speech_mode2, true)
						}
				}
			}


			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				TEHelper.StopChecker()

				return
			}
		}
	}
}

func speakInternal(txt_to_speak string, speech_priority int32, mode int32, auto_gpt bool) {
	if auto_gpt && Utils.IsCommunicatorConnectedSERVER() && GPTComm.SendText("", false) {
		var text string = "Rephrase the following to maintain its meaning but change its wording: \"" + txt_to_speak +
			"\". Current device: user's " + Utils.Gen_settings_GL.Device_settings.Type_ + "."
		if !GPTComm.SendText(text, false) {
			Speech.QueueSpeech("Sorry, the GPT is busy at the moment.", speech_priority,
				SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}

		return
	}

	Speech.QueueSpeech(txt_to_speak, speech_priority, mode, "", 0)
}

func sendToGPT(txt_to_send string) {
	if !Utils.IsCommunicatorConnectedSERVER() {
		var speak string = "GPT unavailable. not connected to the server."
		Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

		return
	}

	if !GPTComm.SendText(txt_to_send, true) {
		Speech.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
			SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
	}
}
