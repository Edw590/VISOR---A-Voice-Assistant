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

package Screens

import (
	"Utils"
	"Utils/ModsFileInfo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"sort"
	"strconv"
)

func ModRSSFeedNotifier() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_RSS_FEED_NOTIFIER

	log.Println("RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR")

	return container.NewAppTabs(
		container.NewTabItem("Feeds list", rssFeedNotifierCreateFeedsListTab()),
		container.NewTabItem("Add feed", rssFeedNotifierCreateAddFeedTab()),
	)
}

func rssFeedNotifierCreateAddFeedTab() *container.Scroll {
	var feed_num int = 1
	for _, feed := range Utils.User_settings_GL.RSSFeedNotifier.Feeds_info {
		if feed.Feed_num == feed_num {
			feed_num++
		}
	}

	var label_id *widget.Label = widget.NewLabel("Feed ID: " + strconv.Itoa(feed_num))

	var check_enabled *widget.Check = widget.NewCheck("Feed enabled", nil)
	check_enabled.SetChecked(true)

	var entry_feed_name *widget.Entry = widget.NewEntry()
	entry_feed_name.PlaceHolder = "Feed name"

	var entry_feed_type *widget.Entry = widget.NewEntry()
	entry_feed_type.PlaceHolder = "Feed type"

	var entry_feed_url *widget.Entry = widget.NewEntry()
	entry_feed_url.PlaceHolder = "Feed URL or YouTube playlist/channel ID"

	var entry_custom_msg_subject *widget.Entry = widget.NewEntry()
	entry_custom_msg_subject.PlaceHolder = "Custom message subject (YT is automatic)"

	var button_save *widget.Button = widget.NewButton("Add", func() {
		Utils.User_settings_GL.RSSFeedNotifier.Feeds_info = append(Utils.User_settings_GL.RSSFeedNotifier.Feeds_info,
			ModsFileInfo.FeedInfo{
			Feed_num:           feed_num,
			Feed_enabled:       check_enabled.Checked,
			Feed_name:          entry_feed_name.Text,
			Feed_type:          entry_feed_type.Text,
			Feed_url:           entry_feed_url.Text,
			Custom_msg_subject: entry_custom_msg_subject.Text,
		})

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})

	return createMainContentScrollUTILS(
		label_id,
		check_enabled,
		entry_feed_name,
		entry_feed_type,
		entry_feed_url,
		entry_custom_msg_subject,
		button_save,
	)
}

func rssFeedNotifierCreateFeedsListTab() *container.Scroll {
	var objects []fyne.CanvasObject = nil
	var feeds_info []ModsFileInfo.FeedInfo = Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	sort.Slice(feeds_info, func(i, j int) bool {
		return feeds_info[i].Feed_num < feeds_info[j].Feed_num
	})
	for i := 0; i < len(feeds_info); i++ {
		objects = append(objects, createFeedInfoSetter(&feeds_info[i], i))
	}

	return createMainContentScrollUTILS(objects...)
}

func createFeedInfoSetter(feed_info *ModsFileInfo.FeedInfo, feed_idx int) *fyne.Container {
	var label_id *widget.Label = widget.NewLabel("Feed ID: " + strconv.Itoa(feed_info.Feed_num))

	var check_enabled *widget.Check = widget.NewCheck("Feed enabled", nil)
	check_enabled.SetChecked(feed_info.Feed_enabled)

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.SetText(feed_info.Feed_name)
	entry_name.PlaceHolder = "Feed name"

	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.SetText(feed_info.Feed_type)
	entry_type.PlaceHolder = "Feed type"

	var entry_url *widget.Entry = widget.NewEntry()
	entry_url.SetText(feed_info.Feed_url)
	entry_url.PlaceHolder = "Feed URL or YouTube playlist/channel ID"

	var entry_custom_msg_subject *widget.Entry = widget.NewEntry()
	entry_custom_msg_subject.SetText(feed_info.Custom_msg_subject)
	entry_custom_msg_subject.PlaceHolder = "Custom message subject (YT is automatic)"

	// Save button
	var button_save *widget.Button = widget.NewButton("Save", func() {
		feed_info.Feed_enabled = check_enabled.Checked
		feed_info.Feed_name = entry_name.Text
		feed_info.Feed_type = entry_type.Text
		feed_info.Feed_url = entry_url.Text
		feed_info.Custom_msg_subject = entry_custom_msg_subject.Text
	})

	var button_delete *widget.Button = widget.NewButton("Delete", func() {
		Utils.DelElemSLICES(&Utils.User_settings_GL.RSSFeedNotifier.Feeds_info, feed_idx)

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})

	return container.NewVBox(
		label_id,
		check_enabled,
		entry_name,
		entry_type,
		entry_url,
		entry_custom_msg_subject,
		button_save,
		button_delete,
	)
}
