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

package DialogMan

import (
	"ACD/ACD"
	"Utils"
	"strconv"
	"strings"
	"time"
)

var last_it_GL string = ""
var last_it_when_GL int64 = 0
var last_and_GL string = ""
var last_and_when_GL int64 = 0

var intents_GL []*Intent = nil
var curr_intent_idx_GL int = -1

// HandleInputResult is the result of HandleInput()
type HandleInputResult struct {
	// Something_detected is set to true if something useful for the intents was detected, false otherwise
	Something_detected bool
	// Response is the response to speak to the user, if any
	Response string
	// next_intent_slot_idx is the index of the intent that will solely be used in the next HandleInput call, or -1 if
	// none
	next_intent_idx int
	// next_intent_slot_idx is the index of the slot of the intent that will solely be used in the next HandleInput
	// call (ignored if next_intent_idx is -1)
	next_intent_slot_idx int
	// Intents is the list of completed intents
	Intents []*Intent
	// Intent0 is the 1st completed intent
	Intent0 *Intent
	// Intent1 is the 2nd completed intent
	Intent1 *Intent
	// Intent2 is the 3rd completed intent
	Intent2 *Intent
	// Intent3 is the 4th completed intent
	Intent3 *Intent
	// Intent4 is the 5th completed intent
	Intent4 *Intent
	// Intent5 is the 6th completed intent
	Intent5 *Intent
	// Intent6 is the 7th completed intent
	Intent6 *Intent
	// Intent7 is the 8th completed intent
	Intent7 *Intent
	// Intent8 is the 9th completed intent
	Intent8 *Intent
	// Intent9 is the 10th completed intent
	Intent9 *Intent
}

/*
AddToIntentList adds an intent to the global intents list.

-----------------------------------------------------------

– Params:
  - intent – the intent to add
 */
func AddToIntentList(intent *Intent) {
	intents_GL = append(intents_GL, intent)
}

/*
ClearIntentsList clears the global intents list.
*/
func ClearIntentsList() {
	intents_GL = nil
	curr_intent_idx_GL = -1
}

/*
HandleInput processes the user input sentence and processes what Intent(s) it contains, if any.

-----------------------------------------------------------

– Params
  - sentence – the input string
  - handle_input_result – the result of the previous HandleInput call, or nil if this is the first call

– Returns:
  - the result, or nil if no intent was detected or completed
 */
func HandleInput(sentence string, handle_input_result *HandleInputResult) *HandleInputResult {
	var ready_intents []*Intent = make([]*Intent, 10)

	sentence = strings.ToLower(strings.TrimSpace(sentence))

	var next_intent_idx int = -1
	var next_intent_slot_idx int = -1
	if handle_input_result != nil {
		next_intent_idx = handle_input_result.next_intent_idx
		next_intent_slot_idx = handle_input_result.next_intent_slot_idx
	}

	if time.Now().UnixMilli() > last_it_when_GL + 60*1000 {
		last_it_GL = ""
	}
	if time.Now().UnixMilli() > last_and_when_GL + 60*1000 {
		last_and_GL = ""
	}

	// Reset all completed intents
	for _, intent := range intents_GL {
		if intent == nil {
			break
		}

		if intent.isReady() {
			// Means it's been completed already - so reset it
			intent.reset()
		}
	}
	// And reset the return array
	if handle_input_result != nil {
		for i := range handle_input_result.Intents {
			handle_input_result.Intents[i] = nil
		}
	}

	var something_detected bool = false
	if next_intent_idx >= 0 && intents_GL[next_intent_idx].Slots[next_intent_slot_idx].Acd_cmd_id == "" {
		// Means HandleInput wants anything as answer for the intent - no need to run the ACD
		var intent *Intent = intents_GL[next_intent_idx]
		var slot *Slot = intent.Slots[next_intent_slot_idx]
		slot.Sentence = sentence
		slot.Value = sentence
		slot.filled = true

		something_detected = true
	} else {
		var cmds_info_str string = ACD.Main(sentence, false, true, last_it_GL + "|" + last_and_GL)
		var cmds_info *ACD.CmdsInfo = ACD.ParseMainResult(cmds_info_str)
		if cmds_info != nil {
			Utils.LogLnDebug("*****************************")
			Utils.LogLnDebug(cmds_info_str)

			if cmds_info.Last_it != "" {
				last_it_GL = cmds_info.Last_it
				last_it_when_GL = time.Now().UnixMilli()
			}
			if cmds_info.Last_and != "" {
				last_and_GL = cmds_info.Last_and
				last_and_when_GL = time.Now().UnixMilli()
			}

			//Utils.LogLnDebug(last_it_GL)
			//Utils.LogLnDebug(last_and_GL)
			Utils.LogLnDebug("*****************************")

			if len(cmds_info.Detected_cmds) == 0 {
				return nil
			}

			for command_idx, command := range cmds_info.Detected_cmds {
				// Step 1: If no active intent, try to detect one
				var where_to_start_idx int = 0
				for intent_idx, intent := range intents_GL {
					if intent.Acd_cmd_id == command.Cmd_id {
						curr_intent_idx_GL = intent_idx
						where_to_start_idx = command_idx + 1

						intent.Sentence = sentence
						intent.Value = command.Cmd_variant
						intent.in_processing = true

						something_detected = true

						break
					}
				}
				if curr_intent_idx_GL == -1 {
					continue
				}

				// Check if there are slots associated
				var intent *Intent = intents_GL[curr_intent_idx_GL]
				if len(intent.Slots) > 0 {
					// Step 2: Fill the unfilled slots with the rest of the detected commands
					for command_idx1 := where_to_start_idx; command_idx1 < len(cmds_info.Detected_cmds); command_idx1++ {
						var command1 ACD.DetectedCmd = cmds_info.Detected_cmds[command_idx1]
						for slot_idx := range intent.Slots {
							var slot *Slot = intent.Slots[slot_idx]
							if !slot.filled && slot.Acd_cmd_id == command1.Cmd_id {
								slot.Sentence = sentence
								if slot.Acd_cmd_id == "" {
									slot.Value = sentence
								} else {
									slot.Value = command1.Cmd_variant
									for _, index := range command1.Indexes {
										slot.Indexes_str += strconv.Itoa(index) + ","
										slot.Indexes = append(slot.Indexes, index)
									}
									slot.Indexes_str = slot.Indexes_str[:len(slot.Indexes_str)-len(ACD.IDXS_SEPARATOR)]
								}
								slot.filled = true

								something_detected = true

								break
							}
						}
					}
				}

				if intent.isReady() {
					curr_intent_idx_GL = -1
				}
			}
		}
	}

	var speak string = ""
	var to_be_next_intent_idx int = -1
	var to_be_next_intent_slot_idx int = -1
	// If there is no active intent, ask for the next unfilled slot of the first intent with unfilled slots (if any).
	// If there is an active intent, ask for the next unfilled slot of it (if any).
	// If there are no unfilled slots, do nothing.
	// Note: if there are multiple intents with unfilled slots, only the first one is considered.
	var valid_next_ids []string = nil
	for intent_idx, intent := range intents_GL {
		if intent == nil {
			break
		}
		if !intent.in_processing {
			continue
		}

		for slot_idx, slot := range intent.Slots {
			if slot.filled {
				continue
			}

			if curr_intent_idx_GL != intent_idx {
				speak = "About the task \"" + strings.ToLower(intent.Task_name) + "\", I need to know: " +
					strings.ToLower(slot.Prompt)
			} else {
				speak = slot.Prompt
			}
			valid_next_ids = append(valid_next_ids, slot.Acd_cmd_id)
			to_be_next_intent_idx = intent_idx
			to_be_next_intent_slot_idx = slot_idx

			curr_intent_idx_GL = intent_idx

			goto almostEndFunc
		}
	}

	almostEndFunc:

	// If there are ready intents, add them to the return structure on the first available space
	for intent_idx, intent := range intents_GL {
		if intent.isReady() {
			for i := range ready_intents {
				if ready_intents[i] == nil {
					ready_intents[i] = intents_GL[intent_idx]

					break
				}
			}
		}
	}

	return &HandleInputResult{
		Something_detected:   something_detected,
		Response:             speak,
		next_intent_idx:      to_be_next_intent_idx,
		next_intent_slot_idx: to_be_next_intent_slot_idx,
		Intents:              ready_intents,
		Intent0:              ready_intents[0],
		Intent1:              ready_intents[1],
		Intent2:              ready_intents[2],
		Intent3:              ready_intents[3],
		Intent4:              ready_intents[4],
		Intent5:              ready_intents[5],
		Intent6:              ready_intents[6],
		Intent7:              ready_intents[7],
		Intent8:              ready_intents[8],
		Intent9:              ready_intents[9],
	}
}
