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

package MOD_11

import (
	"Registry/Registry"
	"Utils"
	"VISOR_Client/ClientRegKeys"
	porcupine "github.com/Picovoice/porcupine/binding/go/v3"
	"github.com/gordonklaus/portaudio"
)

// Speech Recognition //

var (
	in []int16
	stream *portaudio.Stream
)

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		porcupine_ := porcupine.Porcupine{
			AccessKey: Utils.User_settings_GL.PersonalConsts.Picovoice_API_key, // from Picovoice Console (https://console.picovoice.ai/)
			KeywordPaths: []string{
				moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "Hey-Visor_en_windows_v3_0_0.ppn").GPathToStringConversion(),
			},
		}
		err := porcupine_.Init()
		if err != nil {
			panic(err)
		}
		defer porcupine_.Delete()

		portaudio.Initialize()
		if err != nil {
			panic(err)
		}
		defer closeAudio()
		in = make([]int16, porcupine.FrameLength)
		stream, err = portaudio.OpenDefaultStream(1, 0, float64(porcupine.SampleRate), porcupine.FrameLength, in)
		if err != nil {
			panic(err)
		}
		err = stream.Start()
		if err != nil {
			panic(err)
		}

		for {
			keywordIndex, _ := porcupine_.Process(getNextFrameAudio())
			if keywordIndex >= 0 {
				Registry.GetValue(ClientRegKeys.K_SHOW_APP_SIG).SetData(true, false)
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, 0) {
				return
			}
		}
	}
}

func getNextFrameAudio() []int16 {
	err := stream.Read()
	if err != nil {
		panic(err)
	}

	return in
}

func closeAudio() {
	if stream != nil {
		stream.Stop()
		stream.Close()
	}
	portaudio.Terminate()
}
