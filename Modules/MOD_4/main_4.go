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

package MOD_4

import (
	"Utils/ModsFileInfo"
	"context"
	"strconv"
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

// _NewsInfo is the information about news.
type _NewsInfo struct {
	Url   string
	Title string
}

// _MAX_URLS_STORED is the maximum number of URLs stored in the file. This is to avoid having a file with too many URLs.
// 100 because it must be above the number of entries in all the feeds, and 100 is a big number (30 for StackExchange,
// 15 for YT - 100 seems perfect).
const _MAX_URLS_STORED int = 100

const _TIME_SLEEP_S int = 2*60

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		var modUserInfo *ModsFileInfo.Mod4UserInfo = &Utils.User_settings_GL.MOD_4

		for {

			for _, feedInfo := range modUserInfo.Feeds_info {
				// if feedInfo.Feed_num != 8 {
				//	continue
				// }
				//log.Println("__________________________BEGINNING__________________________")

				var feedType _FeedType = getFeedType(feedInfo.Feed_type)

				if !Utils.ContainsSLICES(allowed_feed_types_1_GL, feedType.type_1) {
					//log.Println("Feed type not allowed: " + feedInfo.Feed_type)
					//log.Println("__________________________ENDING__________________________")

					continue
				}

				if feedType.type_1 == _TYPE_1_YOUTUBE {
					// If the feed is a YouTube feed, the feed URL is the channel or playlist ID, so we need to change it to
					// the correct URL.
					if feedType.type_2 == _TYPE_2_YT_CHANNEL {
						feedInfo.Feed_url = "https://www.youtube.com/feeds/videos.xml?channel_id=" + feedInfo.Feed_url
					} else if feedType.type_2 == _TYPE_2_YT_PLAYLIST {
						feedInfo.Feed_url = "https://www.youtube.com/feeds/videos.xml?playlist_id=" + feedInfo.Feed_url
					}
				}

				//log.Println("feed_num: " + strconv.Itoa(feedInfo.Feed_num))
				//log.Println("feed_url: " + feedInfo.Feed_url)
				//log.Println("feed_type: " + feedInfo.Feed_type)
				//log.Println("feedType.type_1: " + feedType.type_1)
				//log.Println("feedType.type_2: " + feedType.type_2)
				//log.Println("feedType.type_3: " + feedType.type_3)

				var notif_news_file_path Utils.GPath = moduleInfo_GL.ModDirsInfo.UserData.Add2(true, "notified_news",
					strconv.Itoa(feedInfo.Feed_num)+".json")
				var newsInfo_list []_NewsInfo = nil
				if notif_news_file_path.Exists() {
					var notified_news_json []byte = notif_news_file_path.ReadFile()
					Utils.FromJsonGENERAL(notified_news_json, &newsInfo_list)
				}

				var new_feed bool = false
				if len(newsInfo_list) == 0 {
					new_feed = true
				}

				ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
				parsed_feed, err := gofeed.NewParser().ParseURLWithContext(feedInfo.Feed_url, ctx)
				cancel()
				if nil != err {
					//log.Println("Error parsing feed: " + err.Error())
					continue
				}

				var notified_news_list_modified bool = false
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

					var email_info Utils.EmailInfo = Utils.EmailInfo{}
					var newsInfo _NewsInfo = _NewsInfo{}

					switch feedType.type_1 {
						case _TYPE_1_YOUTUBE: {
							email_info, newsInfo = youTubeTreatment(feedType, parsed_feed, item_num, new_feed)
						}
						case _TYPE_1_GENERAL: {
							email_info, newsInfo = generalTreatment(parsed_feed, item_num, new_feed,
								feedInfo.Custom_msg_subject)
						}
						default: {
							//log.Println("Unknown feed type_1: " + feedType.type_1)
							continue
						}
					}

					var ignore_video bool = "" == email_info.Html

					if "" == newsInfo.Url { // Some error occurred
						continue
					}

					if check_skipping_later && !isNewNews(newsInfo_list, newsInfo.Title, newsInfo.Url) {
						// If the news is not new, don't notify.
						continue
					}

					var error_notifying bool = false

					//log.Println("New news: " + newsInfo.Title)
					if !new_feed && !ignore_video {
						// If the feed is a newly added one, don't send emails for ALL the items in the feed - which are
						// being treated for the first time.
						//log.Println("Queuing email: " + email_info.Subject)
						error_notifying = !queueEmailAllRecps(email_info.Sender, email_info.Subject, email_info.Html,
							modUserInfo.Mails_to)
					}

					if !error_notifying {
						newsInfo_list = append(newsInfo_list, newsInfo)
						if len(newsInfo_list) > _MAX_URLS_STORED {
							newsInfo_list = newsInfo_list[1:]
						}
						notified_news_list_modified = true
					}
				}
				if notified_news_list_modified {
					_ = notif_news_file_path.WriteTextFile(*Utils.ToJsonGENERAL(newsInfo_list), false)
				}

				//log.Println("__________________________ENDING__________________________")
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				return
			}
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
func isNewNews(newsInfo_list []_NewsInfo, title string, url string) bool {
	//log.Println("Checking if news is new: " + title)
	for _, newsInfo := range newsInfo_list {
		if  newsInfo.Url == url && newsInfo.Title == title {
			return false
		}
	}

	//log.Println("News is new ^^^^^")

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
func queueEmailAllRecps(sender_name string, subject string, html string, mails_to []string) bool {
	for _, mail_to := range mails_to {
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
		_ = moduleInfo_GL.ModDirsInfo.Temp.Add2(false, "last_html_queued.html").WriteTextFile(html, false)

		err := Utils.QueueEmailEMAIL(Utils.EmailInfo{
			Sender:  sender_name,
			Mail_to: mail_to,
			Subject: subject,
			Html:    html,
			Multiparts: nil,
		})
		if nil != err {
			return false
		}
	}

	return true
}
