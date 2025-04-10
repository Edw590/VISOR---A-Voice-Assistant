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

//go:build windows

package Speech

import (
	"github.com/Edw590/sapi-go"
	"github.com/go-ole/go-ole"
	"log"
)

var tts_GL *sapi.Sapi = nil

func initTts() {
	_ = ole.CoInitialize(0)
	defer ole.CoUninitialize()

	if tts, err := sapi.NewSapi(); err != nil {
		panic(err)
	} else {
		tts_GL = tts
	}
	_ = tts_GL.SetRate(0)
	_ = tts_GL.SetVolume(100)

	// Leave this here. It's necessary for the TTS to work on Windows 7. Might be related to bad usage of
	// ole.CoInitialize() which is only for single-threaded applications.
	speak("")
}

func speak(text string) bool {
	_, err := tts_GL.Speak(text, sapi.SVSFDefault)

	return err == nil
}

func stopTts() bool {
	_, err := tts_GL.Skip(50) // Equivalent to stopping all speeches it seems
	if err != nil {
		log.Println("Error stopping speech: ", err)

		return false
	}

	return true
}
