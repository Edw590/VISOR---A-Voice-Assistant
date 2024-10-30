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

import Tcef "github.com/Edw590/TryCatch-go"

var ModsCommsChannels_GL [MODS_ARRAY_SIZE]chan map[string]any = [MODS_ARRAY_SIZE]chan map[string]any{}
var LibsCommsChannels_GL [LIBS_ARRAY_SIZE]chan map[string]any = [LIBS_ARRAY_SIZE]chan map[string]any{}

/*
InitializeCommsChannels initializes the modules and libraries communication channels.
*/
func InitializeCommsChannels() {
	for i := 0; i < MODS_ARRAY_SIZE; i++ {
		ModsCommsChannels_GL[i] = make(chan map[string]any)
	}
	for i := 0; i < LIBS_ARRAY_SIZE; i++ {
		LibsCommsChannels_GL[i] = make(chan map[string]any)
	}
}

/*
CloseCommsChannels closes the modules and libraries communication channels.
 */
func CloseCommsChannels() {
	for i := 0; i < MODS_ARRAY_SIZE; i++ {
		close(ModsCommsChannels_GL[i])
	}
	for i := 0; i < LIBS_ARRAY_SIZE; i++ {
		close(LibsCommsChannels_GL[i])
	}
}

/*
SendToModChannel sends data to a module channel.

In case the module is not supported this function does nothing (to prevent deadlock sending to unused channels).

-----------------------------------------------------------

– Params:
  - mod_num – the module number
  - key – the key of the data
  - data – the data to send
 */
func SendToModChannel(mod_num int, key string, data any) {
	if !IsModSupportedMODULES(mod_num) {
		return
	}

	// Ignore the panic of writing to closed channels (sometimes happens when the app is shutting down).
	Tcef.Tcef{
		Try: func() {
			ModsCommsChannels_GL[mod_num] <- map[string]any{
				key: data,
			}
		},
	}.Do()
}

/*
SendToLibChannel sends data to a library channel.

-----------------------------------------------------------

– Params:
  - lib_num – the library number
  - key – the key of the data
  - data – the data to send
 */
func SendToLibChannel(lib_num int, key string, data any) {
	// Ignore the panic of writing to closed channels (sometimes happens when the app is shutting down).
	Tcef.Tcef{
		Try: func() {
			LibsCommsChannels_GL[lib_num] <- map[string]any{
				key: data,
			}
		},
	}.Do()
}
