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
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

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
