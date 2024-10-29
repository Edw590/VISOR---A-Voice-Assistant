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
	"Utils/UtilsSWA"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

var settings_canvas_object_GL fyne.CanvasObject = nil

func Settings(param any) fyne.CanvasObject {
	Current_screen_GL = settings_canvas_object_GL
	if settings_canvas_object_GL != nil {
		return settings_canvas_object_GL
	}

	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////////////
	// Combine all sections into a vertical box container
	var objects []fyne.CanvasObject = nil
	var values []UtilsSWA.Value = UtilsSWA.GetValuesREGISTRY()
	for i := len(values) - 1; i >= 0; i-- {
		var value *UtilsSWA.Value = &values[i]
		if !value.Auto_set {
			objects = append(objects, createChooser(value.Key))
		}
	}
	var content *fyne.Container = container.NewVBox(objects...)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	settings_canvas_object_GL = main_scroll
	Current_screen_GL = settings_canvas_object_GL

	return settings_canvas_object_GL
}

func createChooser(key string) *fyne.Container {
	var value *UtilsSWA.Value = UtilsSWA.GetValueREGISTRY(key)
	var label *widget.Label = widget.NewLabel(
		"Name: " + value.Pretty_name +
		"\nType: " + strings.ToLower(value.Type_[len("TYPE_"):]) +
		"\nDescription: " + value.Description)
	var content []fyne.CanvasObject = []fyne.CanvasObject{label}

	var entry *widget.Entry = nil
	var check *widget.Check = nil
	switch value.Type_ {
		case UtilsSWA.TYPE_INT: fallthrough
		case UtilsSWA.TYPE_LONG: fallthrough
		case UtilsSWA.TYPE_STRING: fallthrough
		case UtilsSWA.TYPE_FLOAT: fallthrough
		case UtilsSWA.TYPE_DOUBLE:
			entry = widget.NewEntry()
			entry.SetText(value.Curr_data)
			content = append(content, entry)
		case UtilsSWA.TYPE_BOOL:
			check = widget.NewCheck("Check", nil)
			check.SetChecked(value.GetBool(true))
			content = append(content, check)
	}

	// Save button
	content = append(content, widget.NewButton("Save", func() {
		if entry != nil {
			value.SetData(entry.Text, false)
		} else if check != nil {
			value.SetBool(check.Checked, false)
		}
	}))

	return container.NewVBox(
		content...
	)
}
