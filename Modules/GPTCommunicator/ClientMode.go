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

package GPTCommunicator

import (
	"Utils"
	"log"
	"strings"
)

func clientMode() {
	go func() {
		for {
			var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 0)
			if comms_map == nil {
				return
			}

			var map_value []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
			if map_value == nil {
				return
			}

			var request string = Utils.DecompressString(map_value)
			log.Println("request", request)
			if request == "" {
				continue
			}

			var idx_pipe int = strings.Index(request, "|")
			if idx_pipe == -1 {
				continue
			}

			var device_id string = request[:idx_pipe]
			var request_json string = request[idx_pipe+1:]
			if device_id == "" || request_json == "" {
				continue
			}

			sendReceiveOllamaRequest(device_id, []byte(request_json),
				Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id)
		}
	}()

	go func() {
		for {
			var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 1)
			if comms_map == nil {
				return
			}

			var map_value []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)
			if map_value == nil {
				return
			}

			var models string = *Utils.ToJsonGENERAL(getLocalModels())

			var message []byte = []byte("GPT|[models]")
			message = append(message, Utils.CompressString(models)...)
			Utils.QueueNoResponseMessageSERVER(message)
		}
	}()

	for {
		if Utils.WaitWithStopDATETIME(module_stop_GL, 1000000000) {
			return
		}
	}
}
