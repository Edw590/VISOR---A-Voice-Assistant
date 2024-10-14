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

package ClientRegKeys

import (
	"Utils/UtilsSWA"
)

// Modules Manager

// Type: int64
const K_MODULES_ACTIVE string = "MODULES_ACTIVE"

// Speech

// Type: string
const K_LAST_SPEECH string = "LAST_SPEECH"

// Main app

// Type: bool
const K_SHOW_APP_SIG string = "SHOW_APP_SIG"

// System Checker

// Type: int
const K_BATTERY_LEVEL string = "BATTERY_LEVEL"
// Type: bool
const K_POWER_CONNECTED string = "POWER_CONNECTED"
// Type: int
const K_SCREEN_BRIGHTNESS string = "SCREEN_BRIGHTNESS"
// Type: int
const K_SOUND_VOLUME string = "SOUND_VOLUME"
// Type: bool
const K_SOUND_MUTED string = "SOUND_MUTED"

/*
RegisterValues registers the client values in the registry.
 */
func RegisterValues() {
	// Modules Manager
	UtilsSWA.RegisterValueREGISTRY(K_MODULES_ACTIVE, "Modules active", "The modules that are active (in binary)", UtilsSWA.TYPE_LONG)

	// Speech
	UtilsSWA.RegisterValueREGISTRY(K_LAST_SPEECH, "Last speech", "The last speech that was spoken", UtilsSWA.TYPE_STRING)

	// Main app
	UtilsSWA.RegisterValueREGISTRY(K_SHOW_APP_SIG, "Show-app signal", "Signal to show the app", UtilsSWA.TYPE_BOOL)

	// System Checker
	UtilsSWA.RegisterValueREGISTRY(K_BATTERY_LEVEL, "Power - Battery level", "The battery level", UtilsSWA.TYPE_INT)
	UtilsSWA.RegisterValueREGISTRY(K_POWER_CONNECTED, "Power - Power connected", "Whether the power is connected", UtilsSWA.TYPE_BOOL)
	UtilsSWA.RegisterValueREGISTRY(K_SCREEN_BRIGHTNESS, "Screen brightness", "The screen brightness", UtilsSWA.TYPE_INT)
	UtilsSWA.RegisterValueREGISTRY(K_SOUND_VOLUME, "Sound volume", "The sound volume", UtilsSWA.TYPE_INT)
	UtilsSWA.RegisterValueREGISTRY(K_SOUND_MUTED, "Sound muted", "Whether the sound is muted", UtilsSWA.TYPE_BOOL)
}
