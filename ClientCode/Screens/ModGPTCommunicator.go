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
	"GPTComm/GPTComm"
	"Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var mod_gpt_communicator_canvas_object_GL fyne.CanvasObject = nil

func ModGPTCommunicator() fyne.CanvasObject {
	var tabs *container.AppTabs = container.NewAppTabs(
		container.NewTabItem("Communicator", gptCommunicatorCreateCommunicatorTab()),
	)

	mod_gpt_communicator_canvas_object_GL = tabs
	Current_screen_GL = mod_gpt_communicator_canvas_object_GL

	return mod_gpt_communicator_canvas_object_GL
}

func gptCommunicatorCreateCommunicatorTab() *container.Scroll {
	var text_to_send *widget.Entry = widget.NewMultiLineEntry()
	text_to_send.Wrapping = fyne.TextWrapWord
	text_to_send.SetMinRowsVisible(6) // 6 lines, like ChatGPT has
	text_to_send.PlaceHolder = "Text to send to the assistant"
	var btn_send_text *widget.Button = widget.NewButton("Send text", func() {
		Utils.SendToModChannel(Utils.NUM_MOD_CmdsExecutor, "Sentence", text_to_send.Text)
	})
	var btn_send_text_gpt_smart *widget.Button = widget.NewButton("Send text directly to GPT (smart)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Communicator not connected."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		if !GPTComm.SendText(text_to_send.Text, true) {
			Speech.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})
	var btn_send_text_gpt_dumb *widget.Button = widget.NewButton("Send text directly to GPT (dumb)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Communicator not connected."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		if !GPTComm.SendText(text_to_send.Text, false) {
			Speech.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})

	var response_text *widget.Entry = widget.NewMultiLineEntry()
	response_text.PlaceHolder = "Response from the assistant's LLM/GPT"
	response_text.Wrapping = fyne.TextWrapWord
	response_text.SetMinRowsVisible(100)

	go func() {
		time.Sleep(500 * time.Millisecond)
		var old_text string = ""
		for {
			if Current_screen_GL == mod_gpt_communicator_canvas_object_GL {
				var new_text string = GPTComm.GetLastText()
				if old_text != new_text {
					old_text = new_text
					response_text.SetText(old_text)
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(
		text_to_send,
		btn_send_text,
		btn_send_text_gpt_smart,
		btn_send_text_gpt_dumb,
		response_text,
	)
}
