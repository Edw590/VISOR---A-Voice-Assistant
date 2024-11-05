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
*/
func AddLocationLOCATIONS(type_ string, name string, address string, last_detection_s int64, max_distance_m int32,
						  location string) {
	var locs_info []ModsFileInfo.LocInfo = Utils.User_settings_GL.UserLocator.Locs_info
	var loc_id int32 = 1
	for i := 0; i < len(locs_info); i++ {
		if locs_info[i].Id == loc_id {
			loc_id++
		}
	}

	// Add the location to the user settings
	Utils.User_settings_GL.UserLocator.Locs_info = append(Utils.User_settings_GL.UserLocator.Locs_info, ModsFileInfo.LocInfo{
		Id:               loc_id,
		Type:             type_,
		Name:             name,
		Address:          address,
		Last_detection_s: last_detection_s,
		Max_distance_m:   max_distance_m,
		Location:         location,
	})

	sort.Slice(locs_info, func(i, j int) bool {
		return locs_info[i].Location < locs_info[j].Location
	})
}

/*
RemoveLocationLOCATIONS removes a location from the user settings.

-----------------------------------------------------------

– Params:
  - loc_id – the ID of the location to be removed
*/
func RemoveLocationLOCATIONS(loc_id int32) {
	var locs_info []ModsFileInfo.LocInfo = Utils.User_settings_GL.UserLocator.Locs_info
	for i := 0; i < len(locs_info); i++ {
		if locs_info[i].Id == loc_id {
			Utils.DelElemSLICES(&Utils.User_settings_GL.UserLocator.Locs_info, i)

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
	var ids string
	for _, loc_info := range Utils.User_settings_GL.UserLocator.Locs_info {
		ids += strconv.Itoa(int(loc_info.Id)) + "|"
	}

	return ids
}

/*
GetLOCATIONS returns a location by its ID.

-----------------------------------------------------------

– Params:
  - loc_id – the location ID

– Returns:
  - the location or nil if the location was not found
 */
func GetLocationLOCATIONS(loc_id int32) *ModsFileInfo.LocInfo {
	for _, loc_info := range Utils.User_settings_GL.UserLocator.Locs_info {
		if loc_info.Id == loc_id {
			return &loc_info
		}
	}

	return nil
}
