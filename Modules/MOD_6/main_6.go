/*******************************************************************************
 * Copyright 2023-2023 Edw590
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

package MOD_6

import (
	"OnlineInfoChk/OICNews"
	"OnlineInfoChk/OICWeather"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"

	"Utils"
)

// Online Information Checker //

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGIModSpecInfo]
)

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

func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](Utils.NUM_MOD_OnlineInfoChk, realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			////////////////////////////////
			// Prepare the browser
			service, err := selenium.NewChromeDriverService("/usr/bin/chromedriver", 4444)
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

			bypassGoogleCookies(driver)

			OICWeather.UpdateWeather(driver, modUserInfo.Temp_locs)
			OICNews.UpdateNews(driver, modUserInfo.News_locs)

			driver.Quit()
			service.Stop()


			if Utils.WaitWithStop(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}
