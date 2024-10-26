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

package Speech

import (
	"SpeechQueue/SpeechQueue"
	"Utils"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"github.com/itchyny/volume-go"
	"log"
	"strconv"
)

func rightBeforeSpeaking(speech_id string) bool {
	var curr_speech *SpeechQueue.Speech = SpeechQueue.GetSpeech(speech_id)

	var notified bool = false

	var skip_speaking bool = false
	if Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Muted {
		skip_speaking = curr_speech.GetPriority() != SpeechQueue.PRIORITY_CRITICAL &&
			(curr_speech.GetMode() & SpeechQueue.MODE2_BYPASS_NO_SND) == 0
	}

	if skip_speaking {
		// TODO: execute the task of the speech through its ID, if any

		if curr_speech.GetMode() & SpeechQueue.MODE1_NO_NOTIF == 0 {
			Utils.QueueNotificationNOTIFS("Speeches", curr_speech.GetText())
			notified = true
		}
	} else {
		// If it's to speak, prepare the app to speak.
		log.Println("(curr_speech.GetMode() & SpeechQueue.MODE1_NO_NOTIF == 0):", curr_speech.GetMode() & SpeechQueue.MODE1_NO_NOTIF == 0)
		var still_notify = false
		if !volume_muted_done_GL {
			still_notify = setToSpeakChanges(speech_id)
		}
		if still_notify || (curr_speech.GetMode() & SpeechQueue.MODE1_ALWAYS_NOTIFY != 0) {
			Utils.QueueNotificationNOTIFS("Speeches", curr_speech.GetText())
			notified = true
		}

		is_speaking_GL = true
	}

	return notified
}

func setToSpeakChanges(speech_id string) bool {
	var curr_speech *SpeechQueue.Speech = SpeechQueue.GetSpeech(speech_id)

	var still_notify bool = false

	log.Println("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	if curr_speech.GetPriority() == SpeechQueue.PRIORITY_CRITICAL {
		// Set the muted state
		if Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Muted {
			volumeMutedState_GL.was_muted = 1
		} else {
			volumeMutedState_GL.was_muted = 0
		}
		_ = volume.Unmute()

		// Set the volume
		var old_volume int = Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Volume
		volumeMutedState_GL.old_volume = old_volume
		var new_volume int = UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SPEECH_CRITICAL_VOL).GetInt(true)

		setResetWillChangeVolume(true)

		if new_volume != old_volume {
			_ = volume.SetVolume(new_volume)
		}
	} else {
		if curr_speech.GetMode() & SpeechQueue.MODE2_BYPASS_NO_SND != 0 {
			if Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Muted {
				volumeMutedState_GL.was_muted = 1
			} else {
				volumeMutedState_GL.was_muted = 0
			}
			err := volume.Unmute()
			if err != nil {
				still_notify = true
			}
		}

		if !Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Muted {
			// Set the volume
			var curr_volume int = Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Volume
			volumeMutedState_GL.old_volume = curr_volume

			log.Println("GGGGGGGGGGGGGGGGGGGGGGGGG")

			var new_volume int = UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SPEECH_NORMAL_VOL).GetInt(true)
			if curr_volume < new_volume {
				volumeMutedState_GL.old_volume = curr_volume

				log.Println("Setting the volume to speak to " + strconv.Itoa(new_volume) + "...")

				setResetWillChangeVolume(true)

				err := volume.SetVolume(new_volume)
				if err != nil {
					log.Println("Error setting the volume to speak: " + err.Error())
					still_notify = true
				}
			}
		}
	}

	volume_muted_done_GL = true

	log.Println("Still notify:", still_notify)

	return still_notify
}
