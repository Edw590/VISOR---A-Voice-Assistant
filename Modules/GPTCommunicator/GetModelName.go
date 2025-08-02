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
	"GPTComm"
	"Utils"
	"strings"
)

func getModelName(images_needed bool) (string, string) {
	var model_to_use string = ""
	var device_id_to_use string = ""

	var model_type_to_use string = ""
	if images_needed {
		model_type_to_use = GPTComm.MODEL_TYPE_VISION
	} else {
		model_type_to_use = GPTComm.MODEL_TYPE_TEXT
	}

	getServerModels := func() {
		var self_models []string = nil
		for _, model := range getLocalModels().Models {
			self_models = append(self_models, model.Name)
		}
		model_to_use = checkModels(self_models, model_type_to_use)
		if model_to_use != "" {
			device_id_to_use = Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id
			Utils.LogLnDebug(model_to_use)
		}
	}

	getClientsModels := func() {
		var device_models map[string][]string = make(map[string][]string)
		var active_device_ids []string = Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_8.Active_device_IDs
		for _, device_id := range active_device_ids {
			device_models[device_id] = getDeviceLocalModels(device_id)
		}

		for device_id, models := range device_models {
			model_to_use = checkModels(models, model_type_to_use)
			if model_to_use != "" {
				device_id_to_use = device_id
				Utils.LogLnDebug(model_to_use)

				break
			}
		}
	}

	if Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator.Prioritize_clients_models {
		getClientsModels()
		if model_to_use != "" {
			goto end
		}

		getServerModels()
		if model_to_use != "" {
			goto end
		}
	} else {
		getServerModels()
		if model_to_use != "" {
			goto end
		}

		getClientsModels()
		if model_to_use != "" {
			goto end
		}
	}

	end:

	return model_to_use, device_id_to_use
}

func checkModels(device_models []string, model_type_to_use string) string {
	var model_priorities []string = strings.Split(getModUserInfo().Model_priorities, "\n")
	for _, model_name := range model_priorities {
		if model_name == "" {
			continue
		}

		model_info, ok := Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator.Models[model_name]
		if !ok {
			Utils.LogLnInfo(model_name)

			continue
		}

		if strings.Contains(model_info.Type, model_type_to_use) {
			for _, model := range device_models {
				if model == model_name {
					return model_name
				}
			}
		}
	}

	return ""
}

func getDeviceLocalModels(device_id string) []string {
	Utils.QueueMessageBACKEND(true, Utils.NUM_MOD_GPTCommunicator, 1, device_id, nil)

	var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 1, -1)
	if comms_map == nil {
		return nil
	}

	var map_value []byte = []byte(comms_map["Models"].(string))
	if map_value == nil {
		return nil
	}

	var local_models _LocalModels
	err := Utils.FromJsonGENERAL(map_value, &local_models)
	if err != nil {
		Utils.LogLnError(err)

		return nil
	}

	var model_names []string
	for _, model := range local_models.Models {
		model_names = append(model_names, model.Name)
	}

	return model_names
}
