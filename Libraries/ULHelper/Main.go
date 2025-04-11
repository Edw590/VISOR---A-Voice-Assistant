/*******************************************************************************
 * Copyright 2023-2025 The V.I.S.O.R. authors
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

/*
UpdateUserLocation updates the internal user location based on the internal device information.
*/
func UpdateUserLocation() {
	stop_GL = false

	for {
		var curr_location = UNKNOWN_LOCATION

		for _, location_info := range getModUserInfo().Locs_info {
			if !location_info.Enabled {
				continue
			}

			var beacons_list []ModsFileInfo.ExtBeacon
			if location_info.Type == "wifi" {
				beacons_list = getDeviceInfo().System_state.Connectivity_info.Wifi_networks
			} else if location_info.Type == "bluetooth" {
				beacons_list = getDeviceInfo().System_state.Connectivity_info.Bluetooth_devices
			} else {
				continue
			}

			var beacon_found *ModsFileInfo.ExtBeacon = nil
			for _, beacon := range beacons_list {
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

			if beacon_found != nil {
				var distance int32 = int32(UtilsSWA.GetRealDistanceRssiLOCRELATIVE(beacon_found.RSSI, UtilsSWA.DEFAULT_TX_POWER))

				if distance <= location_info.Max_distance_m {
					// If the device is near the beacon, then the user may be near the location.
					curr_location = location_info.Location
					if checkUserLocation(curr_location) {
						getUserLocation().Last_detection_when_s = time.Now().Unix()
					}

					break
				}
			}
		}

		// If no beacon was found, check if the user may still be in the location based on the time the location was
		// last detected.
		if curr_location == UNKNOWN_LOCATION {
			for _, location_info := range getModUserInfo().Locs_info {
				if getUserLocation().Curr_location == location_info.Location &&
					getUserLocation().Last_detection_when_s + location_info.Last_detection_s >= time.Now().Unix() {
					curr_location = location_info.Location

					break
				}
			}
		}

		if checkUserLocation(curr_location) {
			updateUserLocation(curr_location)
		}

		if Utils.WaitWithStopDATETIME(&stop_GL, 1) {
			return
		}
	}
}

func checkUserLocation(location string) bool {
	if location == UNKNOWN_LOCATION {
		return true
	}

	if getModUserInfo().AlwaysWith_device == Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id {
		return true
	}

	var approved bool = false
	if getDeviceInfo().Last_time_used_s + LAST_UNUSED_MAX_S >= time.Now().Unix() {
		approved = true
	}

	return approved
}

func updateUserLocation(new_location string) {
	if new_location != UNKNOWN_LOCATION {
		getUserLocation().Last_known_location = getUserLocation().Curr_location
	}

	if new_location != getUserLocation().Curr_location {
		getUserLocation().Prev_location = getUserLocation().Curr_location
		getUserLocation().Prev_last_time_checked_s = getUserLocation().Last_time_checked_s

		getUserLocation().Curr_location = new_location
	}
	getUserLocation().Last_time_checked_s = time.Now().Unix()
}

/*
StopChecker stops the user location checker.
 */
func StopChecker() {
	stop_GL = true
}

func getUserLocation() *ModsFileInfo.UserLocation {
	return &getModGenSettings().User_location
}

func getDeviceInfo() *ModsFileInfo.DeviceInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_10.Device_info
}

func getModGenSettings() *ModsFileInfo.Mod12GenInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_12
}

func getModUserInfo() *ModsFileInfo.Mod12UserInfo {
	return &Utils.GetUserSettings(Utils.LOCK_UNLOCK).UserLocator
}
