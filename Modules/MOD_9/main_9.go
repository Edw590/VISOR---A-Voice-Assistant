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

package MOD_9

import (
	"GPT/GPT"
	MOD_7 "GPTCommunicator"
	"Registry/Registry"
	MOD_12 "UserLocator"
	"Utils"
	"VISOR_Server/ServerRegKeys"
	"log"
	"strings"
	"time"
)

// Reminders Manager //

const TIME_SLEEP_S int = 1

type _MGIModSpecInfo _ModSpecInfo
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		var prev_curr_last_known_user_loc string = Registry.GetValue(ServerRegKeys.K_CURR_USER_LOCATION).GetString(true)
		var prev_prev_last_known_user_loc string = Registry.GetValue(ServerRegKeys.K_CURR_USER_LOCATION).GetString(false)
		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			// Add each reminder to the internal reminders list
			var list_modified bool = false
			var reminders_info_list map[string]int64 = moduleInfo_GL.ModGenInfo.ModSpecInfo.Reminders_info
			if reminders_info_list == nil {
				reminders_info_list = make(map[string]int64)
				moduleInfo_GL.ModGenInfo.ModSpecInfo.Reminders_info = reminders_info_list
				list_modified = true
			}
			for _, reminder := range modUserInfo.Reminders {
				if _, ok := reminders_info_list[reminder.Id]; !ok {
					reminders_info_list[reminder.Id] = 0
					list_modified = true
				}
			}
			if list_modified {
				_ = moduleInfo_GL.UpdateGenInfo()
			}

			// Location trigger - if the user location changed, check if any reminder is triggered
			var curr_last_known_user_loc string = Registry.GetValue(ServerRegKeys.K_CURR_USER_LOCATION).GetString(true)
			var prev_last_known_user_loc string = Registry.GetValue(ServerRegKeys.K_CURR_USER_LOCATION).GetString(false)
			if curr_last_known_user_loc != prev_curr_last_known_user_loc || prev_last_known_user_loc != prev_prev_last_known_user_loc {
				prev_curr_last_known_user_loc = curr_last_known_user_loc
				prev_prev_last_known_user_loc = prev_last_known_user_loc

				for _, reminder := range modUserInfo.Reminders {
					// If the reminder has a time set or has no location, skip it
					if reminder.Time != "" || reminder.Location == "" {
						continue
					}

					// In case there's a "+", the user must have arrived at the location. If there's a "-", the user
					// must have left the location.
					var condition bool
					if strings.HasPrefix(reminder.Location, "+") {
						var rem_loc string = reminder.Location[1:]
						condition = checkLocation(rem_loc, curr_last_known_user_loc)
					} else if strings.HasPrefix(reminder.Location, "-") {
						var rem_loc string = reminder.Location[1:]
						condition = checkLocation(rem_loc, prev_last_known_user_loc)
					} else {
						// Nothing to do
						continue
					}

					if condition {
						for {
							// Wait until the warning is sent
							if !MOD_7.SpeakOnDevice(GPT.ALL_DEVICES_ID, reminder.Message) {
								time.Sleep(1 * time.Second)
							} else {
								break
							}
						}

						log.Println("Reminder! Message: " + reminder.Message)
					}
				}
			}

			// Time trigger - if the time changed (it always does), check if any reminder is triggered
			for _, reminder := range modUserInfo.Reminders {
				// If the reminder has no time set, skip it
				if reminder.Time == "" {
					continue
				}

				var curr_time int64 = time.Now().Unix() / 60
				var reminder_time string = reminder.Time
				var format = time.RFC3339
				t, _ := time.Parse(format, reminder_time)
				var test_time int64 = t.Unix() / 60
				if reminder.Repeat_each > 0 {
					var repeat_each int64 = reminder.Repeat_each
					for {
						if test_time + repeat_each <= curr_time {
							test_time += repeat_each
						} else {
							break
						}
					}
				}

				// Check if the reminder is due
				var condition1 bool = curr_time >= test_time
				// Check if the reminder was already reminded
				var condition2 bool = reminders_info_list[reminder.Id] != test_time

				var condition_time bool = condition1 && condition2

				var condition_loc bool = true
				if reminder.Location != "" {
					// Check if the reminder has a location and the user is at that location.
					var curr_user_loc string = Registry.GetValue(ServerRegKeys.K_CURR_USER_LOCATION).GetString(true)
					if curr_user_loc != MOD_12.UNKNOWN_LOCATION {
						condition_loc = checkLocation(reminder.Location, curr_user_loc)
					}
				}

				if condition_time && condition_loc {
					for {
						if !MOD_7.SpeakOnDevice(GPT.ALL_DEVICES_ID, reminder.Message) {
							time.Sleep(1 * time.Second)
						} else {
							break
						}
					}

					log.Println("Reminder! Message: " + reminder.Message)

					reminders_info_list[reminder.Id] = test_time

					_ = moduleInfo_GL.UpdateGenInfo()
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, TIME_SLEEP_S) {
				return
			}
		}
	}
}

func checkLocation(reminder_loc string, location string) bool {
	if strings.HasSuffix(reminder_loc, "*") {
		// If the reminder location ends with a "*", it means that the user must be at a location that starts with the
		// reminder location.
		reminder_loc = reminder_loc[:len(reminder_loc) - 1]

		return strings.HasPrefix(location, reminder_loc)
	}

	return reminder_loc == location
}
