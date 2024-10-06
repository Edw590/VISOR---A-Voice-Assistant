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

package ModsFileInfo

// Mod9GenInfo is the format of the custom generated information about this specific module.
type Mod9GenInfo struct {
	// Tasks_info maps the task ID to the last time the task was reminded in Unix minutes
	Tasks_info map[string]int64
	// Tasks is the list of all tasks
	Tasks []Task
}

// Task is the format of a task
type Task struct {
	// Id is the task ID
	Id 		    string
	// Device_IDs are the devices the task is set for
	Device_IDs  []string
	// Message is the task message
	Message     string
	// Command is the command to be executed when the task is triggered on the chosen Device_IDs
	Command     string
	// Time is the time the task is set for in the format "2024-12-31 -- 23:59:59"
	Time        string
	// Repeat_each_min is the time in minutes between each repeatition of the task
	Repeat_each_min int64
	// User_location is the location the user must be in for the task to be triggered
	User_location string
	// Device_condition is an additional "advanced" condition for the task in Go language
	Device_condition string
}
