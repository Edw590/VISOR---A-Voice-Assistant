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
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func DevMode() fyne.CanvasObject {
	type_ := widget.NewEntry()
	text1 := widget.NewEntry()
	text2 := widget.NewEntry()
	text3 := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Type",  Widget: type_},
			{Text: "Text1", Widget: text1},
			{Text: "Text2", Widget: text2},
			{Text: "Text3", Widget: text3},
		},
		OnSubmit: func() {
			Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
				Name:  type_.Text,
				Text1: text1.Text,
				Text2: text2.Text,
				Text3: text3.Text,
			})
		},
	}

	// Entry and Button section
	entry := widget.NewEntry()
	button := widget.NewButton("Submit", func() {
		Utils.QueueSpeechSPEECH(entry.Text)
	})
	entryButtonSection := container.NewVBox(entry, button)

	// Text Display section with vertical scrolling
	text := widget.NewLabel(`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`)
	text.Wrapping = fyne.TextWrapWord // Enable text wrapping
	scroll := container.NewVScroll(text)
	scroll.SetMinSize(fyne.NewSize(300, 200)) // Set the minimum size for the scroll container

	// Combine all sections into a vertical box container
	return container.NewVBox(form, entryButtonSection, scroll)
}
