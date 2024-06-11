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

package Registry

import "strconv"

func (value *Value) GetPrettyName() string {
	return value.pretty_name
}

func (value *Value) GetDescription() string {
	return value.description
}

func (value *Value) GetType() int {
	return value.type_
}

func (value *Value) GetTimeUpdated(curr_data bool) int64 {
	if curr_data {
		return value.time_updated_curr
	} else {
		return value.time_updated_prev
	}
}

/*
GetBool returns the boolean value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the boolean value of the Value
 */
func (value *Value) GetBool(curr_data bool) bool {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseBool(data)
	if err != nil {
		return false
	}

	return i
}

/*
GetInt returns the integer value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the integer value of the Value
 */
func (value *Value) GetInt(curr_data bool) int {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.Atoi(data)
	if err != nil {
		return -1
	}

	return i
}

/*
GetLong returns the long value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the long value of the Value
 */
func (value *Value) GetLong(curr_data bool) int64 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return -1
	}

	return i
}

/*
GetFloat returns the float value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the float value of the Value
 */
func (value *Value) GetFloat(curr_data bool) float32 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseFloat(data, 32)
	if err != nil {
		return -1
	}

	return float32(i)
}

/*
GetDouble returns the double value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the double value of the Value
 */
func (value *Value) GetDouble(curr_data bool) float64 {
	var data string
	if curr_data {
		data = value.curr_data
	} else {
		data = value.prev_data
	}

	i, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return -1
	}

	return i
}

/*
GetString returns the string value of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the string value of the Value
 */
func (value *Value) GetString(curr_data bool) string {
	if curr_data {
		return value.curr_data
	} else {
		return value.prev_data
	}
}
