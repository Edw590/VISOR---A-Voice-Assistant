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

package GPTComm

// _Entry is a struct containing information of a generated text.
type _Entry struct {
	// device_id is the device ID of the entry
	device_id string
	// text is the text generated
	text string
	// time_ms is the Unix time_ms in milliseconds
	time_ms int64
}

/*
getDeviceID gets the device ID of the entry.

-----------------------------------------------------------

– Returns:
  - the device ID
*/
func (entry _Entry) getDeviceID() string {
	return entry.device_id
}

/*
getText gets the text of the entry.

-----------------------------------------------------------

– Returns:
  - the text, ending in END_ENTRY
*/
func (entry _Entry) getText() string {
	return entry.text
}

/*
getTime gets the time_ms of the entry.

-----------------------------------------------------------

– Returns:
  - the time_ms in milliseconds
*/
func (entry _Entry) getTime() int64 {
	return entry.time_ms
}
