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
	"strconv"
	"time"
)

var module_status_canvas_object_GL fyne.CanvasObject = nil

func ModulesStatus(modules []Utils.Module) fyne.CanvasObject {
	Current_screen_GL = module_status_canvas_object_GL
	if module_status_canvas_object_GL != nil {
		return module_status_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	// Text Display section with vertical scrolling
	var module_status *widget.Label = widget.NewLabel("")
	module_status.Wrapping = fyne.TextWrapWord // Enable text wrapping
	var scroll_text *container.Scroll = container.NewVScroll(module_status)
	scroll_text.SetMinSize(fyne.NewSize(300, 400)) // Set the minimum size for the scroll container

	go func() {
		for {
			if Current_screen_GL == module_status_canvas_object_GL {
				var text string = ""
				for i, module := range modules {
					text += "- " + Utils.GetModNameMODULES(i) + " running: " + strconv.FormatBool(!module.Stopped) + "\n\n"
				}
				module_status.SetText(text)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	var canvas_objs []fyne.CanvasObject = []fyne.CanvasObject{
		scroll_text,
	}
	for _, obj := range getCheckBoxes(modules) {
		canvas_objs = append(canvas_objs, obj)
	}

	for _, obj := range canvas_objs {
		if checkbox, ok := obj.(*widget.Check); ok {
			checkbox.SetChecked(true)
		}
	}



	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var content *fyne.Container = container.NewVBox(
		canvas_objs...
	)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(fyne.NewSize(550, 480))

	module_status_canvas_object_GL = main_scroll
	Current_screen_GL = module_status_canvas_object_GL

	return module_status_canvas_object_GL
}

func getCheckBoxes(modules []Utils.Module) []fyne.CanvasObject {
	var check_boxes []fyne.CanvasObject

	// Couldn't do it automatically. So here they are manually...

	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_SMARTChecker), func(b bool) {
		modules[Utils.NUM_MOD_SMARTChecker].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_Speech), func(b bool) {
		modules[Utils.NUM_MOD_Speech].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_RssFeedNotifier), func(b bool) {
		modules[Utils.NUM_MOD_RssFeedNotifier].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_EmailSender), func(b bool) {
		modules[Utils.NUM_MOD_EmailSender].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_OnlineInfoChk), func(b bool) {
		modules[Utils.NUM_MOD_OnlineInfoChk].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_GPTCommunicator), func(b bool) {
		modules[Utils.NUM_MOD_GPTCommunicator].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_WebsiteBackend), func(b bool) {
		modules[Utils.NUM_MOD_WebsiteBackend].Enabled = b
	}))
	check_boxes = append(check_boxes, widget.NewCheck(Utils.GetModNameMODULES(Utils.NUM_MOD_UserLocator), func(b bool) {
		modules[Utils.NUM_MOD_UserLocator].Enabled = b
	}))


	return check_boxes
}
