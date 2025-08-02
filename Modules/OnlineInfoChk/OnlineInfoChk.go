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
	"fmt"
	"github.com/Edw590/go-wolfram"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	gowiki "github.com/trietmn/go-wiki"
)

// Online Information Checker //

const _TIME_SLEEP_S int = 15*60

////////////////////////////////////
// Copilot ideas for the module:

// This is the main function of the OIG library.
// It does not do anything.
// The OIG library is a collection of other libraries and modules.
// It is used to group all the libraries and modules into a single package.
// The libraries and modules are used to perform different tasks.
// The libraries and modules are used to perform tasks such as:
// - Getting the weather
// - Getting the news
// - Getting the time
// - Getting the date
// - Getting the currency exchange rate
// - Getting the stock price
// - Getting the cryptocurrency price
// - Getting the sports scores
// - Getting the traffic information
// - Getting the flight information
// - Getting the hotel information
// - Getting the restaurant information
// - Getting the movie information
// - Getting the music information
// - Getting the book information
// - Getting the game information
// - Getting the TV show information
// - Getting the radio show information
// - Getting the podcast information
// - Getting the audiobook information
// - Getting the e-book information
// - Getting the magazine information
// - Getting the newspaper information
// - Getting the blog information
// - Getting the vlog information
// - Getting the social media information
// - Getting the forum information
// - Getting the chat information
// - Getting the email information
// - Getting the SMS information
// - Getting the MMS information
// - Getting the phone call information
// - Getting the video call information
// - Getting the video conference information
// - Getting the audio call information
// - Getting the audio conference information
// - Getting the video information
// - Getting the audio information
// - Getting the image information
// - Getting the document information
// - Getting the file information
// - Getting the folder information
// - Getting the computer information
// - Getting the smartphone information
// - Getting the tablet information
// - Getting the smartwatch information
// - Getting the smartglasses information
// - Getting the smartclothes information
// - Getting the smartshoes information
// - Getting the smartjewelry information
////////////////////////////////////

// Firefox: log.Sprintf("http://localhost:%d", port)
// Chrome: log.Sprintf("http://127.0.0.1:%d/wd/hub", port) - default

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	for {
		////////////////////////////////
		// Prepare the browser
		service, err := selenium.NewChromeDriverService("chromedriver", 4444)
		if err != nil {
			return
		}

		caps := selenium.Capabilities{}
		caps.AddChrome(chrome.Capabilities{Args: []string{
			// Keep all of these here. They ensure that the driver works properly.
			"start-maximized",
			"enable-automation",
			"--headless",
			"--disable-browser-side-navigation",
			"--disable-dev-shm-usage",
			"--disable-extensions",
			"--disable-gpu",
			"--disable-infobars",
			"--dns-prefetch-disable",
			"--incognito",
			"--no-sandbox",
			"--remote-debugging-port=9222",
			"--window-size=1920,1080",
		}})

		driver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", 4444))
		if err != nil {
			return
		}
		////////////////////////////////

		_ = bypassGoogleCookies(driver)

		getModGenSettings().Weather = UpdateWeather(getModUserInfo().Temp_locs)
		getModGenSettings().News = UpdateNews(driver, getModUserInfo().News_locs)

		_ = driver.Quit()
		_ = service.Stop()


		if Utils.WaitWithStopDATETIME(module_stop, _TIME_SLEEP_S) {
			return
		}
	}
}

/*
RetrieveWolframAlpha retrieves the information from the given query using Wolfram Alpha.

-----------------------------------------------------------

– Params:
  - query – the query to search for

– Returns:
  - the information retrieved
  - whether the information is a direct result or is a combination of results (maybe the LLM should be called to
	summarize it)
*/
func RetrieveWolframAlpha(query string) (string, bool) {
	//Initialize a new client
	c := &wolfram.Client{
		AppID: Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.WolframAlpha_AppID,
	}

	//Get a result without additional parameters
	res, err := c.GetQueryResult(query, nil)
	if err != nil {
		Utils.LogLnError(err)

		return "ERROR", true
	}

	var query_results wolfram.QueryResult = res.QueryResult

	if len(query_results.Pods) < 2 {
		return "ERROR", true
	}

	if query_results.Pods[1].Title == "Result" {
		if len(query_results.Pods[1].SubPods) < 1 {
			return "ERROR", true
		}

		return query_results.Pods[1].SubPods[0].Plaintext, true
	} else {
		// Iterate through the pods and subpods and get the plaintext info from each
		var result string = ""
		for _, i := range res.QueryResult.Pods {
			result += i.Title + ": "

			for _, j := range i.SubPods {
				result += j.Plaintext + " / "
			}
		}
		result = result[:len(result)-3]

		return result, false
	}
}

/*
RetrieveWikipedia retrieves the information from the given query using Wikipedia.

-----------------------------------------------------------

– Params:
  - query – the query to search for

– Returns:
  - the information retrieved
*/
func RetrieveWikipedia(query string) string {
	page, err := gowiki.GetPage(query, -1, true, true)
	if err != nil {
		Utils.LogLnError(err)

		return "ERROR"
	}

	content, err := page.GetContent()
	if err != nil {
		Utils.LogLnError(err)

		return "ERROR"
	}

	return content
}

func getModGenSettings() *ModsFileInfo.Mod6GenInfo {
	return &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_6
}

func getModUserInfo() *ModsFileInfo.Mod6UserInfo {
	return &Utils.GetUserSettings(Utils.LOCK_UNLOCK).OnlineInfoChk
}
