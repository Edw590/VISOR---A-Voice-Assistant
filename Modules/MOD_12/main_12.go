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
	"ULComm/ULComm"
	"Utils"
	"log"
)

// User Locator //

const TIME_SLEEP_S int = 30

type DeviceInfo struct {
	Device_id string
	Curr_location string
}

type _MGIModSpecInfo any
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		log.Println("--------------------------------")
		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			var devices []ULComm.DeviceInfo = nil
			for _, file_info := range moduleInfo_GL.ModDirsInfo.UserData.Add2(true, "devices").GetFileList() {
				var device ULComm.DeviceInfo
				if err := Utils.FromJsonGENERAL(file_info.GPath.ReadFile(), &device); err == nil {
					devices = append(devices, device)
				} else {
					log.Println("Error reading device file", file_info.GPath, ":", err)
				}
			}

			var device_infos []DeviceInfo = nil
			for _, device := range devices {
				var device_info DeviceInfo = DeviceInfo{
					Device_id: device.Device_id,
				}

				for _, wifi_net := range device.System_state.Connectivity_info.Wifi_networks {
					for _, location := range modUserInfo.Locs_info {
						if location.Type == "wifi" {
							if location.Address != "" {
								if wifi_net.BSSID == location.Address {
									device_info.Curr_location = location.Location

									break
								}
							} else {
								if wifi_net.SSID == location.Name {
									device_info.Curr_location = location.Location

									break
								}
							}
						}
					}
				}
				for _, bluetooth_device := range device.System_state.Connectivity_info.Bluetooth_devices {
					for _, location := range modUserInfo.Locs_info {
						if location.Type == "bluetooth" {
							if location.Address != "" {
								if bluetooth_device.Address == location.Address {
									device_info.Curr_location = location.Location

									break
								}
							} else {
								if bluetooth_device.Name == location.Name {
									device_info.Curr_location = location.Location

									break
								}
							}
						}
					}
				}

				log.Println("device_info:", device_info)

				device_infos = append(device_infos, device_info)
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, TIME_SLEEP_S) {
				return
			}
		}
	}
}
