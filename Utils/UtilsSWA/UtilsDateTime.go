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

package UtilsSWA

import "strconv"

/*
GetStartOfDayS returns the timestamp at the start of the day for a given timestamp.

-----------------------------------------------------------

– Params:
  - timestamp – the timestamp in seconds to get the start of the day for

– Returns:
  - the timestamp at the start of the day
*/
func GetStartOfDayS(timestamp int64) int64 {
	return timestamp - timestamp % 86400;
}

/*
GetEventDuration converts a duration in minutes to a human-readable string.

-----------------------------------------------------------

– Params:
  - min – the duration in minutes

– Returns:
  - a human-readable string representing the duration
 */
func GetEventDuration(min int64) string {
	if min >= 60 {
		if min >= 24*60 {
			if min >= 7*24*60 {
				weeks := min / (7 * 24 * 60)
				days := (min % (7 * 24 * 60)) / (24 * 60)
				var week_weeks string = "weeks"
				if weeks == 1 {
					week_weeks = "week"
				}
				var day_days string = "days"
				if days == 1 {
					day_days = "day"
				}
				if days > 0 {
					return strconv.Itoa(int(weeks)) + " " + week_weeks + " and " + strconv.Itoa(int(days)) + " " + day_days
				}
				return strconv.Itoa(int(weeks)) + " " + week_weeks
			}
			days := min / (24 * 60)
			hours := (min % (24 * 60)) / 60
			var day_days string = "days"
			if days == 1 {
				day_days = "day"
			}
			var hour_hours string = "hours"
			if hours == 1 {
				hour_hours = "hour"
			}
			if hours > 0 {
				return strconv.Itoa(int(days)) + " " + day_days + " and " + strconv.Itoa(int(hours)) + " " + hour_hours
			}
			return strconv.Itoa(int(days)) + " " + day_days
		}
		hours := min / 60
		minutes := min % 60
		var hour_hours string = "hours"
		if hours == 1 {
			hour_hours = "hour"
		}
		var minute_minutes string = "minutes"
		if minutes == 1 {
			minute_minutes = "minute"
		}
		if minutes > 0 {
			return strconv.Itoa(int(hours)) + " " + hour_hours + " and " + strconv.Itoa(int(minutes)) + " " + minute_minutes
		}
		return strconv.Itoa(int(hours)) + " " + hour_hours
	}

	var minute_minutes string = "minutes"
	if min == 1 {
		minute_minutes = "minute"
	}

	return strconv.Itoa(int(min)) + " " + minute_minutes
}
