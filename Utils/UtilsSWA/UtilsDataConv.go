/*******************************************************************************
 * Copyright 2023-2023 Edw590
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
	"github.com/cention-sany/utf7"
)

/*
BytesToPrintableDATACONV converts a byte array to printable characters in a string.

Note: all bytes will be attempted to be printed, all based on the platform's default charset (on Android is always the
UTF-8 charset).

-----------------------------------------------------------

– Params:
  - byte_array – the byte array
  - utf7_flag – true if the bytes were encoded using UTF-7, false if they were encoded using UTF-8

– Returns:
  - a string containing printable characters representative of the provided bytes
 */
func BytesToPrintableDATACONV(byte_array []byte, utf7_flag bool) string {
	if utf7_flag {
		utf7_bytes, _ := utf7.UTF7DecodeBytes(byte_array)

		return string(utf7_bytes)
	} else {
		return string(byte_array)
	}
}

const _HEX_ARRAY string = "0123456789ABCDEF";
/*
BytesToHexDATACONV converts the given byte array into a string of the type "XX XX XX", in which the Xs are hexadecimal digits.

-----------------------------------------------------------

– Params:
  - bytes – the byte array

– Returns:
  - the equivalent string
*/
func BytesToHexDATACONV(bytes []byte) string {
	// DO NOT CHANGE THE OUTPUT FORMAT WITHOUT CHECKING ALL USAGES!!!!!! "00 00 00" was chosen because it's easy to
	// replace the spaces by "\x" to write files, for example.

	var bytes_length int = len(bytes)
	var hex_chars string = ""
	for i := 0; i < bytes_length; i++ {
		var positive_byte uint = uint(bytes[i]) & 0xFF // See why it works on byteToIntUnsigned()
		hex_chars += string(_HEX_ARRAY[positive_byte >> 4])
		hex_chars += string(_HEX_ARRAY[positive_byte & 0x0F])
		hex_chars += " "
	}

	return hex_chars[:len(hex_chars) - 1]
}

const _OCT_ARRAY string = "01234567";
/*
BytesToOctalDATACONV converts the given byte array into a string of the type "XXX XXX XXX", in which the Xs are octal digits.

-----------------------------------------------------------

– Params:
  - bytes – the byte array

– Returns:
  - the equivalent string
*/
func BytesToOctalDATACONV(bytes []byte) string {
	// DO NOT CHANGE THE OUTPUT FORMAT WITHOUT CHECKING ALL USAGES!!!!!! "000 000 000" was chosen because it's easy to
	// replace the spaces by "\0" to write files, for example.

	var bytes_length int = len(bytes)
	var hex_chars string = ""
	for i := 0; i < bytes_length; i++ {
		var positive_byte uint = uint(bytes[i]) & 0xFF // See why it works on byteToIntUnsigned()
		hex_chars += string(_OCT_ARRAY[(positive_byte & 0b11000000) >> 6])
		hex_chars += string(_OCT_ARRAY[(positive_byte & 0b00111000 >> 3)])
		hex_chars += string(_OCT_ARRAY[positive_byte & 0b00000111])
		hex_chars += " "
	}

	return hex_chars[:len(hex_chars) - 1]
}
