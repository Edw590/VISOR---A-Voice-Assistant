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

package GoogleManager

import (
	"Utils"
	"Utils/ModsFileInfo"
	"context"
	"net/http"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/tasks/v1"
)

// _SCOPES defines the Google Calendar scope required for read-only access
var _SCOPES = []string{
	calendar.CalendarScope,
	tasks.TasksScope,
	gmail.GmailModifyScope,
}

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	go func() {
		for {
			var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GoogleManager, 0, -1)
			if comms_map == nil {
				break
			}

			var to_process []byte = nil
			var is_event bool = true
			if _, ok := comms_map["Event"]; ok {
				to_process = comms_map["Event"].([]byte)
			} else if _, ok = comms_map["Task"]; ok {
				to_process = comms_map["Task"].([]byte)
				is_event = false
			} else {
				continue
			}
			var to_process_str string = string(to_process)

			client := getClient()
			if client == nil {
				Utils.WaitWithStopDATETIME(module_stop, 15)

				continue
			}

			if is_event {
				var event ModsFileInfo.GEvent
				err := Utils.FromJsonGENERAL([]byte(to_process_str), &event)
				if err != nil {
					Utils.LogfError("Unable to parse calendar event: %v\n", err)

					continue
				}

				if !addEvent(event, client) {
					Utils.LogfError("Unable to add calendar event: %v\n", err)
				}
			} else {
				var task ModsFileInfo.GTask
				err := Utils.FromJsonGENERAL([]byte(to_process_str), &task)
				if err != nil {
					Utils.LogfError("Unable to parse task: %v\n", err)

					continue
				}

				if !addTask(task, client) {
					Utils.LogfError("Unable to add task: %v\n", err)
				}
			}
		}
	}()

	for {
		getModGenSettings().Token_invalid = true

		client := getClient()
		if client == nil {
			Utils.WaitWithStopDATETIME(module_stop, 15)

			continue
		}

		// Store calendar events
		if storeCalendarsEvents(client) {
			setTokenValid()
		}

		// Store tasks
		if storeTasks(client) {
			setTokenValid()
		}

		if getModGenSettings().Token_invalid && !getModGenSettings().Token_invalid_notified {
			var msg_body string = "The saved Google token is invalid. Please re-authenticate."
			var things_replace = map[string]string{
				Utils.MODEL_INFO_DATE_TIME_EMAIL: Utils.GetDateTimeStrDATETIME(-1),
				Utils.MODEL_INFO_MSG_BODY_EMAIL:  msg_body,
			}
			var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_INFO, things_replace)
			email_info.Subject = "Google token is INVALID"
			err := Utils.QueueEmailEMAIL(email_info)
			if err == nil {
				getModGenSettings().Token_invalid_notified = true
			}
		}

		if Utils.WaitWithStopDATETIME(module_stop, 60) {
			return
		}
	}
}

// getClient retrieves a token, saves it, and returns a new client
func getClient() *http.Client {
	// Parse credentials to config
	config, err := ParseConfigJSON()
	if err != nil {
		Utils.LogfError("Unable to parse client secret file to config: %v\n", err)

		return nil
	}

	// Check if the token file exists
	token, err := getToken()
	if err != nil {
		return nil
	}

	return config.Client(context.Background(), token)
}

// getToken retrieves a token from a local file
func getToken() (*oauth2.Token, error) {
	var token oauth2.Token
	err := Utils.FromJsonGENERAL([]byte(getModGenSettings().Token), &token)

	return &token, err
}

func getModGenSettings() *ModsFileInfo.Mod14GenInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14
}

func setTokenValid() {
	getModGenSettings().Token_invalid = false
	getModGenSettings().Token_invalid_notified = false
}
