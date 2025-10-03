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

package GMan

import (
	"Utils"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

/*
setToken sets the token for the Google Manager.

-----------------------------------------------------------

– Params:
  - token – the token to be set
*/
func setToken(token *oauth2.Token) {
	var message []byte = []byte("S_S|GManTok|")
	token_bytes, _ := json.Marshal(token)
	message = append(message, token_bytes...)
	Utils.QueueNoResponseMessageSERVER(message)
}

/*
IsTokenValid checks if the token is valid.

-----------------------------------------------------------

– Returns:
  - true if the token is valid, false otherwise
 */
func IsTokenValid() bool {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GMan, 2, []byte("G_S|true|GManTokVal")) {
		return false
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GMan, 2, 10)
	if comms_map == nil {
		return false
	}

	var response []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	ret, err := strconv.ParseBool(string(response))
	if err != nil {
		ret = false
	}

	return ret
}

/*
ParseConfigJSON parses the Google OAuth2 configuration from the user settings.

-----------------------------------------------------------

– Returns:
  - the OAuth2 configuration
  - error if any
 */
func ParseConfigJSON() (*oauth2.Config, error) {
	var credentials string = Utils.GetUserSettings(Utils.LOCK_UNLOCK).GoogleManager.Credentials_JSON
	if credentials == "" {
		//log.Println("No credentials found in the user settings file")

		return nil, errors.New("no credentials found in the user settings file")
	}

	// Load the credentials from the file
	config, err := google.ConfigFromJSON([]byte(credentials), _SCOPES...)
	if err != nil {
		//log.Printf("Unable to parse client secret file to config: %v\n", err)

		return nil, err
	}

	return config, err
}

/*
GetAuthUrl generates the Google authorization URL.

-----------------------------------------------------------

– Returns:
  - the authorization URL
  - error if any
 */
func GetAuthUrl() (string, error) {
	config, err := ParseConfigJSON()
	if err != nil {
		return "", err
	}

	auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	return auth_url, nil
}

/*
StoreTokenFromAuthCode stores the token obtained from the given authorization code.

-----------------------------------------------------------

– Params:
  - auth_code – the authorization code

– Returns:
  - error if any
*/
func StoreTokenFromAuthCode(auth_code string) error {
	config, err := ParseConfigJSON()
	if err != nil {
		return err
	}

	token, err := config.Exchange(context.Background(), auth_code)
	if err != nil {
		return err
	}

	if !Utils.IsCommunicatorConnectedSERVER() {
		return errors.New("not connected to the server")
	}

	setToken(token)

	return nil
}
