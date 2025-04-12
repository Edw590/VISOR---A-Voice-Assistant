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

func getModelName(model_type_to_use string) (string, string) {
	log.Println("model_type_to_use:", model_type_to_use)
	var self_models []string = nil
	for _, model := range getLocalModels().Models {
		self_models = append(self_models, model.Name)
	}
	log.Println("self_models:", self_models)
	var model_to_use string = checkModels(self_models, model_type_to_use)
	if model_to_use != "" {
		log.Println("Model found in self models:", model_to_use)
		return model_to_use, Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id
	}

	log.Println("No model name found for type:", model_type_to_use)

	return "", Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id

	var device_models map[string][]string = make(map[string][]string)
	var active_device_ids []string = Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_8.Active_device_IDs
	for _, device_id := range active_device_ids {
		device_models[device_id] = getDeviceLocalModels(device_id)
	}
	log.Println("device_models:", device_models)

	for device_id, models := range device_models {
		model_to_use = checkModels(models, model_type_to_use)
		if model_to_use != "" {
			log.Println("Model found in \"" + device_id + "\" models:", model_to_use)
			return model_to_use, device_id
		}
	}

	log.Println("No model name found for type:", model_type_to_use)

	return "", Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id
}

func checkModels(models []string, model_type_to_use string) string {
	for _, model_info := range getModUserInfo().Models_to_use {
		var model_info_split []string = strings.Split(model_info, " - ")
		var model_name string = model_info_split[0]
		var model_type string = model_info_split[1]
		if model_type == model_type_to_use {
			for _, model := range models {
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

	var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_GPTCommunicator, 1)
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
		log.Println("Error parsing local models:", err)

		return nil
	}

	var model_names []string
	for _, model := range local_models.Models {
		model_names = append(model_names, model.Name)
	}

	return model_names
}
