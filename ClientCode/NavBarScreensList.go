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

import (
	"VISOR_Client/Screens"
	"fyne.io/fyne/v2"
)

// _Screen defines the data structure for a tutorial
type _Screen struct {
	Title string
	View func(param any) fyne.CanvasObject
}

// screens_GL defines the metadata for each screen
var screens_GL = map[string]_Screen{
	"home": {"Home", Screens.Home},
	"dev_mode": {"Dev Mode", Screens.DevMode},
	"communicator": {"Communicator", Screens.Communicator},
	"mod_status": {"Modules Status", Screens.ModulesStatus},
	"calendar": {"Calendar", Screens.Calendar},
	"registry": {"Registry", Screens.Registry},
	"tasks": {"Tasks (NOT READY)", Screens.Tasks},
	"add_task": {"Add task", Screens.Tasks},
	"sys_state": {"System State", Screens.SystemState},
	"settings": {"Settings", Screens.Settings},
}

// tree_index defines how the screens should be laid out in the index tree
var tree_index = map[string][]string{
	"": {"home", "dev_mode", "communicator", "mod_status", "calendar", "registry", "tasks", "sys_state", "settings"},
	"tasks": {"add_task"},
}
