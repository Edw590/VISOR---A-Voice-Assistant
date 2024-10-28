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
	"ULHelper/ULHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"strings"
	"time"
)

var tasks_GL []ModsFileInfo.Task
var user_location_GL ModsFileInfo.UserLocation

var tasks_info_list_GL map[int]int64 = make(map[int]int64)

var conditions_were_true_GL map[int]bool = make(map[int]bool)

var prev_curr_last_known_user_loc_GL string = user_location_GL.Curr_location
var prev_prev_last_known_user_loc_GL string = user_location_GL.Prev_location

var stop_GL bool = false

var (
	modUserInfo_GL *ModsFileInfo.Mod9UserInfo = &Utils.User_settings_GL.TasksExecutor
)

/*
CheckDueTasks checks if any Task is due.

This function will block until a Task is due. When that happens, the Task is returned.

-----------------------------------------------------------

– Returns:
  - the Task that is due or nil if the checker was stopped
 */
func CheckDueTasks() *ModsFileInfo.Task {
	for {
		tasks_GL = modUserInfo_GL.Tasks

		// Location trigger - if the user location changed, check if any task is triggered
		var curr_last_known_user_loc string = user_location_GL.Curr_location
		var prev_last_known_user_loc string = user_location_GL.Prev_location
		if curr_last_known_user_loc != prev_curr_last_known_user_loc_GL || prev_last_known_user_loc != prev_prev_last_known_user_loc_GL {
			prev_curr_last_known_user_loc_GL = curr_last_known_user_loc
			prev_prev_last_known_user_loc_GL = prev_last_known_user_loc

			for _, task := range tasks_GL {
				// If the task has a time set or has no location, skip it
				if task.Time != "" || task.User_location == "" {
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

				var condition bool = checkCondition(task, conditions_were_true_GL)

				var device_id_matches bool = checkDeviceID(task)

				if condition_loc && condition && device_id_matches {
					return &task
				}
			}
		}

		// Time/condition trigger - if the time changed (it always does), check if any task is triggered
		for _, task := range tasks_GL {
			var condition_time bool = false
			var test_time int64 = 0
			// If the task has no time set, skip it
			if task.Time == "" {
				condition_time = true
			} else {
				var curr_time int64 = time.Now().Unix() / 60
				var task_time string = task.Time
				var format string = "2006-01-02 -- 15:04:05"
				t, _ := time.ParseInLocation(format, task_time, time.Local)
				test_time = t.Unix() / 60
				if task.Repeat_each_min > 0 {
					for {
						if test_time + task.Repeat_each_min <= curr_time {
							test_time += task.Repeat_each_min
						} else {
							break
						}
					}
				}

				condition_time = curr_time >= test_time && tasks_info_list_GL[task.Id] != test_time
			}

			// Check if the task is due and if it was already reminded

			var condition_loc bool = false
			if task.User_location == "" {
				condition_loc = true
			} else {
				// Check if the task has a location and the user is at that location.
				var curr_user_loc string = user_location_GL.Curr_location
				if curr_user_loc != ULHelper.UNKNOWN_LOCATION {
					condition_loc = checkLocation(task.User_location, curr_user_loc)
				}
			}

			var condition bool = checkCondition(task, conditions_were_true_GL)

			var device_id_matches bool = checkDeviceID(task)

			if condition_time && condition_loc && condition && device_id_matches {
				// Set the last reminded time to the test time
				tasks_info_list_GL[task.Id] = test_time

				return &task
			}
		}

		if Utils.WaitWithStopTIMEDATE(&stop_GL, 1) {
			return nil
		}
	}
}

/*
UpdateUserLocation updates the internal user location.

-----------------------------------------------------------

– Params:
  - user_location – the new user location
 */
func UpdateUserLocation(user_location *ModsFileInfo.UserLocation) {
	user_location_GL = *user_location
}

/*
StopChecker stops the CheckDueTasks function.
 */
func StopChecker() {
	stop_GL = true
}
