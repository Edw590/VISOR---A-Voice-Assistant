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
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"github.com/apaxa-go/eval"
	"strconv"
	"strings"
	"time"
)

func checkTime(task ModsFileInfo.Task) (bool, int64) {
	var test_time_min int64 = 0
	// If the task has no time set, skip it
	if task.Time == "" {
		return true, 0
	} else {
		var curr_time int64 = time.Now().Unix() / 60
		var task_time string = task.Time
		var format string = "2006-01-02 -- 15:04:05"
		t, err := time.ParseInLocation(format, task_time, time.Local)
		if err != nil {
			return false, 0
		}

		test_time_min = t.Unix() / 60
		if task.Repeat_each_min > 0 {
			for {
				if test_time_min + task.Repeat_each_min <= curr_time {
					test_time_min += task.Repeat_each_min
				} else {
					break
				}
			}
		}

		return curr_time >= test_time_min && getTaskInfo(task.Id).Last_time_reminded < test_time_min, test_time_min
	}
}

func checkDeviceActive(task ModsFileInfo.Task) bool {
	if !task.Device_active {
		return true
	}

	return Utils.Gen_settings_GL.MOD_10.Device_info.Last_time_used_s + 5 >= time.Now().Unix()
}

func checkDeviceID(task ModsFileInfo.Task) bool {
	if len(task.Device_IDs) == 0 || task.Device_IDs[0] == "3234_ALL" {
		return true
	}

	for _, device_id := range task.Device_IDs {
		if device_id == Utils.Gen_settings_GL.Device_settings.Id {
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

/*
ComputeCondition computes the result of a condition.

-----------------------------------------------------------

- Params:
  - condition – the condition to compute

- Returns:
  - the result of the condition
  - an error if the condition is invalid
 */
func ComputeCondition(condition string) (bool, error) {
	condition = formatCondition(condition)
	expr, err := eval.ParseString(condition, "")
	if err != nil {
		return false, err
	}
	r, err := expr.EvalToInterface(nil)
	if err != nil {
		return false, err
	}

	return r.(bool), nil
}

func formatCondition(condition string) string {
	var registry_values []*UtilsSWA.Value = UtilsSWA.GetValuesREGISTRY()
	for _, value := range registry_values {
		var value_str string = "ERROR"
		if value.Type_ == UtilsSWA.TYPE_STRING {
			value_str = value.GetString(true)
		} else if value.Type_ == UtilsSWA.TYPE_INT {
			value_str = strconv.FormatInt(int64(value.GetInt(true)), 10)
		} else if value.Type_ == UtilsSWA.TYPE_LONG {
			value_str = strconv.FormatInt(value.GetLong(true), 10)
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

func checkProgrammableCondition(task ModsFileInfo.Task) bool {
	var condition bool = false
	if task.Programmable_condition != "" {
		var cond_was_true *ModsFileInfo.CondWasTrue = getCondWasTrue(task.Id)

		cond_result, _ := ComputeCondition(task.Programmable_condition)
		if cond_result {
			if !cond_was_true.Was_true {
				cond_was_true.Was_true = true

				condition = true
			}
		} else {
			cond_was_true.Was_true = false
		}
	} else {
		condition = true
	}

	return condition
}

func getCondWasTrue(task_id int32) *ModsFileInfo.CondWasTrue {
	for i, cond_was_true := range modGenInfo_GL.Conds_were_true {
		if cond_was_true.Id == task_id {
			return &modGenInfo_GL.Conds_were_true[i]
		}
	}

	modGenInfo_GL.Conds_were_true = append(modGenInfo_GL.Conds_were_true, ModsFileInfo.CondWasTrue{
		Id:       task_id,
	})

	return &modGenInfo_GL.Conds_were_true[len(modGenInfo_GL.Conds_were_true) - 1]
}
