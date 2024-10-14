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

// Speech represents a speech in the speech queue
type Speech struct {
	// id is the id of the speech
	id string
	// text is the text of the speech
	text string
	// time is the time the speech was added in milliseconds
	time int64
	// priority is the priority of the speech
	priority int
	// mode is the mode of the speech - an OR operation of different mode numbers
	mode int
	// task_id is the task id related to the speech
	task_id string
	// interrupted_times is the number of times the speech was interrupted
	interrupted_times int
}

/*
GetID gets the id of the speech

-----------------------------------------------------------

– Returns:
  - the id of the speech
 */
func (speech *Speech) GetID() string {
	return speech.id
}

/*
GetText gets the text of the speech

-----------------------------------------------------------

– Returns:
  - the text of the speech
 */
func (speech *Speech) GetText() string {
	return speech.text
}

/*
GetTime gets the time the speech was added in milliseconds

-----------------------------------------------------------

– Returns:
  - the time the speech was added in milliseconds
 */
func (speech *Speech) GetTime() int64 {
	return speech.time
}

/*
GetPriority gets the priority of the speech

-----------------------------------------------------------

– Returns:
  - the priority of the speech
 */
func (speech *Speech) GetPriority() int {
	return speech.priority
}

/*
GetMode gets the mode of the speech

-----------------------------------------------------------

– Returns:
  - the mode of the speech
 */
func (speech *Speech) GetMode() int {
	return speech.mode
}

/*
GetTaskID gets the task id related to the speech

-----------------------------------------------------------

– Returns:
  - the task id related to the speech
 */
func (speech *Speech) GetTaskID() string {
	return speech.task_id
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
