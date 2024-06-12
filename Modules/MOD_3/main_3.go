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

package MOD_3

import (
	"SpeechQueue/SpeechQueue"
	"Utils"
	"time"

	"github.com/Edw590/sapi-go"
	"github.com/go-ole/go-ole"
	"github.com/itchyny/volume-go"
)

const RATE int = 0
const VOLUME int = 100

const _TIME_SLEEP_S int = 1

var tts_GL *sapi.Sapi = nil
var curr_speech_GL *SpeechQueue.Speech = nil

type _MGIModSpecInfo any
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](realMain, module) }
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		_ = ole.CoInitialize(0)

		if tts, err := sapi.NewSapi(); err != nil {
			panic(err)
		} else {
			tts_GL = tts
		}
		_ = tts_GL.SetRate(RATE)
		_ = tts_GL.SetVolume(VOLUME)

		//log.Println("Waiting for speeches to speak...")

		var higher_priority_came bool = false
		go func() {
			for {
				// FIXME This thread keeps running even after the module stopping. Because this blocks here!
				//  Use timers again instead of channels...
				if curr_speech_GL != nil {
					//log.Println("Speaking speech with priority " + strconv.Itoa(int(speech.GetPriority())) + " and ID " +
					//	speech.GetID()[:10] + "(...)...")

					was_muted, err_muted := volume.GetMuted()
					if err_muted != nil {
						// If there was an error getting the mute status, assume it was muted
						was_muted = true
					}

					// Speak too if there was an error getting the mute status (maybe it's not muted, who knows), if the
					// speech is critical, or if the speech is set to always notify.
					var speak bool = err_muted != nil || curr_speech_GL.GetPriority() == SpeechQueue.PRIORITY_CRITICAL ||
						(curr_speech_GL.GetMode()&SpeechQueue.MODE1_ALWAYS_NOTIFY != 0)

					var speech_mode int32 = curr_speech_GL.GetMode()
					var speech_priority int32 = curr_speech_GL.GetPriority()

					var notify bool = false
					if speech_priority == SpeechQueue.PRIORITY_CRITICAL {
						notify = true
					} else {
						if speech_mode&SpeechQueue.MODE1_ALWAYS_NOTIFY != 0 {
							notify = true
						} else if speech_mode&SpeechQueue.MODE1_NO_NOTIF == 0 {
							// If it's not to not notify, notify if he can't speak
							if was_muted {
								notify = true
							}
						}
					}

					if !was_muted || err_muted != nil || curr_speech_GL.GetPriority() == SpeechQueue.PRIORITY_CRITICAL {
						speak = true
					}

					if notify {
						Utils.QueueNotificationNOTIFS("Speeches", curr_speech_GL.GetText())

						//log.Println("Speech notified.")

						// Remove the speech. This means if he can't speak it in case it's to also speak, he won't retry.
						// But it's notified, so shouldn't be a problem I guess.
						SpeechQueue.RemoveSpeech(curr_speech_GL.GetID())
					}

					if !speak {
						curr_speech_GL = nil

						continue
					}

					old_volume, err := volume.GetVolume()
					if err != nil {
						old_volume = -1
					}
					if curr_speech_GL.GetPriority() == SpeechQueue.PRIORITY_CRITICAL {
						_ = volume.SetVolume(100)
						if curr_speech_GL.GetMode()&SpeechQueue.MODE2_BYPASS_NO_SND != 0 {
							_ = volume.Unmute()
						}
					} else {
						if old_volume < 50 {
							_ = volume.SetVolume(50)
						}
					}
					if err = tts_GL.Speak(curr_speech_GL.GetText(), sapi.SVSFDefault); err == nil {
						if old_volume != -1 {
							_ = volume.SetVolume(old_volume)
						}
						if was_muted {
							_ = volume.Mute()
						}

						if !higher_priority_came {
							//log.Println("Speech spoken successfully.")

							SpeechQueue.RemoveSpeech(curr_speech_GL.GetID())
						} else {
							//log.Println("Speech interrupted successfully.")

							higher_priority_came = false
						}
					} else {
						//log.Println("Error speaking speech: ", err)
					}

					curr_speech_GL = nil
				}

				if Utils.WaitWithStop(module_stop, _TIME_SLEEP_S) {
					return
				}
			}
		}()

		for {
			for i := SpeechQueue.NUM_PRIORITIES - 1; i >= 0; i-- {
				var speech *SpeechQueue.Speech = SpeechQueue.GetNextSpeech(i)
				if speech == nil {
					continue
				}

				if curr_speech_GL == nil {
					curr_speech_GL = speech

					break
				} else if speech.GetPriority() > curr_speech_GL.GetPriority() {
					var old_speech *SpeechQueue.Speech = curr_speech_GL
					if stopTts(tts_GL) {
						higher_priority_came = true
						old_speech.RephraseInterrSpeech()
						curr_speech_GL = speech
					}

					break
				}
			}

			if Utils.WaitWithStop(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}

func QueueSpeech(to_speak string, priority int32, mode int32) {
	SpeechQueue.AddSpeech(to_speak, time.Now().UnixMilli(), priority, mode, "")
}

func SkipCurrentSpeech() bool {
	return stopTts(tts_GL)
}

/*
checkSkipSpeech checks if the current speech should be skipped.

-----------------------------------------------------------

– Returns:
  - true if the current speech should be skipped, false otherwise
 */
func checkSkipSpeech() bool {
	if moduleInfo_GL.ModDirsInfo.UserData.Add2(false, "SKIP").Exists() {
		_ = moduleInfo_GL.ModDirsInfo.UserData.Add2(false, "SKIP").Remove()

		return true
	}

	return false
}

func stopTts(tts *sapi.Sapi) bool {
	err := tts.Skip(50) // Equivalent to stopping all speeches it seems
	if err != nil {
		//log.Println("Error stopping speech: ", err)

		return false
	}

	return true
}
