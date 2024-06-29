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

package MOD_12

// _ModUserInfo is the format of the custom information file about this specific module.
type _ModUserInfo struct {
	// Devices_info is the information about the devices
	Devices_info _DevicesInfo
	// Locs_info is the information about the locations
	Locs_info []_LocInfo
}

type _DevicesInfo struct {
	// AlwaysWith_device_id is the device id of the device that is always with the user
	AlwaysWith_device_id string
}

type _LocInfo struct {
	// Type is the type of the location "detector" (e.g. wifi)
	Type string
	// Name is the name of the detection (e.g. the wifi SSID)
	Name string
	// Address is the address of the detection (e.g. the wifi BSSID) in the format XX:XX:XX:XX:XX:XX
	Address string
	// Last_detection is the maximum amount of time in seconds without checking in which the device may still be in the
	// specified location
	Last_detection int64
	// Max_distance is the maximum distance in meters in which the device is in the specified location
	Max_distance int
	// Location is where the device is (e.g. "home")
	Location string
}
