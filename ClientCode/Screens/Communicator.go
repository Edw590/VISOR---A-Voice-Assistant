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
	"GPTComm/GPTComm"
	MOD_3 "Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var comm_canvas_object_GL fyne.CanvasObject = nil

func Communicator() fyne.CanvasObject {
	Current_screen_GL = comm_canvas_object_GL
	if comm_canvas_object_GL != nil {
		return comm_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// _Entry and Button section
	var entry_txt_to_speech *widget.Entry = widget.NewEntry()
	entry_txt_to_speech.PlaceHolder = "Text to send to the assistant"
	var btn_send_text *widget.Button = widget.NewButton("Send text", func() {
		Utils.ModsCommsChannels_GL[Utils.NUM_MOD_CmdsExecutor] <- map[string]any{
			"Sentence": entry_txt_to_speech.Text,
		}
	})
	var btn_send_text_gpt_smart *widget.Button = widget.NewButton("Send text directly to GPT (smart)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Communicator not connected."
			MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY)

			return
		}

		if !GPTComm.SendText(entry_txt_to_speech.Text, true) {
			MOD_3.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY)
		}
	})
	var btn_send_text_gpt_dumb *widget.Button = widget.NewButton("Send text directly to GPT (dumb)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Communicator not connected."
			MOD_3.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY)

			return
		}

		if !GPTComm.SendText(entry_txt_to_speech.Text, false) {
			MOD_3.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY)
		}
	})

	//////////////////////////////////////////////////////////////////////////////////
	// Text Display section with vertical scrolling
	var response_text *widget.Entry = widget.NewMultiLineEntry()
	response_text.Wrapping = fyne.TextWrapWord // Enable text wrapping
	response_text.SetMinRowsVisible(100)
	var scroll_text *container.Scroll = container.NewVScroll(response_text)
	scroll_text.SetMinSize(response_text.MinSize()) // Set the minimum size for the scroll container

	go func() {
		for {
			if Current_screen_GL == comm_canvas_object_GL {
				response_text.SetText(GPTComm.GetLastText())
			}
			scroll_text.SetMinSize(response_text.MinSize())

			time.Sleep(1 * time.Second)
		}
	}()



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		entry_txt_to_speech,
		btn_send_text,
		btn_send_text_gpt_smart,
		btn_send_text_gpt_dumb,
		scroll_text,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	comm_canvas_object_GL = main_scroll
	Current_screen_GL = comm_canvas_object_GL

	return comm_canvas_object_GL
}
