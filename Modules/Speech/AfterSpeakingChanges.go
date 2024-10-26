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
	"github.com/itchyny/volume-go"
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
		err := volume.SetVolume(volumeMutedState_GL.old_volume)
		if err != nil {
			log.Println("Error setting the volume back to the previous value:", err)
		}
	}

	// Reset the muted state
	if volumeMutedState_GL.was_muted != _DEFAULT_VALUE {
		if volumeMutedState_GL.was_muted == 1 {
			_ = volume.Mute()
		} else {
			_ = volume.Unmute()
		}
	}

	setVoluneMutedStateDefaults()

	volume_muted_done_GL = false
}

func allSpeechesFinished() {
	is_speaking_GL = false

	if volume_muted_done_GL {
		resetToSpeakChanges()
	}

	if assist_will_change_volume_GL {
		setResetWillChangeVolume(false)
	}

	user_changed_volume_GL = false
}
