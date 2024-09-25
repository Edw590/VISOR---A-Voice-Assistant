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

package MOD_10

import (
	"Registry/Registry"
	"ULComm/ULComm"
	"Utils"
	"VISOR_Client/ClientRegKeys"
	"github.com/distatus/battery"
	"github.com/go-vgo/robotgo"
	"github.com/itchyny/volume-go"
	"github.com/schollz/wifiscan"
	"github.com/yusufpapurcu/wmi"
	"runtime"
	"strings"
	"time"
)

// System Checker //

const _TIME_SLEEP_S int = 5

var device_info_GL ULComm.DeviceInfo

type _Battery struct {
	power_connected bool
	level           int
}

type _MousePosition struct {
	x int
	y int
}

// https://learn.microsoft.com/en-us/windows/win32/wmicoreprov/wmimonitorbrightness
type WmiMonitorBrightness struct {
	CurrentBrightness uint8
}

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		device_info_GL = ULComm.DeviceInfo{
			Device_id:    Utils.User_settings_GL.PersonalConsts.Device_ID,
		}
		var curr_mouse_position _MousePosition

		for {
			wifi_on, wifi_networks := getWifiNetworks()

			// Connectivity information
			device_info_GL.System_state.Connectivity_info = ULComm.ConnectivityInfo{
				Airplane_mode_enabled: false,
				Wifi_enabled:          wifi_on,
				Bluetooth_enabled:     false,
				Mobile_data_enabled:   false,
				Wifi_networks:         wifi_networks,
				Bluetooth_devices:     nil,
			}

			// Battery information
			var battery_level int = getBatteryInfo().level
			var power_connected bool = getBatteryInfo().power_connected
			Registry.GetValue(ClientRegKeys.K_BATTERY_LEVEL).SetData(battery_level, false)
			Registry.GetValue(ClientRegKeys.K_POWER_CONNECTED).SetData(power_connected, false)

			device_info_GL.System_state.Battery_info = ULComm.BatteryInfo{
				Level:           battery_level,
				Power_connected: power_connected,
			}

			// Monitor information
			var screen_brightness int = getBrightness()
			Registry.GetValue(ClientRegKeys.K_SCREEN_BRIGHTNESS).SetData(screen_brightness, false)

			device_info_GL.System_state.Monitor_info = ULComm.MonitorInfo{
				Screen_on:  true,
				Brightness: screen_brightness,
			}

			// Sound information
			var sound_volume int = getSoundVolume()
			var sound_muted bool = getSoundMuted()
			Registry.GetValue(ClientRegKeys.K_SOUND_VOLUME).SetData(sound_volume, false)
			Registry.GetValue(ClientRegKeys.K_SOUND_MUTED).SetData(sound_muted, false)

			device_info_GL.System_state.Sound_info = ULComm.SoundInfo{
				Volume: sound_volume,
				Muted:  sound_muted,
			}

			// Check if the device is being used by checking if the mouse is moving
			var x, y int = robotgo.Location()
			if x != curr_mouse_position.x || y != curr_mouse_position.y {
				curr_mouse_position.x = x
				curr_mouse_position.y = y

				device_info_GL.Last_time_used = time.Now().Unix()
			}

			device_info_GL.Last_comm = time.Now().Unix()
			_ = device_info_GL.SendInfo()

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}

func GetDeviceInfoText() string {
	return *Utils.ToJsonGENERAL(device_info_GL)
}

func getBatteryInfo() _Battery {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		return _Battery{}
	}

	var b *battery.Battery = batteries[0]

	return _Battery{
		power_connected: b.State.Raw != battery.Discharging,
		level:           int(b.Current / b.Full * 100),
	}
}

func getBrightness() int {
	if runtime.GOOS != "windows" {
		return -1
	}

	var dst []WmiMonitorBrightness
	err := wmi.QueryNamespace("SELECT CurrentBrightness FROM WmiMonitorBrightness", &dst, "root/wmi")
	if err != nil {
		return -1
	}

	if len(dst) > 0 {
		return int(dst[0].CurrentBrightness)
	}

	return -1
}

func getSoundVolume() int {
	vol, err := volume.GetVolume()
	if err != nil {
		return -1
	}

	return vol
}

func getSoundMuted() bool {
	muted, err := volume.GetMuted()
	if err != nil {
		return false
	}

	return muted
}

func getWifiNetworks() (bool, []ULComm.ExtBeacon) {
	for i := 0; i < 10; i++ {
		wifi_nets, err := wifiscan.Scan()
		if err != nil {
			return false, nil
		}

		var wifi_networks []ULComm.ExtBeacon = nil
		for _, wifi_net := range wifi_nets {
			wifi_networks = append(wifi_networks, ULComm.ExtBeacon{
				Name:    wifi_net.SSID,
				Address: strings.ToUpper(wifi_net.BSSID),
				RSSI:    wifi_net.RSSI,
			})
		}

		if len(wifi_networks) != 0 {
			return true, wifi_networks
		}

		time.Sleep(1 * time.Second)
	}

	return true, nil
}
