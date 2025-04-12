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

package Utils

import "strconv"

/*
QueueMessageBACKEND sends a message to the Website Backend module to be sent to another device.

-----------------------------------------------------------

– Params:
  - is_mod – true if the message is for a module, false if it is for a library
  - num – the number of the module or library
  - channel_num – the number of the channel to send the message to
  - dest_device_id – the ID of the destination device
  - message – the message to send
 */
func QueueMessageBACKEND(is_mod bool, num int, channel_num int, dest_device_id string, message []byte) {
	var msg_to string = ""
	if is_mod {
		msg_to += "M_"
	} else {
		msg_to += "L_"
	}
	msg_to += strconv.Itoa(num) + "_" + strconv.Itoa(channel_num)
	var new_message = []byte(dest_device_id + "|" + msg_to + "|")
	new_message = append(new_message, message...)
	SendToModChannel(NUM_MOD_WebsiteBackend, 0, "Message", new_message)
}
