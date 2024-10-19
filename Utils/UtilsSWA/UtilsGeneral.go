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

package UtilsSWA

import (
	"Utils"
	"strings"
	"unicode"
)

/*
RandStringGENERAL generates a random string with the given length, containing only ASCII letters (upper case
or lower case) and numbers.
RemoveNonGraphicCharsGENERAL removes all the non-graphic characters from a string.

-----------------------------------------------------------

– Params:
  - length – length of the generating string
  - str – the string to remove the non-graphic characters from

– Returns:
  - the generated string
  - the string without the non-graphic characters
*/
func RandStringGENERAL(length int) string {
	return Utils.RandStringGENERAL(length)
}

/*
RemoveNonGraphicCharsGENERAL removes all the non-graphic characters from a string.

-----------------------------------------------------------

– Params:
  - str – the string to remove the non-graphic characters from

– Returns:
  - the string without the non-graphic characters
*/
func RemoveNonGraphicCharsGENERAL(str string) string {
	str = strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}

		return -1
	}, str)

	return str
}

/*
StringHasLetters checks if a string has any letters.

-----------------------------------------------------------

– Params:
  - s – the string to check

– Returns:
  - true if the string has any letters, false otherwise
 */
func StringHasLettersGENERAL(str string) bool {
	for _, r := range str {
		if unicode.IsLetter(r) {
			return true
		}
	}

	return false
}

/*
StringHasNumbers checks if a string has any numbers.

-----------------------------------------------------------

– Params:
  - s – the string to check

– Returns:
  - true if the string has any numbers, false otherwise
 */
func StringHasNumbersGENERAL(str string) bool {
	for _, r := range str {
		if unicode.IsNumber(r) {
			return true
		}
	}

	return false
}
