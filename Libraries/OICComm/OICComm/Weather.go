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

var weathers_GL []ModsFileInfo.Weather = nil

func getAllWeathers() {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_OICComm, 0, []byte("G_S|true|Weather"))
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_OICComm, 0)
	if comms_map == nil {
		return
	}

	var json_bytes []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	if err := Utils.FromJsonGENERAL(json_bytes, &weathers_GL); err != nil {
		return
	}
}

/*
GetWeatherLocationsList returns the weathers locations list separated by "|".

This function will BLOCK FOREVER if there's no Internet connection! Check first with Utils.IsCommunicatorConnectedSERVER().

-----------------------------------------------------------

– Returns:
  - the weathers locations list separated by "|"
*/
func GetWeatherLocationsList() string {
	getAllWeathers()

	var locs_list string = ""
	for _, weather := range weathers_GL {
		locs_list += weather.Location + "|"
	}
	if len(locs_list) > 0 {
		locs_list = locs_list[:len(locs_list)-1]
	}

	return locs_list
}

/*
GetWeather returns the weather for the specified location.

-----------------------------------------------------------

– Returns:
  - the weather or nil if the weather is not found
*/
func GetWeather(weather_location string) *ModsFileInfo.Weather {
	for i := 0; i < len(weathers_GL); i++ {
		var weather *ModsFileInfo.Weather = &weathers_GL[i]
		if weather.Location == weather_location {
			return weather
		}
	}

	return nil
}
