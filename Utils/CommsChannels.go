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

package Utils

import (
	"time"

	Tcef "github.com/Edw590/TryCatch-go"
)

// _COMMS_CH_MUL is the multiplier for the number of communication channels per module/library.
const _COMMS_CH_MUL int = 10

const _MODS_COMMS_CHANNELS_SIZE int = MODS_ARRAY_SIZE * _COMMS_CH_MUL
const _LIBS_COMMS_CHANNELS_SIZE int = LIBS_ARRAY_SIZE * _COMMS_CH_MUL

var mods_comms_channels_GL [_MODS_COMMS_CHANNELS_SIZE]chan map[string]any = [_MODS_COMMS_CHANNELS_SIZE]chan map[string]any{}
var libs_comms_channels_GL [_LIBS_COMMS_CHANNELS_SIZE]chan map[string]any = [_LIBS_COMMS_CHANNELS_SIZE]chan map[string]any{}

/*
InitializeCommsChannels initializes the modules and libraries communication channels.
*/
func InitializeCommsChannels() {
	for i := 0; i < MODS_ARRAY_SIZE; i++ {
		for j := 0; j < _COMMS_CH_MUL; j++ {
			mods_comms_channels_GL[getFullChannelNum(i, j)] = make(chan map[string]any, MOD_NUMS_INFO[i].Chan_size)
		}
	}
	for i := 0; i < LIBS_ARRAY_SIZE; i++ {
		for j := 0; j < _COMMS_CH_MUL; j++ {
			libs_comms_channels_GL[getFullChannelNum(i, j)] = make(chan map[string]any, LIB_NUMS_INFO[i].Chan_size)
		}
	}
}

/*
CloseCommsChannels closes the modules and libraries communication channels.
 */
func CloseCommsChannels() {
	for i := 0; i < _MODS_COMMS_CHANNELS_SIZE; i++ {
		if mods_comms_channels_GL[i] != nil {
			close(mods_comms_channels_GL[i])
			mods_comms_channels_GL[i] = nil
		}
	}
	for i := 0; i < _LIBS_COMMS_CHANNELS_SIZE; i++ {
		if libs_comms_channels_GL[i] != nil {
			close(libs_comms_channels_GL[i])
			libs_comms_channels_GL[i] = nil
		}
	}
}

/*
GetFromCommsChannel gets data from a module or library communication channel.

ALWAYS USE TIMEOUT >= 0 WHENEVER YOU CAN TO AVOID BLOCKING THE APP!!!

-----------------------------------------------------------

– Params:
  - is_mod – whether it's a channel from a module or library
  - num – the number of the module or library
  - ch_num – the channel number
  - timeout_s – the timeout in seconds (>= 0) or -1 for no timeout (blocking)

– Returns:
  - the data from the channel or nil if the timeout is reached
*/
func GetFromCommsChannel(is_mod bool, num int, ch_num int, timeout_s int) map[string]any {
	var full_channel_num int = getFullChannelNum(num, ch_num)
	if timeout_s == -1 {
		if is_mod {
			return <- mods_comms_channels_GL[full_channel_num]
		} else {
			return <- libs_comms_channels_GL[full_channel_num]
		}
	} else if timeout_s == 0 {
		if is_mod {
			select {
				case data := <- mods_comms_channels_GL[full_channel_num]:
					return data
				default:
					return nil
			}
		} else {
			select {
				case data := <- libs_comms_channels_GL[full_channel_num]:
					return data
				default:
					return nil
			}
		}
	} else {
		if is_mod {
			select {
				case <- time.After(time.Duration(timeout_s) * time.Second):
					return nil
				case data := <- mods_comms_channels_GL[full_channel_num]:
					return data
			}
		} else {
			select {
				case <- time.After(time.Duration(timeout_s) * time.Second):
					return nil
				case data := <- libs_comms_channels_GL[full_channel_num]:
					return data
			}
		}
	}
}

/*
GetFullChannelNum returns the full channel number.

-----------------------------------------------------------

– Params:
  - num – the number of the module or library
  - channel_num – the channel number

– Returns:
  - the full channel number
 */
func getFullChannelNum(num int, channel_num int) int {
	return num * _COMMS_CH_MUL + channel_num
}

/*
SendToModChannel sends data to a module channel.

In case the module is not supported this function does nothing (to prevent deadlock sending to unused channels).

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module
  - ch_num – the full channel number
  - key – the key of the data
  - data – the data to send
 */
func SendToModChannel(mod_num int, ch_num int, key string, data any) {
	if !IsModSupportedMODULES(mod_num) {
		return
	}

	// Ignore the panic of writing to closed channels (sometimes happens when the app is shutting down).
	Tcef.Tcef{
		Try: func() {
			mods_comms_channels_GL[getFullChannelNum(mod_num, ch_num)] <- map[string]any{
				key: data,
			}
		},
	}.Do()
}

/*
SendToLibChannel sends data to a library channel.

-----------------------------------------------------------

– Params:
  - ch_num – the full channel number
  - key – the key of the data
  - data – the data to send
 */
func SendToLibChannel(lib_num int, ch_num int, key string, data any) {
	// Ignore the panic of writing to closed channels (sometimes happens when the app is shutting down).
	Tcef.Tcef{
		Try: func() {
			libs_comms_channels_GL[getFullChannelNum(lib_num, ch_num)] <- map[string]any{
				key: data,
			}
		},
	}.Do()
}
