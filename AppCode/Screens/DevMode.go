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

package Screens

import (
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var dev_mode_canvas_object_GL fyne.CanvasObject = nil

func DevMode() fyne.CanvasObject {
	Current_screen_GL = dev_mode_canvas_object_GL
	if dev_mode_canvas_object_GL != nil {
		return dev_mode_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Form section
	var form_type_ *widget.Entry = widget.NewEntry()
	var form_text1 *widget.Entry = widget.NewEntry()
	var form_text2 *widget.Entry = widget.NewEntry()
	var form_text3 *widget.Entry = widget.NewEntry()

	var form *widget.Form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type",  Widget: form_type_},
			{Text: "Text1", Widget: form_text1},
			{Text: "Text2", Widget: form_text2},
			{Text: "Text3", Widget: form_text3},
		},
		OnSubmit: func() {
			Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
				Name:  form_type_.Text,
				Text1: form_text1.Text,
				Text2: form_text2.Text,
				Text3: form_text3.Text,
			})
		},
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Entry and Button section
	var entry_txt_to_speech *widget.Entry = widget.NewEntry()
	entry_txt_to_speech.PlaceHolder = "Enter text to speak"
	entry_txt_to_speech.Text = "This is an example."
	var btn_speak_min *widget.Button = widget.NewButton("Speak (min priority)", func() {
		Utils.QueueSpeechSPEECH(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_LOW)
	})
	var btn_speak_high *widget.Button = widget.NewButton("Speak (high priority)", func() {
		Utils.QueueSpeechSPEECH(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_HIGH)
	})
	var btn_skip_speech *widget.Button = widget.NewButton("Skip current speech", func() {
		Utils.SkipCurrentSpeechSPEECH()
	})



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		form,
		entry_txt_to_speech,
		btn_speak_min,
		btn_speak_high,
		btn_skip_speech,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(fyne.NewSize(550, 480))

	dev_mode_canvas_object_GL = main_scroll
	Current_screen_GL = dev_mode_canvas_object_GL

	return dev_mode_canvas_object_GL
}
