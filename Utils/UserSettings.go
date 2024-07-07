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

package Utils

import "OnlineInfoChk/OICNews"

var User_settings_GL UserSettings

type UserSettings struct {
	PersonalConsts _PersonalConsts
	MOD_2  _MOD_2
	MOD_4  _MOD_4
	MOD_6  _MOD_6
	MOD_7  _MOD_7
	MOD_10 _MOD_10
	MOD_12 _MOD_12
}

///////////////////////////////////////////////////////////////

type _PersonalConsts struct {
	// Device_ID is the device ID of the current device
	Device_ID string

	// VISOR_dir is the full path to the main directory of VISOR.
	VISOR_dir string
	// VISOR_server is an INTERNAL attribute to be filled INTERNALLY that indicates if the version running is the server
	// or the client version
	VISOR_server bool

	// VISOR_email_addr is VISOR's email address
	VISOR_email_addr string
	// VISOR_email_pw is VISOR's email password
	VISOR_email_pw string

	// User_email_addr is the email address of the user, used for all email communication
	User_email_addr string

	// Website_url is the URL of the VISOR website
	Website_url string
	// Website_pw is the password for the VISOR website
	Website_pw string
	// Website_dir is the full path to the directory of the VISOR website
	Website_dir string

	// WolframAlpha_AppID is the app ID for the Wolfram Alpha API
	WolframAlpha_AppID string

	// Picovoice_API_key is the API key for the Picovoice API
	Picovoice_API_key string
}

///////////////////////////////////////////////////////////////

type _MOD_2 struct {
	// Disks_info is the information about the disks. It maps the disk serial number to the disk information struct.
	Disks_info _DiskInfo
}

type _DiskInfo struct {
	// Disk label
	Label string
	// Is the disk an HDD?
	Is_HDD bool
}

///////////////////////////////////////////////////////////////

type _MOD_4 struct {
	// Mails_info is the information about the mails to send the feeds info to
	Mails_to   []string
	// Feed_info is the information about the feeds
	Feeds_info []_FeedInfo
}

// _FeedInfo is the information about a feed.
type _FeedInfo struct {
	// Feed_num is the number of the feed, beginning in 1 (no special reason, but could be useful some time)
	Feed_num int
	// Feed_url is the URL of the feed
	Feed_url string
	// Feed_type is the type of the feed (one of the TYPE_ constants)
	Feed_type string
	// Custom_msg_subject is the custom message subject
	Custom_msg_subject string
}

///////////////////////////////////////////////////////////////

type _MOD_6 struct {
	// Temp_locs is the locations to get the weather from
	Temp_locs   []string
	// News_locs is the locations to get the news from
	News_locs []OICNews.NewsLocs
}

///////////////////////////////////////////////////////////////

type _MOD_7 struct {
	// Model_loc is the location of the model file
	Model_loc string
	// Config_str is the LLM configuration string
	Config_str string
}

///////////////////////////////////////////////////////////////

type _MOD_10 struct {
	// Notifications is the list of notifications
	Notifications []_Notification
}

// _Notification is the format of a notification.
type _Notification struct {
	// Condition is the condition for the notification in Go language
	Condition string
	// Speak is the text to speak when the condition is met
	Speak     string
}

///////////////////////////////////////////////////////////////

type _MOD_12 struct {
	// Devices_info is the information about the devices
	Devices_info _DevicesInfo
	// Locs_info is the information about the locations
	Locs_info []_LocInfo
}

type _DevicesInfo struct {
	// AlwaysWith_device_id is the device id of the device that is always with the user
	AlwaysWith_device_id string
}

type _LocInfo struct {
	// Type is the type of the location "detector" (e.g. wifi)
	Type string
	// Name is the name of the detection (e.g. the wifi SSID)
	Name string
	// Address is the address of the detection (e.g. the wifi BSSID) in the format XX:XX:XX:XX:XX:XX
	Address string
	// Last_detection is the maximum amount of time in seconds without checking in which the device may still be in the
	// specified location
	Last_detection int64
	// Max_distance is the maximum distance in meters in which the device is in the specified location
	Max_distance int
	// Location is where the device is (e.g. "home")
	Location string
}
