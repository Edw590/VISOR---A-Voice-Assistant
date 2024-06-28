/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

package MOD_9

// _ModSpecInfo is the format of the custom generated information about this specific module.
type _ModSpecInfo struct {
	// Reminders_info maps the reminder ID to the last time it was reminded in minutes
	Reminders_info map[string]int64
}

// _ModUserInfo is the format of the custom information file about this specific module.
type _ModUserInfo struct {
	Reminders []_Reminder
}

// _Reminder is the format of a reminder
type _Reminder struct {
	// Id is the reminder id
	Id 		    string
	// Message is the reminder message
	Message     string
	// Time is the time in minutes the reminder is set for
	Time        string
	// Repeat_each is the time in minutes between each repeatition of the reminder
	Repeat_each int64
	// Location is the location of the reminder in case there is any
	Location    string
}
