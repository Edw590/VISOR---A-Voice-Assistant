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

package MOD_4

import (
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/exp/slices"

	"Utils"
)

const _VID_TIME_DEF string = "--:--"
const _VID_TIME_LIVE string = "00:00" // Live videos have 00:00 as duration
const _VID_TITLE_MAX_LEN int = 67
// The max length of the video description on the email preview (YouTube used to trim after 27 chars)
const _VID_DESC_MAX_LEN int = _VID_TITLE_MAX_LEN // Better with 67 chars. 27 is too little.

/*
youTubeTreatment processes the YouTube feed.

-----------------------------------------------------------

– Params:
  - feedType – the type of the feed
  - parsed_feed – the parsed feed
  - item_num – the number of the current item in the feed
  - title_url_only – whether to only get the title and URL of the item through _NewsInfo (can be used for optimization)

– Returns:
  - the email info (without the Mail_to field)
  - the news info (useful especially if it's a playlist and it had to be scraped and the item order reversed internally)
All EmailInfo fields are empty if an error occurs, if the video is to be ignored (like if it's a Short), or if
title_url_only is true. In the 1st case, the _NewsInfo fields are also empty. In the 2nd case, the _NewsInfo fields are
still filled with the video info. To check for errors, check if the video URL is empty on NewsInfo (that one must always
have a value).
*/
func youTubeTreatment(feedType _FeedType, parsed_feed *gofeed.Feed, item_num int, title_url_only bool) (Utils.EmailInfo,
			_NewsInfo) {
	const (
		VIDEO_COLOR string = "#212121" // Default video color (sort of black)
		LIVE_COLOR  string = "#E62117" // Default live color (sort of red)
	)

	var things_replace = map[string]string{
		Utils.MODEL_YT_VIDEO_HTML_TITLE_EMAIL:        _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL:      parsed_feed.Authors[0].Name,
		Utils.MODEL_YT_VIDEO_CHANNEL_CODE_EMAIL:      parsed_feed.Items[0].Extensions["yt"]["channelId"][0].Value,
		Utils.MODEL_YT_VIDEO_CHANNEL_IMAGE_EMAIL:     _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL:       _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL: _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_CODE_EMAIL:        _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_IMAGE_EMAIL:       _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL:        _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_VIDEO_TIME_COLOR_EMAIL:  VIDEO_COLOR,
		Utils.MODEL_YT_VIDEO_PLAYLIST_CODE_EMAIL:     "", // Leave empty if it's not playlist
		Utils.MODEL_YT_VIDEO_SUBSCRIPTION_LINK_EMAIL: _GEN_ERROR,
		Utils.MODEL_YT_VIDEO_SUBSCRIPTION_NAME_EMAIL: parsed_feed.Title,
	}
	if !title_url_only {
		things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_IMAGE_EMAIL] = getChannelImageUrl(things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_CODE_EMAIL])
	}

	if feedType.type_2 == _TYPE_2_YT_CHANNEL {
		// The last part is what YouTube used to put in the URLs (taken from the original model)
		things_replace[Utils.MODEL_YT_VIDEO_SUBSCRIPTION_LINK_EMAIL] = "channel/" + things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_CODE_EMAIL] + "%3Ffeature%3Dem-uploademail"
	} else if feedType.type_2 == _TYPE_2_YT_PLAYLIST {
		things_replace[Utils.MODEL_YT_VIDEO_PLAYLIST_CODE_EMAIL] = parsed_feed.Extensions["yt"]["playlistId"][0].Value
		things_replace[Utils.MODEL_YT_VIDEO_SUBSCRIPTION_LINK_EMAIL] = "playlist?list=" + things_replace[Utils.MODEL_YT_VIDEO_PLAYLIST_CODE_EMAIL]
	}

	if feedType.type_2 == _TYPE_2_YT_PLAYLIST && scrapingNeeded(parsed_feed) {
		// Scraping is only needed for video information. The feed has the rest.
		// For scraping we only use the number of the item to guide through the video array. The rest comes from the
		// playlist page.
		var video_info _VideoInfo = ytPlaylistScraping(things_replace[Utils.MODEL_YT_VIDEO_PLAYLIST_CODE_EMAIL], item_num, len(parsed_feed.Items))
		if video_info.id == _GEN_ERROR {
			return Utils.EmailInfo{}, _NewsInfo{}
		}

		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL] = video_info.title
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_CODE_EMAIL] = video_info.id
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_IMAGE_EMAIL] = video_info.image
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL] = video_info.length

		// No way to get the description from the playlist visual page unless the video appears on the RSS feed.
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL] = _GEN_ERROR
		for _, item := range parsed_feed.Items {
			if item.Extensions["yt"]["videoId"][0].Value == video_info.id {
				things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL] = item.Extensions["media"]["group"][0].Children["description"][0].Value

				break
			}
		}
	} else {
		var feed_item *gofeed.Item = parsed_feed.Items[item_num]
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL] = feed_item.Title
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_CODE_EMAIL] = feed_item.Extensions["yt"]["videoId"][0].Value
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_IMAGE_EMAIL] = feed_item.Extensions["media"]["group"][0].Children["thumbnail"][0].Attrs["url"]
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL] = feed_item.Extensions["media"]["group"][0].Children["description"][0].Value
		if !title_url_only {
			things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL] = getVideoDuration(feed_item.Link)
		}
	}

	var is_short bool = isShort([]string{things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL], things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL]}, things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL])

	// If it's not to include Shorts and the video is a Short, return only the news info (to ignore the notification but
	// memorize that the video is to be ignored).
	if (feedType.type_3 != _TYPE_3_YT_INC_SHORTS && is_short) || title_url_only {
		return Utils.EmailInfo{}, _NewsInfo{
			Title: things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL],
			Url:   "https://www.youtube.com/watch?v=" + things_replace[Utils.MODEL_YT_VIDEO_VIDEO_CODE_EMAIL],
		}
	}

	var vid_title string = things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL]
	var vid_title_original string = vid_title
	if len(vid_title) > _VID_TITLE_MAX_LEN {
		vid_title = vid_title[:_VID_TITLE_MAX_LEN] + "..."
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL] = vid_title
	}

	var video_short string = ""
	if is_short {
		video_short = "Short"
	} else {
		video_short = "vídeo"
	}

	var msg_subject string = _GEN_ERROR
	if feedType.type_2 == _TYPE_2_YT_CHANNEL {
		if things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL] == _VID_TIME_LIVE {
			// Live video
			msg_subject = "🔴 " + things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL] + " está agora em direto: " + vid_title + "!"
			things_replace[Utils.MODEL_YT_VIDEO_HTML_TITLE_EMAIL] = "Em direto no YouTube: " + things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL] + " – " + vid_title + "!"

			// Change the length rectangle
			things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_COLOR_EMAIL] = LIVE_COLOR
			things_replace[Utils.MODEL_YT_VIDEO_VIDEO_TIME_EMAIL] = "LIVE" // Change the video length to "LIVE"
		} else {
			// Normal video
			msg_subject = things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL] + " acabou de carregar um " + video_short
			things_replace[Utils.MODEL_YT_VIDEO_HTML_TITLE_EMAIL] = msg_subject
		}
	} else if feedType.type_2 == _TYPE_2_YT_PLAYLIST {
		// Playlist video
		msg_subject = things_replace[Utils.MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL] + " acabou de adicionar um " + video_short + " a " + parsed_feed.Title
		things_replace[Utils.MODEL_YT_VIDEO_HTML_TITLE_EMAIL] = msg_subject
	}

	if len(things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL]) > _VID_DESC_MAX_LEN {
		things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL] = things_replace[Utils.MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL][:_VID_DESC_MAX_LEN] + "..."
	}

	var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_YT_VIDEO, things_replace)
	email_info.Subject = msg_subject

	return email_info,
	_NewsInfo{
		Title: vid_title_original,
		Url:   "https://www.youtube.com/watch?v=" + things_replace[Utils.MODEL_YT_VIDEO_VIDEO_CODE_EMAIL],
	}
}

/*
getVideoDuration gets the duration of the video by getting the video's page and looking for the duration (scraping).

The format returned is the same as the one from the SecondsToTimeStr() function.

-----------------------------------------------------------

– Params:
  - video_url – the URL of the video

– Returns:
  - the duration of the video if it was found, _VID_TIME_DEF otherwise
*/
func getVideoDuration(video_url string) string {
	var p_page_html *string = Utils.GetPageHtmlWEBPAGES(video_url)
	if p_page_html == nil {
		return _VID_TIME_DEF
	}
	var page_html string = *p_page_html

	// I think the data is in JSON, so I got the lengthSeconds that I found randomly looking for the seconds. It also a
	// double quote after the number ("lengthSeconds":"47" for 47 seconds) --> CAN CHANGE (checked on 2023-07-04).
	text_to_find := "\"lengthSeconds\":\""
	idx_begin := strings.Index(page_html, text_to_find) + len(text_to_find)
	idx_end := strings.Index(page_html[idx_begin:], "\"")
	if idx_begin > 0 && idx_end > 0 {
		return SecondsToTimeStr(page_html[idx_begin : idx_begin+idx_end])
	}

	return _VID_TIME_DEF
}

/*
SecondsToTimeStr converts the seconds to a time string for the video duration.

The format returned is "HH:MM:SS", but if the video is less than an hour long, the hours are removed ("MM:SS").

-----------------------------------------------------------

– Params:
  - seconds_str – the number of seconds as a string

– Returns:
  - the time string
 */
func SecondsToTimeStr(seconds_str string) string {
	var seconds, _ = strconv.Atoi(seconds_str)
	// Note: the location here is useless - I need a duration, not a date. So I chose UTC because yes.
	var length_seconds_time = time.Date(0, 0, 0, 0, 0, seconds, 0, time.UTC)
	var time_str string = length_seconds_time.Format("15:04:05")
	if strings.HasPrefix(time_str, "00:") {
		// Remove the hours if the video is less than an hour long
		time_str = time_str[3:]
	}

	return time_str
}

/*
getChannelImageUrl gets the URL of the channel image of by getting the channel's page and looking for the image
(scraping).

-----------------------------------------------------------

– Params:
  - channel_code – the code of the channel

– Returns:
  - the URL of the channel image if it was found, _GEN_ERROR otherwise
*/
func getChannelImageUrl(channel_code string) string {
	var p_page_html *string = Utils.GetPageHtmlWEBPAGES("https://www.youtube.com/channel/" + channel_code)
	if p_page_html == nil {
		return _GEN_ERROR
	}
	var page_html string = *p_page_html

	// The image URL is on the 3rd occurrence of the "https://yt3.googleusercontent.com/" on HTML of the channel's page
	// --> CAN CHANGE (checked on 2023-07-04).
	// The 1st and 2nd occurrences are the user's image and the channel's background image, respectively.
	var text_to_find string = "https://yt3.googleusercontent.com/"
	var idxs_begin []int = Utils.FindAllIndexesGENERAL(page_html, text_to_find)
	if len(idxs_begin) >= 3 {
		var idx_begin int = idxs_begin[2]
		var idx_end int = strings.Index(page_html[idx_begin:], "\"")

		return page_html[idx_begin : idx_begin+idx_end]
	}

	return _GEN_ERROR
}

/*
isShort checks if the video is a Short.

-----------------------------------------------------------

– Params:
  - video_texts – the texts of the video like title and description
  - video_len – the length of the video from getVideoDuration()

– Returns:
  - true if the video is a short, false otherwise (also false if video_len is _VID_TIME_DEF)
 */
func isShort(video_texts []string, video_len string) bool {
	// If any of the video texts has the #short or #shorts tag, mark as Short.
	for _, video_text := range video_texts {
		video_text_words := strings.Split(strings.ToLower(video_text), " ")
		if slices.Contains(video_text_words, "#short") || slices.Contains(video_text_words, "#shorts") {
			return true
		}
	}

	if video_len == _VID_TIME_DEF {
		// If the video length was not found, mark as not Short (can't know, better safe than sorry).
		return false
	}

	// Lastly, if none of the others worked (a video can be a Short and not have the tags), if the video is 1 minute or
	// less long, mark it as Short.
	var length_seconds = 0
	if len(Utils.FindAllIndexesGENERAL(video_len, ":")) == 1 {
		length_parsed, _ := time.Parse("04:05", video_len)
		length_seconds = length_parsed.Minute()*60 + length_parsed.Second()
	} else {
		length_parsed, _ := time.Parse("15:04:05", video_len)
		length_seconds = length_parsed.Hour()*60*60 + length_parsed.Minute()*60 + length_parsed.Second()
	}

	return length_seconds <= 60
}
