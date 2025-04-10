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

// Mod4GenInfo is the format of the custom generated information about this specific module.
type Mod4GenInfo struct {
	// Notified_news is the information about the notified news from the feeds
	Notified_news []NewsInfo
}

type NewsInfo struct {
	// Id is the ID of the feed
	Id   int32
	// News_info is the information about the news
	News_info []NewsInfo2
}

type NewsInfo2 struct {
	// Url is the URL of the news
	Url   string
	// Title is the title of the news
	Title string
}

///////////////////////////////////////////////////////////////////////////////

// Mod4UserInfo is the format of the custom information file about this specific module.
type Mod4UserInfo struct {
	// Feed_info is the information about the feeds
	Feeds_info []FeedInfo
}

// FeedInfo is the information about a feed.
type FeedInfo struct {
	// Id is the ID of the feed
	Id int32
	// Enabled is whether the feed is enabled
	Enabled bool
	// Name is the user-given name of the feed
	Name string
	// Url is the URL of the feed
	Url string
	// Type_ is the type of the feed
	Type_ string
	// Custom_msg_subject is the custom message subject
	Custom_msg_subject string
}
