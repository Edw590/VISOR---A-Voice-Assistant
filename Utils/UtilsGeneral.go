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

package Utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unsafe"

	"github.com/dchest/jsmin"
	"github.com/ztrue/tracerr"
)

///////////////////////////
// Took from https://stackoverflow.com/a/31832326/8228163, by https://stackoverflow.com/users/1705598/icza.
// Edw590: added numbers to the letterBytes constant (all still fits in the 6 bits).
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var src = rand.NewSource(time.Now().UnixNano())

/*
RandStringGENERAL generates a random string with the given length, containing only ASCII letters (upper case
or lower case) and numbers.

-----------------------------------------------------------

– Params:
  - letters_num – the length of the string to generate

– Returns:
  - the generated string
*/
func RandStringGENERAL(length int) string {
	// Original function name: RandStringBytesMaskImprSrcUnsafe
	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
///////////////////////////

/*
FindAllIndexesGENERAL finds all the indexes of a substring in a string.

-----------------------------------------------------------

– Params:
  - s – the string to search in
  - substr – the substring to search for

– Returns:
  - the indexes of the substring in the string
*/
func FindAllIndexesGENERAL(s string, substr string) []int {
	var indexes []int = nil
	var chars_processed int = 0
	var s_len int = len(s)
	for i := 0; i < s_len; i++ {
		var idx int = strings.Index(s[chars_processed:], substr)
		if idx == -1 {
			break
		}

		indexes = append(indexes, idx + chars_processed)
		chars_processed += idx + len(substr)
	}

	return indexes
}

/*
GetFullErrorMsgGENERAL gets the full error message from an error, including its stacktrace.

-----------------------------------------------------------

– Params:
  - err – the error to get the full message from

– Returns:
  - the full error message
*/
func GetFullErrorMsgGENERAL(err_param any) string {
	var str_error string = ""
	if err, ok := err_param.(error); ok {
		// tracerr only works with errors
		str_error = tracerr.SprintSource(tracerr.Wrap(err), 0)
	} else {
		// If the exception is not an error, get general information about it
		var err_str string = "Invalid type of error information (not a Go \"error\"). " + getVariableInfoGENERAL(err_param)
		str_error = err_str
	}

	return str_error
}

/*
ToJsonGENERAL converts the given data to a JSON string and indents it.

All the needed fields of the struct must be exported like with json.Marshal().

-----------------------------------------------------------

– Params:
  - v – the data to convert to Json. Check the json.Marshal function for more info (used directly here).

– Returns:
  - the JSON string, or nil if the data could not be converted
*/
func ToJsonGENERAL(v any) *string {
	json_data, err := json.Marshal(v)
	if nil != err {
		return nil
	}

	var dst bytes.Buffer
	if json.Indent(&dst, json_data, "", "\t") == nil {
		json_data = dst.Bytes()
	}

	var json_str string = string(json_data)

	return &json_str
}

/*
FromJsonGENERAL minifies and parses the given JSON data.

All the needed fields of the struct must be exported like with json.Marshal().

This function supports minifying the JSON data and removing the last commas of JSON file (possibly rendering the JSON
invalid).

-----------------------------------------------------------

– Params:
  - json_str – the JSON string to parse
  - parsed_data – a pointer of where to write the parsed data to

– Returns:
  - true if the data was parsed correctly, false otherwise
*/
func FromJsonGENERAL(json_data []byte, parsed_data any) error {
	if json_data == nil {
		return errors.New("the JSON data is nil")
	}

	var json_final []byte = json_data
	for {
		var json_min, err = jsmin.Minify(json_final)
		if err == nil {
			json_final = json_min
		} else {
			// If the minifier fails, try to parse the original JSON (probably won't work, but I'll let Unmarshal()
			// decide).
		}

		if err = json.Unmarshal(json_final, parsed_data); err != nil {
			// Remove the last comma after the last element of the JSON array and try again (in case the problem is
			// a trailing comma somewhere).
			if !DelElemSLICES(&json_final, strings.LastIndex(string(json_final), ",")) {
				// No more commas to remove
				return err
			}
		} else {
			return nil
		}
	}
}

/*
getVariableInfoGENERAL gets general information about a variable in a string in a default format.

-----------------------------------------------------------

– Params:
  - v – the variable to get the information about

– Returns:
  - the information about the variable
*/
func getVariableInfoGENERAL(v any) string {
	return "Information about it:" +
		"\n- Type of information (%T): " + fmt.Sprintf("%T", v) +
		"\n- Value of information (%+v): " + fmt.Sprintf("%+v", v) +
		"\n- Go representation of the value (%#v): " + fmt.Sprintf("%#v", v)
}

/*
WasArgUsedGENERAL checks if a wanted argument was used in the arguments list.

-----------------------------------------------------------

– Params:
  - wanted_arg – the argument to check for
  - args – the arguments list

– Returns:
  - true if the argument was used, false otherwise
 */
func WasArgUsedGENERAL(args []string, wanted_arg string) bool {
	for _, curr_arg := range args {
		if wanted_arg == curr_arg {
			return true
		}
	}

	return false
}
