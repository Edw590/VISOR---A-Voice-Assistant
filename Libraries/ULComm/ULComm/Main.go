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

package ULComm

import (
	"Utils"
	"strconv"
	"strings"
)

/*
CreateDeviceInfo creates a DeviceInfo object with the given parameters.
 */
func CreateDeviceInfo(last_comm int64, last_time_used int64, airplane_mode_enabled bool, wifi_enabled bool,
		bluetooth_enabled bool, power_connected bool, battery_level int, screen_on bool, monitor_brightness int,
		wifi_networks string, bluetooth_devices string) *DeviceInfo {
	var wifi_networks_ret []WifiNetwork
	for _, network := range strings.Split(wifi_networks, "\x00") {
		if network == "" {
			continue
		}

		var network_info []string = strings.Split(network, "\x01")
		var wifi_network WifiNetwork = WifiNetwork{
			SSID: network_info[0],
			BSSID: network_info[1],
		}
		rssi, _ := strconv.Atoi(network_info[2])
		wifi_network.RSSI = rssi

		wifi_networks_ret = append(wifi_networks_ret, wifi_network)
	}

	var bluetooth_devices_ret []BluetoothDevice
	for _, device := range strings.Split(bluetooth_devices, "\x00") {
		if device == "" {
			continue
		}

		var device_info []string = strings.Split(device, "\x01")
		var bluetooth_device BluetoothDevice = BluetoothDevice{
			Name: device_info[0],
			Address: device_info[1],
		}
		rssi, _ := strconv.Atoi(device_info[2])
		bluetooth_device.RSSI = rssi

		bluetooth_devices_ret = append(bluetooth_devices_ret, bluetooth_device)
	}

	return &DeviceInfo{
		Device_id: Utils.PersonalConsts_GL.DEVICE_ID,
		Last_comm: last_comm,
		Last_time_used: last_time_used,
		System_state: SystemState{
			Connectivity_info: ConnectivityInfo{
				Airplane_mode_enabled: airplane_mode_enabled,
				Wifi_enabled: wifi_enabled,
				Bluetooth_enabled: bluetooth_enabled,
				Wifi_networks: wifi_networks_ret,
				Bluetooth_devices: bluetooth_devices_ret,
			},
			Battery_info: BatteryInfo{
				Level:           battery_level,
				Power_connected: power_connected,
			},
			Monitor_info: MonitorInfo{
				Screen_on: screen_on,
				Brightness: monitor_brightness,
			},
		},
	}
}
