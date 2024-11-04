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
	)
}

func smartCheckerCreateAddDiskTab() *container.Scroll {
	var entry_id *widget.Entry = widget.NewEntry()
	entry_id.SetPlaceHolder("Disk ID")

	var entry_label *widget.Entry = widget.NewEntry()
	entry_label.SetPlaceHolder("Disk label")

	var check_is_hdd *widget.Check = widget.NewCheck("Is it an HDD? (As opposed to an SSD)", nil)
	check_is_hdd.SetChecked(true)

	var button_save *widget.Button = widget.NewButton("Add", func() {
		for _, disk := range Utils.User_settings_GL.SMARTChecker.Disks_info {
			if disk.Id == entry_id.Text {
				err := errors.New("disk ID already exists")
				dialog.ShowError(err, Current_window_GL)

				return
			}
		}

		Utils.User_settings_GL.SMARTChecker.Disks_info = append(Utils.User_settings_GL.SMARTChecker.Disks_info,
			ModsFileInfo.DiskInfo{
				Id:     entry_id.Text,
				Label:  entry_label.Text,
				Is_HDD: check_is_hdd.Checked,
		})

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})

	return createMainContentScrollUTILS(
		entry_id,
		entry_label,
		check_is_hdd,
		button_save,
	)
}

func smartCheckerCreateDisksListTab() *container.Scroll {
	var objects []fyne.CanvasObject = nil
	var disks []ModsFileInfo.DiskInfo = Utils.User_settings_GL.SMARTChecker.Disks_info
	for i, disk := range disks {
		objects = append(objects, createDiskSetter(&disk, i))
	}

	return createMainContentScrollUTILS(objects...)
}

func createDiskSetter(disk *ModsFileInfo.DiskInfo, disk_idx int) *fyne.Container {
	var label_id *widget.Label = widget.NewLabel("Disk ID: " + disk.Id)

	var entry_label *widget.Entry = widget.NewEntry()
	entry_label.SetPlaceHolder("Disk label")
	entry_label.SetText(disk.Label)

	var check_is_hdd *widget.Check = widget.NewCheck("Is it an HDD? (As opposed to an SSD)", nil)
	check_is_hdd.SetChecked(disk.Is_HDD)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		disk.Label = entry_label.Text
		disk.Is_HDD = check_is_hdd.Checked
	})

	var btn_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationUTILS("Are you sure you want to delete this disk?", func(confirmed bool) {
			if confirmed {
				Utils.DelElemSLICES(&Utils.User_settings_GL.SMARTChecker.Disks_info, disk_idx)

				Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

	var space *widget.Label = widget.NewLabel("")

	return container.NewVBox(
		label_id,
		entry_label,
		check_is_hdd,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
		space,
	)
}
