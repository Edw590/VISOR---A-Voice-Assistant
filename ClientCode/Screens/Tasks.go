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
	"TEHelper/TEHelper"
	"Utils/ModsFileInfo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
	"strings"
)

var tasks_canvas_object_GL fyne.CanvasObject = nil

func Tasks(param any) fyne.CanvasObject {
	Current_screen_GL = tasks_canvas_object_GL

	var objects []fyne.CanvasObject = nil
	var tasks []ModsFileInfo.Task = TEHelper.GetTasks()
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Id > tasks[j].Id
	})
	for i := len(tasks) - 1; i >= 0; i-- {
		objects = append(objects, createTaskSetter(&tasks[i]))
	}
	var content *fyne.Container = container.NewVBox(objects...)

	var main_scroll *container.Scroll = container.NewVScroll(content)
	main_scroll.SetMinSize(screens_size_GL)

	tasks_canvas_object_GL = main_scroll
	Current_screen_GL = tasks_canvas_object_GL

	return tasks_canvas_object_GL
}

func createTaskSetter(task *ModsFileInfo.Task) *fyne.Container {
	var label_id *widget.Label = widget.NewLabel("Task ID: " + strconv.Itoa(task.Id))

	var check_enabled *widget.Check = widget.NewCheck("Task enabled", nil)
	check_enabled.SetChecked(task.Enabled)

	var check_device_active *widget.Check = widget.NewCheck("Device(s) must be active", nil)
	check_device_active.SetChecked(task.Device_active)

	var entry_device_ids *widget.Entry = widget.NewMultiLineEntry()
	entry_device_ids.SetText(strings.Join(task.Device_IDs, "\n"))
	entry_device_ids.SetPlaceHolder("Device IDs (one per line)")
	entry_device_ids.SetMinRowsVisible(3)

	var entry_message *widget.Entry = widget.NewEntry()
	entry_message.SetText(task.Message)
	entry_message.SetPlaceHolder("Message to speak when triggered")

	var entry_command *widget.Entry = widget.NewEntry()
	entry_command.SetText(task.Command)
	entry_command.SetPlaceHolder("Command to execute after speaking")

	var entry_time *widget.Entry = widget.NewEntry()
	entry_time.SetText(task.Time)
	entry_time.SetPlaceHolder("Time trigger (format: 2024-12-31 -- 23:59:59)")
	entry_time.Validator = validation.NewRegexp(`^\d{4}-\d{2}-\d{2} -- \d{2}:\d{2}:\d{2}$`, "wrong format")

	var entry_repeat_each_min *widget.Entry = widget.NewEntry()
	entry_repeat_each_min.SetText(strconv.FormatInt(task.Repeat_each_min, 10))
	entry_repeat_each_min.SetPlaceHolder("Repeat each X minutes")
	entry_repeat_each_min.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}

	var entry_user_location *widget.Entry = widget.NewEntry()
	entry_user_location.SetText(task.User_location)
	entry_user_location.SetPlaceHolder("User location trigger")

	var entry_device_condition *widget.Entry = widget.NewEntry()
	entry_device_condition.SetText(task.Device_condition)
	entry_device_condition.SetPlaceHolder("Programmable condition (in Go)")

	// Save button
	var button_save *widget.Button = widget.NewButton("Save", func() {
		task.Enabled = check_enabled.Checked
		task.Device_active = check_device_active.Checked
		task.Device_IDs = strings.Split(entry_device_ids.Text, "\n")
		task.Message = entry_message.Text
		task.Command = entry_command.Text
		task.Time = entry_time.Text
		task.Repeat_each_min, _ = strconv.ParseInt(entry_repeat_each_min.Text, 10, 64)
		task.User_location = entry_user_location.Text
		task.Device_condition = entry_device_condition.Text
	})

	return container.NewVBox(
		label_id,
		check_enabled,
		check_device_active,
		entry_device_ids,
		entry_message,
		entry_command,
		entry_time,
		entry_repeat_each_min,
		entry_user_location,
		entry_device_condition,
		button_save,
	)
}
