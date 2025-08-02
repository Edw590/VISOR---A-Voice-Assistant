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

package RSSFeedNotifier

import (
	"Utils/ModsFileInfo"
	"context"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"Utils"
)

// RSS Feed Notifier //

// Cut the video title in more than 67 characters. After those 67 put ellipsis. This is what YT used to do.
// //////////////////////////////////////////////////
// (The emails are in PT-PT)
// Família Brisados acabou de carregar um vídeo
// 🔴 SuperHouseTV está agora em direto: [video title here]
// //////////////////////////////////////////////////

//////////////////////////
// Types of feeds:
var allowed_feed_types_1_GL []string = []string{
	_TYPE_1_GENERAL,
	_TYPE_1_YOUTUBE,
}
const (
	_TYPE_1_GENERAL = "General"
	_TYPE_1_YOUTUBE = "YouTube"
)
const (
	_TYPE_2_YT_CHANNEL  = "CH"
	_TYPE_2_YT_PLAYLIST = "PL"
)
const (
	_TYPE_3_YT_INC_SHORTS = "+S"
)
//////////////////////////

const _GEN_ERROR string = "3234_ERROR"

// _FeedType is the type of the feed. Each type is one of the TYPE_x constants, being x the number of the type.
type _FeedType struct {
	type_1 string
	type_2 string
	type_3 string
}

// _MAX_URLS_STORED is the maximum number of URLs stored in the file. This is to avoid having a file with too many URLs.
// 100 because it must be above the number of entries in all the feeds, and 100 is a big number (30 for StackExchange,
// 15 for YT - 100 seems perfect).
const _MAX_URLS_STORED int = 100

const _TIME_SLEEP_S int = 2*60

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	for {
		for _, feedInfo := range Utils.GetUserSettings(Utils.LOCK_UNLOCK).RSSFeedNotifier.Feeds_info {
			if !feedInfo.Enabled {
				continue
			}

			//if feedInfo.Id != 2 {
			//	continue
			//}
			//Utils.LogLnDebug("__________________________BEGINNING__________________________")

			var feedType _FeedType = getFeedType(feedInfo.Type_)

			if !Utils.ContainsSLICES(allowed_feed_types_1_GL, feedType.type_1) {
				//Utils.LogLnDebug("Feed type not allowed: " + feedType.type_1)
				//Utils.LogLnDebug("__________________________ENDING__________________________")

				continue
			}

			if feedType.type_1 == _TYPE_1_YOUTUBE {
				// If the feed is a YouTube feed, the feed URL is the channel or playlist ID, so we need to change it to
				// the correct URL.
				if feedType.type_2 == _TYPE_2_YT_CHANNEL {
					feedInfo.Url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feedInfo.Url
				} else if feedType.type_2 == _TYPE_2_YT_PLAYLIST {
					feedInfo.Url = "https://www.youtube.com/feeds/videos.xml?playlist_id=" + feedInfo.Url
				}
			}

			//Utils.LogLnDebug("feed_id: " + strconv.Itoa(int(feedInfo.Id)))
			//Utils.LogLnDebug("feed_url: " + feedInfo.Url)
			//Utils.LogLnDebug("feed_type: " + feedInfo.Type_)
			//Utils.LogLnDebug("feedType.type_1: " + feedType.type_1)
			//Utils.LogLnDebug("feedType.type_2: " + feedType.type_2)
			//Utils.LogLnDebug("feedType.type_3: " + feedType.type_3)

			var new_feed bool = true
			var newsInfo_list []ModsFileInfo.NewsInfo2 = nil
			for _, newsInfo := range getModGenSettings().Notified_news {
				if newsInfo.Id == feedInfo.Id {
					new_feed = false
					newsInfo_list = newsInfo.News_info

					break
				}
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			parsed_feed, err := gofeed.NewParser().ParseURLWithContext(feedInfo.Url, ctx)
			cancel()
			if nil != err {
				//Utils.LogLnDebug("Error parsing feed: " + err.Error())
				continue
			}

			for item_num, item := range parsed_feed.Items {

				var check_skipping_later bool = true

				// Check if the news is new, and if it's not, skip it. But only if it's not a YouTube playlist, because
				// those may need the order of the items reversed and so the ones got from this loop are wrong. Or if it
				// is playlist, then only if the feed item ordering is correct (no scraping needed).
				// This is also here and not just in the end to prevent useless item processing (optimized).
				if feedType.type_2 != _TYPE_2_YT_PLAYLIST || !scrapingNeeded(parsed_feed) {
					check_skipping_later = false
					if !isNewNews(newsInfo_list, item.Title, item.Link) {
						// If the news is not new, don't notify.
						continue
					}
				}

				var email_info Utils.EmailInfo
				var newsInfo ModsFileInfo.NewsInfo2

				switch feedType.type_1 {
					case _TYPE_1_YOUTUBE: {
						email_info, newsInfo = youTubeTreatment(feedType, parsed_feed, item_num, new_feed)
					}
					case _TYPE_1_GENERAL: {
						email_info, newsInfo = generalTreatment(parsed_feed, item_num, new_feed,
							feedInfo.Custom_msg_subject)
					}
					default: {
						//Utils.LogLnDebug("Unknown feed type_1: " + feedType.type_1)
						continue
					}
				}

				var ignore_video bool = email_info.Html == ""

				if newsInfo.Url == "" { // Some error occurred
					continue
				}

				if check_skipping_later && !isNewNews(newsInfo_list, newsInfo.Title, newsInfo.Url) {
					// If the news is not new, don't notify.
					continue
				}

				var error_notifying bool = false

				//Utils.LogLnDebug("New news: " + newsInfo.Title)
				if !new_feed && !ignore_video {
					// If the feed is a newly added one, don't send emails for ALL the items in the feed - which are
					// being treated for the first time.
					//Utils.LogLnDebug("Queuing email: " + email_info.Subject)
					error_notifying = !queueEmailAllRecps(email_info.Sender, email_info.Subject, email_info.Html)
				}

				if !error_notifying {
					//Utils.LogLnDebug("Adding news to list...")
					newsInfo_list = append(newsInfo_list, newsInfo)
					if len(newsInfo_list) > _MAX_URLS_STORED {
						newsInfo_list = newsInfo_list[1:]
					}
				}
			}

			var found bool = false
			for i, news_info := range getModGenSettings().Notified_news {
				if news_info.Id == feedInfo.Id {
					getModGenSettings().Notified_news[i].News_info = newsInfo_list
					found = true

					break
				}
			}
			if !found {
				getModGenSettings().Notified_news = append(getModGenSettings().Notified_news, ModsFileInfo.NewsInfo{
					Id: feedInfo.Id,
					News_info: newsInfo_list,
				})
			}

			//Utils.LogLnDebug("__________________________ENDING__________________________")
		}

		if Utils.WaitWithStopDATETIME(module_stop, _TIME_SLEEP_S) {
			return
		}
	}
}

/*
getFeedType gets the _FeedType information from _FeedInfo.Feed_type.

-----------------------------------------------------------

– Params:
  - feed_type – the _FeedInfo.Feed_type

– Returns:
  - the _FeedType information
*/
func getFeedType(feed_type string) _FeedType {
	var feed_type_split []string = strings.Split(feed_type, " ")
	var feed_type_split_len int = len(feed_type_split)
	var feedType _FeedType = _FeedType{}
	if feed_type_split_len >= 1 {
		feedType.type_1 = feed_type_split[0]
	}
	if feed_type_split_len >= 2 {
		feedType.type_2 = feed_type_split[1]
	}
	if feed_type_split_len >= 3 {
		feedType.type_3 = feed_type_split[2]
	}

	return feedType
}

/*
isNewNews checks if the news is new.

-----------------------------------------------------------

– Params:
  - newsInfo_list – the list of notified news
  - title – the title of the news
  - url – the URL of the news

– Returns:
  - true if the news is new, false otherwise
 */
func isNewNews(newsInfo_list []ModsFileInfo.NewsInfo2, title string, url string) bool {
	//Utils.LogLnDebug("-------------------------------------------")
	//Utils.LogLnDebug("Checking if news is new: " + title + " - " + url)
	for _, newsInfo := range newsInfo_list {
		//Utils.LogLnDebug("Checking news: " + newsInfo.Title + " - " + newsInfo.Url)
		if  newsInfo.Url == url && newsInfo.Title == title {
			return false
		}
	}

	//Utils.LogLnDebug("News is new ^^^^^")

	return true
}

/*
queueEmailAllRecps queues an email to be sent to all recipients.

-----------------------------------------------------------

– Params:
  - sender_name – the name of the sender
  - subject – the subject of the email
  - html – the HTML of the email

– Returns:
  - true if the email was queued successfully, false otherwise
 */
func queueEmailAllRecps(sender_name string, subject string, html string) bool {
	// This is to add the images to the email using CIDs instead of using URLs which could/can go down at any time.
	// Except most email clients don't support CIDs... So I'll leave this here in case the images stop working with
	// the URLs and then either this or embeded Base64 on the src attribute of the <img> tag or hosted in the server
	// or something.
	// Still, the CID way seems better than the Base64 one. With CIDs, only Gmail Notified Pro wasn't showing them.
	// With Base64, Gmail (web or app) wasn't showing them (don't remember about the notifier). But I didn't test in
	// Hotmail or others.
	//var multiparts []Utils.Multipart = nil
	//if youtube {
	//	// Add the YouTube images to the email instead of using URLs which could/can go down at any time.
	//	var files_add []string = []string{"transparent_pixel.png", "twitter_email_icon_grey.png",
	//		"youtube_email_icon_grey.png", "youtubelogo_60.png"}
	//	for _, file_add := range files_add {
	//		var multipart Utils.Multipart = Utils.Multipart{
	//			Content_type: "image/png",
	//			Content_transfer_encoding: "base64",
	//			Content_id: file_add,
	//		}
	//		data, _ := os.ReadFile(modStartInfo_GL.Dir.Add("", "yt_email_images/", file_add).
	//			GPathToStringConversion())
	//		multipart.Body = base64.StdEncoding.EncodeToString(data)
	//
	//		multiparts = append(multiparts, multipart)
	//	}
	//}

	// Write the HTML to a file in case debugging is needed.
	_ = modDirsInfo_GL.Temp.Add2(false, "last_html_queued.html").WriteTextFile(html, false)

	err := Utils.QueueEmailEMAIL(Utils.EmailInfo{
		Sender:  sender_name,
		Mail_to: Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.User_email_addr,
		Subject: subject,
		Html:    html,
		Multiparts: nil,
	})
	if nil != err {
		return false
	}

	return true
}

func getModGenSettings() *ModsFileInfo.Mod4GenInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_4
}
