/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

import "Utils"

/*
StartCommunicatorSERVER starts the server communicator.

This function does not return until the communicator is stopped.
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
*/
func GetGeneralMessageSERVER() []byte {
	return Utils.GetGeneralMessageSERVER()
}

/*
StopCommunicatorSERVER stops the communicator.
*/
func StopCommunicatorSERVER() {
	Utils.StopCommunicatorSERVER()
}
