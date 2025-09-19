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

package ModsFileInfo

// Mod14GenInfo is the format of the custom generated information about this specific module.
type Mod14GenInfo struct {
	// Token is the cached token
	Token string
	// Token_invalid is whether the token is invalid
	Token_invalid bool
	// Token_invalid_notified is whether the user has been notified about the invalid token
	Token_invalid_notified bool
	// Calendars maps calendar IDs to their corresponding GCalendar structs
	Calendars map[string]GCalendar
	// Events is the list of all events associated with the account
	Events []GEvent
	// Tasks is the list of all tasks associated with the account
	Tasks []GTask
}

// GTask represents a Google Task
type GTask struct {
	// Id is the ID of the task
	Id string
	// Title is the title of the task
	Title string
	// Details are the details of the task
	Details string
	// Date_s is the timestamp of the date of the task in seconds
	Date_s int64
	// Completed is whether the task is completed
	Completed bool
}

// GEvent represents a Google Calendar Event
type GEvent struct {
	// Id is the ID of the event
	Id string
	// Calendar_id is the calendar ID associated with the event
	Calendar_id string
	// Summary is the title of the event
	Summary string
	// Location is the location of the event
	Location string
	// Description is the description of the event
	Description string
	// Start_time_s is the timestamp of the start time of the event in seconds
	Start_time_s int64
	// Duration_min is the duration of the event in minutes
	Duration_min int64
}

// GCalendar represents a Google Calendar
type GCalendar struct {
	// Title is the title of the calendar
	Title string
	// Enabled is whether the calendar is enabled for usage within VISOR
	Enabled bool
}

///////////////////////////////////////////////////////////////////////////////

// Mod14UserInfo is the format of the custom information file about this specific module.
type Mod14UserInfo struct {
	// Credentials_JSON is the text from the credentials.json file obtained from Google
	Credentials_JSON string
}
