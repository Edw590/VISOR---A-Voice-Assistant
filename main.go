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

// Package main provides various examples of Fyne API capabilities.
package main

import (
	"Utils"
	"VISOR/ClientCode/Screens"
	"VISOR/logo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	Utils.PersonalConsts_GL.Init()

	// Create a new application
	myApp := app.NewWithID("com.edw590.visor_c")
	myApp.SetIcon(logo.LogoBlackGmail)
	myWindow := myApp.NewWindow("V.I.S.O.R.")


	// Create the content area with a label to display different screens
	contentLabel := widget.NewLabel("Welcome!")
	contentContainer := container.NewVBox(contentLabel)

	// Create the navigation bar
	navBar := container.NewVBox(
		widget.NewButton("Home", func() {
			contentContainer.Objects = []fyne.CanvasObject{Screens.Home()}
			contentContainer.Refresh()
		}),
		widget.NewButton("Dev Mode", func() {
			contentContainer.Objects = []fyne.CanvasObject{Screens.DevMode()}
			contentContainer.Refresh()
		}),
		/*widget.NewButton("Entry & Button", func() {
			contentContainer.Objects = []fyne.CanvasObject{createEntryButtonScreen()}
			contentContainer.Refresh()
		}),
		widget.NewButton("Progress Bar", func() {
			contentContainer.Objects = []fyne.CanvasObject{createTextScreen()}
			contentContainer.Refresh()
		}),*/
	)


	// Create a split container to hold the navigation bar and the content
	split := container.NewHSplit(navBar, contentContainer)
	split.SetOffset(0.2) // Set the split ratio (20% for nav, 80% for content)

	// Set the content of the window
	myWindow.SetContent(split)

	// Show and run the application
	myWindow.Resize(fyne.NewSize(640, 480))
	myWindow.ShowAndRun()
}
