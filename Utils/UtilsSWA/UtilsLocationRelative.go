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

package UtilsSWA

import (
	"math"
	"strconv"
	"strings"

	"Utils"
)

// Constants to standardize distances inside the app (like maximum distances. Why 4m and not 5m? Unless it's really
// important to be 4).
const ABSTR_DISTANCE_1 int32 = 1
const ABSTR_DISTANCE_5 int32 = 5
const ABSTR_DISTANCE_10 int32 = 10
const ABSTR_DISTANCE_50 int32 = 50
const ABSTR_DISTANCE_100 int32 = 100
const ABSTR_DISTANCE_INFINITY int32 = 999
/*
GetAbstrDistanceRSSILOCATIONRELATIVE returns an abstract distance between the 2 devices using the calculated real
distance.

These distances are solely to standardize values on the app.

Also the distance returned by RSSI calculations may be misleading, so this also discretizes the distances to hopefully
more precise values. For example, _ABSTR_DISTANCE_1 means the devices are right next to each other. _ABSTR_DISTANCE_5
means they're near. A wall or 2 away at most maybe.

---CONSTANTS---
  - ABSTR_DISTANCE_1 – for the returning value: if the distance is <= 1 meter
  - ABSTR_DISTANCE_5 – for the returning value: else if the distance is <= 5 meters
  - ABSTR_DISTANCE_10 – for the returning value: else if the distance is <= 10 meters
  - ABSTR_DISTANCE_50 – for the returning value: else if the distance is <= 50 meters
  - ABSTR_DISTANCE_100 – for the returning value: else if the distance is <= 100 meters
  - ABSTR_DISTANCE_INFINITY – for the returning value: else if the distance is > 100 meters
---CONSTANTS---

-----------------------------------------------------------

– Params:
  - real_distance – same as returned by getRealDistanceRSSI(int, int)

– Returns:
  - one of the constants
 */
func GetAbstrDistanceRSSILOCATIONRELATIVE(real_distance int32) int32 {
	if real_distance <= 1 {
		return ABSTR_DISTANCE_1
	} else if real_distance <= 5 {
		return ABSTR_DISTANCE_5
	} else if real_distance <= 10 {
		return ABSTR_DISTANCE_10
	} else if real_distance <= 50 {
		return ABSTR_DISTANCE_50
	} else if real_distance <= 100 {
		return ABSTR_DISTANCE_100
	} else {
		return ABSTR_DISTANCE_INFINITY
	}
}

/* For now I don't see use for the tx_power parameter, so I'm passing always the default one. Some other time might
implement BLE device detection, on which the value is useful since some devices have a value (not BV9500 it
seems).*/
const DEFAULT_TX_POWER int32 = -60
/*
GetRealDistanceRSSILOCATIONRELATIVE gets the distance in meters between 2 devices from the transmission signal strength
between the 2 devices (RSSI) and the transmission power value at 1 meter for the current device.

-----------------------------------------------------------

– Params:
  - rssi – the RSSI between both devices
  - tx_power – the expected transmission power of the current device for 1 meter distance from another device

– Returns:
  - the distance between the 2 devices rounded to the nearest integer
 */
func GetRealDistanceRSSILOCATIONRELATIVE(rssi int32, tx_power int32) int32 {
	/*
	Copied from: Dong, Q., & Dargie, W. (2012). Evaluation of the reliability of RSSI for indoor localization. 2012
	International Conference on Wireless Communications in Underground and Confined Areas, 1-6.
	Link for IEEE: https://ieeexplore.ieee.org/stamp/stamp.jsp?tp=&arnumber=6402492.

	"RSSI = -10*n*log10(d)-A (4)
	In Equation 4, RSSI is the radio signal strength indicator in dBm, n is the signal propagation constant or
	exponent [, also called Path Loss Exponent, which is dimensionless], d is the relative distance between the
	communicating nodes [in meters], and A is a reference received signal strength in dBm (the RSSI value measured
	when the separation distance between the receiver and the transmitter is one meter)."

	(The formula seems to be from Texas Instruments on the CC2420 chip.)

	So getting d from there, we have this below. Should be RSSI+A, but that term is supposed to be 0 if both values
	are equal, with with a +, that won't happen, so I saw people having +A on the main formula, and then it results
	in:
	d = 10^(-(RSSI-A)/10n)
	*/

	// "The txPower field is supposed to be the expected signal strength in dBm at 1 meter. Every Android phone
	// model has a slightly different BLE transmitter output power, typically measuring from -50 dBm to -70 dBm at
	// 1 meter." - so I'll go for -60 dBm as the default txPower (middle value).

	/*
	According to Wikipedia (https://en.wikipedia.org/wiki/Log-distance_path_loss_model, 2021-07-29):
	----------------
	                         Empirical coefficient values for indoor propagation

	Empirical measurements of coefficients gamma and sigma in dB have shown the following values for a number of
	indoor wave propagation cases.

	Building Type                  Frequency of Transmission    gamma   sigma [dB]
	Vacuum, infinite space                                       2.0    0
	Retail store                         914 MHz                 2.2    8.7
	Grocery store                        914 MHz                 1.8    5.2
	Office with hard partition           1.5 GHz                 3.0    7
	Office with soft partition           900 MHz                 2.4    9.6
	Office with soft partition           1.9 GHz                 2.6    14.1
	Textile or chemical                  1.3 GHz                 2.0    3.0
	Textile or chemical                    4 GHz                 2.1    7.0, 9.7
	Office                                60 GHz                 2.2    3.92
	Commercial                            60 GHz                 1.7    7.9
	----------------
	[Note: our n here is gamma on Wikipedia.]

	So using the Office with soft-partition at 1.9 GHz with 2.6, knowing on 900 MHz it was 2.4, and also that
	Bluetooth operates at 2.4-2.5 GHz, 2.7-2.8 seems a good choice of value at home (so 2.75 - average).
	But that's only with objects in the way. If there's a straight path, 2.1 doesn't seem a bad value (there's the
	phone or computer case). So if the distance is very small, 2.1 will be used (can't be much in the way if the
	distance is like 1 meter with 2.75 as value).
	Another possibility is the object being further away. In case, then other objects might be in its path. Then we
	rise the value, depending on the distance. For example, 5 meters or more could have stuff in the way already.
	10-20 meters and even more stuff (not counting open field here... - can't guess this kind of thing without GPS).
	Else, 2.75 will be used.


	From https://en.wikipedia.org/wiki/Path_loss:
	"In the study of wireless communications, path loss can be represented by the path loss exponent, whose value
	is normally in the range of 2 to 4 (where 2 is for propagation in free space, 4 is for relatively lossy
	environments and for the case of full specular reflection from the earth surface—the so-called flat earth
	model). In some environments, such as buildings, stadiums and other indoor environments, the path loss exponent
	can reach values in the range of 4 to 6. On the other hand, a tunnel may act as a waveguide, resulting in a path
	loss exponent less than 2."

	So, 5 meters or more, could be 2.75 + (4+6)/2 = 3.875 = 3.88. For more than 20 meters, why not 6.0. On more than
	50 meters, I'll use 8.0. A person said n could reach from 2 to 8 in some cases. So I'll use that 8 value here.
	*/

	// This will be used if 2.0 < resulting distance < 7.5.
	var def_ret_value float64 = math.Pow(10.0, float64(-(rssi - tx_power)) / (10.0 * 2.75))
	var ret_value float64 = def_ret_value

	if def_ret_value <= 2.0 {
		// 2 meters here because it's near enough to 1 meter - can't be 1 meter because the idea of the if statements
		// is to correct incorrect calculation, so it's to assume if it gives about 2 meters, it might mean it's less
		// than that actually (like with my phone and tablet, a meter away from each other and this returning more
		// than 1 meter.
		ret_value = math.Pow(10.0, float64(-(rssi - tx_power)) / (10.0 * 2.1))
	} else if def_ret_value >= 50.0 {
		ret_value = math.Pow(10.0, float64(-(rssi - tx_power)) / (10.0 * 8.0))
	} else if def_ret_value >= 20.0 {
		ret_value = math.Pow(10.0, float64(-(rssi - tx_power)) / (10.0 * 6.0))
	} else if def_ret_value >= 7.5 { // 7.5 meters because it's near enough to 5 meters and still away from 20 meters
		ret_value = math.Pow(10.0, float64(-(rssi - tx_power)) / (10.0 * 3.88))
	}

	return int32(math.Round(ret_value))
}

/*
GetAveragePingRTTLOCATIONRELATIVE gets the average round-trip time (RTT) between the current device and the specified IP
address.

This function pings the specified IP address 50 times, eliminates outlier time values, and finally calculates the
average of the resulting values (the round-trip time, RTT).

ATTENTION: this function takes about a minute to return!

It can also return awful results (like 200m when it's in fact 3-4m). That's because it depends on many factors, not
just the speed of light (it doesn't only measure the time between the transmitter and the receptor - electronics delays
are included, and other things).

Though it can too return decent results, like 3-4m when it's 3-4m indeed. But the nonsense cases must be filtered out.

-----------------------------------------------------------

– Params:
  - ip – the IP address to ping

– Returns:
  - the average RTT, or -1 if any error occurred
 */
func GetAveragePingRTTLOCATIONRELATIVE(ip string) (float64, error) {
	const NUM_PACKETS = 50 // 50 packets, each with 0.5 seconds delay, so 25 seconds of waiting time
	// 248 + 8 header = 256 bytes each packet
	command_str := "ping{{EXE}} -c " + strconv.Itoa(NUM_PACKETS) + " -i 0.5 -n -s 248 -t 1 -v " + ip
	output, err := Utils.ExecCmdSHELL([]string{command_str})
	if err != nil {
		return -1.0, err
	}

	output_lines := strings.Split(output.Stdout_str, "\n")

	// Here it gets the time values (excluding the duplicated ones)
	var time_values []float64
	for _, line := range output_lines {
		if !strings.Contains(line, " (DUP!)") && strings.Contains(line, "time=") && strings.Contains(line, " ms") {
			time, _ := strconv.ParseFloat(strings.Split(strings.Split(line, "time=")[1], " ")[0], 64)
			time_values = append(time_values, time)
		}
	}

	const accuracy_parameter = 2.0 // Accuracy parameter to use with UtilsMath.IsOutlier().

	// It will calculate the sum and sum of squares of the first 5 elements of the time_values list (5 seemed a good
	// number), and then will check if any of those first 5 elements is an outlier inside that 5-elements list, and
	// if any is an outlier, it will be removed from the time_values list.
	// If the time values are [2, 20, 18, 16, 34, 129, 21, 23], the average of the first values will still be near
	// the supposed value (between 16 and 20). If I put 4, it will come down to 14 (wrong). If it's less, even
	// worse will get. More than that, could catch the 129, which I don't see appearing in the first 5 elements on
	// some tests.
	const first_n_elements = 5
	var sum float64 = 0.0
	var sum_squares float64 = 0.0
	for i := 0; i < first_n_elements; i++ {
		value := time_values[i]
		sum += value
		sum_squares += value * value
	}

	var any_outlier = false
	for i := 0; i < first_n_elements; i++ {
		if IsOutlierMATH(time_values[i], sum, sum_squares, 5, accuracy_parameter) {
			any_outlier = true
			time_values = append(time_values[:i], time_values[i+1:]...)
			i--
		}
	}

	var summed_elements = first_n_elements
	if any_outlier {
		// Reset the variables since now it will calculate again the new sum and sum of squares (in case there was an
		// outlier and it was removed from the original list). This will be used as the starting parameters for the
		// outlier check.
		sum = 0.0
		sum_squares = 0.0
		for i := 0; i < first_n_elements; i++ {
			value := time_values[i]
			sum += value
			sum_squares += value * value
		}
	}

	// With the remaining elements, the function will check if any of them is an outlier, comparing to the values
	// got from the first 5 elements (those 5 decide the fate of the list xD).
	for i := summed_elements; i < NUM_PACKETS - summed_elements; i++ {
		value := time_values[i]
		if !IsOutlierMATH(value, sum, sum_squares, int32(summed_elements), accuracy_parameter) {
			summed_elements++
			sum += value
			sum_squares += value * value
		}
	}

	// The mean of the time values will be calculated using the first 5 elements (excluding any outliers) and the
	// rest of the list (excluding any outliers).
	return sum / float64(summed_elements), nil
}
