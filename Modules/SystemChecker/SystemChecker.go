/*******************************************************************************
 * Copyright 2023-2024 The V.I.S.O.R. authors
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

package SystemChecker

import (
	"Utils"
	"Utils/ModsFileInfo"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/distatus/battery"
	"github.com/go-vgo/robotgo"
	"github.com/itchyny/volume-go"
	"github.com/schollz/wifiscan"
	"runtime"
	"strings"
	"time"
)

const SCAN_WIFI_EACH_S int64 = 60
var last_check_wifi_when_s int64 = 0

var device_info_GL *ModsFileInfo.DeviceInfo

type _Battery struct {
	power_connected bool
	level           int
}

type _MousePosition struct {
	x int
	y int
}

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
	modGenInfo_GL *ModsFileInfo.Mod10GenInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_10

		var curr_mouse_position _MousePosition

		device_info_GL = &modGenInfo_GL.Device_info

		var wifi_networks []ModsFileInfo.ExtBeacon
		for {
			if time.Now().Unix() >= last_check_wifi_when_s + SCAN_WIFI_EACH_S {
				// Every 3 minutes, update the wifi networks
				wifi_networks = getWifiNetworks()

				last_check_wifi_when_s = time.Now().Unix()
			}


			// Connectivity information
			device_info_GL.System_state.Connectivity_info = ModsFileInfo.ConnectivityInfo{
				Airplane_mode_enabled: false,
				Wifi_enabled:          getWifiEnabled(),
				Bluetooth_enabled:     false,
				Mobile_data_enabled:   false,
				Wifi_networks:         wifi_networks,
				Bluetooth_devices:     nil,
			}


			// Battery information
			battery_level, power_connected := getBatteryInfo(device_info_GL.System_state.Battery_info.Level,
				device_info_GL.System_state.Battery_info.Power_connected)
			UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_BATTERY_LEVEL).SetInt(int32(battery_level), false)
			UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_POWER_CONNECTED).SetBool(power_connected, false)

			device_info_GL.System_state.Battery_info = ModsFileInfo.BatteryInfo{
				Level:           battery_level,
				Power_connected: power_connected,
			}


			// Monitor information
			var screen_brightness int = Utils.GetScreenBrightnessSYSTEM()
			if screen_brightness == -1 {
				screen_brightness = device_info_GL.System_state.Monitor_info.Brightness
			}
			UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SCREEN_BRIGHTNESS).SetInt(int32(screen_brightness), false)

			device_info_GL.System_state.Monitor_info = ModsFileInfo.MonitorInfo{
				Screen_on:  true,
				Brightness: screen_brightness,
			}


			// Sound information
			var sound_volume int = getSoundVolume(device_info_GL.System_state.Sound_info.Volume)
			var sound_muted bool = getSoundMuted(device_info_GL.System_state.Sound_info.Muted)
			UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SOUND_VOLUME).SetInt(int32(sound_volume), false)
			UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_SOUND_MUTED).SetBool(sound_muted, false)

			device_info_GL.System_state.Sound_info = ModsFileInfo.SoundInfo{
				Volume: sound_volume,
				Muted:  sound_muted,
			}


			// Check if the device is being used by checking if the mouse is moving
			var x, y int = robotgo.Location()
			if x != curr_mouse_position.x || y != curr_mouse_position.y {
				curr_mouse_position.x = x
				curr_mouse_position.y = y

				device_info_GL.Last_time_used_s = time.Now().Unix()
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				return
			}
		}
	}
}

func getBatteryInfo(prev1 int, prev2 bool) (int, bool) {
	var batteries []*battery.Battery
	var err error
	var panicked bool = false
	Tcef.Tcef{
		Try: func() {
			batteries, err = battery.GetAll()
		},
		Catch: func(e Tcef.Exception) {
			panicked = true
		},
	}.Do()
	if panicked || err != nil || len(batteries) == 0 {
		return prev1, prev2
	}

	var b *battery.Battery = batteries[0]

	return int(b.Current / b.Full * 100), b.State.Raw != battery.Discharging
}

func getSoundVolume(prev int) int {
	vol, err := volume.GetVolume()
	if err != nil {
		return prev
	}

	return vol
}

func getSoundMuted(prev bool) bool {
	muted, err := volume.GetMuted()
	if err != nil {
		return prev
	}

	return muted
}

func getWifiNetworks() []ModsFileInfo.ExtBeacon {
	var wifi_was_enabled = getWifiEnabled()
	if !wifi_was_enabled {
		setWifiEnabled(true)
	}

	if runtime.GOOS == "windows" {
		// Request a Wi-Fi scan first (wifiscan.Scan() doesn't do it on Windows)
		_, _ = Utils.ExecCmdSHELL([]string{".\\external\\WlanScan.exe"})
	}

	// Then get the cached results on Windows or request a scan and get results
	// on Linux.
	var num_tries int = 1
	if runtime.GOOS == "windows" {
		// I don't know how much time after the scan the results are ready, so 10 seconds seems like a good number. If
		// it's on Linux, there shouldn't be a problem I think.
		num_tries = 10
	}
	var wifi_networks []ModsFileInfo.ExtBeacon = nil
	for i := 0; i < num_tries; i++ {
		wifilist, err := wifiscan.Scan()
		if err != nil {
			break
		}

		for _, wifi_net := range wifilist {
			wifi_networks = append(wifi_networks, ModsFileInfo.ExtBeacon{
				Name:    wifi_net.SSID,
				Address: strings.ToUpper(wifi_net.BSSID),
				RSSI:    wifi_net.RSSI,
			})
		}

		if len(wifi_networks) != 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if !wifi_was_enabled {
		setWifiEnabled(false)
	}

	return wifi_networks
}

func getWifiEnabled() bool {
	if runtime.GOOS == "windows" {
		cmd_output, err := Utils.ExecCmdSHELL([]string{"netsh.exe wlan show networks mode=Bssid"})

		return err == nil && cmd_output.Exit_code == 0
	}

	cmd_output, err := Utils.ExecCmdSHELL([]string{"nmcli radio  wifi"})
	if err != nil {
		return false
	}

	return strings.Contains(cmd_output.Stdout_str, "enabled")
}

func setWifiEnabled(enabled bool) bool {
	var cmd_output Utils.CmdOutput
	var err error
	if runtime.GOOS == "windows" {
		if enabled {
			cmd_output, err = Utils.ExecCmdSHELL([]string{"netsh interface set interface name=Wi-Fi admin=enabled"})
		} else {
			cmd_output, err = Utils.ExecCmdSHELL([]string{"netsh interface set interface name=Wi-Fi admin=disabled"})
		}
	} else {
		if enabled {
			cmd_output, err = Utils.ExecCmdSHELL([]string{"nmcli radio wifi on"})
		} else {
			cmd_output, err = Utils.ExecCmdSHELL([]string{"nmcli radio wifi off"})
		}
	}

	return err == nil && cmd_output.Exit_code == 0
}
