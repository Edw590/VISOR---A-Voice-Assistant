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

package GMan

import (
	"Utils"
	"encoding/json"
	"golang.org/x/oauth2"
	"strconv"
)

/*
SetToken sets the token for the Google Manager.

-----------------------------------------------------------

– Params:
  - token – the token to be set
 */
func SetToken(token *oauth2.Token) {
	var message []byte = []byte("S_S|GManTok|")
	token_bytes, _ := json.Marshal(token)
	message = append(message, Utils.CompressString(string(token_bytes))...)
	Utils.QueueNoResponseMessageSERVER(message)
}

/*
IsTokenValid checks if the token is valid.

-----------------------------------------------------------

– Returns:
  - true if the token is valid, false otherwise
 */
func IsTokenValid() bool {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GMan, 2, []byte("G_S|true|GManTokVal")) {
		return false
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GMan, 2)
	if comms_map == nil {
		return false
	}

	var response []byte = []byte(Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)))

	ret, err := strconv.ParseBool(string(response))
	if err != nil {
		ret = false
	}

	return ret
}
