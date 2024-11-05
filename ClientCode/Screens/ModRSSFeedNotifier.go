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
	"SettingsSync/SettingsSync"
	"Utils"
	"Utils/ModsFileInfo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ModRSSFeedNotifier() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_RSS_FEED_NOTIFIER

	return container.NewAppTabs(
		container.NewTabItem("Feeds list", rssFeedNotifierCreateFeedsListTab()),
		container.NewTabItem("Add feed", rssFeedNotifierCreateAddFeedTab()),
		container.NewTabItem("About", rssFeedNotifierCreateAboutTab()),
	)
}

func rssFeedNotifierCreateAboutTab() *container.Scroll {
	var label_info *widget.Label = widget.NewLabel(RSS_ABOUT)
	label_info.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(label_info)
}

func rssFeedNotifierCreateAddFeedTab() *container.Scroll {
	var check_enabled *widget.Check = widget.NewCheck("Feed enabled", nil)
	check_enabled.SetChecked(true)

	var entry_feed_name *widget.Entry = widget.NewEntry()
	entry_feed_name.SetPlaceHolder("Feed name (just for identification)")

	var entry_feed_type *widget.Entry = widget.NewEntry()
	entry_feed_type.SetPlaceHolder("Feed type (\"General\" or \"YouTube [CH|PL] [+S]\")")

	var entry_feed_url *widget.Entry = widget.NewEntry()
	entry_feed_url.SetPlaceHolder("Feed URL or YouTube playlist/channel ID")

	var entry_custom_msg_subject *widget.Entry = widget.NewEntry()
	entry_custom_msg_subject.SetPlaceHolder("Custom message subject (for YT it's automatic)")

	var btn_add *widget.Button = widget.NewButton("Add", func() {
		SettingsSync.AddFeedRSS(check_enabled.Checked, entry_feed_name.Text, entry_feed_url.Text, entry_feed_type.Text,
			entry_custom_msg_subject.Text)

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})

	return createMainContentScrollUTILS(
		check_enabled,
		entry_feed_name,
		entry_feed_type,
		entry_feed_url,
		entry_custom_msg_subject,
		btn_add,
	)
}

func rssFeedNotifierCreateFeedsListTab() *container.Scroll {
	var accordion *widget.Accordion = widget.NewAccordion()
	accordion.MultiOpen = true
	var feeds_info []ModsFileInfo.FeedInfo = Utils.User_settings_GL.RSSFeedNotifier.Feeds_info
	for i := 0; i < len(feeds_info); i++ {
		var feed_info ModsFileInfo.FeedInfo = feeds_info[i]
		var title string = ""
		if !feed_info.Enabled {
			title += "[X] "
		}
		title += feed_info.Name
		accordion.Append(widget.NewAccordionItem(trimAccordionTitleUTILS(title), createFeedInfoSetter(&feeds_info[i])))
	}

	return createMainContentScrollUTILS(accordion)
}

func createFeedInfoSetter(feed_info *ModsFileInfo.FeedInfo) *fyne.Container {
	var check_enabled *widget.Check = widget.NewCheck("Feed enabled", nil)
	check_enabled.SetChecked(feed_info.Enabled)

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.SetText(feed_info.Name)
	entry_name.SetPlaceHolder("Feed name (just for identification)")

	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.SetText(feed_info.Type_)
	entry_type.SetPlaceHolder("Feed type (\"General\" or \"YouTube [CH|PL] [+S]\")")

	var entry_url *widget.Entry = widget.NewEntry()
	entry_url.SetText(feed_info.Url)
	entry_url.SetPlaceHolder("Feed URL or YouTube playlist/channel ID")

	var entry_custom_msg_subject *widget.Entry = widget.NewEntry()
	entry_custom_msg_subject.SetText(feed_info.Custom_msg_subject)
	entry_custom_msg_subject.SetPlaceHolder("Custom message subject (for YT it's automatic)")

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		feed_info.Enabled = check_enabled.Checked
		feed_info.Name = entry_name.Text
		feed_info.Type_ = entry_type.Text
		feed_info.Url = entry_url.Text
		feed_info.Custom_msg_subject = entry_custom_msg_subject.Text

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})
	btn_save.Importance = widget.SuccessImportance

	var btn_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationUTILS("Are you sure you want to delete this feed?", func(confirmed bool) {
			if confirmed {
				SettingsSync.RemoveFeedRSS(feed_info.Id)

				Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

	return container.NewVBox(
		check_enabled,
		entry_name,
		entry_type,
		entry_url,
		entry_custom_msg_subject,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
	)
}
