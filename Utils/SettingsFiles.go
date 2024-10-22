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

package Utils

import (
	"Utils/ModsFileInfo"
	"errors"
	"os"
)

const DEVICE_SETTINGS_FILE string = "DeviceSettings_EOG.json"
const USER_SETTINGS_FILE string = "UserSettings_EOG.json"
const GEN_SETTINGS_FILE_CLIENT string = "GeneratedSettingsClient_EOG.json"
const _GEN_SETTINGS_FILE_SERVER string = "GeneratedSettingsServer_EOG.json"

var Device_settings_GL DeviceSettings
var User_settings_GL UserSettings
var Gen_settings_GL GenSettings
var VISOR_server_GL bool = false

type DeviceSettings struct {
	// Device_ID is the device ID of the current device
	Device_ID string
	// Device_type is the type of the current device
	Device_type string
	// Device_description is the description of the current device
	Device_description string
}

type UserSettings struct {
	General         _GeneralConsts
	SMARTChecker    ModsFileInfo.Mod2UserInfo
	RSSFeedNotifier ModsFileInfo.Mod4UserInfo
	OnlineInfoChk   ModsFileInfo.Mod6UserInfo
	GPTCommunicator ModsFileInfo.Mod7UserInfo
	WebsiteBackend  ModsFileInfo.Mod8UserInfo
	TasksExecutor   ModsFileInfo.Mod9UserInfo
	UserLocator     ModsFileInfo.Mod12UserInfo
}

type GenSettings struct {
	MOD_2  ModsFileInfo.Mod2GenInfo
	MOD_4  ModsFileInfo.Mod4GenInfo
	MOD_5  ModsFileInfo.Mod5GenInfo
	MOD_6  ModsFileInfo.Mod6GenInfo
	MOD_7  ModsFileInfo.Mod7GenInfo
	MOD_9  ModsFileInfo.Mod9GenInfo
	MOD_10 ModsFileInfo.Mod10GenInfo
	MOD_12 ModsFileInfo.Mod12GenInfo
	Registry []*Value
}

///////////////////////////////////////////////////////////////

type _GeneralConsts struct {
	// VISOR_email_addr is VISOR's email address
	VISOR_email_addr string
	// VISOR_email_pw is VISOR's email password
	VISOR_email_pw string

	// User_email_addr is the email address of the user, used for all email communication
	User_email_addr string

	// Website_domain is the domain of the VISOR website
	Website_domain string
	// Website_pw is the password for the VISOR website
	Website_pw string

	// WolframAlpha_AppID is the app ID for the Wolfram Alpha API
	WolframAlpha_AppID string

	// Picovoice_API_key is the API key for the Picovoice API
	Picovoice_API_key string
}

/*
loadDeviceSettings is the function that initializes the global variables of the DeviceSettings structs.

Call this before SettingsSync.LoadUserSettings.

-----------------------------------------------------------

– Params:
  - server – true if the version running is the server version, false if it is the client version

– Returns:
  - an error if the settings file was not found or if the JSON file could not be parsed, nil otherwise
*/
func loadDeviceSettings(server bool) error {
	bytes, err := os.ReadFile(DEVICE_SETTINGS_FILE)
	if err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "[ERROR]"
		}
		return errors.New("no " + DEVICE_SETTINGS_FILE + " file found in the current working directory: \"" + cwd + "\" - aborting")
	}

	if err = FromJsonGENERAL(bytes, &Device_settings_GL); err != nil {
		return err
	}

	VISOR_server_GL = server

	if Device_settings_GL.Device_ID == "" || Device_settings_GL.Device_type == "" ||
			Device_settings_GL.Device_description == "" {
		return errors.New("some fields in " + DEVICE_SETTINGS_FILE + " are empty or incorrect - aborting")
	}

	return nil
}

/*
SaveUserSettings is the function that saves the global variables of the UserSettings struct.

-----------------------------------------------------------

– Returns:
  - true if the user settings were successfully saved, false otherwise
 */
func SaveUserSettings() bool {
	var p_string *string = ToJsonGENERAL(User_settings_GL)
	if p_string == nil {
		return false
	}

	if err := os.WriteFile(USER_SETTINGS_FILE, []byte(*p_string), 0777); err != nil {
		return false
	}

	return true
}

///////////////////////////////////////////////////////////////

/*
loadGenSettings is the function that initializes the global variables of the GenSettings struct.

-----------------------------------------------------------

– Params:
  - server – true if the version running is the server version, false if it is the client version

– Returns:
  - an error if the settings file was not found or if the JSON file could not be parsed, nil otherwise
*/
func loadGenSettings(server bool) error {
	var settings_file string = GEN_SETTINGS_FILE_CLIENT
	if server {
		settings_file = _GEN_SETTINGS_FILE_SERVER
	}
	bytes, err := os.ReadFile(settings_file)
	if err != nil {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "[ERROR]"
		}
		return errors.New("no " + settings_file + " file found in the current working directory: \"" + cwd + "\" - aborting")
	}

	if err := FromJsonGENERAL(bytes, &Gen_settings_GL); err != nil {
		return err
	}

	return nil
}

/*
saveGenSettings is the function that saves the global variables of the GenSettings struct to the _GEN_SETTINGS_FILE file.

-----------------------------------------------------------

– Params:
  - server – true if the generated settings were successfully saved, false otherwise
*/
func saveGenSettings(server bool) bool {
	var settings_file string = GEN_SETTINGS_FILE_CLIENT
	if server {
		settings_file = _GEN_SETTINGS_FILE_SERVER
	}
	var p_string *string = ToJsonGENERAL(Gen_settings_GL)
	if p_string == nil {
		return false
	}

	if err := os.WriteFile(settings_file, []byte(*p_string), 0777); err != nil {
		return false
	}

	return true
}
