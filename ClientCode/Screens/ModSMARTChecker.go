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
	"SettingsSync"
	"Utils"
	"Utils/ModsFileInfo"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ModSMARTChecker() fyne.CanvasObject {
	Current_screen_GL = ID_SMART_CHECKER

	return container.NewAppTabs(
		container.NewTabItem("Disks list", smartCheckerCreateDisksListTab()),
		container.NewTabItem("Add disk", smartCheckerCreateAddDiskTab()),
		container.NewTabItem("About", smartCheckerCreateAboutTab()),
	)
}

func smartCheckerCreateAboutTab() *container.Scroll {
	var label_info *widget.Label = widget.NewLabel(SMART_ABOUT)
	label_info.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(label_info)
}

func smartCheckerCreateAddDiskTab() *container.Scroll {
	var entry_id *widget.Entry = widget.NewEntry()
	entry_id.SetPlaceHolder("Disk ID")

	var check_enabled *widget.Check = widget.NewCheck("Disk enabled", nil)
	check_enabled.SetChecked(true)

	var entry_label *widget.Entry = widget.NewEntry()
	entry_label.SetPlaceHolder("Disk label")

	var check_is_hdd *widget.Check = widget.NewCheck("Is it an HDD? (As opposed to an SSD)", nil)
	check_is_hdd.SetChecked(true)

	var btn_add *widget.Button = widget.NewButton("Add", func() {
		if !SettingsSync.AddDiskSMART(entry_id.Text, check_enabled.Checked, entry_label.Text, check_is_hdd.Checked) {
			err := errors.New("disk ID already exists")
			dialog.ShowError(err, Current_window_GL)

			return
		}

		reloadScreen()
	})

	return createMainContentScrollUTILS(
		entry_id,
		check_enabled,
		entry_label,
		check_is_hdd,
		btn_add,
	)
}

func smartCheckerCreateDisksListTab() *container.Scroll {
	var accordion *widget.Accordion = widget.NewAccordion()
	accordion.MultiOpen = true
	var disks_info []ModsFileInfo.DiskInfo = Utils.GetUserSettings(Utils.LOCK_UNLOCK).SMARTChecker.Disks_info
	for i := range disks_info {
		var disk_info *ModsFileInfo.DiskInfo = &disks_info[i]
		var title string = disk_info.Label
		if !disk_info.Enabled {
			title = "[X] " + title
		}
		accordion.Append(widget.NewAccordionItem(trimAccordionTitleUTILS(title), createDiskSetter(disk_info)))
	}

	return createMainContentScrollUTILS(accordion)
}

func createDiskSetter(disk *ModsFileInfo.DiskInfo) *fyne.Container {
	var label_id *widget.Label = widget.NewLabel("Disk ID: " + disk.Id)

	var check_enabled *widget.Check = widget.NewCheck("Disk enabled", nil)
	check_enabled.SetChecked(disk.Enabled)

	var entry_label *widget.Entry = widget.NewEntry()
	entry_label.SetPlaceHolder("Disk label")
	entry_label.SetText(disk.Label)

	var check_is_hdd *widget.Check = widget.NewCheck("Is it an HDD? (As opposed to an SSD)", nil)
	check_is_hdd.SetChecked(disk.Is_HDD)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		disk.Enabled = check_enabled.Checked
		disk.Label = entry_label.Text
		disk.Is_HDD = check_is_hdd.Checked

		reloadScreen()
	})
	btn_save.Importance = widget.SuccessImportance

	var btn_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationDialogUTILS("Are you sure you want to delete this disk?", func(confirmed bool) {
			if confirmed {
				SettingsSync.RemoveDiskSMART(disk.Id)

				reloadScreen()
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

	return container.NewVBox(
		label_id,
		check_enabled,
		entry_label,
		check_is_hdd,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
	)
}
