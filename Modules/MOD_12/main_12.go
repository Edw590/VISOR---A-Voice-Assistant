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

import (
	"GPTComm/GPTComm"
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"sort"
	"time"
)

// User Locator //

const TIME_SLEEP_S int = 5

const UNKNOWN_LOCATION string = "3234_UNKNOWN"

const LAST_COMM_MAX_S int64 = 30 + 5 // MUST BE HIGHER THAN MOD_8.PONG_WAIT!!!
// LAST_UNUSED_MAX_S is the maximum time since the last time active to consider the device in a location
const LAST_UNUSED_MAX_S int64 = 5 * 60

type _IntDeviceInfo struct {
	Device_id           string
	Last_comm           int64
	Last_time_used      int64
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

		var device_infos []*_IntDeviceInfo = nil
		for {
			//log.Println("--------------------------------")
			for _, more_device_info := range modGenInfo_GL.More_devices_info {
				var device_info *_IntDeviceInfo
				for _, device_info1 := range device_infos {
					if device_info1.Device_id == more_device_info.Device_id {
						device_info = device_info1

						break
					}
				}
				if device_info == nil {
					device_info = &_IntDeviceInfo{
						Device_id: more_device_info.Device_id,
					}
				}

				// Update some information
				device_info.Last_comm = more_device_info.Last_comm_s
				device_info.Last_time_used = more_device_info.Last_time_used_s
				device_info.Curr_location = UNKNOWN_LOCATION

				for _, location_info := range modUserInfo_GL.Locs_info {
					if location_info.Last_detection_s < int64(TIME_SLEEP_S) * 2 {
						// There must be a minimum. That minimum is the time it takes for the devices to update their
						// status, but double it to be sure they communicated.
						location_info.Last_detection_s = int64(TIME_SLEEP_S) * 2
					}

					var beacon_list []ModsFileInfo.ExtBeacon
					if location_info.Type == "wifi" {
						beacon_list = more_device_info.Device_info.System_state.Connectivity_info.Wifi_networks
					} else if location_info.Type == "bluetooth" {
						beacon_list = more_device_info.Device_info.System_state.Connectivity_info.Bluetooth_devices
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

							if distance <= location_info.Max_distance &&
									device_info.Last_comm + location_info.Last_detection_s >= time.Now().Unix() {
								// If the device was near the location and the last communication was recent, then the
								// user is near the location.
								device_info.Curr_location = location_info.Location
							}

							break
						}
					}
				}

				if device_info.Curr_location != UNKNOWN_LOCATION {
					device_info.Last_known_location = device_info.Curr_location
				}

				//log.Println("device_info:", device_info)

				device_infos = append(device_infos, device_info)
			}

			var curr_user_location string = computeUserLocation(device_infos)
			//log.Println("Current user location:", curr_user_location)
			updateUserLocation(curr_user_location)

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

func IsDeviceActive(device_id string) bool {
	if device_id == "" {
		return false
	}

	var more_devices_info []ModsFileInfo.MoreDeviceInfo = modGenInfo_GL.More_devices_info
	if device_id == GPTComm.ALL_DEVICES_ID {
		// Check if any device is active
		for _, more_device_info := range more_devices_info {
			if time.Now().Unix() - more_device_info.Last_comm_s <= LAST_COMM_MAX_S {
				return true
			}
		}

		return false
	}

	for _, more_device_info := range more_devices_info {
		if more_device_info.Device_id == device_id {
			return time.Now().Unix() - more_device_info.Last_comm_s <= LAST_COMM_MAX_S
		}
	}

	return false
}

func computeUserLocation(devices []*_IntDeviceInfo) string {
	if modUserInfo_GL.Devices_info.AlwaysWith_device_id != "" {
		for _, device := range devices {
			if device.Device_id == modUserInfo_GL.Devices_info.AlwaysWith_device_id &&
					device.Curr_location != UNKNOWN_LOCATION {
				return device.Curr_location
			}
		}
	}

	sortDevicesByLastUsed(devices)
	var curr_location string = UNKNOWN_LOCATION
	for _, device := range devices {
		if device.Curr_location != UNKNOWN_LOCATION && device.Last_time_used + LAST_UNUSED_MAX_S >= time.Now().Unix() {
			curr_location = device.Curr_location

			break
		}
	}

	return curr_location
}

func sortDevicesByLastUsed(devices []*_IntDeviceInfo) {
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Last_time_used > devices[j].Last_time_used
	})
}

func updateUserLocation(new_location string) {
	var user_location *ModsFileInfo.UserLocation = &modGenInfo_GL.User_location
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
