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

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

/*
GetStartOfDayDATETIME returns the timestamp at the start of the day for a given timestamp.

-----------------------------------------------------------

– Params:
  - timestamp – the timestamp in seconds to get the start of the day for

– Returns:
  - the timestamp at the start of the day
*/
func GetStartOfDayDATETIME(timestamp int64) int64 {
	return timestamp - timestamp % 86400;
}

/*
TimeDateToTimestampDATETIME converts a time and/or date string to a timestamp using NLP.

For information on supported formats and examples, check https://github.com/olebedev/when.

-----------------------------------------------------------

– Params:
  - time_date_str – the time and/or date string to convert

– Returns:
  - the timestamp in seconds, or -1 if the string could not be parsed
*/
func TimeDateToTimestampDATETIME(time_date_str string) int64 {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	r, err := w.Parse(time_date_str, time.Now())
	if err != nil || r == nil {
		return -1
	}

	return r.Time.Unix()
}

/*
ToReadableDurationDATETIME converts a duration in minutes to a human-readable string.

The largest unit used is decade(s), and the smallest is second(s).

A years is considered to have 365 days, and a month is considered to have 30 days.

-----------------------------------------------------------

– Params:
  - min – the duration in minutes

– Returns:
  - a human-readable string representing the duration
*/
func ToReadableDurationDATETIME(min int64) string {
	// Units in descending order
	units := []struct {
		nameSingular string
		namePlural   string
		sizeInMin    int64
	}{
		{"decade", "decades", 10 * 365 * 24 * 3600},
		{"year", "years", 365 * 24 * 3600},
		{"month", "months", 30 * 24 * 3600},
		{"week", "weeks", 7 * 24 * 3600},
		{"day", "days", 24 * 3600},
		{"hour", "hours", 3600},
		{"minute", "minutes", 60},
		{"second", "seconds", 1},
	}

	var parts []string
	var remaining int64 = min

	for _, u := range units {
		if remaining >= u.sizeInMin {
			var val int64 = remaining / u.sizeInMin
			remaining = remaining % u.sizeInMin

			var name string = u.namePlural
			if val == 1 {
				name = u.nameSingular
			}
			parts = append(parts, strconv.FormatInt(val, 10) + " " + name)
		}
	}

	if len(parts) == 0 {
		return "0 seconds"
	}

	// Join all parts with commas, but the last with "and"
	if len(parts) == 1 {
		return parts[0]
	}

	return strings.Join(parts[:len(parts)-1], ", ") + " and " + parts[len(parts)-1]
}

/*
ParseDurationDATETIME parses a duration string and returns the duration in seconds.

Supported units range from second(s) to decade(s). The function is case-insensitive.

A year is considered to have 365 days, and a month is considered to have 30 days.

-----------------------------------------------------------

– Params:
  - input – the duration string to parse

– Returns:
  - the duration in seconds or -1 if the input could not be parsed
*/
func ParseDurationDATETIME(input string) int64 {
	input = strings.ToLower(input)
	var re *regexp.Regexp = regexp.MustCompile(`(\d+)\s*(second|minute|hour|day|week|month|year|decade)s?`)
	var matches [][]string = re.FindAllStringSubmatch(input, -1)

	var duration int = 0
	var any_parsed bool = false
	for _, match := range matches {
		val, _ := strconv.Atoi(match[1])
		var unit string = match[2]

		switch unit {
			case "second":
				duration += val
				any_parsed = true
			case "minute":
				duration += val * 60
				any_parsed = true
			case "hour":
				duration += val * 3600
				any_parsed = true
			case "day":
				duration += val * 24 * 3600
				any_parsed = true
			case "week":
				duration += val * 7 * 24 * 3600
				any_parsed = true
			case "month":
				duration += val * 30 * 24 * 3600
				any_parsed = true
			case "year":
				duration += val * 365 * 24 * 3600
				any_parsed = true
			case "decade":
				duration += val * 10 * 365 * 24 * 3600
				any_parsed = true
		}
	}

	if any_parsed {
		return int64(duration)
	}

	return -1
}
