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

package SettingsSync

import "Utils"

/*
GetJsonDeviceSettings returns the device settings in JSON format.

-----------------------------------------------------------

– Returns:
  - the device settings in JSON format
 */
func GetJsonDeviceSettings() string {
	return *Utils.ToJsonGENERAL(Utils.Device_settings_GL)
}

/*
LoadDeviceSettings loads the device settings from the given JSON string.

-----------------------------------------------------------

– Params:
  - json – the JSON string to load the device settings from

– Returns:
  - true if the device settings were successfully loaded, false otherwise
 */
func LoadDeviceSettings(json string) bool {
	if json == "" {
		return false
	}

	if err := Utils.FromJsonGENERAL([]byte(json), &Utils.Device_settings_GL); err != nil {
		return false
	}

	return true
}
