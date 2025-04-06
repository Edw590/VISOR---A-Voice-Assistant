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

package OnlineInfoChk

import (
	"Utils"
	"Utils/ModsFileInfo"
	"github.com/tebeka/selenium"
	"log"
)

/*
UpdateNews updates the news for the given locations.

-----------------------------------------------------------

– Params:
  - driver – the selenium web driver
  - news_data – the news to search for in the following format: "[news search string],[location]" (example: "notícias
    Portugal,Portugal")

– Returns:
  - the error if any
*/
func UpdateNews(driver selenium.WebDriver, news_locs []string) []ModsFileInfo.News {
	var news []ModsFileInfo.News = nil
	for _, news_loc := range news_locs {
		texts, err := findNews(driver, news_loc + " news")
		if err != nil {
			log.Println("Error while searching for news for " + news_loc + ": " + Utils.GetFullErrorMsgGENERAL(err))

			continue
		}
		//log.Println("Current news in " + news_loc + ":")
		//for _, text := range texts {
		//	log.Println(text)
		//}
		//log.Println("")

		// write the info to an json struct
		news = append(news, ModsFileInfo.News{
			Location: news_loc,
			News:     texts,
		})
	}

	return news
}

/*
findNews searches for news on Google.

-----------------------------------------------------------

– Params:
  - driver – the selenium web driver
  - news_str – the news to search for

– Returns:
  - the error if any
  - the news found
 */
func findNews(driver selenium.WebDriver, news_str string) ([]string, error) {
	err := driver.Get("https://www.google.com/search?q=" + news_str + "&tbm=nws&ie=utf-8&oe=utf-8&hl=en")
	if err != nil {
		return nil, err
	}

	elements, err := driver.FindElements(selenium.ByCSSSelector, "div.n0jPhd.ynAwRc.MBeuO.nDgy9d")
	if err != nil {
		return nil, err
	}
	var texts []string = nil
	var length int = len(elements)
	if length > 10 {
		length = 10
	}
	for i := 0; i < length; i++ {
		element := elements[i]
		text, err := element.Text()
		if err != nil {
			continue
		}
		texts = append(texts, text)
	}

	return texts, nil
}
