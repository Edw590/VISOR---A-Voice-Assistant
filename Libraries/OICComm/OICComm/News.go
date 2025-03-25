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

package OICComm

import (
	"Utils"
	"Utils/ModsFileInfo"
)

var news_locs_GL []ModsFileInfo.News = nil

func getAllNews() {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_OICComm, 1, []byte("G_S|true|News")) {
		return
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_OICComm, 1)
	if comms_map == nil {
		return
	}

	var json_bytes []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	if err := Utils.FromJsonGENERAL(json_bytes, &news_locs_GL); err != nil {
		return
	}
}

/*
GetNewsLocationsList returns the news locations list separated by "|".

-----------------------------------------------------------

– Returns:
  - the news locations list separated by "|"
 */
func GetNewsLocationsList() string {
	getAllNews()

	var locs_list string = ""
	for _, news_loc := range news_locs_GL {
		locs_list += news_loc.Location + "|"
	}
	if len(locs_list) > 0 {
		locs_list = locs_list[:len(locs_list)-1]
	}

	return locs_list
}

/*
GetNews returns the news for the specified location.

-----------------------------------------------------------

– Returns:
  - the news or nil if the news are not found
 */
func GetNews(news_location string) *ModsFileInfo.News {
	for i := range news_locs_GL {
		var news_loc *ModsFileInfo.News = &news_locs_GL[i]
		if news_loc.Location == news_location {
			return news_loc
		}
	}

	return nil
}
