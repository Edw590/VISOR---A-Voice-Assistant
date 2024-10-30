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
	"SystemChecker"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var mod_system_checker_canvas_object_GL fyne.CanvasObject = nil

func ModSystemChecker() fyne.CanvasObject {
	var tabs *container.AppTabs = container.NewAppTabs(
		container.NewTabItem("System state", systemCheckerCreateSystemStateTab()),
	)

	mod_system_checker_canvas_object_GL = tabs
	Current_screen_GL = mod_system_checker_canvas_object_GL

	return mod_system_checker_canvas_object_GL
}

func systemCheckerCreateSystemStateTab() *container.Scroll {
	var sys_state_text *widget.Label = widget.NewLabel("")
	sys_state_text.Wrapping = fyne.TextWrapWord

	go func() {
		time.Sleep(500 * time.Millisecond)
		for {
			if Current_screen_GL == mod_system_checker_canvas_object_GL {
				sys_state_text.SetText(SystemChecker.GetDeviceInfoText())
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(sys_state_text)
}
