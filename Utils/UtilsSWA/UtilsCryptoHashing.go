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
	"crypto/sha512"
	"encoding/hex"
)

/*
GetHashBytesOfBytesCRYPTOHASHING gets the byte array of the hash of the given byte array, calculated by the SHA-512
algorithm.

-----------------------------------------------------------

– Params:
  - data – the bytes to calculate the hash from

– Returns:
  - the hash bytes

*/
func GetHashBytesOfBytesCRYPTOHASHING(data []byte) []byte {
	hasher := sha512.New()
	hasher.Write(data)

	return hasher.Sum(nil)
}

/*
GetHashStringOfBytesCRYPTOHASHING gets the string of the hash of the given byte array, calculated by the SHA-512
algorithm (eg. "32B4A667AA8F").

-----------------------------------------------------------

– Params:
  - byteArray – the bytes to calculate the hash from

– Returns:
  - the hash string
 */
func GetHashStringOfBytesCRYPTOHASHING(byteArray []byte) string {
	hash := sha512.Sum512(byteArray)

	return hex.EncodeToString(hash[:])
}
