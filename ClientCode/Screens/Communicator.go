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
	"GPT/GPT"
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
	// Entry and Button section
	var entry_txt_to_speech *widget.Entry = widget.NewEntry()
	entry_txt_to_speech.PlaceHolder = "Text to send to the assistant"
	var btn_send_text *widget.Button = widget.NewButton("Send text", func() {
		GPT.SendText(entry_txt_to_speech.Text)
	})

	//////////////////////////////////////////////////////////////////////////////////
	// Text Display section with vertical scrolling
	var response_text *widget.Label = widget.NewLabel("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed " +
		"do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud " +
		"exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit " +
		"in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non " +
		"proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	response_text.Wrapping = fyne.TextWrapWord // Enable text wrapping
	var scroll_text *container.Scroll = container.NewVScroll(response_text)
	scroll_text.SetMinSize(response_text.MinSize()) // Set the minimum size for the scroll container

	go func() {
		for {
			if Current_screen_GL == comm_canvas_object_GL {
				response_text.SetText(GPT.GetTextFromEntry(GPT.GetEntry(-1)))
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
		scroll_text,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(fyne.NewSize(550, 480))

	comm_canvas_object_GL = main_scroll
	Current_screen_GL = comm_canvas_object_GL

	return comm_canvas_object_GL
}
