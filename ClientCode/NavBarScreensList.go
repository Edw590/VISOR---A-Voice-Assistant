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

package main

// screens_GL defines the metadata for each screen
var screens_GL = map[string]string{
	"home": "Home",
	"mod_mod_manager": "Modules Manager",
	"mod_speech": "Speech",
	"mod_gpt_comm": "GPT Communicator",
	"tasks_executor": "Tasks Executor",
		"add_task": "Add task",
	"sys_checker": "System Checker",
	"registry": "Registry",
	"settings": "Settings",
}

// tree_index defines how the screens should be laid out in the index tree
var tree_index = map[string][]string{
	"": {
		"home",
		"mod_mod_manager",
		"mod_speech",
		"mod_gpt_comm",
		"tasks_executor",
		"sys_checker",
		"registry",
		"settings",
	},
	"tasks": {
		"add_task",
	},
}
