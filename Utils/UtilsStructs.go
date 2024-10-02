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
	"reflect"
)

/*
CompareSTRUCTS compares two structs and returns true if they are equal in all their fields.

-----------------------------------------------------------

– Params:
  - a – the first struct to compare
  - b – the second struct to compare

– Returns:
  - true if the structs are equal in all their fields, false otherwise
 */
func CompareSTRUCTS[T any](a T, b T) bool {
	valA := reflect.ValueOf(a)
	valB := reflect.ValueOf(b)
	typA := reflect.TypeOf(a)
	typB := reflect.TypeOf(b)

	if typA.Kind() != reflect.Struct || typB.Kind() != reflect.Struct {
		return false
	}

	if typA != typB {
		return false
	}

	result := true

	for i := 0; i < valA.NumField(); i++ {
		fieldA := valA.Field(i)
		fieldB := valB.Field(i)

		if fieldA.Kind() == reflect.Struct {
			if !CompareSTRUCTS(fieldA.Interface(), fieldB.Interface()) {
				result = false

				break
			}
		} else if fieldA.Kind() == reflect.Slice || fieldA.Kind() == reflect.Array {
			if fieldA.Len() != fieldB.Len() {
				result = false

				break
			} else {
				for j := 0; j < fieldA.Len(); j++ {
					if fieldA.Index(j).Interface() != fieldB.Index(j).Interface() {
						result = false

						break
					}
				}
			}
		} else {
			if fieldA.Interface() != fieldB.Interface() {
				result = false

				break
			}
		}
	}

	return result
}
