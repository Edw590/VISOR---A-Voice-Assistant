/*******************************************************************************
 * Copyright 2023-2023 Edw590
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

package main

import (
	"GPT/GPT"
	"Utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

func main() {
	Utils.PersonalConsts_GL.Init()

	GPT.SetWebsiteInfo(Utils.PersonalConsts_GL.WEBSITE_URL, Utils.PersonalConsts_GL.WEBSITE_PW)

	a := app.New()
	w := a.NewWindow("Entry Widget")

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	//text := canvas.NewText("Text Object", color.White)
	var text = widget.NewLabel("Text Object")
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}

	content := container.NewVBox(input,
		widget.NewButton("Save", func() {
			log.Println("Content was:", input.Text)
		}),
		text,
	)

	fmt.Println(GPT.GetTextFromEntry(GPT.GetEntry(-1)))

	go func() {
		for {
			var num_entries int = GPT.GetNumEntries()
			fmt.Println(num_entries)
			for i := num_entries - 1; i >= 0; i-- {
				if GPT.GetTextFromEntry(GPT.GetEntry(i)) != GPT.GEN_ERROR {
					text.SetText(GPT.GetTextFromEntry(GPT.GetEntry(i)))
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()

	w.SetContent(content)
	w.ShowAndRun()
}
