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
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"github.com/apaxa-go/eval"
	"strconv"
	"strings"
)

func checkDeviceID(task ModsFileInfo.Task) bool {
	var device_IDs_split []string = strings.Split(task.Device_IDs, "|")
	if len(device_IDs_split) == 0 || device_IDs_split[0] == "3234_ALL" {
		return true
	}

	for _, device_id := range device_IDs_split {
		if device_id == Utils.Device_settings_GL.Device_ID {
			return true
		}
	}

	return false
}

func checkLocation(task_loc string, location string) bool {
	if strings.HasSuffix(task_loc, "*") {
		// If the task location ends with a "*", it means that the user must be at a location that starts with the
		// task location.
		task_loc = task_loc[:len(task_loc) - 1]

		return strings.HasPrefix(location, task_loc)
	}

	return task_loc == location
}

func computeCondition(condition string) bool {
	condition = formatCondition(condition)
	//log.Println("Condition:", condition)
	expr, err := eval.ParseString(condition, "")
	if err != nil {
		return false
	}
	r, err := expr.EvalToInterface(nil)
	if err != nil {
		return false
	}

	return r.(bool)
}

func formatCondition(condition string) string {
	var registry_values []UtilsSWA.Value = UtilsSWA.GetValuesREGISTRY()
	for _, value := range registry_values {
		var value_str string = "ERROR"
		if value.Type_ == UtilsSWA.TYPE_STRING {
			value_str = value.GetString(true)
		} else if value.Type_ == UtilsSWA.TYPE_INT {
			value_str = strconv.Itoa(value.GetInt(true))
		} else if value.Type_ == UtilsSWA.TYPE_LONG {
			value_str = strconv.Itoa(int(value.GetLong(true)))
		} else if value.Type_ == UtilsSWA.TYPE_BOOL {
			value_str = strconv.FormatBool(value.GetBool(true))
		} else if value.Type_ == UtilsSWA.TYPE_FLOAT {
			value_str = strconv.FormatFloat(float64(value.GetFloat(true)), 'f', -1, 32)
		} else if value.Type_ == UtilsSWA.TYPE_DOUBLE {
			value_str = strconv.FormatFloat(value.GetDouble(true), 'f', -1, 64)
		}

		condition = strings.Replace(condition, strings.ToLower(value.Key), value_str, -1)
	}

	return condition
}

func checkCondition(task ModsFileInfo.Task, conditions_were_true map[int]bool) bool {
	var condition bool = false
	if task.Device_condition != "" {
		if ok := conditions_were_true[task.Id]; !ok {
			conditions_were_true[task.Id] = false
		}

		if computeCondition(task.Device_condition) {
			if !conditions_were_true[task.Id] {
				conditions_were_true[task.Id] = true

				condition = true
			}
		} else {
			conditions_were_true[task.Id] = false
		}
	} else {
		condition = true
	}

	return condition
}
