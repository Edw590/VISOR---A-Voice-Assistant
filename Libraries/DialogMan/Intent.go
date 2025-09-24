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

// Intent has an ID and a list of slots that need to be filled
type Intent struct {
	// Ad_cmd_id is the ID of the ACD command to use to detect the intent beginning
	Acd_cmd_id string
	// Task_name is the name of the task associated with the intent
	Task_name string
	// Value is the value of the intent
	Value string
	// Sentence is the sentence that was used to trigger the intent command
	Sentence string
	// in_processing is true if the intent is being processed, false otherwise
	in_processing bool
	// Slots is the list of slots that need to be filled to complete the intent
	Slots []*Slot
	// Slot0 is the 1st slot of the intent
	Slot0 *Slot
	// Slot1 is the 2nd slot of the intent
	Slot1 *Slot
	// Slot2 is the 3rd slot of the intent
	Slot2 *Slot
	// Slot3 is the 4th slot of the intent
	Slot3 *Slot
	// Slot4 is the 5th slot of the intent
	Slot4 *Slot
	// Slot5 is the 6th slot of the intent
	Slot5 *Slot
	// Slot6 is the 7th slot of the intent
	Slot6 *Slot
	// Slot7 is the 8th slot of the intent
	Slot7 *Slot
	// Slot8 is the 9th slot of the intent
	Slot8 *Slot
	// Slot9 is the 10th slot of the intent
	Slot9 *Slot
}

// Slot is a piece of information the assistant needs to complete an intent
type Slot struct {
	// Prompt is the question to ask the user to fill the slot
	Prompt string
	// Acd_cmd_id is the ID of the ACD command to use to detect the response to the Prompt
	Acd_cmd_id string
	// Value is the value of the slot
	Value string
	// Indexes_str is the list of word indexes in the original sentence sent to the ACD split by " " that triggered the
	// slot command, separated by ","
	Indexes_str string
	// Indexes is the same as Indexes_str but as a slice of ints
	Indexes []int
	// Sentence is the sentence that was used to trigger the slot command
	Sentence string

	// filled is true if the slot has been filled, false otherwise
	filled bool
}

/*
ManualSlotsReady prepares the Slots slice from the individual Slot fields.
 */
func (intent *Intent) ManualSlotsReady() {
	intent.Slots = []*Slot{intent.Slot0, intent.Slot1, intent.Slot2, intent.Slot3, intent.Slot4,
		intent.Slot5, intent.Slot6, intent.Slot7, intent.Slot8, intent.Slot9}
}

func (intent *Intent) reset() {
	for i := range intent.Slots {
		intent.Slots[i].reset()
	}
	intent.Value = ""
	intent.Sentence = ""
	intent.in_processing = false
}

func (slot *Slot) reset() {
	slot.Value = ""
	slot.filled = false
	slot.Indexes = nil
	slot.Indexes_str = ""
	slot.Sentence = ""
}

func (intent *Intent) isReady() bool {
	var any_slot bool = false

	var slots []*Slot = intent.Slots
	if slots == nil {
		slots = []*Slot{intent.Slot0, intent.Slot1, intent.Slot2, intent.Slot3, intent.Slot4, intent.Slot5,
			intent.Slot6, intent.Slot7, intent.Slot8, intent.Slot9}
	}

	// An intent is ready for processing if all its mandatory slots are filled
	for _, slot := range intent.Slots {
		if slot == nil {
			continue
		}

		any_slot = true
		if !slot.filled {
			return false
		}
	}

	if !any_slot {
		// Or if there are no slots, then the intent must have a value
		return intent.Value != ""
	}

	return true
}
