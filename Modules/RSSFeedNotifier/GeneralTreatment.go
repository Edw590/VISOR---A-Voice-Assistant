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

package RSSFeedNotifier

import (
	"Utils/ModsFileInfo"
	"time"

	"github.com/mmcdole/gofeed"

	"Utils"
)

/*
generalTreatment does the general treatment of an RSS feed item.

-----------------------------------------------------------

– Params:
  - parsed_feed – the parsed feed
  - item_num – the number of the item to get
  - title_url_only – whether to only get the title and URL of the item through _NewsInfo (can be used for optimization)

– Returns:
  - the email info (without the Mail_to field) or all fields empty if title_url_only is true
  - the news info
 */
func generalTreatment(parsed_feed *gofeed.Feed, item_num int, title_url_only bool, custom_msg_subject string) (
					  Utils.EmailInfo, ModsFileInfo.NewsInfo) {
	var feed_item *gofeed.Item = parsed_feed.Items[item_num]

	var author string = ""
	if len(feed_item.Authors) > 0 {
		author = feed_item.Authors[0].Name
	}

	var things_replace = map[string]string{
		Utils.MODEL_RSS_ENTRY_TITLE_EMAIL:       feed_item.Title,
		Utils.MODEL_RSS_ENTRY_AUTHOR_EMAIL:      author,
		Utils.MODEL_RSS_ENTRY_DESCRIPTION_EMAIL: feed_item.Description,
		Utils.MODEL_RSS_ENTRY_URL_EMAIL:         feed_item.Link,
		Utils.MODEL_RSS_ENTRY_PUB_DATE_EMAIL:    feed_item.Published,
		Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL:    feed_item.Updated,
	}
	var newsInfo ModsFileInfo.NewsInfo = ModsFileInfo.NewsInfo{
		Title: things_replace[Utils.MODEL_RSS_ENTRY_TITLE_EMAIL],
		Url:   things_replace[Utils.MODEL_RSS_ENTRY_URL_EMAIL],
	}

	if title_url_only {
		return Utils.EmailInfo{}, newsInfo
	}

	if things_replace[Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL] != "" {
		if things_replace[Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL] == things_replace[Utils.MODEL_RSS_ENTRY_PUB_DATE_EMAIL] {
			things_replace[Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL] = "[new]"
		}
	}

	things_replace[Utils.MODEL_RSS_ENTRY_PUB_DATE_EMAIL] = convertDate(things_replace[Utils.MODEL_RSS_ENTRY_PUB_DATE_EMAIL])
	things_replace[Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL] = convertDate(things_replace[Utils.MODEL_RSS_ENTRY_UPD_DATE_EMAIL])

	var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_RSS, things_replace)
	email_info.Subject = custom_msg_subject

	return email_info, newsInfo
}

/*
convertDate converts a date from RFC3339 to DATE_FORMAT, also correcting the timezone to the local one.

-----------------------------------------------------------

– Params:
  - date – the date to convert

– Returns:
  - the converted date or the original date if it couldn't be converted
 */
func convertDate(date string) string {
	var date_time, err = time.Parse(time.RFC3339, date)
	if nil != err {
		return date
	}

	return date_time.Local().Format(Utils.DATE_TIME_FORMAT)
}
