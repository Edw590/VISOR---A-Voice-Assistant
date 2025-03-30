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
	"SettingsSync/SettingsSync"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"net/url"
	"time"
)

func Home() fyne.CanvasObject {
	Current_screen_GL = ID_HOME

	return container.NewAppTabs(
		container.NewTabItem("Main", homeCreateHomeTab()),
		container.NewTabItem("Settings", homeCreateSettingsTab()),
		container.NewTabItem("Local settings", homeCreateLocalSettingsTab()),
	)
}

func homeCreateLocalSettingsTab() *container.Scroll {
	var entry_password *widget.Entry = widget.NewPasswordEntry()
	entry_password.SetPlaceHolder("Settings encryption password or empty to disable")
	entry_password.SetText(Utils.GetPasswordCREDENTIALS())

	var btn_save_temp *widget.Button = widget.NewButton("Save for this session", func() {
		Utils.Password_GL = entry_password.Text
		_ = Utils.DeletePasswordCREDENTIALS()
	})
	btn_save_temp.Importance = widget.SuccessImportance

	var btn_save_perm *widget.Button = widget.NewButton("Save permanently", func() {
		Utils.Password_GL = entry_password.Text
		if entry_password.Text == "" {
			_ = Utils.DeletePasswordCREDENTIALS()
		} else {
			_ = Utils.SavePasswordCREDENTIALS(entry_password.Text)
		}
	})
	btn_save_perm.Importance = widget.HighImportance

	var entry_device_id *widget.Entry = widget.NewEntry()
	entry_device_id.SetPlaceHolder("Unique device ID (for example \"MyComputer\")")
	entry_device_id.SetText(Utils.Gen_settings_GL.Device_settings.Id)

	var entry_device_type *widget.Entry = widget.NewEntry()
	entry_device_type.SetPlaceHolder("Device type (for example \"computer\")")
	entry_device_type.SetText(Utils.Gen_settings_GL.Device_settings.Type_)

	var entry_device_description *widget.Entry = widget.NewEntry()
	entry_device_description.SetPlaceHolder("Device description (for example the model, \"Legion Y520\")")
	entry_device_description.SetText(Utils.Gen_settings_GL.Device_settings.Description)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.Gen_settings_GL.Device_settings.Id = entry_device_id.Text
		Utils.Gen_settings_GL.Device_settings.Type_ = entry_device_type.Text
		Utils.Gen_settings_GL.Device_settings.Description = entry_device_description.Text
	})
	btn_save.Importance = widget.SuccessImportance

	return createMainContentScrollUTILS(
		entry_password,
		container.New(layout.NewGridLayout(2), btn_save_temp, btn_save_perm),
		entry_device_id,
		entry_device_type,
		entry_device_description,
		btn_save,
		createValuesChooserAccordionUTILS("General - "),
	)
}

func homeCreateSettingsTab() *container.Scroll {
	var entry_pin *widget.Entry = widget.NewPasswordEntry()
	entry_pin.SetPlaceHolder("App protection PIN (any number of digits or empty to disable)")
	entry_pin.SetText(Utils.User_settings_GL.General.Pin)
	entry_pin.Validator = validation.NewRegexp(`^(\d+)?$`, "PIN must be numberic")

	var entry_visor_email_addr *widget.Entry = widget.NewEntry()
	entry_visor_email_addr.SetPlaceHolder("V.I.S.O.R. email address")
	entry_visor_email_addr.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"Invalid email address")
	entry_visor_email_addr.SetText(Utils.User_settings_GL.General.VISOR_email_addr)

	var entry_visor_email_pw *widget.Entry = widget.NewPasswordEntry()
	entry_visor_email_pw.SetPlaceHolder("V.I.S.O.R. email password (2FA password if enabled)")
	entry_visor_email_pw.SetText(Utils.User_settings_GL.General.VISOR_email_pw)

	var entry_user_email_addr *widget.Entry = widget.NewEntry()
	entry_user_email_addr.SetPlaceHolder("User email address (used for all email communication)")
	entry_user_email_addr.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		"Invalid email address")
	entry_user_email_addr.SetText(Utils.User_settings_GL.General.User_email_addr)

	var entry_server_domain *widget.Entry = widget.NewEntry()
	entry_server_domain.SetPlaceHolder("Server domain or IP (example: localhost)")
	entry_server_domain.SetText(Utils.User_settings_GL.General.Website_domain)

	var entry_server_pw *widget.Entry = widget.NewPasswordEntry()
	entry_server_pw.SetPlaceHolder("Server password (strong letters and numbers password)")
	entry_server_pw.SetText(Utils.User_settings_GL.General.Website_pw)

	var entry_wolframalpha_appid *widget.Entry = widget.NewEntry()
	entry_wolframalpha_appid.SetPlaceHolder("WolframAlpha App ID")
	entry_wolframalpha_appid.SetText(Utils.User_settings_GL.General.WolframAlpha_AppID)

	var entry_picovoice_api_key *widget.Entry = widget.NewEntry()
	entry_picovoice_api_key.SetPlaceHolder("Picovoice API key")
	entry_picovoice_api_key.SetText(Utils.User_settings_GL.General.Picovoice_API_key)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.User_settings_GL.General.Pin = entry_pin.Text
		Utils.User_settings_GL.General.VISOR_email_addr = entry_visor_email_addr.Text
		Utils.User_settings_GL.General.VISOR_email_pw = entry_visor_email_pw.Text
		Utils.User_settings_GL.General.User_email_addr = entry_user_email_addr.Text
		Utils.User_settings_GL.General.Website_domain = entry_server_domain.Text
		Utils.User_settings_GL.General.Website_pw = entry_server_pw.Text
		Utils.User_settings_GL.General.WolframAlpha_AppID = entry_wolframalpha_appid.Text
		Utils.User_settings_GL.General.Picovoice_API_key = entry_picovoice_api_key.Text
	})
	btn_save.Importance = widget.SuccessImportance

	link, _ := url.Parse("https://developer.wolframalpha.com/")
	var link_wolframalpha *widget.Hyperlink = widget.NewHyperlink("Get WolframAlpha App ID from here", link)

	link, _ = url.Parse("https://console.picovoice.ai/")
	var link_picovoice *widget.Hyperlink = widget.NewHyperlink("Get Picovoice API key from here", link)

	return createMainContentScrollUTILS(
		entry_pin,
		entry_visor_email_addr,
		entry_visor_email_pw,
		entry_user_email_addr,
		entry_server_domain,
		entry_server_pw,
		entry_wolframalpha_appid,
		entry_picovoice_api_key,
		btn_save,
		link_wolframalpha,
		link_picovoice,
	)
}

func homeCreateHomeTab() *container.Scroll {
	var text *canvas.Text = canvas.NewText("V.I.S.O.R. Systems", color.RGBA{
		R: 34,
		G: 177,
		B: 76,
		A: 255,
	})
	text.TextSize = 40
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle.Bold = true

	var communicator_checkbox *widget.Check = widget.NewCheck("Connected to the server", nil)

	var no_website_info_label *widget.Label = widget.NewLabel("")

	go func() {
		for {
			if Current_screen_GL == ID_HOME {
				communicator_checkbox.SetChecked(Utils.IsCommunicatorConnectedSERVER())

				if SettingsSync.IsWebsiteInfoEmpty() {
					no_website_info_label.SetText("No server info exists. Enter it to activate full functionality.")
				} else {
					no_website_info_label.SetText("Server info exists")
				}
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(
		text,
		communicator_checkbox,
		no_website_info_label,
	)
}
