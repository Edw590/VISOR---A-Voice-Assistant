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

package SettingsSync

import (
	"Utils"
	"Utils/ModsFileInfo"
	"sort"
)

/*
AddDiskSMART adds a disk to the user settings.

-----------------------------------------------------------

– Params:
  - id – the disk serial number
  - label – the disk label
  - is_hdd – true if the disk is an HDD, false otherwise

– Returns:
  - true if the disk was added, false if the ID already exists
 */
func AddDiskSMART(id string, enabled bool, label string, is_hdd bool) bool {
	var disks_info *[]ModsFileInfo.DiskInfo = &Utils.GetUserSettings().SMARTChecker.Disks_info
	for _, disk_info := range *disks_info {
		if disk_info.Id == id {
			return false
		}
	}

	// Add the disk to the user settings
	*disks_info = append(*disks_info, ModsFileInfo.DiskInfo{
		Id:      id,
		Enabled: enabled,
		Label:   label,
		Is_HDD:  is_hdd,
	})

	sort.SliceStable(*disks_info, func(i, j int) bool {
		return (*disks_info)[i].Label < (*disks_info)[j].Label
	})

	return true
}

/*
RemoveDiskSMART removes a disk from the user settings.

-----------------------------------------------------------

– Params:
  - id – the disk serial number
 */
func RemoveDiskSMART(id string) {
	var disks_info *[]ModsFileInfo.DiskInfo = &Utils.GetUserSettings().SMARTChecker.Disks_info
	for i := range *disks_info {
		if (*disks_info)[i].Id == id {
			Utils.DelElemSLICES(disks_info, i)

			break
		}
	}
}

/*
GetIdsListSMART returns a list of all disks' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all disks' IDs separated by "|"
 */
func GetIdsListSMART() string {
	var ids_list string = ""
	for _, disk_info := range Utils.GetUserSettings().SMARTChecker.Disks_info {
		ids_list += disk_info.Id + "|"
	}
	ids_list = ids_list[:len(ids_list)-1]

	return ids_list
}

/*
GetDiskSMART returns a disk by its ID.

-----------------------------------------------------------

– Params:
  - id – the disk serial number

– Returns:
  - the disk or nil if the disk was not found
 */
func GetDiskSMART(id string) *ModsFileInfo.DiskInfo {
	var disks_info []ModsFileInfo.DiskInfo = Utils.GetUserSettings().SMARTChecker.Disks_info
	for i := range disks_info {
		var disk_info *ModsFileInfo.DiskInfo = &disks_info[i]
		if disk_info.Id == id {
			return disk_info
		}
	}

	return nil
}
