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

package ModsFileInfo

const MOD_7_STATE_STOPPED int32 = 0
const MOD_7_STATE_STARTING int32 = 1
const MOD_7_STATE_READY int32 = 2
const MOD_7_STATE_BUSY int32 = 3

// Mod12GenInfo is the format of the custom generated information about this specific module.
type Mod7GenInfo struct {
	// State is the state of the module
	State int32
	// N_mems_when_last_memorized is the number of memories when the last session was memorized
	N_mems_when_last_memorized int
	// Memories is the list of memories the GPT has
	Memories []string
	// Sessions is the list of sessions of the user with the GPT indexed by their ID
	Sessions []Session
}

// Session is the format of a chat session with the GPT.
type Session struct {
	// Id is the ID of the session
	Id string
	// Name is the name of the session
	Name string
	// Created_time_s is the timestamp of the creation of the session
	Created_time_s int64
	// History is the chat history of the session
	History []OllamaMessage
	// Last_interaction_s is the timestamp of the last interaction with the session
	Last_interaction_s int64
	// Memorized is whether the session has been memorized since the last interaction
	Memorized bool
}

///////////////////////////////////////////////////////////////////////////////

// Mod7UserInfo is the format of the custom information file about this specific module.
type Mod7UserInfo struct {
	// Model_name is the name of the LLM model to use
	Model_name string
	// Model_has_tool_role indicates if the tool role of the LLM model exists or not
	Model_has_tool_role bool
	// Context_size is the context size to use
	Context_size int32
	// Temperature is the temperature to use
	Temperature float32
	// System_info is the LLM's system information, like the cutting knowledge date and today's date
	System_info string
	// User_nickname is the user nickname to be used by the LLM
	User_nickname string
}
