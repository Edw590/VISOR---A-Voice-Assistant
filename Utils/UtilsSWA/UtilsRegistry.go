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

package UtilsSWA

import (
	"Utils"
)

const TYPE_BOOL string = Utils.TYPE_BOOL
const TYPE_INT string = Utils.TYPE_INT
const TYPE_LONG string = Utils.TYPE_LONG
const TYPE_FLOAT string = Utils.TYPE_FLOAT
const TYPE_DOUBLE string = Utils.TYPE_DOUBLE
const TYPE_STRING string = Utils.TYPE_STRING

var registry_GL *[]*Utils.Value = &Utils.Gen_settings_GL.Registry

// There's an init() function on Keys.go

/*
RegisterValueREGISTRY registers a value in the registry.

-----------------------------------------------------------

– Params:
  - key – the key of the value
  - pretty_name – the pretty name of the value
  - description – the description of the value
  - value_type – the type of the value

– Returns:
  - the created value or nil if the value already exists
*/
func RegisterValueREGISTRY(key string, pretty_name string, description string, value_type string) *Utils.Value {
	if value := GetValueREGISTRY(key); value != nil {
		return value
	}

	var value *Utils.Value = &Utils.Value{
		Key:          key,
		Pretty_name:  pretty_name,
		Description:  description,
		Type:         value_type,
	}

	switch value.Type {
		case TYPE_BOOL:
			value.Prev_data = "false"
			value.Curr_data = "false"
		case TYPE_INT: fallthrough
		case TYPE_LONG: fallthrough
		case TYPE_FLOAT: fallthrough
		case TYPE_DOUBLE:
			value.Prev_data = "-1"
			value.Curr_data = "-1"
		case TYPE_STRING:
			value.Prev_data = ""
	}


	*registry_GL = append(*registry_GL, value)

	return value
}

func GetValueREGISTRY(key string) *Utils.Value {
	for _, value := range *registry_GL {
		if value.Key == key {
			return value
		}
	}

	return nil
}

/*
RemoveValueREGISTRY removes a value from the registry based on its key.

-----------------------------------------------------------

– Params:
  - key – the key of the value
*/
func RemoveValueREGISTRY(key string) {
	for i, value := range *registry_GL {
		if value.Key == key {
			*registry_GL = append((*registry_GL)[:i], (*registry_GL)[i+1:]...)

			return
		}
	}
}

func GetRegistryTextREGISTRY() string {
	var text string = ""

	for _, value := range *registry_GL {
		text += "Name: " + value.Pretty_name + "\n" +
				"Type: " + value.Type + "\n" +
				"Prev time: " + Utils.GetDateTimeStrTIMEDATE(value.Time_updated_prev) + "\n" +
				"Prev data: " + value.Prev_data + "\n" +
				"Curr time: " + Utils.GetDateTimeStrTIMEDATE(value.Time_updated_curr) + "\n" +
				"Curr data: " + value.Curr_data + "\n" +
				"Description: " + value.Description + "\n\n"
	}

	return text
}
