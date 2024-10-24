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

// Current_screen_GL is the current app screen. It's currently used to let threads specific to each screen know if they
// should continue processing data or not (they don't stop, they just keep waiting for the screen to become active again).
var Current_screen_GL fyne.CanvasObject = nil

var screens_size_GL fyne.Size = fyne.NewSize(550, 480)

var home_canvas_object_GL fyne.CanvasObject = nil

func Home() fyne.CanvasObject {
	Current_screen_GL = home_canvas_object_GL
	if home_canvas_object_GL != nil {
		return home_canvas_object_GL
	}

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
			if Current_screen_GL == home_canvas_object_GL {
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
			}

			time.Sleep(1 * time.Second)
		}
	}()



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		container.NewVBox(text),
		communicator_checkbox,
		no_website_info_label,
		domain_entry,
		password_entry,
		save_button,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	home_canvas_object_GL = main_scroll
	Current_screen_GL = home_canvas_object_GL

	return home_canvas_object_GL
}
