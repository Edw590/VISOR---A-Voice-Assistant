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
	"GPT/GPT"
	"ULComm/ULComm"
	"Utils"
	"Utils/UtilsSWA"
	"log"
	"sort"
	"time"
)

// User Locator //

const TIME_SLEEP_S int = 5

const UNKNOWN_LOCATION string = "3234_UNKNOWN"

const LAST_COMM_MAX int64 = 20

var Device_infos_ULComm_GL []ULComm.DeviceInfo = nil

type _IntDeviceInfo struct {
	Device_id           string
	Last_comm           int64
	Last_time_used      int64
	Curr_location       string
	Last_known_location string
}

type _MGI any
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGI]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGI](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGI])

		var user_location ULComm.UserLocation
		var user_location_json Utils.GPath = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "user_location.json")
		if user_location_json.Exists() {
			_ = Utils.FromJsonGENERAL(user_location_json.ReadFile(), &user_location)
		}

		var device_infos []*_IntDeviceInfo = nil
		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			Device_infos_ULComm_GL = nil
			for _, file_info := range moduleInfo_GL.ModDirsInfo.UserData.Add2(true, "devices").GetFileList() {
				var device ULComm.DeviceInfo
				if err := Utils.FromJsonGENERAL(file_info.GPath.ReadFile(), &device); err == nil {
					Device_infos_ULComm_GL = append(Device_infos_ULComm_GL, device)
				} else {
					log.Println("Error reading device file", file_info.GPath, ":", err)
				}
			}

			log.Println("--------------------------------")
			for _, device_ULComm := range Device_infos_ULComm_GL {
				var device_info *_IntDeviceInfo
				for _, device_info1 := range device_infos {
					if device_info1.Device_id == device_ULComm.Device_id {
						device_info = device_info1

						break
					}
				}
				if device_info == nil {
					device_info = &_IntDeviceInfo{
						Device_id: device_ULComm.Device_id,
					}
				}

				// Update some information
				device_info.Last_comm = device_ULComm.Last_comm
				device_info.Last_time_used = device_ULComm.Last_time_used
				device_info.Curr_location = UNKNOWN_LOCATION

				for _, location_info := range modUserInfo.Locs_info {
					if location_info.Last_detection < int64(TIME_SLEEP_S) * 2 {
						// There must be a minimum. That minimum is the time it takes for the devices to update their
						// status, but double it to be sure they communicated.
						location_info.Last_detection = int64(TIME_SLEEP_S) * 2
					}

					var beacon_list []ULComm.ExtBeacon
					if location_info.Type == "wifi" {
						beacon_list = device_ULComm.System_state.Connectivity_info.Wifi_networks
					} else if location_info.Type == "bluetooth" {
						beacon_list = device_ULComm.System_state.Connectivity_info.Bluetooth_devices
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
									device_info.Last_comm + location_info.Last_detection >= time.Now().Unix() {
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

				log.Println("device_info:", device_info)

				device_infos = append(device_infos, device_info)
			}

			var curr_user_location string = getUserLocation(modUserInfo, device_infos)
			log.Println("Current user location:", curr_user_location)
			updateUserLocation(&user_location, curr_user_location)
			_ = user_location_json.WriteTextFile(*Utils.ToJsonGENERAL(user_location), false)

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

	var device_infos []*_IntDeviceInfo = getIntDeviceInfos()
	if device_id == GPT.ALL_DEVICES_ID {
		// Check if any device is active
		for _, device_info := range device_infos {
			if time.Now().Unix() - device_info.Last_comm <= LAST_COMM_MAX {
				return true
			}
		}

		return false
	}

	for _, device_info := range device_infos {
		if device_info.Device_id == device_id {
			return time.Now().Unix() - device_info.Last_comm <= LAST_COMM_MAX
		}
	}

	return false
}

func getIntDeviceInfos() []*_IntDeviceInfo {
	var device_infos []*_IntDeviceInfo = nil
	for _, file_info := range moduleInfo_GL.ModDirsInfo.UserData.Add2(true, "devices").GetFileList() {
		var device ULComm.DeviceInfo
		if err := Utils.FromJsonGENERAL(file_info.GPath.ReadFile(), &device); err == nil {
			device_infos = append(device_infos, &_IntDeviceInfo{
				Device_id: device.Device_id,
				Last_comm: device.Last_comm,
				Last_time_used: device.Last_time_used,
				Curr_location: UNKNOWN_LOCATION,
			})
		} else {
			log.Println("Error reading device file", file_info.GPath, ":", err)
		}
	}

	return device_infos
}

func getUserLocation(modUserInfo _ModUserInfo, devices []*_IntDeviceInfo) string {
	if modUserInfo.Devices_info.AlwaysWith_device_id != "" {
		for _, device := range devices {
			if device.Device_id == modUserInfo.Devices_info.AlwaysWith_device_id &&
					device.Curr_location != UNKNOWN_LOCATION {
				return device.Curr_location
			}
		}
	}

	sortDevicesByLastUsed(devices)
	var curr_location string = UNKNOWN_LOCATION
	for _, device := range devices {
		if device.Curr_location != UNKNOWN_LOCATION && device.Last_time_used + 5*60 >= time.Now().Unix() {
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

func updateUserLocation(user_location *ULComm.UserLocation, new_location string) {
	if new_location != UNKNOWN_LOCATION {
		user_location.Last_known_location = user_location.Curr_location
	}

	if new_location != user_location.Curr_location {
		user_location.Prev_location = user_location.Curr_location
		user_location.Prev_last_time_checked = user_location.Last_time_checked

		user_location.Curr_location = new_location
	}
	user_location.Last_time_checked = time.Now().Unix()
}
