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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func ModUserLocator() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_USER_LOCATOR

	return container.NewAppTabs(
		container.NewTabItem("Locations list", userLocatorCreateLocationsListTab()),
		container.NewTabItem("Add location", userLocatorCreateAddLocationTab()),
		container.NewTabItem("Settings", userLocatorCreateSettingsTab()),
		container.NewTabItem("About", userLocatorCreateAboutTab()),
	)
}

func userLocatorCreateAboutTab() *container.Scroll {
	var label_info *widget.Label = widget.NewLabel(LOCATOR_ABOUT)
	label_info.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(label_info)
}

func userLocatorCreateSettingsTab() *container.Scroll {
	var entry_always_with_device *widget.Entry = widget.NewEntry()
	entry_always_with_device.SetPlaceHolder("ID of the device always with the user (user's phone for example) or " +
		"empty if none")
	entry_always_with_device.SetText(Utils.GetUserSettings().UserLocator.AlwaysWith_device)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.GetUserSettings().UserLocator.AlwaysWith_device = entry_always_with_device.Text
	})
	btn_save.Importance = widget.SuccessImportance

	return createMainContentScrollUTILS(
		entry_always_with_device,
		btn_save,
	)
}

func userLocatorCreateAddLocationTab() *container.Scroll {
	var check_enabled *widget.Check = widget.NewCheck("Location enabled", nil)
	check_enabled.SetChecked(true)

	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.SetPlaceHolder("Beacon type (\"wifi\" or \"bluetooth\")")
	entry_type.Validator = validation.NewRegexp(`^(wifi|bluetooth)$`, "The location type must be either \"wifi\" or \"bluetooth\"")

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.SetPlaceHolder("Beacon name (Wi-Fi SSID or Bluetooth device name)")

	var entry_address *widget.Entry = widget.NewEntry()
	entry_address.SetPlaceHolder("Beacon address (Wi-Fi BSSID or Bluetooth device address)")
	entry_address.Validator = validation.NewRegexp(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, "The address must be in the format XX:XX:XX:XX:XX:XX")

	var entry_last_detection_s *widget.Entry = widget.NewEntry()
	entry_last_detection_s.SetPlaceHolder("How long the beacon is not found but user may still be in the location (in seconds)")
	entry_last_detection_s.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}

	var entry_max_distance *widget.Entry = widget.NewEntry()
	entry_max_distance.SetPlaceHolder("Maximum distance from the beacon to the user (in meters)")
	entry_max_distance.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 32)

		return err
	}

	var entry_location_name *widget.Entry = widget.NewEntry()
	entry_location_name.SetPlaceHolder("Location name")


	last_detection_s, _ := strconv.ParseInt(entry_last_detection_s.Text, 10, 64)
	max_distance, _ := strconv.ParseInt(entry_max_distance.Text, 10, 32)
	var btn_add *widget.Button = widget.NewButton("Add", func() {
		SettingsSync.AddLocationLOCATIONS(check_enabled.Checked, entry_type.Text, entry_name.Text, entry_address.Text,
			last_detection_s, int32(max_distance), entry_location_name.Text)

		reloadScreen()
	})

	return createMainContentScrollUTILS(
		check_enabled,
		entry_type,
		entry_name,
		entry_address,
		entry_last_detection_s,
		entry_max_distance,
		entry_location_name,
		btn_add,
	)
}

func userLocatorCreateLocationsListTab() *container.Scroll {
	var accordion *widget.Accordion = widget.NewAccordion()
	accordion.MultiOpen = true
	var locs_info []ModsFileInfo.LocInfo = Utils.GetUserSettings().UserLocator.Locs_info
	for i := range locs_info {
		var loc_info *ModsFileInfo.LocInfo = &locs_info[i]
		var title = loc_info.Name
		if title == "" {
			title = loc_info.Address
		}
		title = loc_info.Location + " - " + title
		if !loc_info.Enabled {
			title = "[X] " + title
		}
		accordion.Append(widget.NewAccordionItem(trimAccordionTitleUTILS(title), createLocationSetter(loc_info)))
	}

	return createMainContentScrollUTILS(accordion)
}

func createLocationSetter(loc_info *ModsFileInfo.LocInfo) *fyne.Container {
	var check_enabled *widget.Check = widget.NewCheck("Location enabled", nil)
	check_enabled.SetChecked(loc_info.Enabled)

	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.SetPlaceHolder("Beacon type (\"wifi\" or \"bluetooth\")")
	entry_type.Validator = validation.NewRegexp(`^(wifi|bluetooth)$`, "The location type must be either \"wifi\" or \"bluetooth\"")
	entry_type.SetText(loc_info.Type)

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.SetPlaceHolder("Beacon name (Wi-Fi SSID or Bluetooth device name)")
	entry_name.SetText(loc_info.Name)

	var entry_address *widget.Entry = widget.NewEntry()
	entry_address.SetPlaceHolder("Beacon address (Wi-Fi BSSID or Bluetooth device address)")
	entry_address.Validator = validation.NewRegexp(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, "The address must be in the format XX:XX:XX:XX:XX:XX")
	entry_address.SetText(loc_info.Address)

	var entry_last_detection_s *widget.Entry = widget.NewEntry()
	entry_last_detection_s.SetPlaceHolder("How long the beacon is not found but user may still be in the location (in seconds)")
	entry_last_detection_s.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}
	entry_last_detection_s.Text = strconv.FormatInt(loc_info.Last_detection_s, 10)

	var entry_max_distance *widget.Entry = widget.NewEntry()
	entry_max_distance.SetPlaceHolder("Maximum distance from the beacon to the user (in meters)")
	entry_max_distance.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 32)

		return err
	}
	entry_max_distance.SetText(strconv.Itoa(int(loc_info.Max_distance_m)))

	var entry_location_name *widget.Entry = widget.NewEntry()
	entry_location_name.SetPlaceHolder("Location name")
	entry_location_name.SetText(loc_info.Location)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		loc_info.Enabled = check_enabled.Checked
		loc_info.Type = entry_type.Text
		loc_info.Name = entry_name.Text
		loc_info.Address = entry_address.Text
		loc_info.Last_detection_s, _ = strconv.ParseInt(entry_last_detection_s.Text, 10, 64)
		max_distance_m, _ := strconv.ParseInt(entry_max_distance.Text, 10, 32)
		loc_info.Max_distance_m = int32(max_distance_m)
		loc_info.Location = entry_location_name.Text

		reloadScreen()
	})
	btn_save.Importance = widget.SuccessImportance

	var btn_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationDialogUTILS("Are you sure you want to delete this location?", func(confirmed bool) {
			if confirmed {
				SettingsSync.RemoveLocationLOCATIONS(loc_info.Id)

				reloadScreen()
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

	return container.NewVBox(
		check_enabled,
		entry_type,
		entry_name,
		entry_address,
		entry_last_detection_s,
		entry_max_distance,
		entry_location_name,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
	)
}
