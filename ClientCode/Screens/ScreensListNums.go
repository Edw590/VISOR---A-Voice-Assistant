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

package Screens

import "fyne.io/fyne/v2"

// Current_screen_GL is the current app screen. It's currently used to let threads specific to each screen know if they
// should continue processing data or not (they don't stop, they just keep waiting for the screen to become active again).
var Current_screen_GL int = -1

var screens_size_GL fyne.Size = fyne.NewSize(550, 480)

const (
	NUM_HOME = iota
	NUM_MODULES_MANAGER
	NUM_SPEECH
	NUM_RSS_FEED_NOTIFIER
	NUM_GPT_COMMUNICATOR
	NUM_TASKS_EXECUTOR
	NUM_SYSTEM_CHECKER
	NUM_REGISTRY
)
