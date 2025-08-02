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

package SpeechRecognition

import (
	"Utils"
	porcupine "github.com/Picovoice/porcupine/binding/go/v3"
	"github.com/gordonklaus/portaudio"
	"log"
)

// TODO: Go find a real way of making him listen to us
//  For the moment this module is disabled - less work with the new Program Data assets directory.
//  The module doesn't do anything useful anyway.

var in_GL []int16
var stream_GL *portaudio.Stream

var (
	modDirsInfo_GL Utils.ModDirsInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	porcupine_ := porcupine.Porcupine{
		// AccessKey from Picovoice Console (https://console.picovoice.ai/)
		AccessKey: Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Picovoice_API_key,
		KeywordPaths: []string{
			modDirsInfo_GL.ProgramData.Add2(false, "Hey-Visor_en_windows_v3_0_0.ppn").GPathToStringConversion(),
		},
	}
	err := porcupine_.Init()
	if err != nil {
		Utils.LogLnError(err)

		return
	}
	defer porcupine_.Delete()

	err = portaudio.Initialize()
	if err != nil {
		Utils.LogLnError(err)

		return
	}
	defer closeAudio()
	in_GL = make([]int16, porcupine.FrameLength)
	stream_GL, err = portaudio.OpenDefaultStream(1, 0, float64(porcupine.SampleRate), porcupine.FrameLength, in_GL)
	if err != nil {
		Utils.LogLnError(err)

		return
	}
	err = stream_GL.Start()
	if err != nil {
		Utils.LogLnError(err)

		return
	}

	for {
		keywordIndex, _ := porcupine_.Process(getNextFrameAudio())
		if keywordIndex >= 0 {
			Utils.SendToModChannel(Utils.NUM_MOD_VISOR, 0, "ShowApp", nil)
		}

		if Utils.WaitWithStopDATETIME(module_stop, 0) {
			return
		}
	}
}

func getNextFrameAudio() []int16 {
	err := stream_GL.Read()
	if err != nil {
		Utils.LogLnError(err)

		return nil
	}

	return in_GL
}

func closeAudio() {
	if stream_GL != nil {
		_ = stream_GL.Stop()
		_ = stream_GL.Close()
	}
	_ = portaudio.Terminate()
}
