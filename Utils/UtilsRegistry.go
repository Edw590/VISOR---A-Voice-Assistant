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

package Utils

import (
	"strconv"
	"time"
)

const TYPE_BOOL string = "TYPE_BOOL"
const TYPE_INT string = "TYPE_INT"
const TYPE_LONG string = "TYPE_LONG"
const TYPE_FLOAT string = "TYPE_FLOAT"
const TYPE_DOUBLE string = "TYPE_DOUBLE"
const TYPE_STRING string = "TYPE_STRING"

// Value represents a value in the registry
type Value struct {
	// Key is the Key of the value
	Key string
	// Pretty_name is the pretty name of the value
	Pretty_name string
	// Description is the Description of the value
	Description string
	// Type is the type of the value
	Type string

	// Prev_data is the previous data of the value
	Prev_data string
	// Time_updated_prev is the time the previous data was updated in milliseconds
	Time_updated_prev int64
	// Curr_data is the current data of the value
	Curr_data string
	// Time_updated_curr is the time the data was updated in milliseconds
	Time_updated_curr int64
}

/////////////////////////////////////////////////////////
// Getters

/*
getInternal returns whether the data is internal and the data to return.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data
  - no_data – the data to return if there's no data or nil to return the default values

– Returns:
  - whether the no_data parameter was used
  - the data to return
*/
func (value *Value) getInternal(curr_data bool, no_data any) (bool, any) {
	if no_data != nil {
		if curr_data {
			if value.Time_updated_curr == 0 {
				return true, no_data
			}
		} else {
			if value.Time_updated_prev == 0 {
				return true, no_data
			}
		}
	}

	return false, nil
}

/*
GetTimeUpdated returns the time the data was updated in milliseconds.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data

– Returns:
  - the time the data was updated in milliseconds
*/
func (value *Value) GetTimeUpdated(curr_data bool) int64 {
	if curr_data {
		return value.Time_updated_curr
	} else {
		return value.Time_updated_prev
	}
}

/*
GetType returns the type of the Value.

-----------------------------------------------------------

– Returns:
  - the type of the Value
*/
func (value *Value) GetType() string {
	return value.Type
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
		data = value.Curr_data
	} else {
		data = value.Prev_data
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
		data = value.Curr_data
	} else {
		data = value.Prev_data
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
		data = value.Curr_data
	} else {
		data = value.Prev_data
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
		data = value.Curr_data
	} else {
		data = value.Prev_data
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
		data = value.Curr_data
	} else {
		data = value.Prev_data
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
		return value.Curr_data
	} else {
		return value.Prev_data
	}
}

/*
GetData returns the data of the Value.

-----------------------------------------------------------

– Params:
  - curr_data – true to get the current data, false to get the previous data
  - no_data – the data to return if there's no data or nil to return the default values
*/
func (value *Value) GetData(curr_data bool, no_data any) any {
	no_data_used, no_data_ret := value.getInternal(curr_data, no_data)
	if no_data_used {
		return no_data_ret
	}

	switch value.Type {
	case TYPE_BOOL:
		return value.GetBool(curr_data)
	case TYPE_INT:
		return value.GetInt(curr_data)
	case TYPE_LONG:
		return value.GetLong(curr_data)
	case TYPE_FLOAT:
		return value.GetFloat(curr_data)
	case TYPE_DOUBLE:
		return value.GetDouble(curr_data)
	case TYPE_STRING:
		return value.GetString(curr_data)
	}

	// Won't happen
	return nil
}

/////////////////////////////////////////////////////////
// Setters

/*
setInternal sets the internal variables for the value.
*/
func (value *Value) setInternal(new_data string) {
	if value.Curr_data != new_data {
		value.Prev_data = value.Curr_data
		value.Time_updated_prev = value.Time_updated_curr
	}

	value.Time_updated_curr = time.Now().UnixMilli()
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
	if value.Type != TYPE_BOOL {
		return false
	}

	var data_str string = strconv.FormatBool(data)
	if !update_if_same && value.Curr_data == data_str {
		return false
	}

	var new_data string
	if data {
		new_data = "true"
	} else {
		new_data = "false"
	}

	value.setInternal(new_data)
	value.Curr_data = new_data

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
	if value.Type != TYPE_INT {
		return false
	}

	var data_str string = strconv.Itoa(data)
	if !update_if_same && value.Curr_data == data_str {
		return false
	}

	var new_data string = strconv.Itoa(data)

	value.setInternal(new_data)
	value.Curr_data = new_data

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
	if value.Type != TYPE_LONG {
		return false
	}

	var data_str string = strconv.FormatInt(data, 10)
	if !update_if_same && value.Curr_data == data_str {
		return false
	}

	var new_data string = strconv.FormatInt(data, 10)

	value.setInternal(new_data)
	value.Curr_data = new_data

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
	if value.Type != TYPE_FLOAT {
		return false
	}

	var data_str string = strconv.FormatFloat(float64(data), 'f', -1, 32)
	if !update_if_same && value.Curr_data == data_str {
		return false
	}

	var new_data string = strconv.FormatFloat(float64(data), 'f', -1, 32)

	value.setInternal(new_data)
	value.Curr_data = new_data

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
	if value.Type != TYPE_DOUBLE {
		return false
	}

	var data_str string = strconv.FormatFloat(data, 'f', -1, 64)
	if !update_if_same && value.Curr_data == data_str {
		return false
	}

	var new_data string = strconv.FormatFloat(data, 'f', -1, 64)

	value.setInternal(new_data)
	value.Curr_data = new_data

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
	if value.Type != TYPE_STRING {
		return false
	}

	if !update_if_same && value.Curr_data == data {
		return false
	}

	value.setInternal(data)
	value.Curr_data = data

	return true
}

/*
SetData sets the value and converts it to the right type automatically.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
*/
func (value *Value) SetData(data any, update_if_same bool) bool {
	switch value.Type {
	case TYPE_BOOL:
		return value.SetBool(data.(bool), update_if_same)
	case TYPE_INT:
		return value.SetInt(data.(int), update_if_same)
	case TYPE_LONG:
		return value.SetLong(data.(int64), update_if_same)
	case TYPE_FLOAT:
		return value.SetFloat(data.(float32), update_if_same)
	case TYPE_DOUBLE:
		return value.SetDouble(data.(float64), update_if_same)
	case TYPE_STRING:
		return value.SetString(data.(string), update_if_same)
	}

	// Won't happen
	return false
}
