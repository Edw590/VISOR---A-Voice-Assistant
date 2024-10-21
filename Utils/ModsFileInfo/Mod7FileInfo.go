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

package ModsFileInfo

const MOD_7_STATE_READY int = 0
const MOD_7_STATE_STARTING int = 1
const MOD_7_STATE_BUSY int = 2
const MOD_7_STATE_STOPPING int = 3

// Mod12GenInfo is the format of the custom generated information about this specific module.
type Mod7GenInfo struct {
	// State is the state of the module
	State int
	// Memories is the list of memories the GPT has
	Memories []string
}

///////////////////////////////////////////////////////////////////////////////

// Mod7UserInfo is the format of the custom information file about this specific module.
type Mod7UserInfo struct {
	// Model_smart_loc is the location of the model file for the smart LLM
	Model_smart_loc string
	// Model_dumb_loc is the location of the model file for the dumb LLM
	Model_dumb_loc string
	// System_info is the LLM's system information, like the cutting knowledge date and today's date
	System_info string
}
