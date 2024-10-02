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

package main

import (
	"ULComm/ULComm"
	"Utils"
	"log"
	"time"
)

func main() {
	Utils.LoadUserSettings(false)
	Utils.InitializeCommsChannels()

	go func() {
		Utils.StartCommunicatorSERVER()
	}()
	time.Sleep(2 * time.Second)

	var device_info ULComm.DeviceInfo = ULComm.DeviceInfo{
		Device_id:    "Test Device",
		Last_comm:    time.Now().Unix(),
		System_state: ULComm.SystemState{
			Connectivity_info: ULComm.ConnectivityInfo{
				Airplane_mode_enabled: false,
				Wifi_enabled:          true,
				Bluetooth_enabled:     false,
				Mobile_data_enabled:   true,
				Wifi_networks:         nil,
				Bluetooth_devices:     nil,
			},
			Battery_info:      ULComm.BatteryInfo{
				Level:           54,
				Power_connected: true,
			},
			Monitor_info:      ULComm.MonitorInfo{
				Screen_on:  true,
				Brightness: 30,
			},
			Sound_info:        ULComm.SoundInfo{
				Volume:  50,
				Muted:   false,
			},
		},
	}
	if device_info.Device_id == "" {}

	var device_info2 *ULComm.DeviceInfo = ULComm.CreateDeviceInfo(0, 0, false, false, false, false, 0, false, -1,
		"test\x01XX:XX:XX:XX:XX:XX\x01-50\x00test2\x01YY:YY:YY:YY:YY:YY\x01-60\x00",
		"test\x01XX:XX:XX:XX:XX:XX\x01-23\x00test2\x01YY:YY:YY:YY:YY:YY\x01-14\x00", 100, false)
	log.Println(*device_info2)

	device_info2.SendInfo()
	device_info2.SendInfo()

	log.Println(device_info)
	device_info.SendInfo()

	time.Sleep(5 * time.Second)
}
