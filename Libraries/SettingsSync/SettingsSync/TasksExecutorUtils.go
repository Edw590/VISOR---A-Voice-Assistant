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

package SettingsSync

import (
	"Utils"
	"Utils/ModsFileInfo"
	"sort"
	"strconv"
	"strings"
)

/*
AddTask adds a task to the user settings.

-----------------------------------------------------------

– Params:
  - enabled – whether the task is enabled
  - device_active – whether the device is active
  - device_ids – the IDs of the devices separated by "\n"
  - message – the message of the task
  - command – the command to be executed
  - time – the time the task is set for
  - repeat_each_min – the time in minutes between each repeatition
  - user_location – the location the user must be in
  - programmable_condition – an additional condition for the task, in Go
 */
func AddTaskTASKS(enabled bool, device_active bool, device_ids string, message string, command string, time string,
			 repeat_each_min int64, user_location string, programmable_condition string) {
	var tasks []ModsFileInfo.Task = Utils.User_settings_GL.TasksExecutor.Tasks
	var task_id int = 1
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == task_id {
			task_id++
		}
	}

	// Add the task to the user settings
	Utils.User_settings_GL.TasksExecutor.Tasks = append(Utils.User_settings_GL.TasksExecutor.Tasks, ModsFileInfo.Task{
		Id:                     task_id,
		Enabled:                enabled,
		Device_active:          device_active,
		Device_IDs:             strings.Split(device_ids, "\n"),
		Message:                message,
		Command:                command,
		Time:                   time,
		Repeat_each_min:        repeat_each_min,
		User_location:          user_location,
		Programmable_condition: programmable_condition,
	})

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id < tasks[j].Id
	})
}

/*
RemoveTask removes a task from the user settings.

-----------------------------------------------------------

– Params:
  - id – the task ID
 */
func RemoveTaskTASKS(id int) {
	var tasks []ModsFileInfo.Task = Utils.User_settings_GL.TasksExecutor.Tasks
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == id {
			Utils.DelElemSLICES(&Utils.User_settings_GL.TasksExecutor.Tasks, i)

			break
		}
	}
}

/*
GetIdsList returns a list of all tasks' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all tasks' IDs separated by "|"
*/
func GetIdsListTASKS() string {
	var ids string
	for _, task := range Utils.User_settings_GL.TasksExecutor.Tasks {
		ids += strconv.Itoa(task.Id) + "|"
	}

	return ids
}

/*
GetTaskById returns a task by its ID.

-----------------------------------------------------------

– Params:
  - id – the task ID

– Returns:
  - the task or nil if the task was not found
*/
func GetTaskByIdTASKS(id int) *ModsFileInfo.Task {
	for _, task := range Utils.User_settings_GL.TasksExecutor.Tasks {
		if task.Id == id {
			return &task
		}
	}

	return nil
}
