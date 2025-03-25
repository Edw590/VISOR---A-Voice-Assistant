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
	"strconv"
	"strings"
)

var sessions_GL []ModsFileInfo.Session = nil

/*
RetrieveSessions retrieves the list of sessions, ready to be used by the other functions.
 */
func retrieveSessions() {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, 4, []byte("G_S|true|GPTSessions")) {
		return
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 4)
	if comms_map == nil {
		return
	}

	var response []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	var json_bytes []byte = []byte(Utils.DecompressString(response))

	if err := Utils.FromJsonGENERAL(json_bytes, &sessions_GL); err != nil {
		return
	}
}

/*
GetSessionIdsList updates and gets the list of session IDs.

-----------------------------------------------------------

– Returns:
  - the list of session IDs separated by "|"
 */
func GetSessionIdsList() string {
	retrieveSessions()

	var ids_list string = ""
	for _, session := range sessions_GL {
		ids_list += session.Id + "|"
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetSessionName gets the name of the session.

-----------------------------------------------------------

– Returns:
  - the name of the session
*/
func GetSessionName(session_id string) string {
	for _, session := range sessions_GL {
		if session_id == session.Id {
			return session.Name
		}
	}

	return ""
}

/*
GetSessionCreatedTime gets the creation timestamp of the session.

-----------------------------------------------------------

– Returns:
  - the creation timestamp
 */
func GetSessionCreatedTime(session_id string) int64 {
	for _, session := range sessions_GL {
		if session_id == session.Id {
			return session.Created_time_s
		}
	}

	return -1
}

/*
GetSessionLastInteraction gets the last interaction timestamp of the session.

-----------------------------------------------------------

– Returns:
  - the last interaction timestamp
 */
func GetSessionLastInteraction(session_id string) int64 {
	for _, session := range sessions_GL {
		if session_id == session.Id {
			return session.Last_interaction_s
		}
	}

	return -1
}

/*
GetSessionHistory gets the history of the session.

-----------------------------------------------------------

– Returns:
  - the list of ModsFileInfo.OllamaMessage messages of the session in the following format: "Role/Timestamp|Content"
    separated by "\0"
 */
func GetSessionHistory(session_id string) string {
	for _, session := range sessions_GL {
		if session_id == session.Id {
			var session_history string = ""
			for _, message := range session.History {
				var msg_content string = message.Content
				if message.Role == "user" {
					msg_content = msg_content[strings.Index(msg_content, "]") + 1:]
				}
				session_history += message.Role + "/" + strconv.Itoa(int(message.Timestamp_s)) + "|" + msg_content +
					"\000"
			}
			if len(session_history) > 0 {
				session_history = session_history[:len(session_history)-1]
			}

			return session_history
		}
	}

	return ""
}

/*
DeleteSession deletes the session.

-----------------------------------------------------------

– Params:
  - session_id – the ID of the session
 */
func DeleteSession(session_id string) {
	var message []byte = []byte("S_S|GPTSession|")
	message = append(message, Utils.CompressString(session_id + "\000" + "delete")...)
	Utils.QueueNoResponseMessageSERVER(message)
}

/*
SetSessionName sets the name of the session.

-----------------------------------------------------------

– Params:
  - session_id – the ID of the session
  - name – the name of the session
 */
func SetSessionName(session_id string, name string) {
	var message []byte = []byte("S_S|GPTSession|")
	message = append(message, Utils.CompressString(session_id + "\000" + "rename" + "\000" + name)...)
	Utils.QueueNoResponseMessageSERVER(message)
}
