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
	"Utils"
	"strconv"
	"strings"
	"time"
)

const TYPE_BOOL string = "TYPE_BOOL"
const TYPE_INT string = "TYPE_INT"
const TYPE_LONG string = "TYPE_LONG"
const TYPE_FLOAT string = "TYPE_FLOAT"
const TYPE_DOUBLE string = "TYPE_DOUBLE"
const TYPE_STRING string = "TYPE_STRING"

var keys_added_GL []string = nil

// Value represents a value in the registry
type Value struct {
	// Key is the Key of the value
	Key string
	// Pretty_name is the pretty name of the value
	Pretty_name string
	// Description is the Description of the value
	Description string
	// Type_ is the type of the value
	Type_ string
	// Auto_set is true if the value is automatically set by VISOR, false if it is set by the user
	Auto_set bool

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
// Normal functions

/*
RegisterValueREGISTRY registers a Value in the Registry.

In case the value already exists, the pretty name and the description will be updated with the given ones.

After registering all necessary Values, call CleanRegistryREGISTRY().

-----------------------------------------------------------

– Params:
  - key – the key of the Value
  - pretty_name – the pretty name of the Value
  - description – the description of the Value
  - value_type – the type of the Value
  - init_data – the initial current data of the Value

– Returns:
  - the created value or nil if the Value already exists
*/
func RegisterValueREGISTRY(key string, pretty_name string, description string, value_type string, init_data string,
						   auto_set bool) *Value {
	keys_added_GL = append(keys_added_GL, key)
	if value := GetValueREGISTRY(key); value != nil {
		if value.Time_updated_curr == 0 {
			// If the value was never changed, update the initial data in case there was a change.
			switch value.Type_ {
				case TYPE_BOOL:
					if init_data == "" {
						value.Curr_data = "false"
					} else {
						value.Curr_data = init_data
					}
				case TYPE_INT: fallthrough
				case TYPE_LONG: fallthrough
				case TYPE_FLOAT: fallthrough
				case TYPE_DOUBLE:
					if init_data == "" {
						value.Curr_data = "-1"
					} else {
						value.Curr_data = init_data
					}
				case TYPE_STRING:
					value.Curr_data = init_data
			}
		}
		value.Pretty_name = pretty_name
		value.Description = description

		return nil
	}

	var value *Value = &Value{
		Key:         key,
		Pretty_name: pretty_name,
		Description: description,
		Type_:       value_type,
		Auto_set:    auto_set,
	}

	switch value.Type_ {
		case TYPE_BOOL:
			value.Prev_data = "false"
			if init_data == "" {
				value.Curr_data = "false"
			} else {
				value.Curr_data = init_data
			}
		case TYPE_INT: fallthrough
		case TYPE_LONG: fallthrough
		case TYPE_FLOAT: fallthrough
		case TYPE_DOUBLE:
			value.Prev_data = "-1"
			if init_data == "" {
				value.Curr_data = "-1"
			} else {
				value.Curr_data = init_data
			}
		case TYPE_STRING:
			value.Prev_data = ""
			value.Curr_data = init_data
	}

	Utils.GetGenSettings().Registry = append(Utils.GetGenSettings().Registry, (*Utils.Value) (value))

	return value
}

/*
CleanRegistryREGISTRY cleans the registry by removing old unused Values that were not registered in the current session.

Useful for when a Value is removed from the code but is still on users' Gen Settings.

Call this after registering all necessary Values.
 */
func CleanRegistryREGISTRY() {
	var registry []*Utils.Value
	for _, value := range Utils.GetGenSettings().Registry {
		for _, key := range keys_added_GL {
			if value.Key == key {
				registry = append(registry, value)
				break
			}
		}
	}
	Utils.GetGenSettings().Registry = registry
}

/*
GetValueREGISTRY gets a value from the registry based on its key.

-----------------------------------------------------------

– Params:
  - key – the key of the value

– Returns:
  - the value or nil if the value doesn't exist
 */
func GetValueREGISTRY(key string) *Value {
	for _, value := range Utils.GetGenSettings().Registry {
		if value.Key == key {
			return (*Value) (value)
		}
	}

	return nil
}

/*
GetValuesREGISTRY gets all the values in the registry.

-----------------------------------------------------------

– Returns:
  - all the values in the registry
 */
func GetValuesREGISTRY() []*Value {
	var values []*Value
	for _, value := range Utils.GetGenSettings().Registry {
		values = append(values, (*Value)(value))
	}

	return values
}

/*
GetKeysREGISTRY gets all the keys in the registry.

-----------------------------------------------------------

– Returns:
  - all the keys in the registry separated by "|"
 */
func GetKeysREGISTRY() string {
	var keys string = ""

	for _, value := range Utils.GetGenSettings().Registry {
		keys += value.Key + "|"
	}
	keys = keys[:len(keys) - 1]

	return keys
}

/*
RemoveValueREGISTRY removes a value from the registry based on its key.

-----------------------------------------------------------

– Params:
  - key – the key of the value
*/
func RemoveValueREGISTRY(key string) {
	for i, value := range Utils.GetGenSettings().Registry {
		if value.Key == key {
			Utils.DelElemSLICES(&Utils.GetGenSettings().Registry, i)

			break
		}
	}
	for i, key_added := range keys_added_GL {
		if key_added == key {
			Utils.DelElemSLICES(&keys_added_GL, i)

			break
		}
	}
}

/*
GetRegistryTextREGISTRY returns a text representation of the Registry.

-----------------------------------------------------------

– Params:
  - type_ – the type of the values to include (0 for all, 1 for only auto-set values, 2 for only manually-set values)

– Returns:
  - a text representation of the Registry
 */
func GetRegistryTextREGISTRY(type_ int) string {
	var text string = ""

	for _, value := range Utils.GetGenSettings().Registry {
		if type_ == 1 && !value.Auto_set {
			continue
		} else if type_ == 2 && value.Auto_set {
			continue
		}
		text += "Name: " + value.Pretty_name + "\n" +
			"Key: " + value.Key + "\n" +
			"Auto set: " + strconv.FormatBool(value.Auto_set) + "\n" +
			"Type: " + strings.ToLower(value.Type_[len("TYPE_"):]) + "\n" +
			"Prev time: " + Utils.GetDateTimeStrTIMEDATE(value.Time_updated_prev / 1000) + "\n" +
			"Prev data: " + value.Prev_data + "\n" +
			"Curr time: " + Utils.GetDateTimeStrTIMEDATE(value.Time_updated_curr / 1000) + "\n" +
			"Curr data: " + value.Curr_data + "\n" +
			"Description: " + value.Description + "\n\n"
	}

	return text
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
	return value.Type_
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
func (value *Value) GetInt(curr_data bool) int32 {
	var data string
	if curr_data {
		data = value.Curr_data
	} else {
		data = value.Prev_data
	}

	i, err := strconv.ParseInt(data, 10, 32)
	if err != nil {
		return -1
	}

	return int32(i)
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

	switch value.Type_ {
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
	if value.Type_ != TYPE_BOOL {
		return false
	}

	var new_data string = strconv.FormatBool(data)
	if !update_if_same && value.Curr_data == new_data {
		return false
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
func (value *Value) SetInt(data int32, update_if_same bool) bool {
	if value.Type_ != TYPE_INT {
		return false
	}

	var new_data string = strconv.FormatInt(int64(data), 10)
	if !update_if_same && value.Curr_data == new_data {
		return false
	}

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
	if value.Type_ != TYPE_LONG {
		return false
	}

	var new_data string = strconv.FormatInt(data, 10)
	if !update_if_same && value.Curr_data == new_data {
		return false
	}

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
	if value.Type_ != TYPE_FLOAT {
		return false
	}

	var new_data string = strconv.FormatFloat(float64(data), 'f', -1, 32)
	if !update_if_same && value.Curr_data == new_data {
		return false
	}

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
	if value.Type_ != TYPE_DOUBLE {
		return false
	}

	var new_data string = strconv.FormatFloat(data, 'f', -1, 64)
	if !update_if_same && value.Curr_data == new_data {
		return false
	}

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
	if value.Type_ != TYPE_STRING {
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
SetData sets the value directly in string without conversion.

-----------------------------------------------------------

- Params:
  - data – the data to set
  - update_if_same – whether to still update if the data is the same

- Returns:
  - whether the data was set
*/
func (value *Value) SetData(data string, update_if_same bool) bool {
	if !update_if_same && value.Curr_data == data {
		return false
	}

	// Check if the data is valid
	switch value.Type_ {
		case TYPE_BOOL:
			if _, err := strconv.ParseBool(data); err != nil {
				return false
			}
		case TYPE_INT:
			if _, err := strconv.ParseInt(data, 10, 32); err != nil {
				return false
			}
		case TYPE_LONG:
			if _, err := strconv.ParseInt(data, 10, 64); err != nil {
				return false
			}
		case TYPE_FLOAT:
			if _, err := strconv.ParseFloat(data, 32); err != nil {
				return false
			}
		case TYPE_DOUBLE:
			if _, err := strconv.ParseFloat(data, 64); err != nil {
				return false
			}
	}

	value.setInternal(data)
	value.Curr_data = data

	return true
}
