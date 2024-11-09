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
	"GPTComm/GPTComm"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

type _VolumeDndState struct {
	audio_stream int
	old_volume int
	was_muted int
}

const _DEFAULT_VALUE int = -3234

var volumeMutedState_GL _VolumeDndState = _VolumeDndState{
	audio_stream: _DEFAULT_VALUE,
	old_volume:   _DEFAULT_VALUE,
	was_muted:    _DEFAULT_VALUE,
}

// For slow devices, maybe 250 is good?
// EDIT: I think if we reset and set the volume too quickly, Android will mess up somewhere and one of the changes
// won't be detected (put a LOW and a CRITICAL one in row and that might happen). The detected change may happen
// about 500ms after the volume was set. More specifically, in a test, it gave 530ms. For slower devices, I've put
// 750ms at most. I think this time it should be enough...
// EDIT 2: this is on the computer, with no broadcasts of volume changed. Must be busy-waiting checking. And it's every
// second. So this must be greater than a second. I'll just add it to the old value.
// EDIT 3: Went up by 500ms until it worked. + 2 seconds it seems. So I'll add 500ms more just to be sure.
const VOLUME_CHANGE_INTERVAL int64 = 750 + 2500

var assist_changed_volume_time_ms_GL int64 = math.MaxInt64 - VOLUME_CHANGE_INTERVAL
var assist_will_change_volume_GL bool = false
var user_changed_volume_GL bool = false
var is_speaking_GL bool = false
var volume_mute_done_GL bool = false
var higher_priority_came_GL bool = false

var curr_speech_GL *SpeechQueue.Speech = nil

var mutex sync.Mutex

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		initTts()

		processVolumeChanges()

		//log.Println("Waiting for speeches to speak...")

		go func() {
			for {
				if curr_speech_GL != nil {
					var curr_speech *SpeechQueue.Speech = curr_speech_GL
					//log.Println("Speaking speech with priority " + strconv.Itoa(int(curr_speech.GetPriority())) + " and ID " +
					//	curr_speech.GetID()[:10] + "(...)...")

					var notified bool = rightBeforeSpeaking(curr_speech.GetID())
					log.Println("Notified:", notified)

					if err := speak(curr_speech.GetText()); err != nil {
						log.Println("Error speaking speech:", err)
						if !notified {
							Utils.QueueNotificationNOTIFS("Speeches", curr_speech.GetText())
						}
					}

					UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_LAST_SPEECH).SetString(curr_speech.GetText(), true)

					var speech_id string = curr_speech.GetID()
					if higher_priority_came_GL {
						higher_priority_came_GL = false

						curr_speech.RephraseInterrSpeech()
						speech_id = ""
					}

					curr_speech_GL = nil

					speechTreatment(speech_id)
				}

				if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
					return
				}
			}
		}()

		go func() {
			GPTComm.SetPreparations(time.Now().UnixMilli())
			for {
				// Keep getting the next sentence to speak from the server
				var speak string = GPTComm.GetNextSpeechSentence()
				if speak == "" || speak == GPTComm.END_ENTRY {
					time.Sleep(1 * time.Second)

					continue
				}

				if *module_stop {
					break
				}

				QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE_DEFAULT, "", 0)
			}
		}()

		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				stop_volume_processing_GL = true
				SpeechQueue.ClearQueue()

				return
			}
		}
	}
}

/*
QueueSpeech adds a speech to the speech queue.

Note: this function enforces the SpeechQueue.MODE1_ALWAYS_NOTIFY mode (if there's music playing, the speech won't be
heard).

-----------------------------------------------------------

– Params:
  - to_speak – the text to speak
  - priority – the priority of the speech (one of the constants in SpeechQueue)
  - mode – the mode of the speech (one of the constants in SpeechQueue)
 */
func QueueSpeech(to_speak string, priority int32, mode int32, speech_id string, task_id int32) {
	mutex.Lock()
	defer mutex.Unlock()

	if UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SPEECH_ALWAYS_NOTIFY).GetBool(true) {
		mode = mode | SpeechQueue.MODE1_ALWAYS_NOTIFY
	}

	var speech_id_to_use string = ""
	if speech_id != "" {
		speech_id_to_use = speech_id
	}
	if SpeechQueue.GetSpeech(speech_id_to_use) == nil {
		// If it's a new speech, add to the lists.
		speech_id_to_use = SpeechQueue.AddSpeech(to_speak, "", time.Now().UnixMilli(), priority, mode, 0, task_id)
	}

	if curr_speech_GL == nil {
		log.Println("No speech in progress. Speaking speech with ID " + speech_id_to_use[:10] + "(...)...")
		// If there's no speech in progress, send it to be spoken.
		curr_speech_GL = SpeechQueue.GetSpeech(speech_id_to_use)
	} else {
		log.Println("Speech in progress. Adding speech with ID " + speech_id_to_use[:10] + "(...) to the queue...")
		// If there's a speech already being spoken, the new one is just added to the list (when the current one
		// stops, it will take care of starting the next ones on the queues).
		// Except if the new speech has a higher priority than the current one. In that case, the current one
		// stops temporarily to give place to the new one.
		if priority > curr_speech_GL.GetPriority() {
			log.Println("Priority: " + strconv.Itoa(int(priority)) + " > " + strconv.Itoa(int(curr_speech_GL.GetPriority())))
			if stopTts() {
				higher_priority_came_GL = true
			}
		}
	}
}

/*
SkipCurrentSpeech skips the current speech.

-----------------------------------------------------------

– Returns:
  - true if the speech was skipped successfully, false otherwise
 */
func SkipCurrentSpeech() bool {
	return stopTts()
}
