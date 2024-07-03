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

package RRComm

// Reminder is the format of a reminder
type Reminder struct {
	// Id is the reminder ID
	Id 		    string
	// Devices is the devices the reminder is set for
	Devices     []string
	// Message is the reminder message
	Message     string
	// Command is the command to be executed when the reminder is triggered on the chosen Devices
	Command     string
	// Time is the time in minutes the reminder is set for
	Time        string
	// Repeat_each is the time in minutes between each repeatition of the reminder
	Repeat_each int64
	// User_location is the location the user must be in for the reminder to be triggered
	User_location string
	// Device_condition is an additional "advanced" condition for the reminder in Go language
	Device_condition string
}
