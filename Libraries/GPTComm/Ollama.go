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
	"Utils/ModsFileInfo"
)

/*
StartReportingNoModelsOLLAMA starts the reporting of no available LLM models to the server.

Call this function if the client version of the GPT Communicator module is not available (for example on Android).
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

			var message []byte = []byte("GPT|models|")
			Utils.QueueNoResponseMessageSERVER(message)
		}
	}()
}

/*
GetModelOLLAMA gets a model from the list of models.

-----------------------------------------------------------

– Params:
  - model_name – the name of the model to get

– Returns:
  - the model or nil if the model does not exist
 */
func GetModelOLLAMA(model_name string) *ModsFileInfo.Model {
	// Get the model from the list of models
	model, ok := Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator.Models[model_name]
	if !ok {
		return nil
	}

	return model
}

/*
AddUpdateModelOLLAMA adds or updates a model in the list of models.

-----------------------------------------------------------

– Params:
  - model_name – the name of the model to add or update
  - model_type – the type of the model one of the GPTComm.MODEL_TYPE_* constants
  - has_tool_role – whether the model has the tool role available or not
  - context_size – the context size of the model
  - temperature – the temperature of the model
  - system_info – the system info of the model
 */
func AddUpdateModelOLLAMA(model_name string, model_type string, has_tool_role bool, context_size int32,
						  temperature float32, system_info string) {
	// Add the model to the list of models
	Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator.Models[model_name] = &ModsFileInfo.Model{
		Type:          model_type,
		Has_tool_role: has_tool_role,
		Context_size:  context_size,
		Temperature:   temperature,
		System_info:   system_info,
	}
}

/*
DeleteModelOLLAMA deletes a model from the list of models.

-----------------------------------------------------------

– Params:
  - model_name – the name of the model to delete
 */
func DeleteModelOLLAMA(model_name string) {
	delete(Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator.Models, model_name)
}
