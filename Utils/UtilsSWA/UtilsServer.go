/*******************************************************************************
 * Copyright 2023-2024 The V.I.S.O.R. authors
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

package UtilsSWA

import (
	"Utils"
)

/*
StartCommunicatorSERVER keeps the server communicator running in background, unless it's requested to stop.
*/
func StartCommunicatorSERVER() {
	Utils.StartCommunicatorSERVER()
}

/*
QueueGeneralMessageSERVER queues a general message to be sent to the server.

It is received by GetGeneralMessageSERVER().

-----------------------------------------------------------

– Params:
  - message – the message to be sent
*/
func QueueGeneralMessageSERVER(message []byte) {
	Utils.QueueGeneralMessageSERVER(message)
}

/*
GetGeneralMessageSERVER gets a general message from the server.

The message is sent by QueueGeneralMessageSERVER().

If no message is available, the function will wait until a message is received.

-----------------------------------------------------------

– Returns:
  - the message received or nil if the communicator is stopping or stopped
  - true if a message was received, false otherwise
*/
func GetGeneralMessageSERVER() ([]byte, bool) {
	return Utils.GetGeneralMessageSERVER()
}

/*
IsCommunicatorConnectedSERVER checks if the communicator is connected.

-----------------------------------------------------------

– Returns:
  - true if the communicator is connected, false otherwise
*/
func IsCommunicatorConnectedSERVER() bool {
	return Utils.IsCommunicatorConnectedSERVER()
}
