/*******************************************************************************
 * Copyright 2023-2024 Edw590
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
	"OnlineInfoChk/OICNews"
	"OnlineInfoChk/OICWeather"
	"Utils"
)

/*
GetNews gets the news from the given page contents.

-----------------------------------------------------------

– Params:
  - page_contents – the page contents

– Returns:
  - the news separated by " ||| " and each news location separated by "\n"
 */
func GetNews() string {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, []byte("File|false|news.json"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GPTComm]
	var file_contents []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	var news_list []OICNews.News
	if err := Utils.FromJsonGENERAL(file_contents, &news_list); err != nil {
		return ""
	}

	var ret string = ""
	for _, news := range news_list {
		ret += news.Location + " ||| "
		for _, s := range news.News {
			ret += s + " ||| "
		}
		ret += "\n"
	}

	return ret
}

/*
GetWeather gets the weather from the given page contents.

Weather data in order:
  - Location
  - Temperature
  - Precipitation
  - Humidity
  - Wind
  - Status
  - Max_temp
  - Min_temp

-----------------------------------------------------------

– Params:
  - page_contents – the page contents

– Returns:
  - the weather separated by " ||| " and each weather location separated by "\n"
 */
func GetWeather() string {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, []byte("File|false|weather.json"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GPTComm]
	var file_contents []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	var weather []OICWeather.Weather
	if err := Utils.FromJsonGENERAL(file_contents, &weather); err != nil {
		return ""
	}

	var ret string = ""
	for _, w := range weather {
		ret += w.Location + " ||| " + w.Temperature + " ||| " + w.Precipitation + " ||| " + w.Humidity + " ||| " +
			w.Wind + " ||| " + w.Status + " ||| " + w.Max_temp + " ||| " + w.Min_temp + "\n"
	}

	return ret
}
