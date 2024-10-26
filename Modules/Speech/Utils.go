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

package Speech

import (
	"Utils"
	"math"
	"time"
)

var stop_volume_processing_GL bool = false
var processing_volume_GL bool = false

func setResetWillChangeVolume(set bool) {
	if set {
		assist_changed_volume_time_ms_GL = time.Now().UnixMilli()
		assist_will_change_volume_GL = true
	} else {
		assist_will_change_volume_GL = false
		assist_changed_volume_time_ms_GL = math.MaxInt64 - VOLUME_CHANGE_INTERVAL
	}
}

func setVoluneMutedStateDefaults() {
	volumeMutedState_GL.audio_stream = _DEFAULT_VALUE
	volumeMutedState_GL.old_volume = _DEFAULT_VALUE
	volumeMutedState_GL.was_muted = _DEFAULT_VALUE
}

func processVolumeChanges() {
	// Equivalent to the setUserChangedVolumeTrue() Android function
	go func() {
		if processing_volume_GL {
			return
		}

		processing_volume_GL = true

		// Detect user changes only after some time after the assistant changed the volume to speak, since the first
		// volume change to be detected would be the assistant himself changing the volume - in case he changed the
		// volume (otherwise the first and next changes will be user changes). Also, detect only if the assistant is
		// speaking, of course.

		var prev_sound_volume int = Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Volume
		for {
			if Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Volume != prev_sound_volume {
				prev_sound_volume = Utils.Gen_settings_GL.MOD_10.Device_info.System_state.Sound_info.Volume

				var carry_on bool = false
				if is_speaking_GL {
					carry_on = true
					if assist_will_change_volume_GL {
						if time.Now().UnixMilli() <= assist_changed_volume_time_ms_GL + VOLUME_CHANGE_INTERVAL {
							// If the assistant will change the volume and it's detected here a volume change before the
							// maximum allowed waiting time for the assistant to change the volume, reset the will
							// change volume variables.
							setResetWillChangeVolume(false)

							continue
						} else {
							// Else, if the assistant will change the volume but the first volume change detection was
							// after the maximum allowed waiting period, reset the will change volume variables and
							// check anyways. This, as a start, shouldn't happen. But if it does, assume it's a user
							// change, and assume there was some error and the assistant didn't get to change the volume.
							setResetWillChangeVolume(false);
						}
					}
				} else {
					// If the assistant is not speaking, discard any volume changes.
					carry_on = false;
				}

				if carry_on {
					// As soon as a user volume change is detected, set the variable to true to indicate the user
					// changed the volume.
					user_changed_volume_GL = true;
				}
			}

			if Utils.WaitWithStopTIMEDATE(&stop_volume_processing_GL, 1) {
				processing_volume_GL = false

				return
			}
		}
	}()
}
