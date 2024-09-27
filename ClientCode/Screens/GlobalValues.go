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
	"Utils/UtilsSWA"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var global_values_canvas_object_GL fyne.CanvasObject = nil

func GlobalValues() fyne.CanvasObject {
	Current_screen_GL = global_values_canvas_object_GL
	if global_values_canvas_object_GL != nil {
		return global_values_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Text Display section with vertical scrolling
	var registry_text *widget.Label = widget.NewLabel("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed " +
		"do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud " +
		"exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit " +
		"in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non " +
		"proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	registry_text.Wrapping = fyne.TextWrapWord // Enable text wrapping
	var scroll_text *container.Scroll = container.NewVScroll(registry_text)
	scroll_text.SetMinSize(screens_size_GL) // Set the minimum size for the scroll container

	go func() {
		for {
			if Current_screen_GL == global_values_canvas_object_GL {
				registry_text.SetText(UtilsSWA.GetRegistryTextREGISTRY())
			}

			time.Sleep(1 * time.Second)
		}
	}()



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		scroll_text,
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	global_values_canvas_object_GL = main_scroll
	Current_screen_GL = global_values_canvas_object_GL

	return global_values_canvas_object_GL
}
