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

package MOD_9

import (
	"Utils"
	"github.com/schollz/wifiscan"
	"log"
	"strings"
)

const _TIME_SLEEP_S int = 5

const _DEFAULT_LOCATION string = "unknown"

const _LOCATION_FILE string = "curr_loc.txt"

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](Utils.NUM_MOD_UserLocator, realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			// Get the current wifi SSID and check if it's in the list of locations
			// (The "device" parameter of GetSSID is ignored when the OS is Windows)
			wifi_nets, _ := wifiscan.Scan()
			log.Println("Nearby wifi networks:", wifi_nets)
			for _, wifi_net := range wifi_nets {
				if wifi_net.SSID == "" {
					continue
				}

				var wifi_bssid string = strings.ToUpper(wifi_net.SSID)

				log.Println(modUserInfo.LocsInfo)

				for _, locInfo := range modUserInfo.LocsInfo {
					log.Println(locInfo.Address == wifi_bssid)
					if locInfo.Type == "wifi" && locInfo.Address == wifi_bssid {
						log.Println("Location found:", locInfo.Location)
						_ = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, _LOCATION_FILE).WriteTextFile(locInfo.Location,false)

						break
					}
				}
			}

			if Utils.WaitWithStop(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}
