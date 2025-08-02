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
	"Utils/ModsFileInfo"
	"strconv"
	"time"
)

const _START_CMD string = "[3234_START:"
const _END_CMD string = "[3234_END]"

type _LocalModels struct {
	Models []_LocalModel `json:"models"`
}

type _LocalModel struct {
	Name string `json:"name"`
}

func getStartString(device_id string) string {
	return _START_CMD + strconv.FormatInt(time.Now().UnixMilli(), 10) + "|" + device_id + "|]"
}

func getEndString() string {
	return "\n" + _END_CMD + "\n"
}

func getModUserInfo() *ModsFileInfo.Mod7UserInfo {
	return &Utils.GetUserSettings(Utils.LOCK_UNLOCK).GPTCommunicator
}

func getLocalModels() _LocalModels {
	body, err := Utils.MakeGetRequest("http://localhost:11434/api/tags")
	if err != nil {
		Utils.LogLnError(err)
	}

	var local_models _LocalModels
	if err = Utils.FromJsonGENERAL(body, &local_models); err != nil {
		Utils.LogLnError(err)
	}

	return local_models
}
