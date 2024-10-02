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

package GPTComm

import (
	"Utils"
	"strings"
	"time"
)

var time_begin_GL int64 = -1
var curr_entry_time_GL int64 = -1
var curr_idx_GL int = 0
var last_speech_GL string = ""

const END_ENTRY string = "[3234_END]"
const ALL_DEVICES_ID string = "3234_ALL"

/*
SetTimeBegin sets the time to begin searching for the next speech.

-----------------------------------------------------------

– Params:
  - time_begin – the time to begin searching for the next speech
 */
func SetTimeBegin(time_begin int64) {
	time_begin_GL = time_begin
}

/*
GetNextSpeechSentence gets the next sentence to be spoken of the most recent speech.

Each time the function is called, a new sentence is returned, until the end of the text file is reached, in which case
the function will return END_ENTRY.

In case a new speech is added to the text file, the function will continue the speech it was on until its end.

The function will wait until the time of the next speech is reached.

-----------------------------------------------------------

– Returns:
  - the next sentence to be spoken or END_ENTRY if the end of the text file is reached
 */
func GetNextSpeechSentence() string {
	if curr_entry_time_GL == -1 {
		var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GPTComm]
		if comms_map == nil {
			return END_ENTRY
		}

		if string(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)) == "start" {
			var entry *_Entry = getEntry(-1, -1)
			var device_id string = entry.getDeviceID()
			if entry.getTime() >= time_begin_GL && (device_id == Utils.User_settings_GL.PersonalConsts.Device_ID || device_id == ALL_DEVICES_ID) {
				curr_entry_time_GL = entry.getTime()
				time_begin_GL = curr_entry_time_GL + 1
				last_speech_GL = ""
			}
		}
	}

	var sentence string = ""
	for {
		var entry *_Entry = getEntry(curr_entry_time_GL, -1)
		var text_split []string = strings.Split(entry.getText(), " ")

		//log.Println("--------------------------")
		if curr_idx_GL >= len(text_split) {
			break
		}
		//log.Println("len(text_split):", len(text_split))
		for i := curr_idx_GL; i < len(text_split); i++ {
			var word string = text_split[i]
			//log.Println("curr_idx:", i)
			//log.Println("word:", word)

			if strings.Contains(word, END_ENTRY) {
				// If the word contains END_ENTRY, remove it and add a period at the end in case there's not already
				// one. Example: "peers[3234_END]" --> "peers.". Or "peers.[3234_END]" --> "peers.".

				// Add one more to go out of bounds next time the function is called. Will make it break the loop
				// instantly.
				curr_idx_GL++

				// But if the word is END_ENTRY alone, just break the loop and return whatever there is - including
				// nothing, which is taken care of below.
				if word == END_ENTRY {
					break
				}

				word = strings.Replace(word, END_ENTRY, "", -1)
				if !strings.HasSuffix(word, ".") && !strings.HasSuffix(word, "!") && !strings.HasSuffix(word, "?") {
					word += "."
				}
			} else if word != "" {
				curr_idx_GL++
			} else {
				continue
			}

			if sentence == "" {
				sentence = word
			} else {
				sentence += " " + word
			}

			if strings.HasSuffix(word, ".") || strings.HasSuffix(word, "!") || strings.HasSuffix(word, "?") {
				//log.Println("sentence: \"" + sentence + "\"")
				//log.Println("word: \"" + word + "\"")
				sentence = strings.TrimSuffix(sentence, " ")

				break
			}
		}

		if strings.HasSuffix(sentence, ".") || strings.HasSuffix(sentence, "!") || strings.HasSuffix(sentence, "?") {
			break
		}

		time.Sleep(1 * time.Second)
	}

	//log.Println("sentence: \"" + sentence + "\"")

	if sentence == "" {
		sentence = END_ENTRY
		curr_entry_time_GL = -1
		curr_idx_GL = 0
	}

	last_speech_GL += " " + sentence

	return sentence
}

func GetLastSpeech() string {
	return last_speech_GL
}
