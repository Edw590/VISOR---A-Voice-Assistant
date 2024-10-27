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

const _SETTINGS_PREFIX string = "MANUAL_"

// Type: int64
const K_MODULES_ACTIVE string = "MODULES_ACTIVE"

// Speech

// Type: string
const K_LAST_SPEECH string = "LAST_SPEECH"
// Type: int
const K_SPEECH_NORMAL_VOL string = _SETTINGS_PREFIX + "SPEECH_NORMAL_VOL"
// Type: int
const K_SPEECH_CRITICAL_VOL string = _SETTINGS_PREFIX + "SPEECH_CRITICAL_VOL"
// Type: bool
const K_SPEECH_ALWAYS_NOTIFY string = _SETTINGS_PREFIX + "SPEECH_ALWAYS_NOTIFY"

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
	/////////////////////////////////////////////
	// Automatic values

	// Modules Manager
	UtilsSWA.RegisterValueREGISTRY(K_MODULES_ACTIVE, "General - Modules active",
		"The modules that are active (in binary)", UtilsSWA.TYPE_LONG, "", true)

	// Speech
	UtilsSWA.RegisterValueREGISTRY(K_LAST_SPEECH, "Speech - Last speech",
		"The last speech that was spoken", UtilsSWA.TYPE_STRING, "", true)

	// System Checker
	UtilsSWA.RegisterValueREGISTRY(K_BATTERY_LEVEL, "System Checker - Battery level",
		"The battery level", UtilsSWA.TYPE_INT, "", true)
	UtilsSWA.RegisterValueREGISTRY(K_POWER_CONNECTED, "System Checker - Power connected",
		"Whether the power is connected", UtilsSWA.TYPE_BOOL, "", true)
	UtilsSWA.RegisterValueREGISTRY(K_SCREEN_BRIGHTNESS, "System Checker - Screen brightness",
		"The screen brightness", UtilsSWA.TYPE_INT, "", true)
	UtilsSWA.RegisterValueREGISTRY(K_SOUND_VOLUME, "System Checker - Sound volume",
		"The sound volume", UtilsSWA.TYPE_INT, "", true)
	UtilsSWA.RegisterValueREGISTRY(K_SOUND_MUTED, "System Checker - Sound muted",
		"Whether the sound is muted", UtilsSWA.TYPE_BOOL, "", true)

	/////////////////////////////////////////////
	// Manual values

	// Speech
	UtilsSWA.RegisterValueREGISTRY(K_SPEECH_NORMAL_VOL, "Speech - Normal speech volume",
		"The volume at which to speak non-critical speeches", UtilsSWA.TYPE_INT, "25", false)
	UtilsSWA.RegisterValueREGISTRY(K_SPEECH_CRITICAL_VOL, "Speech - Critical speech volume",
		"The volume at which to speak critical speeches", UtilsSWA.TYPE_INT, "100", false)
	UtilsSWA.RegisterValueREGISTRY(K_SPEECH_ALWAYS_NOTIFY, "Speech - Always notify",
		"Always notify speeches (useful if media is playing)", UtilsSWA.TYPE_BOOL, "true", false)
}
