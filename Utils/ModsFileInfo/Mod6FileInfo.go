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

package ModsFileInfo

// Mod6GenInfo is the format of the custom generated information about this specific module.
type Mod6GenInfo struct {
	// News is the news information
	News []News
	// Weather is the weather information
	Weather []Weather
}

// News is the format of the news information.
type News struct {
	Location string
	News []string
}

// Weather is the format of the weather information.
type Weather struct {
	Location string
	Temperature string
	Max_temp string
	Min_temp string
	Precipitation string
	Humidity string
	Wind   string
	Status string
}

/*
GetNewsList returns the news list separated by "\n".

-----------------------------------------------------------

– Returns:
  - the news list separated by "\n"
 */
func (news *News) GetNewsList() string {
	var news_list string = ""
	for _, news_item := range news.News {
		news_list += news_item + "\n"
	}
	if len(news_list) > 0 {
		news_list = news_list[:len(news_list)-1]
	}

	return news_list
}

///////////////////////////////////////////////////////////////////////////////

// Mod6UserInfo is the format of the custom information file about this specific module.
type Mod6UserInfo struct {
	// Temp_locs is the locations to get the weather from
	Temp_locs   []string
	// News_locs is the locations to get the news from
	News_locs []string
}
