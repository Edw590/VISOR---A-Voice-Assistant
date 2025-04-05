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

package Utils

import (
	"time"
)

const TIME_FORMAT string = "15:04:05 (MST)"
const DATE_FORMAT string = "Monday 2006-01-02"
const DATE_TIME_FORMAT string = DATE_FORMAT + " -- " + TIME_FORMAT

/*
GetDateTimeStrTIMEDATE gets the current time and date in the format DATE_TIME_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the time in seconds or -1 for the current time

– Returns:
  - the current time and date in the default format
*/
func GetDateTimeStrTIMEDATE(s int64) string {
	return getTimeDateInFormatTIMEDATE(s, DATE_TIME_FORMAT)
}

/*
GetDateStrTIMEDATE gets the current date in the format DATE_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the time in seconds or -1 for the current time

– Returns:
  - the current time in the default format
*/
func GetDateStrTIMEDATE(s int64) string {
	return getTimeDateInFormatTIMEDATE(s, DATE_FORMAT)
}

/*
GetTimeStrTIMEDATE gets the current time in the format TIME_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the time in seconds or -1 for the current time

– Returns:
  - the current date in the default format
*/
func GetTimeStrTIMEDATE(s int64) string {
	return getTimeDateInFormatTIMEDATE(s, TIME_FORMAT)
}

/*
getTimeDateInFormatTIMEDATE gets the time and/or date in the given format.

-----------------------------------------------------------

– Params:
  - millis – the time in milliseconds or -1 for the current time
  - format – the format to use

– Returns:
  - the time and/or date in the given format
*/
func getTimeDateInFormatTIMEDATE(millis int64, format string) string {
	if millis == -1 {
		return time.Now().Format(format)
	} else {
		return time.Unix(0, millis*1e6).Format(format)
	}
}

/*
WaitWithStopTIMEDATE waits for a certain amount of time or until a stop signal is received (checked every second).

-----------------------------------------------------------

– Params:
  - stop – the stop signal
  - time_sleep_s – the time to wait in seconds

– Returns:
  - true if the loop was stopped, false if it reached the end time
*/
func WaitWithStopTIMEDATE(stop *bool, time_wait_s int) bool {
	if *stop {
		// In case time_wait_s is 0, it returns immediately in case the stop signal has been given
		return true
	}

	var time_end int64 = time.Now().Unix() + int64(time_wait_s)
	var stopped bool = false
	for time.Now().Unix() < time_end {
		if *stop {
			stopped = true

			break
		}

		time.Sleep(1 * time.Second)
	}

	return stopped
}
