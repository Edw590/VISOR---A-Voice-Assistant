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

package ULHelper

import (
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"time"
)

const UNKNOWN_LOCATION string = "3234_UNKNOWN"

// LAST_UNUSED_MAX_S is the maximum time since the last time active to consider the device in a location
const LAST_UNUSED_MAX_S int64 = 5 * 60

var stop_GL bool = false

var device_info_GL *ModsFileInfo.DeviceInfo = &Utils.Gen_settings_GL.MOD_10.Device_info

var (
	modGenInfo_GL  *ModsFileInfo.Mod12GenInfo = &Utils.Gen_settings_GL.MOD_12
	modUserInfo_GL *ModsFileInfo.Mod12UserInfo = &Utils.User_settings_GL.UserLocator
)

/*
UpdateUserLocation updates the internal user location based on the internal device information.
*/
func UpdateUserLocation() {
	stop_GL = false

	for {
		//log.Println("--------------------------------")

		// Update some information
		var curr_location = UNKNOWN_LOCATION

		for _, location_info := range modUserInfo_GL.Locs_info {
			var beacon_list []ModsFileInfo.ExtBeacon
			if location_info.Type == "wifi" {
				beacon_list = device_info_GL.System_state.Connectivity_info.Wifi_networks
			} else if location_info.Type == "bluetooth" {
				beacon_list = device_info_GL.System_state.Connectivity_info.Bluetooth_devices
			} else {
				continue
			}

			var beacon_found *ModsFileInfo.ExtBeacon = nil
			for _, beacon := range beacon_list {
				var address_match bool = true
				var name_match bool = true
				if location_info.Address != "" {
					address_match = beacon.Address == location_info.Address
				}
				if location_info.Name != "" {
					name_match = beacon.Name == location_info.Name
				}

				if address_match && name_match {
					beacon_found = &beacon

					break
				}
			}

			var distance_match bool = false
			var may_still_be_in_location bool = false
			if beacon_found != nil {
				var distance int = UtilsSWA.GetRealDistanceRssiLOCRELATIVE(beacon_found.RSSI, UtilsSWA.DEFAULT_TX_POWER)

				if distance <= location_info.Max_distance {
					// If the device is near the location, then the user is near the location.
					distance_match = true
				}
			} else {
				if curr_location != UNKNOWN_LOCATION {
					if modGenInfo_GL.User_location.Last_detection_when_s+ location_info.Last_detection_s >= time.Now().Unix() {
						may_still_be_in_location = true
					}
				}
			}

			if distance_match || may_still_be_in_location {
				curr_location = location_info.Location
			}
		}

		if curr_location != UNKNOWN_LOCATION {
			modGenInfo_GL.User_location.Last_known_location = curr_location
			modGenInfo_GL.User_location.Last_detection_when_s = time.Now().Unix()
		}

		updateUserLocation(computeUserLocation(curr_location))

		if Utils.WaitWithStopTIMEDATE(&stop_GL, 1) {
			return
		}
	}
}

func computeUserLocation(location string) string {
	if modUserInfo_GL.AlwaysWith_device == Utils.Device_settings_GL.Device_ID && location != UNKNOWN_LOCATION {
		return location
	}

	var curr_location string = UNKNOWN_LOCATION
	if location != UNKNOWN_LOCATION && device_info_GL.Last_time_used_s + LAST_UNUSED_MAX_S >= time.Now().Unix() {
		curr_location = location
	}

	return curr_location
}

func updateUserLocation(new_location string) {
	var user_location *ModsFileInfo.UserLocation = &modGenInfo_GL.User_location

	if new_location != UNKNOWN_LOCATION {
		user_location.Last_detection_when_s = time.Now().Unix()
		user_location.Last_known_location = user_location.Curr_location
	}

	if new_location != user_location.Curr_location {
		user_location.Prev_location = user_location.Curr_location
		user_location.Prev_last_time_checked_s = user_location.Last_time_checked_s

		user_location.Curr_location = new_location
	}
	user_location.Last_time_checked_s = time.Now().Unix()
}

/*
StopChecker stops the user location checker.
 */
func StopChecker() {
	stop_GL = true
}
