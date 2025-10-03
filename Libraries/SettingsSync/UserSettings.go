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

package SettingsSync

import (
	"Utils"
	"bytes"
	"errors"
	"time"
)

const _GET_SETTINGS_EACH_S int64 = 30

var last_remote_crc16_GL []byte = nil
var stop_GL bool = false

/*
SyncUserSettings keeps synchronizing the remote user settings file with the local settings in memory in background.

-----------------------------------------------------------

– Params:
  - loop – if true, the function will keep running until it's stopped with StopUserSettingsSyncer()
*/
func SyncUserSettings() {
	go func() {
		var last_get_settings_when_s int64 = 0
		var last_user_settings_json string = GetJsonUserSettings()
		for {
			var new_user_settings_json string = GetJsonUserSettings()
			if last_user_settings_json != new_user_settings_json {
				last_user_settings_json = new_user_settings_json

				var message []byte = []byte("S_S|US|")
				message = append(message, last_user_settings_json...)
				Utils.QueueNoResponseMessageSERVER(message)
			} else {
				var get_settings bool = false
				if time.Now().Unix() >= last_get_settings_when_s + _GET_SETTINGS_EACH_S && Utils.IsCommunicatorConnectedSERVER() {
					get_settings = true

					last_get_settings_when_s = time.Now().Unix()
				}

				if get_settings && areRemoteSettingsChanged() {
					if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, 0, []byte("G_S|true|US")) {
						goto end
					}
					var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_SettingsSync, 0, 10)
					if comms_map == nil {
						goto end
					}

					var json []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

					_ = Utils.FromJsonGENERAL(json, Utils.GetUserSettings(Utils.LOCK_UNLOCK))
					last_user_settings_json = GetJsonUserSettings()
				}
			}

			end:

			if Utils.WaitWithStopDATETIME(&stop_GL, 1) {
				return
			}
		}
	}()
}

func areRemoteSettingsChanged() bool {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_SettingsSync, 1, []byte("G_S|false|US")) {
		return false
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_SettingsSync, 1, 10)
	if comms_map == nil {
		return false
	}

	var map_value []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
	if map_value == nil {
		return false
	}

	if !bytes.Equal(map_value, last_remote_crc16_GL) {
		last_remote_crc16_GL = map_value

		return true
	}

	return false
}

/*
StopUserSettingsSyncer stops the user settings synchronizer.
 */
func StopUserSettingsSyncer() {
	stop_GL = true
}

/*
GetUserSettings returns the user settings in JSON format.

-----------------------------------------------------------

– Returns:
  - the user settings in JSON format
 */
func GetJsonUserSettings() string {
	return *Utils.ToJsonGENERAL(*Utils.GetUserSettings(Utils.LOCK_UNLOCK))
}

/*
LoadUserSettings loads the user settings from the given JSON string.

-----------------------------------------------------------

– Params:
  - json – the JSON string to load the user settings from

– Returns:
  - an error if the user settings could not be loaded, nil otherwise
 */
func LoadUserSettings(json string) error {
	if json == "" {
		return errors.New("empty JSON string")
	}

	var user_settings Utils.UserSettings
	if err := Utils.FromJsonGENERAL([]byte(json), &user_settings); err != nil {
		return err
	}

	*Utils.GetUserSettings(Utils.LOCK_UNLOCK) = user_settings

	return nil
}

/*
IsWebsiteInfoEmpty returns true if the website domain and password are empty, false otherwise.

-----------------------------------------------------------

– Returns:
  - true if the website domain and password are empty, false otherwise
 */
func IsWebsiteInfoEmpty() bool {
	return Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Website_domain == "" && Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Website_pw == ""
}
