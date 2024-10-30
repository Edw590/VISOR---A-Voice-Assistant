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
var Current_screen_GL string = ""

var screens_size_GL fyne.Size = fyne.NewSize(550, 480)

const ID_HOME string = "home"
const ID_MOD_MOD_MANAGER string = "mod_mod_manager"
const ID_MOD_SPEECH string = "mod_speech"
const ID_MOD_RSS_FEED_NOTIFIER string = "mod_rss_feed_notifier"
const ID_MOD_GPT_COMM string = "mod_gpt_comm"
const ID_MOD_TASKS_EXECUTOR string = "tasks_executor"
const ID_MOD_SYS_CHECKER string = "sys_checker"
const ID_REGISTRY string = "registry"
