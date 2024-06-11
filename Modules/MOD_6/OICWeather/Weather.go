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

package OICWeather

import (
	"Utils"
	"github.com/tebeka/selenium"
)

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
UpdateWeather updates the weather for the given locations to a file called "weather.json".

-----------------------------------------------------------

– Params:
  - driver – the selenium web driver
  - locations – the locations to search for the weather

– Returns:
  - the error if any
*/
func UpdateWeather(driver selenium.WebDriver, locations []string) error {
	var weather []Weather = nil
	for _, location := range locations {
		if location == "" {
			continue
		}

		temperature, max_temp, min_temp, precipitation, humidity, wind, status, err := findWeather(driver, location)
		if err != nil {
			panic(err)
		}

		//log.Println("Current temperature in " + location + ": " + temperature + "ºC")
		//log.Println("Maximum temperature in " + location + ": " + max_temp + "ºC")
		//log.Println("Minimum temperature in " + location + ": " + min_temp + "ºC")
		//log.Println("Current precipitation in " + location + ": " + precipitation)
		//log.Println("Current humidity in " + location + ": " + humidity)
		//log.Println("Current wind in " + location + ": " + wind)
		//log.Println("Current status in " + location + ": " + status)
		//log.Println("")

		// write the info to an json struct
		weather = append(weather, Weather{
			Location:      location,
			Temperature:   temperature,
			Max_temp:      max_temp,
			Min_temp:      min_temp,
			Precipitation: precipitation,
			Humidity:      humidity,
			Wind:          wind,
			Status:        status,
		})
	}

	return Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, "weather.json").WriteTextFile(*Utils.ToJsonGENERAL(weather), false)
}

/*
findWeather finds the weather for the given location.

The output strings will be "ERROR" if an error occurs with that specific information.

-----------------------------------------------------------

– Params:
  - driver – the selenium web driver
  - location – the location to search for the weather

– Returns:
  - the error if any
  - the temperature
  - the precipitation
  - the humidity
  - the wind
  - the status
 */
func findWeather(driver selenium.WebDriver, location string) (string, string, string, string, string, string, string, error) {
	err := driver.Get("https://www.google.com/search?q=tempo " + location + "&ie=utf-8&oe=utf-8&hl=en")
	if err != nil {
		return "", "", "", "", "", "", "", err
	}

	temperature := "ERROR"
	element, err := driver.FindElement(selenium.ByCSSSelector, "#wob_tm")
	if err == nil {
		temperature, _ = element.Text()
	}

	max_temp := "ERROR"
	min_temp := "ERROR"
	element, err = driver.FindElement(selenium.ByCSSSelector, ".wob_df.wob_ds")
	if err == nil {
		elements, err := element.FindElements(selenium.ByCSSSelector, ".wob_t")
		if err == nil {
			if len(elements) > 0 {
				max_temp, _ = elements[0].Text()
			}
			if len(elements) >= 3 {
				min_temp, _ = elements[2].Text()
			}
		}
	}

	precipitation := "ERROR"
	element, err = driver.FindElement(selenium.ByCSSSelector, "#wob_pp")
	if err == nil {
		precipitation, _ = element.Text()
	}

	humidity := "ERROR"
	element, err = driver.FindElement(selenium.ByCSSSelector, "#wob_hm")
	if err == nil {
		humidity, _ = element.Text()
	}

	wind := "ERROR"
	element, err = driver.FindElement(selenium.ByCSSSelector, "#wob_ws")
	if err == nil {
		wind, _ = element.Text()
	}

	status := "ERROR"
	element, err = driver.FindElement(selenium.ByCSSSelector, "#wob_dc")
	if err == nil {
		status, _ = element.Text()
	}

	return temperature, max_temp, min_temp, precipitation, humidity, wind, status, nil
}
