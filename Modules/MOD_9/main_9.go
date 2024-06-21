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
	"strconv"
)

const _TIME_SLEEP_S int = 5

const _DEFAULT_LOCATION string = "unknown"

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		for {
			// Get the current wifi SSID and check if it's in the list of locations
			wifi_nets, err := wifiscan.Scan()
			if err != nil {
				// TODO: do something with this. Can mean the adapter is off
			}
			//log.Println("Nearby wifi networks:", wifi_nets)
			var final_string string = ""
			for _, wifi_net := range wifi_nets {
				final_string += wifi_net.BSSID + "|" + strconv.Itoa(wifi_net.RSSI) + "|" + wifi_net.SSID + "\n"
			}
			//log.Println("Final string:", final_string)

			// TODO: send the string to the server
			//  And be careful processing it! It can have ANY characters, including \n...

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}
