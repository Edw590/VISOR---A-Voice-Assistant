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

package ModsFileInfo

import "strings"

// Mod9GenInfo is the format of the custom generated information about this specific module.
type Mod9GenInfo struct {
	// Tasks_info is the information about the tasks
	Tasks_info []TaskInfo
	// Conds_were_true is the information about if programmable conditions that were true
	Conds_were_true []CondWasTrue
}

type TaskInfo struct {
	// Id is the task ID
	Id int32
	// Last_time_reminded is the last time the task was reminded in Unix minutes
	Last_time_reminded int64
}

type CondWasTrue struct {
	// Id is the task ID
	Id int32
	// Was_true is whether the programmable condition was true
	Was_true bool
}

///////////////////////////////////////////////////////////////////////////////

// Mod9UserInfo is the format of the custom information file about this specific module.
type Mod9UserInfo struct {
	// Tasks is the list of all tasks
	Tasks []Task
}

// Task is the format of a task
type Task struct {
	// Id is the task ID
	Id 		    int32
	// Enabled is whether the task is enabled
	Enabled     bool
	// Device_active is whether the device must be active to trigger the task
	Device_active bool
	// Device_IDs are the devices the task is set for, separated by "|"
	Device_IDs  []string
	// Message is the task message
	Message     string
	// Command is the command to be executed when the task is triggered on the chosen Device_IDs
	Command     string
	// Time_s is the time the task is set for in the format "2024-12-31 -- 23:59:59"
	Time_s int64
	// Repeat_each_min is the time in minutes between each repetition of the task
	Repeat_each_min int64
	// User_location is the location the user must be in for the task to be triggered
	User_location string
	// Programmable_condition is an additional "advanced" condition for the task in Go language
	Programmable_condition string
}

/*
GetDeviceIDs returns the Device_IDs separated by "\n".

-----------------------------------------------------------

– Returns:
  - the Device_IDs separated by "\n"
 */
func (task *Task) GetDeviceIDs() string {
	return strings.Join(task.Device_IDs, "\n")
}

/*
SetDeviceIDs sets the Device_IDs from a string separated by "\n".

-----------------------------------------------------------

– Params:
  - device_ids – the Device_IDs separated by "\n"
 */
func (task *Task) SetDeviceIDs(device_ids string) {
	task.Device_IDs = strings.Split(device_ids, "\n")
}
