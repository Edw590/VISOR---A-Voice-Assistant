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
	"encoding/json"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"Utils"
)

type _PlaylistPage struct {
	id               string
	videos_info_json []string
}

type _VideoInfo struct {
	id     string
	title  string
	length string
	image  string
}

const _YT_TIME_DATE_FORMAT string = "2006-01-02T15:04:05-07:00"

var playlistPage_GL _PlaylistPage = _PlaylistPage{}

/*
ytPlaylistScraping scrapes the YT playlist page to get the video information and reads the video list backwards to get
the latest videos - because this function is to be used *only* if scrapingNeeded() returns true.

-----------------------------------------------------------

– Params:
  - playlist_id – the ID of the playlist
  - item_num – the number of the item to get (0 is the last, 1 is the second to last, etc.)
  - item_count – the number of items in the playlist RSS feed

– Returns:
  - the video info or all fields with _GEN_ERROR on them if any error occurs (check that with the video ID)
*/
func ytPlaylistScraping(playlist_id string, item_num int, item_count int) _VideoInfo {
	var videoInfo _VideoInfo = _VideoInfo{
		id:     _GEN_ERROR,
		title:  _GEN_ERROR,
		length: _GEN_ERROR,
		image:  _GEN_ERROR,
	}

	var videos_info_json []string = playlistPage_GL.videos_info_json
	// This is here to make sure the page is only scraped once
	if playlistPage_GL.id != playlist_id {
		var playlist_url string = "https://www.youtube.com/playlist?list=" + playlist_id
		var page_html *string = Utils.GetPageHtmlWEBPAGES(playlist_url)
		if page_html == nil {
			playlistPage_GL.id = ""

			return videoInfo
		}

		videos_info_json = strings.Split(*page_html, "{\"playlistVideoRenderer\":")

		// Separate and remove the rest of the page from the last item
		var last_json string = videos_info_json[len(videos_info_json)-1]
		videos_info_json[len(videos_info_json)-1] = strings.Split(last_json, "],\"playlistId\":")[0]

		// Remove the first too (it's also the rest of the page)
		videos_info_json = videos_info_json[1:]
		for i := 0; i < len(videos_info_json); i++ {
			// Remove the last "}," from the string (it's part of the main JSON object)
			videos_info_json[i] = videos_info_json[i][:strings.LastIndex(videos_info_json[i], "}")]
		}

		playlistPage_GL.id = playlist_id
		playlistPage_GL.videos_info_json = videos_info_json
	}

	var index int = len(videos_info_json) - item_count + item_num
	if index < 0 {
		// This should never happen - but it has, somehow xD (len was 0, item_num 0 and item_count 15...). So here is
		// the prevention.
		return videoInfo
	}
	var video_info_json string = videos_info_json[index]
	var json_decoded any = nil
	err := json.NewDecoder(strings.NewReader(video_info_json)).Decode(&json_decoded)
	if nil != err {
		return videoInfo
	}

	// Video ID
	var val, ok = toMap(json_decoded)["videoId"]
	if ok {
		videoInfo.id = val.(string)

		// toMap(json_decoded)["videoId"].(string)
	}

	// Video title
	val, ok = toMap(json_decoded)["title"]
	if ok {
		val, ok = toMap(val)["runs"]
		if ok {
			if len(toArr(val)) > 0 {
				val, ok = toMap(toArr(val)[0])["text"]
				if ok {
					videoInfo.title = val.(string)

					// toMap(toArr(toMap(toMap(json_decoded)["title"])["runs"])[0])["text"].(string)
				}
			}
		}
	}

	// Video length
	val, ok = toMap(json_decoded)["lengthSeconds"]
	if ok {
		videoInfo.length = SecondsToTimeStr(val.(string))

		// SecondsToTimeStr(toStr(toMap(json_decoded)["lengthSeconds"]))
	}

	// Video thumbnail
	val, ok = toMap(json_decoded)["thumbnail"]
	if ok {
		val, ok = toMap(val)["thumbnails"]
		if ok {
			var array []any = toArr(val)

			// toArr(toMap(toMap(json_decoded)["thumbnail"])["thumbnails"])

			if len(array) > 0 {
				// The last element is the highest quality thumbnail
				val, ok = toMap(array[len(array)-1])["url"]
				if ok {
					videoInfo.image = val.(string)

					// toMap(array[len(array)-1])["url"].(string)
				}
			}
		}
	}

	return videoInfo
}

/*
scrapingNeeded checks if the playlist visual page needs to be scraped to get the videos information or not.

Some playliists have the videos listed ascending order, others in descending order. If they're listed in descending
order, the RSS feed page which contains always the "last" 15 videos, will have the actual lastest videos. On the other
hand, if the videos are listed in ascending order, the RSS feed page will have the oldest videos instead and will never
update them --> scraping the visual playlist page for the videos is needed.

-----------------------------------------------------------

– Params:
  - parsed_feed – the parsed feed

– Returns:
  - true if scraping is needed, false otherwise
 */
func scrapingNeeded(parsed_feed *gofeed.Feed) bool {
	var num_items int = len(parsed_feed.Items)
	// 15 is the maximum number of items in the YouTube RSS feed page. So if there are 15 ("or more" - just to not put
	// == which is too strict... just a precaution) items, we need to scrape the playlist page to check if there are
	// actually more than 15 items.
	if num_items >= 2 && num_items >= 15 {
		first_date, err1 := time.Parse(_YT_TIME_DATE_FORMAT, parsed_feed.Items[0].Published)
		last_date, err2 := time.Parse(_YT_TIME_DATE_FORMAT, parsed_feed.Items[num_items-1].Published)
		if err1 == nil && err2 == nil {
			// If the first date is before the last date, the playlist is in ascending order - need to scrape.
			return first_date.Before(last_date)
		} else {
			// Default to true in this case
			return true
		}
	}

	return false
}

func toMap(m any) map[string]any {
	return m.(map[string]any)
}

func toArr(m any) []any {
	return m.([]any)
}
