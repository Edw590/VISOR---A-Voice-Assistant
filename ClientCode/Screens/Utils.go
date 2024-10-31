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
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

func createValueChooserUTILS(value *UtilsSWA.Value) *fyne.Container {
	var label *widget.Label = widget.NewLabel(
		"Name: " + value.Pretty_name[strings.Index(value.Pretty_name, "-") + 1:] +
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
		entry.Validator = func(s string) error {
			if value.Type_ == UtilsSWA.TYPE_INT {
				if _, err := strconv.Atoi(s); err != nil {
					return errors.New("not an int")
				}
			} else if value.Type_ == UtilsSWA.TYPE_LONG {
				if _, err := strconv.ParseInt(s, 10, 64); err != nil {
					return errors.New("not a long")
				}
			} else if value.Type_ == UtilsSWA.TYPE_FLOAT {
				if _, err := strconv.ParseFloat(s, 32); err != nil {
					return errors.New("not a float")
				}
			} else if value.Type_ == UtilsSWA.TYPE_DOUBLE {
				if _, err := strconv.ParseFloat(s, 64); err != nil {
					return errors.New("not a double")
				}
			}
			return nil
		}
		content = append(content, entry)
	case UtilsSWA.TYPE_BOOL:
		check = widget.NewCheck("Check", nil)
		check.SetChecked(value.GetBool(true))
		content = append(content, check)
	}

	content = append(content, widget.NewButton("Save", func() {
		if entry != nil {
			value.SetData(entry.Text, false)
		} else if check != nil {
			value.SetBool(check.Checked, false)
		}
	}))

	var space *widget.Label = widget.NewLabel("")

	content = append(content, space)

	return container.NewVBox(
		content...
	)
}

func createMainContentScrollUTILS(objects... fyne.CanvasObject) *container.Scroll {
	var content *fyne.Container = container.NewVBox(objects...)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	return main_scroll
}

func createConfirmationUTILS(message string, callback func(bool)) *dialog.ConfirmDialog {
	cnf := dialog.NewConfirm("Confirmation", message, callback, Current_window_GL)
	cnf.SetDismissText("Cancel")
	cnf.SetConfirmText("Yes")
	cnf.Show()

	return cnf
}
