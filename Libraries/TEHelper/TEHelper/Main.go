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

package TEHelper

import (
	MOD_3 "Speech"
	"SpeechQueue/SpeechQueue"
	MOD_12 "UserLocator"
	"Utils"
	"Utils/ModsFileInfo"
	"bytes"
	"log"
	"strings"
	"time"
)

const _TIME_SLEEP_S int = 1

const GET_TASKS_EACH_S int64 = 1 * 60
var last_get_tasks_when_s int64 = 0

var last_crc16_GL []byte = nil

var tasks_GL []ModsFileInfo.Task
var user_location_GL ModsFileInfo.UserLocation

var tasks_info_list map[string]int64 = make(map[string]int64)

var conditions_were_true_GL map[string]bool = make(map[string]bool)

var prev_curr_last_known_user_loc_GL string = user_location_GL.Curr_location
var prev_prev_last_known_user_loc_GL string = user_location_GL.Prev_location

var stop_GL bool = false

func CheckDueTasks() *ModsFileInfo.Task {
	for {
		if time.Now().Unix() >= last_get_tasks_when_s + GET_TASKS_EACH_S {
			UpdateLocalTasks() // TODO: RUN THIS IN ANOTHER THREAD!!!

			last_get_tasks_when_s = time.Now().Unix()
		}

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
					MOD_3.QueueSpeech(task.Message, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)

					log.Println("Task! -->", task.Message)

					// TODO: Execute command here

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

				condition_time = curr_time >= test_time && tasks_info_list[task.Id] != test_time
			}

			// Check if the task is due and if it was already reminded

			var condition_loc bool = false
			if task.User_location == "" {
				condition_loc = true
			} else {
				// Check if the task has a location and the user is at that location.
				var curr_user_loc string = user_location_GL.Curr_location
				if curr_user_loc != MOD_12.UNKNOWN_LOCATION {
					condition_loc = checkLocation(task.User_location, curr_user_loc)
				}
			}

			var condition bool = checkCondition(task, conditions_were_true_GL)

			var device_id_matches bool = checkDeviceID(task)

			if condition_time && condition_loc && condition && device_id_matches {
				MOD_3.QueueSpeech(task.Message, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)

				log.Println("Task! -->", task)

				// TODO: Execute command here

				// Set the last reminded time to the test time
				tasks_info_list[task.Id] = test_time

				// TODO: This must RETURN the task that just went off. Think about how to do it with multiple tasks
				//  triggering at the same time over and over again.

				return &task
			}
		}

		if Utils.WaitWithStopTIMEDATE(&stop_GL, _TIME_SLEEP_S) {
			return nil
		}
	}
}

func LoadLocalTasks(json string) {
	var p_tasks []ModsFileInfo.Task
	if err := Utils.FromJsonGENERAL([]byte(json), &p_tasks); err != nil {
		return
	}

	tasks_GL = p_tasks
}

func UpdateLocalTasks() string {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_TEHelper, []byte("File|true|tasks.json"))
	// TODO: This must be in another thread - will block if there's no Internet connection
	//  Like any other communication on the project!
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_TEHelper]
	if comms_map == nil {
		return ""
	}

	var new_crc16 []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
	if !bytes.Equal(new_crc16, last_crc16_GL) {
		last_crc16_GL = new_crc16

		var p_tasks *[]ModsFileInfo.Task = GetTasksList()
		if p_tasks == nil {
			return ""
		}

		tasks_GL = *p_tasks

		return *Utils.ToJsonGENERAL(tasks_GL)
	}

	return ""
}

func UpdateUserLocation(user_location *ModsFileInfo.UserLocation) {
	user_location_GL = *user_location
}

func StopChecker() {
	stop_GL = true
}
