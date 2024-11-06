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

package SettingsSync

import (
	"Utils"
	"Utils/ModsFileInfo"
	"sort"
	"strconv"
)

/*
AddFeedRSS adds a feed to the user settings.

-----------------------------------------------------------

– Params:
  - enabled – whether the feed is enabled
  - name – the name of the feed
  - url – the URL of the feed
  - type_ – the type of the feed
  - custom_msg_subject – the custom message subject

– Returns:
  - the ID of the feed
 */
func AddFeedRSS(enabled bool, name string, url string, type_ string, custom_msg_subject string) int32 {
	var feeds_info *[]ModsFileInfo.FeedInfo = &Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	var id int32 = 1
	for i := 0; i < len(*feeds_info); i++ {
		if (*feeds_info)[i].Id == id {
			id++
			i = -1
		}
	}

	// Add the feed to the user settings
	*feeds_info = append(*feeds_info, ModsFileInfo.FeedInfo{
		Id:                 id,
		Enabled:            enabled,
		Name:               name,
		Url:                url,
		Type_:              type_,
		Custom_msg_subject: custom_msg_subject,
	})

	sort.SliceStable(*feeds_info, func(i, j int) bool {
		return (*feeds_info)[i].Name < (*feeds_info)[j].Name
	})

	return id
}

/*
RemoveFeedRSS removes a feed from the user settings.

-----------------------------------------------------------

– Params:
  - feed_id – the feed ID
 */
func RemoveFeedRSS(feed_id int32) {
	var feeds_info *[]ModsFileInfo.FeedInfo = &Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	for i := 0; i < len(*feeds_info); i++ {
		if (*feeds_info)[i].Id == feed_id {
			Utils.DelElemSLICES(feeds_info, i)

			break
		}
	}
}

/*
GetIdsListRSS returns a list of all feeds' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all feeds' IDs separated by "|"
 */
func GetIdsListRSS() string {
	var feeds_info *[]ModsFileInfo.FeedInfo = &Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	var ids_list string = ""
	for _, feed_info := range *feeds_info {
		ids_list += strconv.Itoa(int(feed_info.Id)) + "|"
	}
	ids_list = ids_list[:len(ids_list) - 1]

	return ids_list
}

/*
GetFeedRSS returns a feed by its ID.

-----------------------------------------------------------

– Params:
  - feed_id – the feed ID

– Returns:
  - the feed or nil if the feed was not found
 */
func GetFeedRSS(feed_id int32) *ModsFileInfo.FeedInfo {
	var feeds_info []ModsFileInfo.FeedInfo = Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	for i := 0; i < len(feeds_info); i++ {
		var feed_info *ModsFileInfo.FeedInfo = &feeds_info[i]
		if feed_info.Id == feed_id {
			return feed_info
		}
	}

	return nil
}
