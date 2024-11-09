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

package ModsFileInfo

// Mod14GenInfo is the format of the custom generated information about this specific module.
type Mod14GenInfo struct {
	// Token is the cached token
	Token string
	// Events is the list of events, with the keys being the events IDs
	Events []Event
}

type Event struct {
	// Id is the ID of the event
	Id string
	// Summary is the title of the event
	Summary string
	// Location is the location of the event
	Location string
	// Description is the description of the event
	Description string
	// Start_time is the time of the event in RFC3339 format
	Start_time string
	// Duration_min is the duration of the event in minutes
	Duration_min int64
}

///////////////////////////////////////////////////////////////////////////////

// Mod14UserInfo is the format of the custom information file about this specific module.
type Mod14UserInfo struct {
	// Credentials_JSON is the text from the credentials.json file obtained from Google
	Credentials_JSON string
}