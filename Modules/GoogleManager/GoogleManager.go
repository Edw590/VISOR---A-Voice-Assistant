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

package GoogleManager

import (
	"Utils"
	"Utils/ModsFileInfo"
	"context"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/tasks/v1"
	"log"
	"net/http"
)

// _SCOPES defines the Google Calendar scope required for read-only access
var _SCOPES = []string{
	calendar.CalendarScope,
	tasks.TasksScope,
	gmail.GmailModifyScope,
}

var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
	modGenInfo_GL  *ModsFileInfo.Mod14GenInfo = &Utils.Gen_settings_GL.MOD_14
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_14

		for {
			// Parse credentials to config
			config, err := ParseConfigJSON()
			if err != nil {
				log.Printf("Unable to parse client secret file to config: %v\n", err)

				return
			}
			client := getClient(config)
			if client == nil {
				//log.Println("No token saved")

				return
			}

			// Store calendar events
			storeCalendarsEvents(client)

			// Store tasks
			storeTasks(client)

			if Utils.WaitWithStopTIMEDATE(module_stop, 60) {
				return
			}
		}
	}
}

// getClient retrieves a token, saves it, and returns a new client
func getClient(config *oauth2.Config) *http.Client {
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
	err := Utils.FromJsonGENERAL([]byte(modGenInfo_GL.Token), &token)

	return &token, err
}
