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

package GPTComm

// Entry is a struct containing information of a generated text.
type Entry struct {
	// device_id is the device ID of the entry
	device_id string
	// text is the text generated
	text string
	// time is the Unix time in milliseconds
	time int64
}

/*
GetDeviceID gets the device ID of the entry.

-----------------------------------------------------------

– Returns:
  - the device ID
 */
func (entry Entry) GetDeviceID() string {
	return entry.device_id
}

/*
GetText gets the text of the entry.

-----------------------------------------------------------

– Returns:
  - the text, ending in END_ENTRY
*/
func (entry Entry) GetText() string {
	return entry.text
}

/*
GetTime gets the time of the entry.

-----------------------------------------------------------

– Returns:
  - the time in milliseconds
*/
func (entry Entry) GetTime() int64 {
	return entry.time
}
