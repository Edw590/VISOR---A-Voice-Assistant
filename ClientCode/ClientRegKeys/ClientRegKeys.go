/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

import "Registry/Registry"

// Type: int64
const K_MODULES_ACTIVE string = "MODULES_ACTIVE"

// Type: string
const K_LAST_SPEECH string = "LAST_SPEECH"

// Type: bool
const K_SHOW_APP_SIG string = "SHOW_APP_SIG"

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
	Registry.RegisterValue(K_MODULES_ACTIVE, "Modules active", "The modules that are active (in binary)", Registry.TYPE_LONG)

	Registry.RegisterValue(K_LAST_SPEECH, "Last speech", "The last speech that was spoken", Registry.TYPE_STRING)

	Registry.RegisterValue(K_SHOW_APP_SIG, "Show-app signal", "Signal to show the app", Registry.TYPE_BOOL)

	Registry.RegisterValue(K_BATTERY_LEVEL, "Battery level", "The battery level", Registry.TYPE_INT)
	Registry.RegisterValue(K_POWER_CONNECTED, "Power connected", "Whether the power is connected", Registry.TYPE_BOOL)
	Registry.RegisterValue(K_SCREEN_BRIGHTNESS, "Screen brightness", "The screen brightness", Registry.TYPE_INT)
	Registry.RegisterValue(K_SOUND_VOLUME, "Sound volume", "The sound volume", Registry.TYPE_INT)
	Registry.RegisterValue(K_SOUND_MUTED, "Sound muted", "Whether the sound is muted", Registry.TYPE_BOOL)
}
