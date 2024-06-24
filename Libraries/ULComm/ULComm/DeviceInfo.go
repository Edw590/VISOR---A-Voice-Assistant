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

import "Utils"

const TYPE_DEVICE_PHONE string = "phone"
const TYPE_DEVICE_TABLET string = "tablet"
const TYPE_DEVICE_LAPTOP string = "laptop"
const TYPE_DEVICE_DESKTOP string = "desktop"

type DeviceInfo struct {
	// Device_id is the unique identifier of the device
	Device_id string
	// Device_type is the type of the device
	Device_type string
	// Last_comm is the last time the device communicated with the server in Unix time
	Last_comm int64
	// System_state is the information about the system state of the device
	System_state SystemState
}

func (device_info *DeviceInfo) SendInfo() error {
	return Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
		Type:  "UserLocator",
		Text1: Utils.PersonalConsts_GL.DEVICE_ID,
		Text2: *Utils.ToJsonGENERAL(device_info),
	})
}

type SystemState struct {
	// Connectivity_info is the information about the connectivity of the device
	Connectivity_info ConnectivityInfo
	// Battery_info is the information about the battery of the device
	Battery_info BatteryInfo
	// Monitor_info is the information about the main monitor of the device
	Monitor_info MonitorInfo
}

type ConnectivityInfo struct {
	// Airplane_mode_enabled is whether the airplane mode is enabled
	Airplane_mode_enabled bool
	// Wifi_enabled is whether the Wi-Fi is enabled
	Wifi_enabled bool
	// Bluetooth_enabled is whether the bluetooth is enabled
	Bluetooth_enabled bool
	// Mobile_data_enabled is whether the mobile data is enabled
	Mobile_data_enabled bool
	// Wifi_info is the information about the Wi-Fi network the device is connected to
	Wifi_networks []WifiNetwork
	// Bluetooth_info is the information about the bluetooth devices the device has in range
	Bluetooth_devices []BluetoothDevice
}

type WifiNetwork struct {
	// SSID is the name of the Wi-Fi network
	SSID string
	// BSSID is the address of the Wi-Fi network in the format XX:XX:XX:XX:XX:XX
	BSSID string
	// Signal is the signal strength of the Wi-Fi network in dBm
	Signal int
}

type BluetoothDevice struct {
	// Name is the name of the bluetooth device
	Name string
	// Address is the address of the bluetooth device in the format XX:XX:XX:XX:XX:XX
	Address string
	// Signal is the signal strength of the bluetooth device in dBm
	Signal int
}

type BatteryInfo struct {
	// Battery_level is the battery level of the device in percentage
	Level int
	// Power_connected is whether the device is connected to power
	Power_connected bool
}

type MonitorInfo struct {
	// Screen_on is whether the screen is on
	Screen_on bool
	// Brightness is the brightness of the screen in percentage
	Brightness int
}
