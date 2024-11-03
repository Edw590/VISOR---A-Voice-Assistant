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
	"log"
	"time"
)

func speechTreatment(speech_id string) {
	var last_speech *SpeechQueue.Speech = SpeechQueue.RemoveSpeech(speech_id)

	// If there's an ID of a Task to run after the speech is finished, run it in a new thread.
	if last_speech != nil && last_speech.GetTaskID() != 0 {
		go func() {
			// TODO: execute the task
		}()
	}

	var next_speech *SpeechQueue.Speech = SpeechQueue.GetNextSpeech(-1)
	log.Println("next_speech != nil?", next_speech != nil)
	if next_speech != nil {
		if volume_mute_done_GL && ((last_speech == nil || last_speech.GetPriority() != SpeechQueue.PRIORITY_CRITICAL) &&
					next_speech.GetPriority() == SpeechQueue.PRIORITY_CRITICAL) ||
				(last_speech != nil && last_speech.GetPriority() == SpeechQueue.PRIORITY_CRITICAL) {
			// This if statement is for when a CRITICAL speech comes after or came before the current speech. In that
			// case, the volume/mute state may be changed because of the CRITICAL speech. So we have to reset it before
			// speaking the CRITICAL speech (because before came a normal speech that changed the volume), and reset it
			// after speaking it to then speak normal speeches.
			// This is because the volume/mute state is saved only once and not in a queue. So to know the previous
			// states, we have to reset it before changing it again.
			resetToSpeakChanges()
		}

		QueueSpeech("", 0, 0, next_speech.GetID(), 0)

		// This is a break between speeches so they're not all at once without a small break in between (which is awkward,
		// and doesn't help the brain process when one ends and the other one starts). 500 milliseconds should suffice, I
		// guess.
		time.Sleep(500 * time.Millisecond)

		log.Println("Returning from speechTreatment() after speech queued...")

		return
	}

	allSpeechesFinished()
}

func resetToSpeakChanges() {
	setResetWillChangeVolume(false)

	// Reset the volume
	log.Println("volumeMutedState_GL.old_volume:", volumeMutedState_GL.old_volume)
	log.Println("user_changed_volume_GL:", user_changed_volume_GL)
	if volumeMutedState_GL.old_volume != _DEFAULT_VALUE && !user_changed_volume_GL {
		log.Println("Setting the volume back to the previous value...")
		if !Utils.SetVolumeVOLUME(volumeMutedState_GL.old_volume) {
			log.Println("Error setting the volume back to the previous value")
		}
	}

	// Reset the muted state
	if volumeMutedState_GL.was_muted != _DEFAULT_VALUE {
		if volumeMutedState_GL.was_muted == 1 {
			Utils.SetMutedVOLUME(true)
		} else {
			Utils.SetMutedVOLUME(false)
		}
	}

	setVoluneMutedStateDefaults()

	volume_mute_done_GL = false
}

func allSpeechesFinished() {
	is_speaking_GL = false

	if volume_mute_done_GL {
		resetToSpeakChanges()
	}

	if assist_will_change_volume_GL {
		setResetWillChangeVolume(false)
	}

	user_changed_volume_GL = false
}
