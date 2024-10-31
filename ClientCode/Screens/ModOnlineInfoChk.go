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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func ModOnlineInfoChk() fyne.CanvasObject {
	Current_screen_GL = ID_ONLINE_INFO_CHK

	return container.NewAppTabs(
		container.NewTabItem("Settings", onlineInfoChkCreateSettingsTab()),
	)
}

func onlineInfoChkCreateSettingsTab() *container.Scroll {
	var entry_weather_locs *widget.Entry = widget.NewMultiLineEntry()
	entry_weather_locs.SetPlaceHolder("The weather locationos to check, one per line")
	entry_weather_locs.SetMinRowsVisible(3)
	entry_weather_locs.SetText(strings.Join(Utils.User_settings_GL.OnlineInfoChk.Temp_locs, "\n"))

	var entry_news_locs *widget.Entry = widget.NewMultiLineEntry()
	entry_news_locs.SetPlaceHolder("The news locations to check, one per line")
	entry_news_locs.SetMinRowsVisible(3)
	entry_news_locs.SetText(strings.Join(Utils.User_settings_GL.OnlineInfoChk.News_locs, "\n"))

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.User_settings_GL.OnlineInfoChk.Temp_locs = strings.Split(entry_weather_locs.Text, "\n")
		Utils.User_settings_GL.OnlineInfoChk.News_locs = strings.Split(entry_news_locs.Text, "\n")
	})

	return createMainContentScrollUTILS(
		entry_weather_locs,
		entry_news_locs,
		btn_save,
	)
}