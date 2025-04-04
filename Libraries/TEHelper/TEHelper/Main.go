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

package TEHelper

import (
	"ULHelper/ULHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"strings"
	"time"
)

var prev_curr_last_known_user_loc_GL string = getUserLocation().Curr_location
var prev_prev_last_known_user_loc_GL string = getUserLocation().Prev_location

var stop_GL bool = false

/*
CheckDueTasks checks if any Task is due.

This function will block until a Task is due. When that happens, the Task is returned.

-----------------------------------------------------------

– Returns:
  - the Task that is due or nil if the checker was stopped
 */
func CheckDueTasks() *ModsFileInfo.Task {
	var task_return ModsFileInfo.Task
	stop_GL = false
	for {
		// Location trigger - if the user location changed, check if any task is triggered
		var curr_last_known_user_loc string = getUserLocation().Curr_location
		var prev_last_known_user_loc string = getUserLocation().Prev_location
		if curr_last_known_user_loc != prev_curr_last_known_user_loc_GL || prev_last_known_user_loc != prev_prev_last_known_user_loc_GL {
			prev_curr_last_known_user_loc_GL = curr_last_known_user_loc
			prev_prev_last_known_user_loc_GL = prev_last_known_user_loc

			for _, task := range getModUserSettings().Tasks {
				// If the task has a time set or has no location, skip it
				if !task.Enabled || task.Time_s != 0 || task.User_location == "" {
					continue
				}

				// In case there's a "+", the user must have arrived at the location. If there's a "-", the user
				// must have left the location.
				var condition_loc bool
				if strings.HasPrefix(task.User_location, "+") {
					var rem_loc string = task.User_location[1:]
					condition_loc = checkLocation(rem_loc, curr_last_known_user_loc)
				} else if strings.HasPrefix(task.User_location, "-") {
					var rem_loc string = task.User_location[1:]
					condition_loc = checkLocation(rem_loc, prev_last_known_user_loc)
				} else {
					// Nothing to do
					continue
				}

				var programmable_condition bool = checkProgrammableCondition(task)

				var device_id_matches bool = checkDeviceID(task)

				var condition_device_active bool = checkDeviceActive(task)

				var task_info *ModsFileInfo.TaskInfo = getTaskInfo(task.Id)
				if condition_loc && programmable_condition && device_id_matches && condition_device_active {
					if task_info.Last_time_reminded == 0 {
						if task_return.Id == 0 {
							// Only set the last reminded time if no other task was triggered
							task_info.Last_time_reminded = time.Now().Unix() / 60

							task_return = task
						}
					}
				} else {
					task_info.Last_time_reminded = 0
				}
			}
		}

		// Time/condition trigger - if the time changed (it always does), check if any task is triggered
		for _, task := range getModUserSettings().Tasks {
			if !task.Enabled {
				continue
			}

			condition_time, test_time_min := checkTime(task)

			var condition_loc bool = false
			if task.User_location == "" {
				condition_loc = true
			} else {
				// Check if the task has a location and the user is at that location.
				var curr_user_loc string = getUserLocation().Curr_location
				if curr_user_loc != ULHelper.UNKNOWN_LOCATION {
					condition_loc = checkLocation(task.User_location, curr_user_loc)
				}
			}

			var programmable_condition bool = checkProgrammableCondition(task)

			var device_id_matches bool = checkDeviceID(task)

			var condition_device_active bool = checkDeviceActive(task)

			if condition_time && condition_loc && programmable_condition && device_id_matches && condition_device_active {
				if task_return.Id == 0 {
					// Only set the last reminded time if no other task was triggered
					getTaskInfo(task.Id).Last_time_reminded = test_time_min

					task_return = task
				}
			}
		}

		if task_return.Id != 0 {
			return &task_return
		}

		if Utils.WaitWithStopTIMEDATE(&stop_GL, 1) {
			return nil
		}
	}
}

func getTaskInfo(task_id int32) *ModsFileInfo.TaskInfo {
	for i, task_info := range getModGenSettings().Tasks_info {
		if task_info.Id == task_id {
			return &getModGenSettings().Tasks_info[i]
		}
	}

	getModGenSettings().Tasks_info = append(getModGenSettings().Tasks_info, ModsFileInfo.TaskInfo{
		Id: task_id,
	})

	return &getModGenSettings().Tasks_info[len(getModGenSettings().Tasks_info) - 1]
}

/*
StopChecker stops the CheckDueTasks function.
 */
func StopChecker() {
	stop_GL = true
}

func getUserLocation() *ModsFileInfo.UserLocation {
	return &Utils.GetGenSettings().MOD_12.User_location
}

func getModGenSettings() *ModsFileInfo.Mod9GenInfo {
	return &Utils.GetGenSettings().MOD_9
}

func getModUserSettings() *ModsFileInfo.Mod9UserInfo {
	return &Utils.GetUserSettings().TasksExecutor
}
