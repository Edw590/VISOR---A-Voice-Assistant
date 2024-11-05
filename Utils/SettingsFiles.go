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

const USER_SETTINGS_FILE string = "UserSettings_EOG.dat"
const GEN_SETTINGS_FILE_CLIENT string = "GeneratedSettingsClient_EOG.dat"
const _GEN_SETTINGS_FILE_SERVER string = "GeneratedSettingsServer_EOG.dat"

// User_settings_GL is the global variable that holds the user settings. It is saved to the USER_SETTINGS_FILE file
// every 5 seconds.
var User_settings_GL UserSettings
// Gen_settings_GL is the global variable that holds the general settings. It is saved to the GEN_SETTINGS_FILE_CLIENT
// file every 5 seconds.
var Gen_settings_GL GenSettings

var VISOR_server_GL bool = false

var Password_GL string = ""

type UserSettings struct {
	General         ModsFileInfo.GeneralConsts
	SMARTChecker    ModsFileInfo.Mod3UserInfo
	RSSFeedNotifier ModsFileInfo.Mod4UserInfo
	OnlineInfoChk   ModsFileInfo.Mod6UserInfo
	GPTCommunicator ModsFileInfo.Mod7UserInfo
	WebsiteBackend  ModsFileInfo.Mod8UserInfo
	TasksExecutor   ModsFileInfo.Mod9UserInfo
	UserLocator     ModsFileInfo.Mod12UserInfo
}

///////////////////////////////////////////////////////////////

/*
WriteUserSettings is the function that saves the global variables of the UserSettings struct.

-----------------------------------------------------------

– Returns:
  - true if the user settings were successfully saved, false otherwise
*/
func WriteUserSettings() bool {
	var p_string *string = ToJsonGENERAL(User_settings_GL)
	if p_string == nil {
		return false
	}

	var to_write []byte = []byte(*p_string)
	if Password_GL != "" {
		to_write = EncryptBytesCRYPTOENDECRYPT([]byte(Password_GL), []byte(Password_GL), to_write, nil)
	}

	if err := os.WriteFile(USER_SETTINGS_FILE, to_write, 0777); err != nil {
		return false
	}

	return true
}

///////////////////////////////////////////////////////////////

/*
readGenSettings is the function that initializes the global variables of the GenSettings struct.

-----------------------------------------------------------

– Params:
  - server – true if the version running is the server version, false if it is the client version

– Returns:
  - an error if the settings file was not found or if the JSON file could not be parsed, nil otherwise
*/
func readGenSettings(server bool) error {
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

	var to_read []byte = bytes
	if Password_GL != "" {
		to_read = DecryptBytesCRYPTOENDECRYPT([]byte(Password_GL), []byte(Password_GL), bytes, nil)
	}

	if err = FromJsonGENERAL(to_read, &Gen_settings_GL); err != nil {
		return err
	}

	return nil
}

/*
writeGenSettings is the function that saves the global variables of the GenSettings struct to the _GEN_SETTINGS_FILE file.

-----------------------------------------------------------

– Params:
  - server – true if the generated settings were successfully saved, false otherwise
*/
func writeGenSettings(server bool) bool {
	var settings_file string = GEN_SETTINGS_FILE_CLIENT
	if server {
		settings_file = _GEN_SETTINGS_FILE_SERVER
	}
	var p_string *string = ToJsonGENERAL(Gen_settings_GL)
	if p_string == nil {
		return false
	}

	var to_write []byte = []byte(*p_string)
	if Password_GL != "" {
		to_write = EncryptBytesCRYPTOENDECRYPT([]byte(Password_GL), []byte(Password_GL), to_write, nil)
	}

	if err := os.WriteFile(settings_file, to_write, 0777); err != nil {
		return false
	}

	return true
}
