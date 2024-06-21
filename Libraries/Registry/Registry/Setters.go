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

import (
	"strconv"
	"time"
)

/*
setInternal sets the internal variables for the value.
 */
func (value *Value) setInternal() {
	value.prev_data = value.curr_data
	value.time_updated_prev = value.time_updated_curr

	value.time_updated_curr = time.Now().UnixMilli()
}

/*
SetBool sets the value to a boolean.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetBool(data bool, update_if_same bool) bool {
	if value.type_ != TYPE_BOOL {
		return false
	}

	var data_str string = strconv.FormatBool(data)
	if !update_if_same && value.curr_data == data_str {
		return false
	}

	value.setInternal()

	if data {
		value.curr_data = "true"
	} else {
		value.curr_data = "false"
	}

	return true
}

/*
SetInt sets the value to an integer.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetInt(data int, update_if_same bool) bool {
	if value.type_ != TYPE_INT {
		return false
	}

	var data_str string = strconv.Itoa(data)
	if !update_if_same && value.curr_data == data_str {
		return false
	}

	value.setInternal()

	value.curr_data = strconv.Itoa(data)

	return true
}

/*
SetLong sets the value to a long.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetLong(data int64, update_if_same bool) bool {
	if value.type_ != TYPE_LONG {
		return false
	}

	var data_str string = strconv.FormatInt(data, 10)
	if !update_if_same && value.curr_data == data_str {
		return false
	}

	value.setInternal()

	value.curr_data = strconv.FormatInt(data, 10)

	return true
}

/*
SetFloat sets the value to a float.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetFloat(data float32, update_if_same bool) bool {
	if value.type_ != TYPE_FLOAT {
		return false
	}

	var data_str string = strconv.FormatFloat(float64(data), 'f', -1, 32)
	if !update_if_same && value.curr_data == data_str {
		return false
	}

	value.setInternal()

	value.curr_data = strconv.FormatFloat(float64(data), 'f', -1, 32)

	return true
}

/*
SetDouble sets the value to a double.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetDouble(data float64, update_if_same bool) bool {
	if value.type_ != TYPE_DOUBLE {
		return false
	}

	var data_str string = strconv.FormatFloat(data, 'f', -1, 64)
	if !update_if_same && value.curr_data == data_str {
		return false
	}

	value.setInternal()

	value.curr_data = strconv.FormatFloat(data, 'f', -1, 64)

	return true
}

/*
SetString sets the value to a string.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
 */
func (value *Value) SetString(data string, update_if_same bool) bool {
	if value.type_ != TYPE_STRING {
		return false
	}

	if !update_if_same && value.curr_data == data {
		return false
	}

	value.setInternal()

	value.curr_data = data

	return true
}
