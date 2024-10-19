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

import (
	"Utils"
	"bytes"
	"time"
)

const _GET_SETTINGS_EACH_S int64 = 30

var last_crc16_GL []byte = nil
var stop_GL bool = false

/*
SyncUserSettings keeps synchronizing the remote user settings file with the local one.

This function only returns when it's stopped with StopUserSettingsSyncer().
*/
func SyncUserSettings() {
	var last_get_settings_when_s int64 = 0
	for {
		var update_settings bool = false
		if time.Now().Unix() >= last_get_settings_when_s + _GET_SETTINGS_EACH_S && Utils.IsCommunicatorConnectedSERVER() {
			update_settings = true

			last_get_settings_when_s = time.Now().Unix()
		}

		if update_settings {
			Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, []byte("US|false"))
			var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_SettingsSync]
			if comms_map == nil {
				return
			}

			var new_crc16 []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
			if !bytes.Equal(new_crc16, last_crc16_GL) {
				last_crc16_GL = new_crc16

				Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, []byte("US|true"))
				comms_map = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_SettingsSync]
				if comms_map == nil {
					return
				}

				var json []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

				_ = Utils.FromJsonGENERAL(json, &Utils.User_settings_GL)
			}
		}

		if Utils.WaitWithStopTIMEDATE(&stop_GL, 1) {
			return
		}
	}
}

/*
StopUserSettingsSyncer stops the user settings synchronizer.
 */
func StopUserSettingsSyncer() {
	stop_GL = true
}
