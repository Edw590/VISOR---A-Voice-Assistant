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

import "time"

/*
WaitForNetwork waits for the network to be connected by checking the server communicator connection status.

-----------------------------------------------------------

– Params:
  - timeout_s – the timeout in seconds

– Returns:
  - true if the network is connected, false if it's not or if the timeout expired
 */
func WaitForNetwork(timeout_s int64) bool {
	var start_time int64 = time.Now().Unix()
	for !IsCommunicatorConnectedSERVER() {
		if time.Now().Unix() - start_time >= timeout_s {
			return false
		}

		time.Sleep(1 * time.Second)
	}

	return IsCommunicatorConnectedSERVER()
}
