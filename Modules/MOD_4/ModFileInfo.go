/*******************************************************************************
 * Copyright 2023-2023 Edw590
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

package MOD_4

// _ModUserInfo is the format of the custom information file about this specific module.
type _ModUserInfo struct {
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
