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

package SettingsSync

import (
	"Utils"
	"Utils/ModsFileInfo"
	"sort"
	"strconv"
	"strings"
)

/*
AddTaskTASKS adds a task to the user settings.

-----------------------------------------------------------

– Params:
  - enabled – whether the task is enabled
  - device_active – whether the device is active
  - device_ids – the IDs of the devices separated by "\n"
  - message – the message of the task
  - command – the command to be executed
  - time – the time the task is set for
  - repeat_each_min – the time in minutes between each repetition
  - user_location – the location the user must be in
  - programmable_condition – an additional condition for the task, in Go

– Returns:
  - the ID of the task
 */
func AddTaskTASKS(enabled bool, device_active bool, device_ids string, message string, command string, time_s int64,
			 	  repeat_each_min int64, user_location string, programmable_condition string) int32 {
	var tasks *[]ModsFileInfo.Task = &Utils.GetUserSettings().TasksExecutor.Tasks
	var id int32 = 1 // 0 is reserved
	for i := 0; i < len(*tasks); i++ {
		if (*tasks)[i].Id == id {
			id++
			i = -1
		}
	}

	// Add the task to the user settings
	*tasks = append(*tasks, ModsFileInfo.Task{
		Id:                     id,
		Enabled:                enabled,
		Device_active:          device_active,
		Device_IDs:             strings.Split(device_ids, "\n"),
		Message:                message,
		Command:                command,
		Time_s:                 time_s,
		Repeat_each_min:        repeat_each_min,
		User_location:          user_location,
		Programmable_condition: programmable_condition,
	})

	sort.SliceStable(*tasks, func(i, j int) bool {
		return (*tasks)[i].Id < (*tasks)[j].Id
	})

	return id
}

/*
RemoveTaskTASKS removes a task from the user settings.

-----------------------------------------------------------

– Params:
  - id – the task ID
 */
func RemoveTaskTASKS(id int32) {
	var tasks *[]ModsFileInfo.Task = &Utils.GetUserSettings().TasksExecutor.Tasks
	for i := range *tasks {
		if (*tasks)[i].Id == id {
			Utils.DelElemSLICES(tasks, i)

			break
		}
	}
}

/*
GetIdsListTASKS returns a list of all tasks' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all tasks' IDs separated by "|"
*/
func GetIdsListTASKS() string {
	var ids_list string
	for _, task := range Utils.GetUserSettings().TasksExecutor.Tasks {
		ids_list += strconv.Itoa(int(task.Id)) + "|"
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetTaskTASKS returns a task by its ID.

-----------------------------------------------------------

– Params:
  - id – the task ID

– Returns:
  - the task or nil if the task was not found
*/
func GetTaskTASKS(id int32) *ModsFileInfo.Task {
	var tasks []ModsFileInfo.Task = Utils.GetUserSettings().TasksExecutor.Tasks
	for i := range tasks {
		var task *ModsFileInfo.Task = &tasks[i]
		if task.Id == id {
			return task
		}
	}

	return nil
}
