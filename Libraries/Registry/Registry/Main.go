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
	"Utils"
)

const TYPE_BOOL string = "TYPE_BOOL"
const TYPE_INT string = "TYPE_INT"
const TYPE_LONG string = "TYPE_LONG"
const TYPE_FLOAT string = "TYPE_FLOAT"
const TYPE_DOUBLE string = "TYPE_DOUBLE"
const TYPE_STRING string = "TYPE_STRING"

// Value represents a value in the registry
type Value struct {
	// key is the key of the value
	key string
	// pretty_name is the pretty name of the value
	pretty_name string
	// description is the description of the value
	description string
	// type_ is the type of the value
	type_ string

	// prev_data is the previous data of the value
	prev_data string
	// time_updated_prev is the time the previous data was updated in milliseconds
	time_updated_prev int64
	// curr_data is the current data of the value
	curr_data string
	// time_updated_curr is the time the data was updated in milliseconds
	time_updated_curr int64
}

var registry_GL []*Value = nil

// There's an init() function on Keys.go

/*
RegisterValue registers a value in the registry.

-----------------------------------------------------------

– Params:
  - key – the key of the value
  - pretty_name – the pretty name of the value
  - description – the description of the value
  - value_type – the type of the value

– Returns:
  - the created value or nil if the value already exists
*/
func RegisterValue(key string, pretty_name string, description string, value_type string) *Value {
	if value := GetValue(key); value != nil {
		return value
	}

	var value *Value = &Value{
		key:          key,
		pretty_name:  pretty_name,
		description:  description,
		type_:        value_type,
	}

	switch value.type_ {
		case TYPE_BOOL:
			value.prev_data = "false"
			value.curr_data = "false"
		case TYPE_INT: fallthrough
		case TYPE_LONG: fallthrough
		case TYPE_FLOAT: fallthrough
		case TYPE_DOUBLE:
			value.prev_data = "-1"
			value.curr_data = "-1"
		case TYPE_STRING:
			value.prev_data = ""
	}


	registry_GL = append(registry_GL, value)

	return value
}

func GetValue(key string) *Value {
	for _, value := range registry_GL {
		if value.key == key {
			return value
		}
	}

	return nil
}

/*
RemoveValue removes a value from the registry based on its key.

-----------------------------------------------------------

– Params:
  - key – the key of the value
 */
func RemoveValue(key string) {
	for i, value := range registry_GL {
		if value.key == key {
			registry_GL = append(registry_GL[:i], registry_GL[i+1:]...)

			return
		}
	}
}

func GetRegistryText() string {
	var text string = ""

	for _, value := range registry_GL {
		text += "Name: " + value.pretty_name + "\n" +
				"Type: " + value.type_ + "\n" +
				"Prev time: " + Utils.GetDateTimeStrTIMEDATE(value.time_updated_prev) + "\n" +
				"Prev data: " + value.prev_data + "\n" +
				"Curr time: " + Utils.GetDateTimeStrTIMEDATE(value.time_updated_curr) + "\n" +
				"Curr data: " + value.curr_data + "\n" +
				"Description: " + value.description + "\n\n"
	}

	return text
}
