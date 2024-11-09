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

package GMan

import (
	"Utils"
	"Utils/ModsFileInfo"
)

var tasks_GL []ModsFileInfo.GTask = nil

func getTasks() {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GMan, []byte("JSON|true|GManTasks"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GMan]
	if comms_map == nil {
		return
	}

	var json_bytes []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	if err := Utils.FromJsonGENERAL(json_bytes, &tasks_GL); err != nil {
		return
	}
}

/*
GetEventsIdsListGMAN returns a list of all events' IDs.

This function will BLOCK FOREVER if there's no Internet connection! Check first with Utils.IsCommunicatorConnectedSERVER().

DON'T REQUEST LISTS FROM THIS LIBRARY IN DIFFERENT THREADS!!! All in the same thread always.

-----------------------------------------------------------

– Returns:
  - a list of all events' IDs separated by "|"
*/
func GetTasksIdsList() string {
	getTasks()

	var ids_list string = ""
	for _, task := range tasks_GL {
		ids_list += task.Id + "|"
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetEventGMAN returns an event by its ID.

-----------------------------------------------------------

– Params:
  - event_id – the event ID

– Returns:
  - the event or nil if the event was not found
*/
func GetTask(task_id string) *ModsFileInfo.GTask {
	for i := 0; i < len(tasks_GL); i++ {
		var task *ModsFileInfo.GTask = &tasks_GL[i]
		if task.Id == task_id {
			return task
		}
	}

	return nil
}
