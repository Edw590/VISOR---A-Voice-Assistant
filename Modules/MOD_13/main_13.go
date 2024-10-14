/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

package MOD_13

import (
	"ACD/ACD"
	"GPTComm/GPTComm"
	MOD_3 "Speech"
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

			var sentence_str string = comms_map["Sentence"].(string)
			var cmds_info_str = ACD.Main(sentence_str, false, true, last_it + "|" + last_and)
			log.Println("*****************************")
			fmt.Println(sentence_str)
			log.Println(cmds_info_str)
			var cmds_info []string = strings.Split(cmds_info_str, ACD.INFO_CMDS_SEPARATOR)
			if len(cmds_info) < 2 {
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

			if strings.HasPrefix(cmds_info_str, ACD.ERR_CMD_DETECT) {
				var speak string = "WARNING! There was a problem processing the commands sir. This needs a fix. The " +
					"error was the following: " + cmds_info_str + ". You said: " + sentence_str
				MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY)
				log.Println("EXECUTOR - ERR_PROC_CMDS")

				time.Sleep(1 * time.Second)

				continue
			}

			if len(detected_cmds) == 0 || detected_cmds[0] == "" {
				if !Utils.IsCommunicatorConnectedSERVER() {
					var speak string = "GPT unavailable. Communicator not connected."
					MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY)

					return
				}

				if !GPTComm.SendText(sentence_str, true) {
					MOD_3.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.",
						SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY)
				}

				return
			}

			for _, command := range detected_cmds {
				var dot_index int = strings.Index(command, ".")
				if dot_index == -1 {
					// No command.
					continue
				}

				var cmd_id string = command[:dot_index] // "14.3" --> "14"
				//var cmd_variant string = command[dot_index:] // "14.3" --> ".3"

				var speech_mode2 = SpeechQueue.MODE_DEFAULT
				if cmdi_info[cmd_id] == CMDi_INF1_ONLY_SPEAK {
					speech_mode2 = SpeechQueue.MODE2_BYPASS_NO_SND
				}

				switch cmd_id {
					case CMD_ASK_TIME:
						var speak string = "It's " + Utils.GetTimeStrTIMEDATE(-1)
						speakInternal(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2, true)

					case CMD_ASK_DATE:
						var speak string = "Today's " + Utils.GetDateStrTIMEDATE(-1)
						speakInternal(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2, true)

					case CMD_ASK_BATTERY_PERCENT:
						var battery_percentage int = UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).
							GetData(true, nil).(int)
						var speak string = "Battery percentage: " + strconv.Itoa(battery_percentage) + "%"
						speakInternal(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2, true)
				}
			}


			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				TEHelper.StopChecker()

				return
			}
		}
	}
}

func speakInternal(txt_to_speak string, speech_priority int, mode int, auto_gpt bool) {
	if auto_gpt && Utils.IsCommunicatorConnectedSERVER() {
		var text string = "Sent from my " + Utils.Device_settings_GL.Device_type + ": write ONE concise sentence " +
			"saying \"" + txt_to_speak + "\"."
		if !GPTComm.SendText(text, false) {
			MOD_3.QueueSpeech("Sorry, the GPT is busy at the moment.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY)
		}

		return
	}

	MOD_3.QueueSpeech(txt_to_speak, speech_priority, mode)
}
