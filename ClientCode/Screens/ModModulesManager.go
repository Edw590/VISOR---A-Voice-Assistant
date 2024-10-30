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
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
)

var module_status_canvas_object_GL fyne.CanvasObject = nil

func ModulesStatus(param any) fyne.CanvasObject {
	var modules []Utils.Module = param.([]Utils.Module)

	var tabs *container.AppTabs = container.NewAppTabs(
		container.NewTabItem("Modules status", modulesManagerCreateModulesStatusTab(modules)),
	)

	module_status_canvas_object_GL = tabs
	Current_screen_GL = module_status_canvas_object_GL

	return module_status_canvas_object_GL
}

func modulesManagerCreateModulesStatusTab(modules []Utils.Module) *container.Scroll {
	var module_status_text *widget.Label = widget.NewLabel("")
	module_status_text.Wrapping = fyne.TextWrapWord

	go func() {
		time.Sleep(500 * time.Millisecond)
		for {
			if Current_screen_GL == module_status_canvas_object_GL {
				var text string = ""
				for i, module := range modules {
					if Utils.MOD_NUMS_SUPPORT[i] & Utils.MOD_CLIENT != 0 {
						text += "- " + Utils.GetModNameMODULES(i) + " running: " + strconv.FormatBool(!module.Stopped) +
							"\n\n"
					}
				}
				text = text[:len(text)-2] // Remove the last 2 newlines
				module_status_text.SetText(text)
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(module_status_text)
}
