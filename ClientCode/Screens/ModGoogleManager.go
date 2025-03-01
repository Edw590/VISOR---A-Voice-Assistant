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

package Screens

import (
	"GoogleManager"
	"Utils"
	"context"
	"encoding/json"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/oauth2"
	"net/url"
	"time"
)

type date struct {
	instruction *widget.Label
	dateChosen  *widget.Label
}

func (d *date) onSelected(t time.Time) {
	// use time object to set text on label with given format
	d.instruction.SetText("Date Selected:")
	d.dateChosen.SetText(t.Format("Mon 02 Jan 2006"))
}

func ModGoogleManager() fyne.CanvasObject {
	Current_screen_GL = ID_GOOGLE_MANAGER

	return container.NewAppTabs(
		container.NewTabItem("Settings", googleManagerCreateSettingsTab()),
	)
}

func googleManagerCreateSettingsTab() *container.Scroll {
	link, _ := url.Parse("https://console.cloud.google.com/projectcreate")
	var link_google *widget.Hyperlink = widget.NewHyperlink("Click here and watch the video below", link)

	link, _ = url.Parse("https://youtu.be/B2E82UPUnOY?si=TIHV5U1kxY5mCKsD&t=95")
	var link_video *widget.Hyperlink = widget.NewHyperlink("How to obtain the Google credentials JSON", link)

	var label_additional_info *widget.Label = widget.NewLabel("Activate the Calendar, Gmail and Tasks APIs by " +
		"looking them up in the Search bar and in the Scopes, choose \"auth/calendar\", \"auth/tasks\" and " +
		"\"auth/gmail.modify\".")
	label_additional_info.Wrapping = fyne.TextWrapWord

	var entry_credentials_json *widget.Entry = widget.NewEntry()
	entry_credentials_json.SetPlaceHolder("Google credentials JSON file contents")
	entry_credentials_json.SetText(Utils.User_settings_GL.GoogleManager.Credentials_JSON)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.User_settings_GL.GoogleManager.Credentials_JSON = entry_credentials_json.Text
	})
	btn_save.Importance = widget.SuccessImportance

	var label_additional_info2 *widget.Label = widget.NewLabel("To get the authorization code, when you get to an " +
		"error page (it's normal - Google stuff), look at the URL bar. Look for \"code=\" and copy what's after the " +
		"= sign until just before the next & sign.")
	label_additional_info2.Wrapping = fyne.TextWrapWord

	var btn_authorize *widget.Button = widget.NewButton("Authorize", func() {
		if Utils.User_settings_GL.GoogleManager.Credentials_JSON == "" {
			dialog.ShowError(errors.New("no credentials JSON saved"), Current_window_GL)
		}

		if !Utils.IsCommunicatorConnectedSERVER() {
			dialog.ShowError(errors.New("not connected to the server"), Current_window_GL)

			return
		}

		config, err := GoogleManager.ParseConfigJSON()
		if err != nil {
			dialog.ShowError(err, Current_window_GL)

			return
		}

		auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

		var entry_auth_code *widget.Entry = widget.NewEntry()
		dialog.ShowForm("Google authorization code", "Enter", "Cancel", []*widget.FormItem{
			{Text: "Code", Widget: entry_auth_code},
		}, func(ok bool) {
			if (!ok) || (entry_auth_code.Text == "") {
				return
			}

			token, err := config.Exchange(context.Background(), entry_auth_code.Text)
			if err != nil {
				dialog.ShowError(err, Current_window_GL)

				return
			}

			if !Utils.IsCommunicatorConnectedSERVER() {
				dialog.ShowError(errors.New("not connected to the server"), Current_window_GL)

				return
			}

			var message []byte = []byte("S_S|GManTok|")
			token_bytes, _ := json.Marshal(token)
			message = append(message, Utils.CompressString(string(token_bytes))...)
			Utils.QueueNoResponseMessageSERVER(message)

			dialog.ShowInformation("Information", "Authorization code saved. You're all set!", Current_window_GL)
		}, Current_window_GL)

		link, _ = url.Parse(auth_url)
		var link_authorize *widget.Hyperlink = widget.NewHyperlink("External authorization prompt", link)
		dialog.ShowCustom("Open the following Google link", "Close", link_authorize, Current_window_GL)
	})
	btn_authorize.Importance = widget.HighImportance

	return createMainContentScrollUTILS(
		link_google,
		link_video,
		label_additional_info,
		entry_credentials_json,
		btn_save,
		label_additional_info2,
		btn_authorize,
	)
}
