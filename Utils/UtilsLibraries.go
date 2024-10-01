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

package Utils

const (
	NUM_LIB_ACD               int = iota
	NUM_LIB_OICComm
	NUM_LIB_GPTComm
	NUM_LIB_SpeechQueue
	NUM_LIB_ULComm
	NUM_LIB_RRComm

	LIBS_ARRAY_SIZE
)
// LIB_NUMS_NAMES is a map of the numbers of the libraries and their names. Use with the NUM_LIB_ constants.
var LIB_NUMS_NAMES map[int]string = map[int]string{
	NUM_LIB_ACD:         "Advanced Commands Detection",
	NUM_LIB_OICComm:     "Online Information Checker Communicator",
	NUM_LIB_GPTComm:     "GPT Communicator",
	NUM_LIB_SpeechQueue: "Speech Queue",
	NUM_LIB_ULComm:      "User Locator Communicator",
	NUM_LIB_RRComm:      "Reminders Reminder Communicator",
}
