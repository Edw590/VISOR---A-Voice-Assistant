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

package UserLocator

import (
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"log"
	"time"
)

// User Locator //

const TIME_SLEEP_S int = 5

const UNKNOWN_LOCATION string = "3234_UNKNOWN"

const LAST_COMM_MAX_S int64 = 30 + 5 // MUST BE HIGHER THAN MOD_8.PONG_WAIT!!!
// LAST_UNUSED_MAX_S is the maximum time since the last time active to consider the device in a location
const LAST_UNUSED_MAX_S int64 = 5 * 60

type _IntDeviceInfo struct {
	Last_time_used_s    int64
	Curr_location       string
	Last_known_location string
}

var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
	modGenInfo_GL  *ModsFileInfo.Mod12GenInfo
	modUserInfo_GL *ModsFileInfo.Mod12UserInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_12
		modUserInfo_GL = &Utils.User_settings_GL.MOD_12

		var device_info *ModsFileInfo.DeviceInfo = &Utils.Gen_settings_GL.MOD_10.Device_info

		for {
			//log.Println("--------------------------------")
			var int_device_info _IntDeviceInfo

			// Update some information
			int_device_info.Curr_location = UNKNOWN_LOCATION
			int_device_info.Last_time_used_s = device_info.Last_time_used_s

			for _, location_info := range modUserInfo_GL.Locs_info {
				if location_info.Last_detection_s < int64(TIME_SLEEP_S) * 2 {
					// There must be a minimum. That minimum is the time it takes for the devices to update their
					// status, but double it to be sure they communicated.
					location_info.Last_detection_s = int64(TIME_SLEEP_S) * 2
				}

				var beacon_list []ModsFileInfo.ExtBeacon
				if location_info.Type == "wifi" {
					beacon_list = device_info.System_state.Connectivity_info.Wifi_networks
				} else if location_info.Type == "bluetooth" {
					beacon_list = device_info.System_state.Connectivity_info.Bluetooth_devices
				} else {
					continue
				}

				for _, beacon := range beacon_list {
					var location_matches bool = false
					if location_info.Address != "" {
						if beacon.Address == location_info.Address {
							location_matches = true
						}
					} else {
						if beacon.Name == location_info.Name {
							location_matches = true
						}
					}

					if location_matches {
						var distance int = UtilsSWA.GetRealDistanceRssiLOCRELATIVE(beacon.RSSI, UtilsSWA.DEFAULT_TX_POWER)

						if distance <= location_info.Max_distance {
							// If the device is near the location, then the user is near the location.
							int_device_info.Curr_location = location_info.Location
						}

						break
					}
				}
			}

			if int_device_info.Curr_location != UNKNOWN_LOCATION {
				int_device_info.Last_known_location = int_device_info.Curr_location
			}

			log.Println("int_device_info:", int_device_info)

			var curr_user_location string = computeUserLocation(int_device_info)
			log.Println("Current user location:", curr_user_location)
			updateUserLocation(&modGenInfo_GL.User_location, curr_user_location)

			// TODO: Also check if the location changed on some device. The user must be with it then, even if not using
			//  it.

			// TODO: Give priorities to devices. You're always with the phone even if not using it, but not with the
			//  computer.

			if Utils.WaitWithStopTIMEDATE(module_stop, TIME_SLEEP_S) {
				return
			}
		}
	}
}

func computeUserLocation(int_device_info _IntDeviceInfo) string {
	if modUserInfo_GL.AlwaysWith_device && int_device_info.Curr_location != UNKNOWN_LOCATION {
		return int_device_info.Curr_location
	}

	var curr_location string = UNKNOWN_LOCATION
	if int_device_info.Curr_location != UNKNOWN_LOCATION && int_device_info.Last_time_used_s + LAST_UNUSED_MAX_S >= time.Now().Unix() {
		curr_location = int_device_info.Curr_location
	}

	return curr_location
}

func updateUserLocation(user_location *ModsFileInfo.UserLocation, new_location string) {
	if new_location != UNKNOWN_LOCATION {
		user_location.Last_known_location = user_location.Curr_location
	}

	if new_location != user_location.Curr_location {
		user_location.Prev_location = user_location.Curr_location
		user_location.Prev_last_time_checked_s = user_location.Last_time_checked_s

		user_location.Curr_location = new_location
	}
	user_location.Last_time_checked_s = time.Now().Unix()
}
