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

const TYPE_BOOL int = 0
const TYPE_INT int = 1
const TYPE_LONG int = 2
const TYPE_FLOAT int = 3
const TYPE_DOUBLE int = 4
const TYPE_STRING int = 5

// Value represents a value in the registry.
type Value struct {
	// key is the key of the value.
	key string
	// pretty_name is the pretty name of the value.
	pretty_name string
	// description is the description of the value.
	description string
	// type_ is the type of the value.
	type_ int

	// prev_data is the previous data of the value.
	prev_data string
	// time_updated_prev is the time the previous data was updated in milliseconds.
	time_updated_prev int64
	// curr_data is the current data of the value.
	curr_data string
	// time_updated_curr is the time the data was updated in milliseconds.
	time_updated_curr int64
}

var registry_GL map[string]*Value = make(map[string]*Value)

/*
AddValue adds a value to the registry.

-----------------------------------------------------------

– Params:
  - key – the key of the value
  - pretty_name – the pretty name of the value
  - description – the description of the value
  - value_type – the type of the value

– Returns:
  - the created value or nil if the value already exists
*/
func AddValue(key string, pretty_name string, description string, value_type int) *Value {
	if _, ok := registry_GL[key]; ok {
		return nil
	}

	registry_GL[key] = &Value{
		key:          key,
		pretty_name:  pretty_name,
		description:  description,
		type_:        value_type,
	}

	return registry_GL[key]
}

/*
RemoveValue removes a value from the registry based on its key.

-----------------------------------------------------------

– Params:
  - key – the key of the value
 */
func RemoveValue(key string) {
	delete(registry_GL, key)
}

type TestFunction func() bool
type TestInterface interface {
	Test() bool
	TestBool (bool) bool
}

func TestFunc(function TestFunction) bool {
	return true
}
func TestFunc2(func() bool) bool {
	return true
}
func TestInterfaceFunc(test_interface TestInterface) bool {
	return test_interface.Test()
}

type SomeCallback interface {
	DoSomething()
}

func TestCallback(callback SomeCallback) {
	callback.DoSomething()
}
