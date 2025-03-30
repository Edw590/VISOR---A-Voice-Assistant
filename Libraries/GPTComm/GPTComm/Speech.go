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

package GPTComm

import (
	"Utils"
	"Utils/UtilsSWA"
	"strings"
	"time"
)

var time_begin_ms_GL int64 = -1
var curr_entry_time_ms_GL int64 = -1
var last_speech_GL string = ""
var last_idx_begin_GL int = 0
var last_text_GL string = ""
var ignore_sentence_GL bool = false

var start_checking_GL []bool = nil

const END_ENTRY string = "[3234_END]"
const ALL_DEVICES_ID string = "3234_ALL"

/*
SetPreparations sets the time to begin searching for the next speech and starts the process of getting the next speech
beginning.

Set only for the first time. The next times, the time is automatic and based on the time of the last speech.

-----------------------------------------------------------

– Params:
  - time_begin_ms – the time to begin searching for the next speech in milliseconds
*/
func SetPreparations(time_begin_ms int64) {
	time_begin_ms_GL = time_begin_ms
	go func() {
		for {
			// Get the start message but ignore it (we just need to know *when* to start).
			Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 0)

			// Keep waiting until the start message is received. If multiple are received, stack them up.
			start_checking_GL = append(start_checking_GL, true)
		}
	}()
}

/*
GetNextSpeechSentence gets the next sentence to be spoken of the most recent speech.

Each time the function is called, a new sentence is returned, until the end of the text file is reached, in which case
the function will return END_ENTRY.

In case a new speech is added to the text file, the function will continue the speech it was on until its end.

The function will wait until the time of the next speech is reached.

-----------------------------------------------------------

– Returns:
  - the next sentence to be spoken (sometimes may return an empty string - ignore) or END_ENTRY if the end of the text
    file is reached
 */
func GetNextSpeechSentence() string {
	if curr_entry_time_ms_GL == -1 {
		for {
			if len(start_checking_GL) == 0 {
				time.Sleep(1 * time.Second)

				continue
			}

			Utils.DelElemSLICES(&start_checking_GL, 0)

			break
		}

		var entry *_Entry = getEntry(-1, -1)
		var device_id string = entry.getDeviceID()
		if entry.getTime() >= time_begin_ms_GL && (device_id == Utils.GetGenSettings().Device_settings.Id ||
				device_id == ALL_DEVICES_ID) {
			curr_entry_time_ms_GL = entry.getTime()
			if curr_entry_time_ms_GL != 1 {
				time_begin_ms_GL = curr_entry_time_ms_GL + 1
			}
			last_speech_GL = ""
			ignore_sentence_GL = false
		}
	}
	if curr_entry_time_ms_GL == -1 {
		// If no Entry was found, return END_ENTRY
		return END_ENTRY
	}

	//log.Println("JJJJJJJJJJJJJJJJJJJJJJJJJJJJJ")
	//log.Println("curr_entry_time_ms_GL:", curr_entry_time_ms_GL)
	//log.Println("time_begin_ms_GL:", time_begin_ms_GL)

	var sentence string = ""
	for {
		var entry *_Entry = getEntry(curr_entry_time_ms_GL, -1)
		if entry.getTime() == -1 {
			//log.Println("No entry found")
			// Maybe no Internet connection, so it returns an empty Entry. Just wait until there is connection again.
			time.Sleep(1 * time.Second)

			continue
		}
		var text = entry.getText()
		last_text_GL = text

		//log.Println("--------------------------")
		//log.Println("text: \"" + text + "\"")

		text = strings.Replace(text, "\n", ". ", -1)
		text = strings.Replace(text, END_ENTRY, ". " + END_ENTRY, 1)
		text = strings.Replace(text, "...", ".", -1)
		//log.Println("text: \"" + text + "\"")
		if last_idx_begin_GL != 0 && last_idx_begin_GL >= len(text) {
			sentence = ""

			break
		}

		var dot_idx = strings.Index(text[last_idx_begin_GL:], ". ")
		var dot_idx2 = strings.IndexAny(text[last_idx_begin_GL:], "!?")
		if dot_idx2 != -1 && (dot_idx == -1 || dot_idx2 < dot_idx) {
			dot_idx = dot_idx2
		}

		//log.Println("dot_idx:", dot_idx)
		//log.Println("last_idx_begin_GL:", last_idx_begin_GL)
		//log.Println("text[last_idx_begin_GL:]:", text[last_idx_begin_GL:])

		// If the last dot index is not found, it means that the sentence is not finished yet. So, we must wait for the
		// next entry to be added to the text file.
		if dot_idx != -1 {
			sentence = text[last_idx_begin_GL : last_idx_begin_GL + dot_idx + 2]
			sentence = strings.Trim(sentence, " ")

			last_idx_begin_GL += dot_idx + 2

			// Ignore code
			if strings.Contains(sentence, "```") {
				ignore_sentence_GL = !ignore_sentence_GL
			}

			if ignore_sentence_GL {
				time.Sleep(1 * time.Second)

				continue
			}

			break
		}

		if strings.Contains(text[last_idx_begin_GL:], END_ENTRY) {
			//log.Println("RRRRRRRRRRRRRRRRRRR")
			sentence = END_ENTRY
			curr_entry_time_ms_GL = -1
			last_idx_begin_GL = 0

			break
		} else {
			//fmt.Println("text[last_idx_begin_GL:]:", text[last_idx_begin_GL:])
			time.Sleep(1 * time.Second)
		}
	}

	if !UtilsSWA.StringHasLettersGENERAL(sentence) {
		sentence = ""
	}

	//log.Println("sentence: \"" + sentence + "\"")

	if sentence != "" {
		if last_speech_GL != "" {
			last_speech_GL += " "
		}
		last_speech_GL += sentence
	}

	//log.Println("last_speech_GL: \"" + last_speech_GL + "\"")

	return sentence
}

/*
GetLastSpeech gets the last speech that was spoken, optimized for speech by a TTS engine.

-----------------------------------------------------------

– Returns:
  - the last speech that was spoken
 */
func GetLastSpeech() string {
	return last_speech_GL
}

/*
GetLastText gets the last text that was spoken, exactly as written by the LLM.

-----------------------------------------------------------

– Returns:
  - the last text that was spoken
 */
func GetLastText() string {
	return last_text_GL
}
