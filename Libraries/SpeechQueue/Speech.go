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

package SpeechQueue

// Speech represents a speech in the speech queue
type Speech struct {
	// id is the ID of the speech
	id string
	// text is the text of the speech
	text string
	// time is the time the speech was added (any chosen resolution)
	time int64
	// priority is the priority of the speech
	priority int32
	// mode is the mode of the speech - an OR operation of different mode numbers
	mode int32
	// task_id is the task ID related to the speech
	task_id int32
	// interrupted_times is the number of times the speech was interrupted
	interrupted_times int32
	// audio_stream is the stream on which to speak the speech, if applicable
	audio_stream int32
	// stopped is a flag that indicates if the speech was stopped or not
	stopped bool
}

/*
GetID gets the ID of the speech.

-----------------------------------------------------------

– Returns:
  - the id of the speech
 */
func (speech *Speech) GetID() string {
	return speech.id
}

/*
GetText gets the text of the speech.

-----------------------------------------------------------

– Returns:
  - the text of the speech
 */
func (speech *Speech) GetText() string {
	return speech.text
}

/*
GetTime gets the time the speech was added.

-----------------------------------------------------------

– Returns:
  - the time the speech was added
 */
func (speech *Speech) GetTime() int64 {
	return speech.time
}

/*
GetPriority gets the priority of the speech.

-----------------------------------------------------------

– Returns:
  - the priority of the speech
 */
func (speech *Speech) GetPriority() int32 {
	return speech.priority
}

/*
GetMode gets the mode of the speech.

-----------------------------------------------------------

– Returns:
  - the mode of the speech
 */
func (speech *Speech) GetMode() int32 {
	return speech.mode
}

/*
GetTaskID gets the task ID related to the speech.

-----------------------------------------------------------

– Returns:
  - the task id related to the speech
 */
func (speech *Speech) GetTaskID() int32 {
	return speech.task_id
}

/*
GetAudioStream gets the audio stream of the speech.

-----------------------------------------------------------

– Returns:
  - the task id related to the speech
*/
func (speech *Speech) GetAudioStream() int32 {
	return speech.audio_stream
}

/*
GetStopped gets the stopped flag of the speech.

-----------------------------------------------------------

– Returns:
  - the value of the stopped flag
*/
func (speech *Speech) GetStopped() bool {
	return speech.stopped
}

/*
SetStopped sets the stopped flag of the speech.

-----------------------------------------------------------

– Params:
  - stopped – the new value for the stopped flag
*/
func (speech *Speech) SetStopped(stopped bool) {
	speech.stopped = stopped
}

/*
RephraseInterrSpeech rephrases an interrupted speech depending on the number of attempts.
*/
func (speech *Speech) RephraseInterrSpeech() {
	const PREFIX_1 = "As I was saying, "
	const PREFIX_2 = "Once again, as I was saying, "
	const PREFIX_3 = "And again, as I was saying, "

	if speech.interrupted_times == 0 {
		speech.text = PREFIX_1 + speech.text
	} else if speech.interrupted_times == 1 {
		speech.text = PREFIX_2 + speech.text[len(PREFIX_1):]
	} else if speech.interrupted_times == 2 {
		speech.text = PREFIX_3 + speech.text[len(PREFIX_2):]
	} else {
		return
	}

	speech.interrupted_times++
}
