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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func ModUserLocator() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_USER_LOCATOR

	return container.NewAppTabs(
		container.NewTabItem("Locations list", userLocatorCreateLocationsListTab()),
		container.NewTabItem("Add location", userLocatorCreateAddLocationTab()),
	)
}

func userLocatorCreateAddLocationTab() *container.Scroll {
	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.PlaceHolder = "Beacon type (\"wifi\" or \"bluetooth\")"
	entry_type.Validator = validation.NewRegexp(`^(wifi|bluetooth)$`, "The location type must be either \"wifi\" or \"bluetooth\"")

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.PlaceHolder = "Beacon name (Wi-Fi SSID or Bluetooth device name)"

	var entry_address *widget.Entry = widget.NewEntry()
	entry_address.PlaceHolder = "Beacon address (Wi-Fi BSSID or Bluetooth device address)"
	entry_address.Validator = validation.NewRegexp(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, "The address must be in the format XX:XX:XX:XX:XX:XX")

	var entry_last_detection_s *widget.Entry = widget.NewEntry()
	entry_last_detection_s.PlaceHolder = "How long the beacon is not found but user may still be in the location (in seconds)"
	entry_last_detection_s.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}

	var entry_max_distance *widget.Entry = widget.NewEntry()
	entry_max_distance.PlaceHolder = "Maximum distance from the beacon to the user (in meters)"
	entry_max_distance.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 32)

		return err
	}

	var entry_location_name *widget.Entry = widget.NewEntry()
	entry_location_name.PlaceHolder = "Location name"


	last_detection_s, _ := strconv.ParseInt(entry_last_detection_s.Text, 10, 64)
	max_distance, _ := strconv.ParseInt(entry_max_distance.Text, 10, 32)
	var button_save *widget.Button = widget.NewButton("Add", func() {
		Utils.User_settings_GL.UserLocator.Locs_info = append(Utils.User_settings_GL.UserLocator.Locs_info, ModsFileInfo.LocInfo{
			Type: entry_type.Text,
			Name: entry_name.Text,
			Address: entry_address.Text,
			Last_detection_s: last_detection_s,
			Max_distance_m: int(max_distance),
			Location: entry_location_name.Text,
		})

		Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
	})

	return createMainContentScrollUTILS(
		entry_type,
		entry_name,
		entry_address,
		entry_last_detection_s,
		entry_max_distance,
		entry_location_name,
		button_save,
	)
}

func userLocatorCreateLocationsListTab() *container.Scroll {
	var objects []fyne.CanvasObject = nil
	var locs_info []ModsFileInfo.LocInfo = Utils.User_settings_GL.UserLocator.Locs_info
	for i := 0; i < len(locs_info); i++ {
		objects = append(objects, createLocationSetter(&locs_info[i], i))
	}

	return createMainContentScrollUTILS(objects...)
}

func createLocationSetter(loc_info *ModsFileInfo.LocInfo, loc_info_idx int) *fyne.Container {
	var entry_type *widget.Entry = widget.NewEntry()
	entry_type.PlaceHolder = "Beacon type (\"wifi\" or \"bluetooth\")"
	entry_type.Validator = validation.NewRegexp(`^(wifi|bluetooth)$`, "The location type must be either \"wifi\" or \"bluetooth\"")
	entry_type.Text = loc_info.Type

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.PlaceHolder = "Beacon name (Wi-Fi SSID or Bluetooth device name)"
	entry_name.Text = loc_info.Name

	var entry_address *widget.Entry = widget.NewEntry()
	entry_address.PlaceHolder = "Beacon address (Wi-Fi BSSID or Bluetooth device address)"
	entry_address.Validator = validation.NewRegexp(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`, "The address must be in the format XX:XX:XX:XX:XX:XX")
	entry_address.Text = loc_info.Address

	var entry_last_detection_s *widget.Entry = widget.NewEntry()
	entry_last_detection_s.PlaceHolder = "How long the beacon is not found but user may still be in the location (in seconds)"
	entry_last_detection_s.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)

		return err
	}
	entry_last_detection_s.Text = strconv.FormatInt(loc_info.Last_detection_s, 10)

	var entry_max_distance *widget.Entry = widget.NewEntry()
	entry_max_distance.PlaceHolder = "Maximum distance from the beacon to the user (in meters)"
	entry_max_distance.Validator = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 32)

		return err
	}
	entry_max_distance.Text = strconv.Itoa(loc_info.Max_distance_m)

	var entry_location_name *widget.Entry = widget.NewEntry()
	entry_location_name.PlaceHolder = "Location name"
	entry_location_name.Text = loc_info.Location

	var button_save *widget.Button = widget.NewButton("Save", func() {
		loc_info.Type = entry_type.Text
		loc_info.Name = entry_name.Text
		loc_info.Address = entry_address.Text
		loc_info.Last_detection_s, _ = strconv.ParseInt(entry_last_detection_s.Text, 10, 64)
		loc_info.Max_distance_m, _ = strconv.Atoi(entry_max_distance.Text)
		loc_info.Location = entry_location_name.Text
	})

	var button_delete *widget.Button = widget.NewButton("Delete", func() {
		createConfirmationUTILS("Are you sure you want to delete this location?", func(confirmed bool) {
			if confirmed {
				Utils.DelElemSLICES(&Utils.User_settings_GL.UserLocator.Locs_info, loc_info_idx)

				Utils.SendToModChannel(Utils.NUM_MOD_VISOR, "Redraw", nil)
			}
		})
	})

	var space *widget.Label = widget.NewLabel("")

	return container.NewVBox(
		entry_type,
		entry_name,
		entry_address,
		entry_last_detection_s,
		entry_max_distance,
		entry_location_name,
		button_save,
		button_delete,
		space,
	)
}
