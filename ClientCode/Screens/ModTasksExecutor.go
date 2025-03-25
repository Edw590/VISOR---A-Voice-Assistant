/*******************************************************************************
 * Copyright 2023-2025 The V.I.S.O.R. authors
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
	"SettingsSync/SettingsSync"
	"TEHelper/TEHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
)

func ModTasksExecutor() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_TASKS_EXECUTOR

	return container.NewAppTabs(
		container.NewTabItem("Tasks list", tasksExecutorCreateTasksListTab()),
		container.NewTabItem("Add task", tasksExecutorCreateAddTaskTab()),
		container.NewTabItem("About", tasksExecutorCreateAboutTab()),
	)
}

func tasksExecutorCreateAboutTab() *container.Scroll {
	var label_info *widget.Label = widget.NewLabel(TASKS_ABOUT)
	label_info.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(label_info)
}

func tasksExecutorCreateAddTaskTab() *container.Scroll {
	var check_enabled *widget.Check = widget.NewCheck("Task enabled", nil)
	check_enabled.SetChecked(true)

	var check_device_active *widget.Check = widget.NewCheck("Device(s) must be in use", nil)
	check_device_active.SetChecked(false)

	var entry_device_ids *widget.Entry = widget.NewMultiLineEntry()
	entry_device_ids.SetPlaceHolder("Device IDs where the task is triggered (one per line)")
	entry_device_ids.SetMinRowsVisible(3)

	var entry_message *widget.Entry = widget.NewEntry()
	entry_message.SetPlaceHolder("Message to speak when triggered")

	var entry_command *widget.Entry = widget.NewEntry()
	entry_command.SetPlaceHolder("Command to execute after speaking")

	var entry_time *widget.Entry = widget.NewEntry()
	entry_time.SetPlaceHolder("Time trigger (format: \"2024-12-31 -- 23:59:59\")")
	entry_time.Validator = validation.NewRegexp(`^(\d{4}-\d{2}-\d{2} -- \d{2}:\d{2}:\d{2})?$`, "wrong format")

	var entry_repeat_each_min *widget.Entry = widget.NewEntry()
	entry_repeat_each_min.SetPlaceHolder("Repeat each X minutes")
	entry_repeat_each_min.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}

	var entry_user_location *widget.Entry = widget.NewEntry()
	entry_user_location.SetPlaceHolder("User location trigger")

	var entry_programmable_condition *widget.Entry = widget.NewEntry()
	entry_programmable_condition.SetPlaceHolder("Programmable condition (in Go)")
	entry_programmable_condition.Validator = func(s string) error {
		if s == "" {
			return nil
		}

		_, err := TEHelper.ComputeCondition(s)

		return err
	}

	repeat_each_min, _ := strconv.ParseInt(entry_repeat_each_min.Text, 10, 64)
	var btn_add *widget.Button = widget.NewButton("Add", func() {
		SettingsSync.AddTaskTASKS(check_enabled.Checked, check_device_active.Checked, entry_device_ids.Text,
			entry_message.Text, entry_command.Text, entry_time.Text, repeat_each_min, entry_user_location.Text,
			entry_programmable_condition.Text)

		reloadScreen()
	})

	return createMainContentScrollUTILS(
		check_enabled,
		check_device_active,
		entry_device_ids,
		entry_message,
		entry_command,
		entry_time,
		entry_repeat_each_min,
		entry_user_location,
		entry_programmable_condition,
		btn_add,
	)
}

func tasksExecutorCreateTasksListTab() *container.Scroll {
	var accordion *widget.Accordion = widget.NewAccordion()
	accordion.MultiOpen = true
	var tasks []ModsFileInfo.Task = Utils.User_settings_GL.TasksExecutor.Tasks
	for i := range tasks {
		var task *ModsFileInfo.Task = &tasks[i]
		var title = task.Message
		if title == "" {
			title = task.Command
		}
		if !task.Enabled {
			title = "[X] " + title
		}
		accordion.Append(widget.NewAccordionItem(trimAccordionTitleUTILS(title), createTaskSetter(task)))
	}

	return createMainContentScrollUTILS(accordion)
}

func createTaskSetter(task *ModsFileInfo.Task) *fyne.Container {
	var label_id *widget.Label = widget.NewLabel("Task ID: " + strconv.Itoa(int(task.Id)))

	var check_enabled *widget.Check = widget.NewCheck("Task enabled", nil)
	check_enabled.SetChecked(task.Enabled)

	var check_device_active *widget.Check = widget.NewCheck("Device(s) must be in use", nil)
	check_device_active.SetChecked(task.Device_active)

	var entry_device_ids *widget.Entry = widget.NewMultiLineEntry()
	entry_device_ids.SetText(strings.Join(task.Device_IDs, "\n"))
	entry_device_ids.SetPlaceHolder("Device IDs where the task is triggered (one per line)")
	entry_device_ids.SetMinRowsVisible(3)

	var entry_message *widget.Entry = widget.NewEntry()
	entry_message.SetText(task.Message)
	entry_message.SetPlaceHolder("Message to speak when triggered")

	var entry_command *widget.Entry = widget.NewEntry()
	entry_command.SetText(task.Command)
	entry_command.SetPlaceHolder("Command to execute after speaking")

	var entry_time *widget.Entry = widget.NewEntry()
	entry_time.SetText(task.Time)
	entry_time.SetPlaceHolder("Time trigger (format: \"2024-12-31 -- 23:59:59\")")
	entry_time.Validator = validation.NewRegexp(`^(\d{4}-\d{2}-\d{2} -- \d{2}:\d{2}:\d{2})?$`, "wrong format")

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

	var entry_programmable_condition *widget.Entry = widget.NewEntry()
	entry_programmable_condition.SetText(task.Programmable_condition)
	entry_programmable_condition.SetPlaceHolder("Programmable condition (in Go)")
	entry_programmable_condition.Validator = func(s string) error {
		if s == "" {
			return nil
		}

		_, err := TEHelper.ComputeCondition(s)

		return err
	}

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		task.Enabled = check_enabled.Checked
		task.Device_active = check_device_active.Checked
		task.Device_IDs = strings.Split(entry_device_ids.Text, "\n")
		task.Message = entry_message.Text
		task.Command = entry_command.Text
		task.Time = entry_time.Text
		task.Repeat_each_min, _ = strconv.ParseInt(entry_repeat_each_min.Text, 10, 64)
		task.User_location = entry_user_location.Text
		task.Programmable_condition = entry_programmable_condition.Text

		reloadScreen()
	})
	btn_save.Importance = widget.SuccessImportance

	var btn_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationDialogUTILS("Are you sure you want to delete this task?", func(confirmed bool) {
			if confirmed {
				SettingsSync.RemoveTaskTASKS(task.Id)

				reloadScreen()
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

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
		entry_programmable_condition,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
	)
}
