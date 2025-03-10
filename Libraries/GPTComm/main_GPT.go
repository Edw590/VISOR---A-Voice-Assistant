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

package main

import (
	"GPTComm/GPTComm"
	"Utils"
	"log"
	"time"
)

func main() {
	Utils.LoadDeviceUserSettings(false)
	Utils.InitializeCommsChannels()

	Utils.StartCommunicatorSERVER()
	time.Sleep(4 * time.Second)
	//Utils.StopCommunicatorSERVER()

	/*for {
		time.Sleep(1 * time.Second)
	}*/

	/*for {
		log.Println(GPT.GetEntry(-1, -1))

		time.Sleep(1 * time.Second)
	}*/

	go func() {for {
		var sentence string = GPTComm.GetNextSpeechSentence()
		if sentence == "" {
			continue
		}

		log.Println("sentence: " + sentence)

		time.Sleep(1 * time.Second)
	}
	}()

	time.Sleep(1 * time.Second)

	log.Println(GPTComm.SendText("hello", GPTComm.SESSION_TYPE_NEW))

	for {
		time.Sleep(1 * time.Second)
	}
}
