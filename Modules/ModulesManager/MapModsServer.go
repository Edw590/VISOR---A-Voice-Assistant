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

//go:build server

package ModulesManager

import (
	"EmailSender"
	"GPTCommunicator"
	"GoogleManager"
	"OnlineInfoChk"
	"RSSFeedNotifier"
	"SMARTChecker"
	"Utils"
	"WebsiteBackend"
)

// Make sure to add the modules support check for each new module too...
var _MAP_MOD_NUM_START = map[int]func(modules *Utils.Module){
	Utils.NUM_MOD_SMARTChecker:    SMARTChecker.Start,
	Utils.NUM_MOD_RssFeedNotifier: RSSFeedNotifier.Start,
	Utils.NUM_MOD_EmailSender:     EmailSender.Start,
	Utils.NUM_MOD_OnlineInfoChk:   OnlineInfoChk.Start,
	Utils.NUM_MOD_GPTCommunicator: GPTCommunicator.Start,
	Utils.NUM_MOD_WebsiteBackend:  WebsiteBackend.Start,
	Utils.NUM_MOD_GoogleManager:   GoogleManager.Start,
}
