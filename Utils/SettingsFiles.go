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

package Utils

import (
	"Utils/ModsFileInfo"
	"errors"
	"sync"
)

const USER_SETTINGS_FILE string = "UserSettings_EOG.dat"
const GEN_SETTINGS_FILE_CLIENT string = "GeneratedSettingsClient_EOG.dat"
const _GEN_SETTINGS_FILE_SERVER string = "GeneratedSettingsServer_EOG.dat"

// User_settings_GL is the global variable that holds the user settings. It is saved to the USER_SETTINGS_FILE file
// every 5 seconds.
var user_settings_GL UserSettings
// Gen_settings_GL is the global variable that holds the general settings. It is saved to the GEN_SETTINGS_FILE_CLIENT
// or _GEN_SETTINGS_FILE_SERVER file every 5 seconds.
var gen_settings_GL GenSettings

var VISOR_server_GL bool = false

var Password_GL string = ""

var mutex_US_GL sync.Mutex
var mutex_GS_GL sync.Mutex

type UserSettings struct {
	General         ModsFileInfo.GeneralConsts
	SMARTChecker    ModsFileInfo.Mod3UserInfo
	RSSFeedNotifier ModsFileInfo.Mod4UserInfo
	OnlineInfoChk   ModsFileInfo.Mod6UserInfo
	GPTCommunicator ModsFileInfo.Mod7UserInfo
	WebsiteBackend  ModsFileInfo.Mod8UserInfo
	TasksExecutor   ModsFileInfo.Mod9UserInfo
	UserLocator     ModsFileInfo.Mod12UserInfo
	GoogleManager   ModsFileInfo.Mod14UserInfo
}

type GenSettings struct {
	Device_settings ModsFileInfo.DeviceSettings
	MOD_3           ModsFileInfo.Mod3GenInfo
	MOD_4           ModsFileInfo.Mod4GenInfo
	MOD_5           ModsFileInfo.Mod5GenInfo
	MOD_6           ModsFileInfo.Mod6GenInfo
	MOD_7           ModsFileInfo.Mod7GenInfo
	MOD_9           ModsFileInfo.Mod9GenInfo
	MOD_10          ModsFileInfo.Mod10GenInfo
	MOD_12          ModsFileInfo.Mod12GenInfo
	MOD_14		    ModsFileInfo.Mod14GenInfo
	Registry        []*Value
}

/*
GetUserSettings is the function that returns the User settings.

This function is thread-safe. It uses a mutex to lock and unlock the access to the global variable.

DO NOT SET THIS POINTER ON VARIABLES!!! Each time you want to use the settings, call this function and use its pointer
directly - not indirectly! It's like if it was a temporary pointer that only lasts for one usage.

-----------------------------------------------------------

– Returns:
  - a pointer to the user_settings_GL variable
 */
func GetUserSettings() *UserSettings {
	mutex_US_GL.Lock()
	defer mutex_US_GL.Unlock()

	return &user_settings_GL
}

/*
GetGenSettings is the function that returns the Generated settings.

This function is thread-safe. It uses a mutex to lock and unlock the access to the global variable.

DO NOT SET THIS POINTER ON VARIABLES!!! Each time you want to use the settings, call this function and use its pointer
directly - not indirectly! It's like if it was a temporary pointer that only lasts for one usage.

-----------------------------------------------------------

– Returns:
  - a pointer to the gen_settings_GL variable
 */
func GetGenSettings() *GenSettings {
	mutex_GS_GL.Lock()
	defer mutex_GS_GL.Unlock()

	return &gen_settings_GL
}

/*
ReadSettingsFile is the function that reads the User and Generated settings from disk.

-----------------------------------------------------------

– Params:
  - user_settings – true if the user settings should be read, false if the generated settings should be read

– Returns:
  - an error if the settings file was not found or if the JSON file could not be parsed, nil otherwise
*/
func ReadSettingsFile(user_settings bool) error {
	var settings_file string = USER_SETTINGS_FILE
	if !user_settings {
		settings_file = GEN_SETTINGS_FILE_CLIENT
		if VISOR_server_GL {
			settings_file = _GEN_SETTINGS_FILE_SERVER
		}
	}
	var backup_file string = settings_file + ".bak"

	var bin_dir GPath = GetBinDirFILESDIRS()
	var bytes []byte = bin_dir.Add2(false, settings_file).ReadFile()

	decryptToJson := func() error {
		if Password_GL != "" {
			bytes = DecryptBytesCRYPTOENDECRYPT([]byte(Password_GL), []byte(Password_GL), bytes, nil)
		}

		var p_settings any = GetGenSettings()
		if user_settings {
			p_settings = GetUserSettings()
		}
		if err := FromJsonGENERAL(bytes, p_settings); err != nil {
			return err
		}

		return nil
	}

	// Try to decrypt and parse the obtained JSON file (normal or backup)
	if err := decryptToJson(); err != nil {
		// If the decryption and/or parsing failed, maybe the file was empty or corrupted. So try to read the backup
		// file.
		bytes = bin_dir.Add2(false, backup_file).ReadFile()
		if bytes == nil {
			var user_generated string = "generated"
			if user_settings {
				user_generated = "user"
			}

			return errors.New("no valid " + user_generated + " settings file found in the directory: \"" +
				bin_dir.GPathToStringConversion() + "\" - aborting")
		}

		if err = decryptToJson(); err != nil {
			// If not even the backup file could be decrypted and/or parsed, return the error
			return err
		}
	}

	return nil
}

/*
WriteSettingsFile is the function that writes the User and Generated settings to disk.

-----------------------------------------------------------

– Params:
  - user_settings – true if the user settings should be saved, false if the generated settings should be saved

– Returns:
  - true if the settings were successfully saved, false otherwise
 */
func WriteSettingsFile(user_settings bool) bool {
	var settings any = *GetGenSettings()
	if user_settings {
		settings = *GetUserSettings()
	}
	var p_string *string = ToJsonGENERAL(settings)
	if p_string == nil {
		return false
	}

	var to_write []byte = []byte(*p_string)
	if Password_GL != "" {
		to_write = EncryptBytesCRYPTOENDECRYPT([]byte(Password_GL), []byte(Password_GL), to_write, nil)
	}

	var settings_file string = USER_SETTINGS_FILE
	if !user_settings {
		settings_file = GEN_SETTINGS_FILE_CLIENT
		if VISOR_server_GL {
			settings_file = _GEN_SETTINGS_FILE_SERVER
		}
	}
	var backup_file string = settings_file + ".bak"

	var bin_dir GPath = GetBinDirFILESDIRS()
	_ = bin_dir.Add2(false, settings_file).WriteFile(to_write, false)
	_ = bin_dir.Add2(false, backup_file).WriteFile(to_write, false)

	return true
}
