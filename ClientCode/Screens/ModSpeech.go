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
	"Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var file_path_GL string = ""

var mod_speech_canvas_object_GL fyne.CanvasObject = nil

func ModSpeech(param any) fyne.CanvasObject {
	Current_screen_GL = mod_speech_canvas_object_GL
	if mod_speech_canvas_object_GL != nil {
		return mod_speech_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Entry and Button section
	var entry_txt_to_speech *widget.Entry = widget.NewEntry()
	entry_txt_to_speech.PlaceHolder = "Enter text to speak"
	entry_txt_to_speech.Text = "This is an example."
	var btn_speak_min *widget.Button = widget.NewButton("Speak (min priority)", func() {
		Speech.QueueSpeech(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_LOW, SpeechQueue.MODE_DEFAULT, "", 0)
	})
	var btn_speak_high *widget.Button = widget.NewButton("Speak (high priority)", func() {
		Speech.QueueSpeech(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE_DEFAULT, "", 0)
	})
	var btn_speak_critical *widget.Button = widget.NewButton("Speak (critical priority)", func() {
		Speech.QueueSpeech(entry_txt_to_speech.Text, SpeechQueue.PRIORITY_CRITICAL, SpeechQueue.MODE_DEFAULT, "", 0)
	})
	var btn_skip_speech *widget.Button = widget.NewButton("Skip current speech", func() {
		Speech.SkipCurrentSpeech()
	})
	var btn_config_tts *widget.Button = widget.NewButton("Configure Windows SAPI TTS", func() {
		_, _ = Utils.ExecCmdSHELL([]string{"control.exe C:\\Windows\\System32\\Speech\\SpeechUX\\sapi.cpl"})
	})



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		entry_txt_to_speech,
		btn_speak_min,
		btn_speak_high,
		btn_speak_critical,
		btn_skip_speech,
		btn_config_tts,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	mod_speech_canvas_object_GL = main_scroll
	Current_screen_GL = mod_speech_canvas_object_GL

	return mod_speech_canvas_object_GL
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
