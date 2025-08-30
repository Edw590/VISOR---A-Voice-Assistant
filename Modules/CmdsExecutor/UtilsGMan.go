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

package CmdsExecutor

import (
	"GMan"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"time"
)

func getTasksList(tasks_ids []string, cmd_variant string) string {
	var speak string = ""

	for _, task_id := range tasks_ids {
		var task *ModsFileInfo.GTask = GMan.GetTask(task_id)
		if task == nil {
			continue
		}

		var add_task bool = false
		if task.Date_s == 0 {
			// If the task has no date, we add it to the list (it's to be done every day)
			add_task = true
		} else {
			// Else we check the date
			var task_date time.Time = time.Unix(task.Date_s, 0)
			switch cmd_variant {
				case RET_31_TODAY:
					if task_date.Day() == time.Now().Day() {
						add_task = true
					}
				case RET_31_TOMORROW:
					if task_date.Day() == time.Now().AddDate(0, 0, 1).Day() {
						add_task = true
					}
			}
		}

		if add_task {
			speak += "\"" + task.Title + "\"; "
		}
	}

	var when string
	if cmd_variant == RET_31_TODAY {
		when = "today"
	} else if cmd_variant == RET_31_TOMORROW {
		when = "tomorrow"
	}
	if speak == "" {
		speak = "You have no tasks found for " + when + "."
	} else {
		speak = "Your list of tasks for " + when + ": " + speak + "."
	}

	return speak
}

func getEventsList(events_ids []string, cmd_variant string) string {
	var speak string = ""

	for _, event_id := range events_ids {
		var event *ModsFileInfo.GEvent = GMan.GetEvent(event_id)
		if event == nil {
			continue
		}

		var event_end_time_s int64 = event.Start_time_s + event.Duration_min*60
		if event_end_time_s < time.Now().Unix() {
			// Event already ended
			continue
		}

		var now time.Time = time.Now()

		var start_of_day_s int64 = UtilsSWA.GetStartOfDayS(now.Unix())
		var end_of_day_s int64 = start_of_day_s + 86400 - 1 // 86400 seconds in a day (24*60*60)

		var start_of_next_day_s int64 = start_of_day_s + 86400
		var end_of_next_day_s int64 = start_of_next_day_s + 86400 - 1

		var start_of_week_s int64 = UtilsSWA.GetStartOfDayS(now.Unix() - int64(now.Weekday())*86400)
		var end_of_week_s int64 = start_of_week_s + 7*86400 - 1

		var start_of_next_week_s int64 = start_of_week_s + 7*86400
		var end_of_next_week_s int64 = start_of_next_week_s + 7*86400 - 1

		var add_event bool = false
		switch cmd_variant {
			case RET_31_TODAY:
				if (event.Start_time_s >= start_of_day_s && event.Start_time_s <= end_of_day_s) ||
					(event_end_time_s >= start_of_day_s && event_end_time_s <= end_of_day_s) ||
					(start_of_day_s >= event.Start_time_s && end_of_day_s <= event_end_time_s) {
					add_event = true
				}
			case RET_31_TOMORROW:
				if (event.Start_time_s >= start_of_next_day_s && event.Start_time_s <= end_of_next_day_s) ||
					(event_end_time_s >= start_of_next_day_s && event_end_time_s <= end_of_next_day_s) ||
					(start_of_next_day_s >= event.Start_time_s && end_of_next_day_s <= event_end_time_s) {
					add_event = true
				}
			case RET_31_THIS_WEEK:
				if (event.Start_time_s >= start_of_week_s && event.Start_time_s <= end_of_week_s) ||
					(event_end_time_s >= start_of_week_s && event_end_time_s <= end_of_week_s) ||
					(start_of_week_s >= event.Start_time_s && end_of_week_s <= event_end_time_s) {
					add_event = true
				}
			case RET_31_NEXT_WEEK:
				if (event.Start_time_s >= start_of_next_week_s && event.Start_time_s <= end_of_next_week_s) ||
					(event_end_time_s >= start_of_next_week_s && event_end_time_s <= end_of_next_week_s) ||
					(start_of_next_week_s >= event.Start_time_s && end_of_next_week_s <= event_end_time_s) {
					add_event = true
				}
		}
		if add_event {
			var event_date_time time.Time = time.Unix(event.Start_time_s, 0)

			var event_on string = ""
			if cmd_variant == RET_31_THIS_WEEK || cmd_variant == RET_31_NEXT_WEEK {
				event_on = " on " + event_date_time.Weekday().String()
			}

			var event_began_today bool = event.Start_time_s >= start_of_day_s && event.Start_time_s <= end_of_day_s
			var event_at string = ""
			if event_began_today {
				event_at = "at " + event_date_time.Format("15:04")
			}

			var curr_duration int64 = event.Start_time_s/60 + event.Duration_min - now.Unix()/60

			speak += "\"" + event.Summary + "\"" + event_on + " " + event_at + " for " +
				UtilsSWA.GetEventDuration(curr_duration) + "; "
		}
	}

	var when string
	if cmd_variant == RET_31_TODAY {
		when = "today"
	} else if cmd_variant == RET_31_TOMORROW {
		when = "tomorrow"
	} else if cmd_variant == RET_31_THIS_WEEK {
		when = "this week"
	} else if cmd_variant == RET_31_NEXT_WEEK {
		when = "next week"
	}
	if speak == "" {
		speak = "You have no events found for " + when + "."
	} else {
		speak = "Your list of events for " + when + ": " + speak + "."
	}

	return speak
}
