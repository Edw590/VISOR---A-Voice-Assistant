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
	"regexp"
	"strconv"
	"time"
)

func parseDuration(input string) time.Duration {
	re := regexp.MustCompile(`(\d+)\s*(hour|minute|second|day|week)s?`)
	matches := re.FindAllStringSubmatch(input, -1)

	var duration time.Duration
	for _, match := range matches {
		val, _ := strconv.Atoi(match[1])
		unit := match[2]

		switch unit {
			case "second":
				duration += time.Duration(val) * time.Second
			case "minute":
				duration += time.Duration(val) * time.Minute
			case "hour":
				duration += time.Duration(val) * time.Hour
			case "day":
				duration += time.Duration(val) * 24 * time.Hour
			case "week":
				duration += time.Duration(val) * 7 * 24 * time.Hour
		}
	}

	return duration
}
