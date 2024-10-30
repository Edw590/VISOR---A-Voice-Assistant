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

// Mod4GenInfo is the format of the custom generated information about this specific module.
type Mod4GenInfo struct {
	// Tasks_info maps the task ID to the last time the task was reminded in Unix minutes
	Notified_news map[int][]NewsInfo
}

// NewsInfo is the information about news.
type NewsInfo struct {
	Url   string
	Title string
}

///////////////////////////////////////////////////////////////////////////////

// Mod4UserInfo is the format of the custom information file about this specific module.
type Mod4UserInfo struct {
	// Mails_info is the information about the mails to send the feeds info to
	Mails_to   []string
	// Feed_info is the information about the feeds
	Feeds_info []FeedInfo
}

// FeedInfo is the information about a feed.
type FeedInfo struct {
	// Feed_num is the number of the feed, beginning in 1 (no special reason, but could be useful some time)
	Feed_num int
	// Feed_name is the user-given name of the feed
	Feed_name string
	// Feed_url is the URL of the feed
	Feed_url string
	// Feed_type is the type of the feed (one of the TYPE_ constants)
	Feed_type string
	// Custom_msg_subject is the custom message subject
	Custom_msg_subject string
}
