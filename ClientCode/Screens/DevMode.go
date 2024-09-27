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
	MOD_3 "Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var file_path_GL string = ""

var dev_mode_canvas_object_GL fyne.CanvasObject = nil

func DevMode(window fyne.Window) fyne.CanvasObject {
	Current_screen_GL = dev_mode_canvas_object_GL
	if dev_mode_canvas_object_GL != nil {
		return dev_mode_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Form section
	var form_type *widget.Entry = widget.NewEntry()
	var form_text1 *widget.Entry = widget.NewEntry()
	var form_text2 *widget.Entry = widget.NewEntry()
	var file_button *widget.Button = widget.NewButton("File", func() {
		showFilePicker(window)
	})

	var form *widget.Form = &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type",  Widget: form_type},
			{Text: "Text1", Widget: form_text1},
			{Text: "Text2", Widget: form_text2},
			{Text: "File",  Widget: file_button},
		},
		OnSubmit: func() {
			Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
				Type:  form_type.Text,
				Text1: form_text1.Text,
				Text2: form_text2.Text,
				File:  Utils.CompressString(*Utils.PathFILESDIRS(false, "", file_path_GL).ReadTextFile()),
			})
		},
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Entry and Button section
	var entry_txt_to_speech *widget.Entry = widget.NewEntry()
	entry_txt_to_speech.PlaceHolder = "Enter text to speak"
	entry_txt_to_speech.Text = "This is an example."
	var btn_speak_min *widget.Button = widget.NewButton("Speak (min priority)", func() {
		MOD_3.QueueSpeech(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_LOW, SpeechQueue.MODE_DEFAULT)
	})
	var btn_speak_high *widget.Button = widget.NewButton("Speak (high priority)", func() {
		MOD_3.QueueSpeech(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)
	})
	var btn_skip_speech *widget.Button = widget.NewButton("Skip current speech", func() {
		MOD_3.SkipCurrentSpeech()
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
	main_scroll.SetMinSize(screens_size_GL)

	dev_mode_canvas_object_GL = main_scroll
	Current_screen_GL = dev_mode_canvas_object_GL

	return dev_mode_canvas_object_GL
}

func showFilePicker(w fyne.Window) {
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		file_path_GL = ""
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if f == nil {
			return
		}
		file_path_GL = f.URI().Path()
	}, w)
}
