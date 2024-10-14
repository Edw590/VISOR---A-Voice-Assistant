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

package TEHelper

import (
	"Utils"
	"Utils/ModsFileInfo"
)

/*
GetTasksList returns a list of all reminders.

-----------------------------------------------------------

– Returns:
  - a list of all reminders
*/
func GetTasksList() *[]ModsFileInfo.Task {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_TEHelper, []byte("File|false|tasks.json"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_TEHelper]
	if comms_map == nil {
		return nil
	}

	var file_contents []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	var reminders []ModsFileInfo.Task
	if err := Utils.FromJsonGENERAL(file_contents, &reminders); err != nil {
		return nil
	}

	return &reminders
}

/*
GetIdsList returns a list of all reminders' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all reminders' IDs separated by "|"
*/
func GetIdsList() string {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_TEHelper, []byte("File|false|tasks.json"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_TEHelper]
	if comms_map == nil {
		return ""
	}

	var file_contents []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	var user_location []ModsFileInfo.Task
	if err := Utils.FromJsonGENERAL(file_contents, &user_location); err != nil {
		return ""
	}

	var ids string
	for _, reminder := range user_location {
		ids += reminder.Id + "|"
	}

	return ids
}

/*
GetTaskById returns a reminder by its ID.

-----------------------------------------------------------

– Params:
  - id – the reminder ID

– Returns:
  - the reminder or nil if the reminder was not found
*/
func GetTaskById(id string) *ModsFileInfo.Task {
	var p_reminders *[]ModsFileInfo.Task = GetTasksList()
	if p_reminders == nil {
		return nil
	}

	for _, reminder := range *p_reminders {
		if reminder.Id == id {
			return &reminder
		}
	}

	return nil
}
