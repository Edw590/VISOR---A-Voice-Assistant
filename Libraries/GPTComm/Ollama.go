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

package GPTComm

import (
	"Utils"
)

/*
StartReportingNoModelsOLLAMA starts the reporting of no Ollama LLM models for the server.

Call this function if the client version of the GPT Communicator module does not exist (for example on Android).
 */
func StartReportingNoModelsOLLAMA() {
	go func() {
		for {
			var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 1, -1)
			if comms_map == nil {
				return
			}

			var map_value []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
			if map_value == nil {
				return
			}

			var message []byte = []byte("GPT|[models]")
			Utils.QueueNoResponseMessageSERVER(message)
		}
	}()
}
