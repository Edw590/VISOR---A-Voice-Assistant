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
	"log"
	"math"
	"strconv"
	"strings"
)

type OpenMeteoWeather struct {
	Current_units Current_units `json:"current_units"`
	Current       Current       `json:"current"`
	Daily_units   Daily_units   `json:"daily_units"`
	Daily         Daily         `json:"daily"`
}

type Current_units struct {
	Temperature_2m string `json:"temperature_2m"`
}

type Daily_units struct {
	Temperature_2m_max             string `json:"temperature_2m_max"`
	Temperature_2m_min             string `json:"temperature_2m_min"`
	Precipitation_probability_mean string `json:"precipitation_probability_mean"`
	Relative_humidity_2m_mean      string `json:"relative_humidity_2m_mean"`
	Wind_speed_10m_mean            string `json:"wind_speed_10m_mean"`
}

type Current struct {
	Temperature_2m float32 `json:"temperature_2m"`
}

type Daily struct {
	Temperature_2m_max             []float32 `json:"temperature_2m_max"`
	Temperature_2m_min             []float32 `json:"temperature_2m_min"`
	Precipitation_probability_mean []float32 `json:"precipitation_probability_mean"`
	Relative_humidity_2m_mean      []float32 `json:"relative_humidity_2m_mean"`
	Wind_speed_10m_mean            []float32 `json:"wind_speed_10m_mean"`
}

/*
UpdateWeather updates the weather for the given locations.

-----------------------------------------------------------

– Params:
  - driver – the selenium web driver
  - locations – the locations to search for the weather

– Returns:
  - the error if any
*/
func UpdateWeather(locations_info []string) []ModsFileInfo.Weather {
	var weathers []ModsFileInfo.Weather = nil
	for _, location_info := range locations_info {
		if location_info == "" {
			continue
		}

		var location_parts []string = strings.Split(location_info, ": ")
		if len(location_parts) != 2 {
			log.Println("Invalid location format: " + location_info)

			continue
		}

		var location string = location_parts[0]
		var location_coords []string = strings.Split(location_parts[1], ", ")
		if len(location_coords) != 2 {
			log.Println("Invalid location coordinates format: " + location_info)

			continue
		}

		latitude, err := strconv.ParseFloat(location_coords[0], 32)
		if err != nil {
			log.Println("Invalid location latitude format: " + location_info)

			continue
		}
		longitude, err := strconv.ParseFloat(location_coords[1], 32)
		if err != nil {
			log.Println("Invalid location longitude format: " + location_info)

			continue
		}

		weather, err := findWeather(location, float32(latitude), float32(longitude))
		if err != nil {
			log.Println("Error getting weather for " + location_info + " --> " + Utils.GetFullErrorMsgGENERAL(err))

			continue
		}

		//log.Println("Current temperature in " + location_info + ": " + weather.Temperature + "ºC")
		//log.Println("Maximum temperature in " + location_info + ": " + weather.Max_temp + "ºC")
		//log.Println("Minimum temperature in " + location_info + ": " + weather.Min_temp + "ºC")
		//log.Println("Current precipitation in " + location_info + ": " + weather.Precipitation)
		//log.Println("Current humidity in " + location_info + ": " + weather.Humidity)
		//log.Println("Current wind in " + location_info + ": " + weather.Wind)
		//log.Println("Current status in " + location_info + ": " + weather.Status)
		//log.Println("")

		weathers = append(weathers, weather)
	}

	return weathers
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
func findWeather(location string, latitude float32, longitude float32) (ModsFileInfo.Weather, error) {
	var latitude_str string = strconv.FormatFloat(float64(latitude), 'f', -1, 32)
	var longitude_str string = strconv.FormatFloat(float64(longitude), 'f', -1, 32)
	source, err := Utils.MakeGetRequest("https://api.open-meteo.com/v1/forecast?" +
		"latitude=" + latitude_str + "&longitude=" + longitude_str +
		"&current=temperature_2m" +
		"&daily=temperature_2m_max,temperature_2m_min,precipitation_probability_mean,relative_humidity_2m_mean,wind_speed_10m_mean" +
		"&forecast_days=1" +
		"&timeformat=unixtime")
	if err != nil {
		return ModsFileInfo.Weather{}, err
	}

	var weather OpenMeteoWeather
	err = Utils.FromJsonGENERAL(source, &weather)
	if err != nil {
		return ModsFileInfo.Weather{}, err
	}

	// Don't add the units to the temperature. Useful in case we want VISOR to say just "degrees" and not "degrees
	// Celsius".
	var temperature string = float32ToIntToString(weather.Current.Temperature_2m)
	var max_temp string = float32ToIntToString(weather.Daily.Temperature_2m_max[0])
	var min_temp string = float32ToIntToString(weather.Daily.Temperature_2m_min[0])
	var precipitation string = float32ToIntToString(weather.Daily.Precipitation_probability_mean[0]) +
		weather.Daily_units.Precipitation_probability_mean
	var humidity string = float32ToIntToString(weather.Daily.Relative_humidity_2m_mean[0]) +
		weather.Daily_units.Relative_humidity_2m_mean
	var wind string = float32ToIntToString(weather.Daily.Wind_speed_10m_mean[0]) +
		weather.Daily_units.Wind_speed_10m_mean

	var status string = ""
	source, err = Utils.MakeGetRequest("https://wttr.in/" + location + "?lang=en&format=%C")
	if err != nil {
		status = "ERROR"
	} else {
		status = strings.TrimSpace(string(source))
	}

	return ModsFileInfo.Weather{
		Location:      location,
		Temperature:   temperature,
		Max_temp:      max_temp,
		Min_temp:      min_temp,
		Precipitation: precipitation,
		Humidity:      humidity,
		Wind:          wind,
		Status:        status,
	}, nil
}

func float32ToIntToString(value float32) string {
	var int_value int = int(math.Round(float64(value)))

	return strconv.Itoa(int_value)
}
