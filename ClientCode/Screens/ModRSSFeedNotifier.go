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
	"sort"
	"strconv"
)

var mod_rss_feed_notifier_canvas_object_GL fyne.CanvasObject = nil

func ModRSSFeedNotifier() fyne.CanvasObject {
	var tabs *container.AppTabs = container.NewAppTabs(
		container.NewTabItem("Feeds list", rssFeedNotifierCreateFeedsListTab()),
	)

	mod_rss_feed_notifier_canvas_object_GL = tabs
	Current_screen_GL = mod_rss_feed_notifier_canvas_object_GL

	return mod_rss_feed_notifier_canvas_object_GL
}

func rssFeedNotifierCreateFeedsListTab() *container.Scroll {
	var objects []fyne.CanvasObject = nil
	var feeds_info []ModsFileInfo.FeedInfo = Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	sort.Slice(feeds_info, func(i, j int) bool {
		return feeds_info[i].Feed_num > feeds_info[j].Feed_num
	})
	for i := len(feeds_info) - 1; i >= 0; i-- {
		objects = append(objects, createFeedInfo(&feeds_info[i]))
	}

	return createMainContentScrollUTILS(objects...)
}

func createFeedInfo(feed_info *ModsFileInfo.FeedInfo) *fyne.Container {
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
	entry_custom_msg_subject.PlaceHolder = "Custom message subject"

	// Save button
	var button_save *widget.Button = widget.NewButton("Save", func() {
		feed_info.Feed_enabled = check_enabled.Checked
		feed_info.Feed_name = entry_name.Text
		feed_info.Feed_type = entry_type.Text
		feed_info.Feed_url = entry_url.Text
		feed_info.Custom_msg_subject = entry_custom_msg_subject.Text
	})

	return container.NewVBox(
		label_id,
		check_enabled,
		entry_name,
		entry_type,
		entry_url,
		entry_custom_msg_subject,
		button_save,
	)
}
