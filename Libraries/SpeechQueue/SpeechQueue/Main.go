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

package SpeechQueue

import (
	"Utils"
	"math"
	"time"
)

const PRIORITY_LOW int = 0
const PRIORITY_MEDIUM int = 1
const PRIORITY_USER_ACTION int = 2
const PRIORITY_HIGH int = 3
const PRIORITY_CRITICAL int = 4

const NUM_PRIORITIES int = PRIORITY_CRITICAL + 1

// MODE_DEFAULT represents all default modes. The default of MODE1 is to only notify if he can't speak. The default of
//MODE2 is to not bypass the sound.
const MODE_DEFAULT int = 0;
// MODE1_NO_NOTIF doesn't notify even if he can't speak
const MODE1_NO_NOTIF int = 1 << 0;
// MODE1_ALWAYS_NOTIFY always notifies, even if he can speak
const MODE1_ALWAYS_NOTIFY int = 1 << 1;
// MODE2_BYPASS_NO_SND bypasses the no-sound state in case the device is in a no-sound state
const MODE2_BYPASS_NO_SND int = 1 << 2;

var speech_queue_GL []*Speech = nil

/*
AddSpeech adds a speech to the speech queue.

-----------------------------------------------------------

– Params:
  - text – the text of the speech
  - millis – the time at which the speech was added in milliseconds
  - priority – the priority of the speech
  - mode – the mode of the speech - an OR operation of different mode numbers
  - task_id – the task id related to the speech

– Returns:
  - the id of the speech
 */
func AddSpeech(text string, millis int64, priority int, mode int, task_id string) string {
	var id = Utils.RandStringGENERAL(2048)

	if millis == 0 {
		millis = time.Now().UnixMilli()
	}

	speech := &Speech{
		id: id,
		text: text,
		time: millis,
		priority: priority,
		mode: mode,
		task_id: task_id,
	}

	speech_queue_GL = append(speech_queue_GL, speech)

	return id
}

/*
GetSpeech gets a speech from the speech queue.

-----------------------------------------------------------

– Params:
  - id – the id of the speech

– Returns:
  - the speech or nil if the speech does not exist
 */
func GetSpeech(id string) *Speech {
	for _, speech := range speech_queue_GL {
		if speech.id == id {
			return speech
		}
	}

	return nil
}

/*
RemoveSpeech removes a speech from the speech queue.

-----------------------------------------------------------

– Params:
  - id – the id of the speech

– Returns:
  - true if the speech was removed, false if the speech does not exist
 */
func RemoveSpeech(id string) bool {
	for i, speech := range speech_queue_GL {
		if speech.id == id {
			Utils.DelElemSLICES(&speech_queue_GL, i)

			return true
		}
	}

	return false
}

/*
GetNextSpeech gets the next/oldest speech in the speech queue based on the priority and time.

-----------------------------------------------------------

– Params:
  - priority – the priority of the speech

– Returns:
  - the next speech or nil if there are no speeches with the priority
 */
func GetNextSpeech(priority int) *Speech {
	var oldest_time int64 = math.MaxInt64
	var oldest_speech *Speech = nil
	for _, speech := range speech_queue_GL {
		if speech.priority == priority && speech.time <= oldest_time {
			oldest_time = speech.time
			oldest_speech = speech
		}
	}

	return oldest_speech
}
