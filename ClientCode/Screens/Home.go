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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

func Home() fyne.CanvasObject {
	Current_screen_GL = NUM_HOME

	return container.NewAppTabs(
		container.NewTabItem("Home", homeCreateHomeTab()),
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

	var communicator_checkbox *widget.Check = widget.NewCheck("Communicator connected", func(checked bool) {
	})

	var no_website_info_label *widget.Label = widget.NewLabel("")
	var domain_entry *widget.Entry = widget.NewEntry()
	domain_entry.SetPlaceHolder("Website domain or IP (example: localhost)")
	var password_entry *widget.Entry = widget.NewPasswordEntry()
	password_entry.SetPlaceHolder("Website password")
	var save_button *widget.Button = widget.NewButton("Save", func() {
		SettingsSync.SetWebsiteInfo(domain_entry.Text, password_entry.Text)
		domain_entry.SetText("")
		password_entry.SetText("")
	})

	go func() {
		for {
			if Current_screen_GL == NUM_HOME {
				communicator_checkbox.SetChecked(Utils.IsCommunicatorConnectedSERVER())

				if SettingsSync.IsWebsiteInfoEmpty() {
					domain_entry.Enable()
					password_entry.Enable()
					save_button.Enable()
					no_website_info_label.SetText("No website info exists. Please enter it to activate full functionality.")
				} else {
					domain_entry.Disable()
					password_entry.Disable()
					save_button.Disable()
					no_website_info_label.SetText("Website info exists")
				}
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(
		container.NewVBox(text),
		communicator_checkbox,
		no_website_info_label,
		domain_entry,
		password_entry,
		save_button,
	)
}
