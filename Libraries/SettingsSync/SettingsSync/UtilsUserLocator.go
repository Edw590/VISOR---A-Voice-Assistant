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
	"strconv"
)

/*
AddLocationLOCATIONS adds a location to the user settings.

-----------------------------------------------------------

– Params:
  - type_ – the type of the location "detector"
  - name – the name of the detection
  - address – the address of the location
  - last_detection_s – the last time the location was detected in Unix time
  - max_distance_m – the maximum distance in meters from the location
  - location – the location of the user

– Returns:
  - the ID of the location
*/
func AddLocationLOCATIONS(enabled bool, type_ string, name string, address string, last_detection_s int64,
						  max_distance_m int32, location string) int32 {
	var locs_info *[]ModsFileInfo.LocInfo = &Utils.GetUserSettings().UserLocator.Locs_info
	var id int32 = 1
	for i := 0; i < len(*locs_info); i++ {
		if (*locs_info)[i].Id == id {
			id++
			i = -1
		}
	}

	// Add the location to the user settings
	*locs_info = append(*locs_info, ModsFileInfo.LocInfo{
		Id:               id,
		Enabled:          enabled,
		Type:             type_,
		Name:             name,
		Address:          address,
		Last_detection_s: last_detection_s,
		Max_distance_m:   max_distance_m,
		Location:         location,
	})

	sort.SliceStable(*locs_info, func(i, j int) bool {
		return (*locs_info)[i].Location < (*locs_info)[j].Location
	})

	return id
}

/*
RemoveLocationLOCATIONS removes a location from the user settings.

-----------------------------------------------------------

– Params:
  - loc_id – the ID of the location to be removed
*/
func RemoveLocationLOCATIONS(id int32) {
	var locs_info *[]ModsFileInfo.LocInfo = &Utils.GetUserSettings().UserLocator.Locs_info
	for i := range *locs_info {
		if (*locs_info)[i].Id == id {
			Utils.DelElemSLICES(locs_info, i)

			break
		}
	}
}

/*
GetIdsListLOCATIONS returns a list of all locations' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all locations' IDs separated by "|"
 */
func GetIdsListLOCATIONS() string {
	var ids_list string
	for _, loc_info := range Utils.GetUserSettings().UserLocator.Locs_info {
		ids_list += strconv.Itoa(int(loc_info.Id)) + "|"
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetLOCATIONS returns a location by its ID.

-----------------------------------------------------------

– Params:
  - loc_id – the location ID

– Returns:
  - the location or nil if the location was not found
 */
func GetLocationLOCATIONS(id int32) *ModsFileInfo.LocInfo {
	var locs_info []ModsFileInfo.LocInfo = Utils.GetUserSettings().UserLocator.Locs_info
	for i := range locs_info {
		var loc_info *ModsFileInfo.LocInfo = &locs_info[i]
		if loc_info.Id == id {
			return loc_info
		}
	}

	return nil
}
