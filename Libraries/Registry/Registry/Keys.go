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

package Registry

// Type: int64
const K_MODULES_ACTIVE string = "MODULES_ACTIVE"

// Type: string
const K_LAST_SPEECH string = "LAST_SPEECH"

// Type: bool
const K_SHOW_APP_SIG string = "SHOW_APP_SIG"

// Type: bool
const K_BATTERY_PRESENT string = "BATTERY_PRESENT"
// Type: int32
const K_BATTERY_PERCENT string = "BATTERY_PERCENT"
// Type: bool
const K_POWER_CONNECTED string = "POWER_CONNECTED"

func init() {
	// Modules Manager
	RegisterValue(K_MODULES_ACTIVE, "Modules active", "The modules that are active (in binary)", TYPE_LONG)
	// Speech
	RegisterValue(K_LAST_SPEECH, "Last speech", "The last speech that was said", TYPE_STRING)
	// Speech Recognition
	RegisterValue(K_SHOW_APP_SIG, "Show app signal", "Signal to show the app", TYPE_BOOL)
	// Power Processor
	//RegisterValue(K_BATTERY_PRESENT, "Battery present", "Whether the battery is present", TYPE_BOOL)
	//RegisterValue(K_BATTERY_PERCENT, "Battery percentage", "The battery percentage", TYPE_INT)
	//RegisterValue(K_POWER_CONNECTED, "Power connected", "Whether the power is connected", TYPE_BOOL)
}
