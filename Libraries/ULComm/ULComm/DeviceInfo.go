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
	"log"
)

var prev_device_info DeviceInfo

type DeviceInfo struct {
	// Device_id is the unique identifier of the device
	Device_id string
	// Last_comm is the last time the device communicated with the server in Unix time
	Last_comm int64
	// Last_time_used is the last time the device was used in Unix time
	Last_time_used int64
	// System_state is the information about the system state of the device
	System_state SystemState
}

func (device_info *DeviceInfo) SendInfo() {
	if Utils.CompareSTRUCTS[DeviceInfo](*device_info, prev_device_info) {
		log.Println("No changes in device info")

		return
	}

	var message []byte = []byte("UserLocator|" + Utils.User_settings_GL.PersonalConsts.Device_ID + "|")
	message = append(message, Utils.CompressString(*Utils.ToJsonGENERAL(*device_info))...)
	Utils.QueueNoResponseMessageSERVER(message)

	log.Println("Sent device info")
	prev_device_info = *device_info
}

type SystemState struct {
	// Connectivity_info is the information about the connectivity of the device
	Connectivity_info ConnectivityInfo
	// Battery_info is the information about the battery of the device
	Battery_info BatteryInfo
	// Monitor_info is the information about the main monitor of the device
	Monitor_info MonitorInfo
	// Sound_info is the information about the sound of the device
	Sound_info SoundInfo
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
	Wifi_networks []ExtBeacon
	// Bluetooth_info is the information about the bluetooth devices the device has in range
	Bluetooth_devices []ExtBeacon
}

type ExtBeacon struct {
	// Name is the name of the device (bluetooth device name or SSID for Wi-Fi, for example)
	Name string
	// Address is the address of the device (MAC address for bluetooth devices, BSSID for Wi-Fi networks, for example)
	Address string
	// RSSI is the signal strength of the device in dBm
	RSSI int
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

type SoundInfo struct {
	// Volume is the volume of the sound in percentage
	Volume int
	// Muted is whether the sound is muted
	Muted bool
}
