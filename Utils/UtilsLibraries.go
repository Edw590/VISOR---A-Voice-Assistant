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

package Utils

const (
	NUM_LIB_ACD               int = iota
	NUM_LIB_OICComm
	NUM_LIB_GPTComm
	NUM_LIB_SpeechQueue
	NUM_LIB_TEHelper
	NUM_LIB_SettingsSync
	NUM_LIB_ULHelper
	NUM_LIB_SCLink
	NUM_LIB_GMan
	NUM_LIB_DialogMan

	LIBS_ARRAY_SIZE
)

// LIB_NUMS_INFO is a map of the numbers of the libraries and their respective ModuleInfo.
var LIB_NUMS_INFO map[int]ModuleInfo = map[int]ModuleInfo{
	NUM_LIB_ACD: {
		Name: "Advanced Commands Detection",
	},
	NUM_LIB_OICComm: {
		Name: "Online Information Checker Communicator",
	},
	NUM_LIB_GPTComm: {
		Name: "GPT Communicator",
	},
	NUM_LIB_SpeechQueue: {
		Name: "Speech Queue",
	},
	NUM_LIB_TEHelper: {
		Name: "Tasks Executor Helper",
	},
	NUM_LIB_SettingsSync: {
		Name: "Settings Synchronizer",
	},
	NUM_LIB_ULHelper: {
		Name: "User Locator Helper",
	},
	NUM_LIB_SCLink: {
		Name: "System Checker Link",
	},
	NUM_LIB_GMan: {
		Name: "Google Manager",
	},
	NUM_LIB_DialogMan: {
		Name: "Dialogue Manager",
	},
}

/*
GetLibNameMODULES gets the name of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the name of the module or an empty string if the module number is invalid
*/
func GetLibNameMODULES(lib_num int) string {
	if library_info, ok := LIB_NUMS_INFO[lib_num]; ok {
		return library_info.Name
	}

	return ""
}
