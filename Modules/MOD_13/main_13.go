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
						MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2)
					case CMD_ASK_DATE:
						var speak string = "Today's " + Utils.GetDateStrTIMEDATE(-1)
						MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2)
					case CMD_ASK_BATTERY_PERCENT:
						var battery_percentage int = UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).
							GetData(true, nil).(int)
						var speak string = "Battery percentage: " + strconv.Itoa(battery_percentage) + "%"
						MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, speech_mode2)
				}
			}


			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				TEHelper.StopChecker()

				return
			}
		}
	}
}

func prepareCommandsString() string {
	var commands = [...][]string{
		//{CMD_TOGGLE_FLASHLIGHT, ACD.CMDi_TYPE_TURN_ONFF, "", "", "flashlight/lantern"},
		{CMD_ASK_TIME, ACD.CMDi_TYPE_ASK, "", "", "time"},
		{CMD_ASK_DATE, ACD.CMDi_TYPE_ASK, "", "", "date"},
		//{CMD_TOGGLE_WIFI, ACD.CMDi_TYPE_TURN_ONFF, "", "", "wifi"},
		//{CMD_TOGGLE_MOBILE_DATA, ACD.CMDi_TYPE_TURN_ONFF, "", "", "mobile data"},
		//{CMD_TOGGLE_BLUETOOTH, ACD.CMDi_TYPE_TURN_ONFF, "", "", "bluetooth"},
		//{CMD_ANSWER_CALL, ACD.CMDi_TYPE_ANSWER, "", "", "call"},
		//{CMD_END_CALL, ACD.CMDi_TYPE_STOP, "", "", "call"},
		//{CMD_TOGGLE_SPEAKERS, ACD.CMDi_TYPE_TURN_ONFF, "", "", "speaker/speakers"},
		//{CMD_TOGGLE_AIRPLANE_MODE, ACD.CMDi_TYPE_TURN_ONFF, "", "", "airplane mode"},
		{CMD_ASK_BATTERY_PERCENT, ACD.CMDi_TYPE_ASK, "", "", "battery percentage", "battery status", "battery level"},
		//{CMD_SHUT_DOWN_DEVICE, ACD.CMDi_TYPE_SHUT_DOWN, "", "", "device/phone"},
		//{CMD_REBOOT_DEVICE, ACD.CMDi_TYPE_REBOOT, "fast", "fast|;4; -fast", "reboot/restart device/phone|device/phone|device/phone recovery|device/phone safe mode|device/phone bootloader"},
		//{CMD_TAKE_PHOTO, ACD.CMDi_TYPE_NONE, "take", "", "picture/photo|frontal picture/photo"},
		//{CMD_RECORD_MEDIA, ACD.CMDi_TYPE_START, "record", "record|record|;4; -record", "audio/sound|video/camera|recording audio/sound|recording video/camera"},
		//{CMD_SAY_AGAIN, ACD.CMDi_TYPE_REPEAT_SPEECH, "", "", "again", "say", "said"},
		//{CMD_MAKE_CALL, ACD.CMDi_TYPE_NONE, "make place", "", "call"},
		//{CMD_TOGGLE_POWER_SAVER_MODE, ACD.CMDi_TYPE_TURN_ONFF, "", "", "power/battery saver"},
		//{CMD_STOP_RECORD_MEDIA, ACD.CMDi_TYPE_STOP, "", "", "recording audio/sound|recording video/camera"},
		//{CMD_CONTROL_MEDIA, ACD.CMDi_TYPE_NONE, "play continue resume pause stop next previous", "play continue resume|pause|stop|next|previous", "media/song/songs/music/audio/musics/video/videos"},
		//{CMD_CONFIRM, ACD.CMDi_TYPE_NONE, "i", "", "do/confirm/approve/certify"},
		//{CMD_REJECT, ACD.CMDi_TYPE_NONE, "i", "", "don't/reject/disapprove"},
		//{CMD_STOP_LISTENING, ACD.CMDi_TYPE_STOP, "", "", "listening"},
		//{CMD_START_LISTENING, ACD.CMDi_TYPE_START, "", "", "listening"},
		//{CMD_TELL_WEATHER, ACD.CMDi_TYPE_ASK, "", "", "weather"},
		//{CMD_TELL_NEWS, ACD.CMDi_TYPE_ASK, "", "", "news"},
		//{CMD_GONNA_SLEEP, ACD.CMDi_TYPE_WILL_GO, "", "", "sleep"},
	}

	var commands_almost_str []string = nil
	for _, array := range commands {
		commands_almost_str = append(commands_almost_str, strings.Join(array, "||"))
	}

	return strings.Join(commands_almost_str, "\\")
}
