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

import "VISOR_Client/Screens"

// screens_GL defines the metadata for each screen
var screens_GL = map[string]string{
	Screens.ID_HOME:                  "Home",
	Screens.ID_MOD_MOD_MANAGER:       "Modules Manager",
	Screens.ID_MOD_SPEECH:            "Speech",
	Screens.ID_MOD_RSS_FEED_NOTIFIER: "RSS Feed Notifier",
	Screens.ID_MOD_GPT_COMM:          "GPT Communicator",
	Screens.ID_MOD_TASKS_EXECUTOR:    "Tasks Executor",
	Screens.ID_MOD_USER_LOCATOR:      "User Locator",
	Screens.ID_ONLINE_INFO_CHK:       "Online Info Checker",
	Screens.ID_MOD_SYS_CHECKER:       "System Checker",
	Screens.ID_SMART_CHECKER:         "S.M.A.R.T. Checker",
	Screens.ID_REGISTRY:              "Registry",
}

// tree_index defines how the screens should be laid out in the index tree
var tree_index = map[string][]string{
	"": {
		Screens.ID_HOME,
		Screens.ID_MOD_MOD_MANAGER,
		Screens.ID_MOD_SPEECH,
		Screens.ID_MOD_GPT_COMM,
		Screens.ID_MOD_RSS_FEED_NOTIFIER,
		Screens.ID_MOD_TASKS_EXECUTOR,
		Screens.ID_MOD_USER_LOCATOR,
		Screens.ID_ONLINE_INFO_CHK,
		Screens.ID_MOD_SYS_CHECKER,
		Screens.ID_SMART_CHECKER,
		Screens.ID_REGISTRY,
	},
}
