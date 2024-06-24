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
	MOD_3 "Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"github.com/apaxa-go/eval"
	"github.com/distatus/battery"
	"github.com/itchyny/volume-go"
	"github.com/yusufpapurcu/wmi"
	"log"
	"runtime"
	"strconv"
	"strings"
)

// System State //

const _TIME_SLEEP_S int = 5

type _Battery struct {
	charging bool
	percentage int
}

// https://learn.microsoft.com/en-us/windows/win32/wmicoreprov/wmimonitorbrightness
type WmiMonitorBrightness struct {
	CurrentBrightness uint8
}

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		var notifs_were_true []bool = nil

		for {
			var modUserInfo _ModUserInfo
			if err := moduleInfo_GL.GetModUserInfo(&modUserInfo); err != nil {
				panic(err)
			}

			if len(notifs_were_true) != len(modUserInfo.Notifications) {
				notifs_were_true = make([]bool, len(modUserInfo.Notifications))
			}

			for i, notification := range modUserInfo.Notifications {
				if computeCondition(notification.Condition) {
					if !notifs_were_true[i] {
						notifs_were_true[i] = true

						log.Println(notification.Speak)
						MOD_3.QueueSpeech(notification.Speak, SpeechQueue.PRIORITY_HIGH, SpeechQueue.MODE1_ALWAYS_NOTIFY)
					}
				} else {
					notifs_were_true[i] = false
				}
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}

func computeCondition(condition string) bool {
	condition = formatCondition(condition)
	log.Println("Condition:", condition)
	expr, err := eval.ParseString(condition, "")
	if err != nil {
		log.Println(err)
	}
	r, err := expr.EvalToInterface(nil)
	if err != nil {
		log.Println(err)
	}

	return r.(bool)
}

func formatCondition(condition string) string {
	var battery_info _Battery = getBatteryInfo()
	var monitor_brightness int = getBrightness()
	var sound_volume int = getSoundVolume()
	var sound_muted bool = getSoundMuted()

	condition = strings.Replace(condition, "power_connected", strconv.FormatBool(battery_info.charging), -1)
	condition = strings.Replace(condition, "battery_percent", strconv.Itoa(battery_info.percentage), -1)
	condition = strings.Replace(condition, "brightness", strconv.Itoa(monitor_brightness), -1)
	condition = strings.Replace(condition, "sound_volume", strconv.Itoa(sound_volume), -1)
	condition = strings.Replace(condition, "sound_muted", strconv.FormatBool(sound_muted), -1)

	return condition
}

func getBatteryInfo() _Battery {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		// TODO: handle error

		return _Battery{}
	}

	var b *battery.Battery = batteries[0]

	return _Battery{
		charging:   b.State.Raw != battery.Discharging,
		percentage: int(b.Current / b.Full * 100),
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
