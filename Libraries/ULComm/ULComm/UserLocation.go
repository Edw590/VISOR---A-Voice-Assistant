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

package ULComm

type UserLocation struct {
	// Last_known_location is the last known location of the user
	Last_known_location string
	// Curr_location is the current location of the user
	Curr_location string
	// Last_time_checked is the last time the current location was checked in Unix time
	Last_time_checked int64
	// Prev_location is the previous location of the user
	Prev_location string
	// Prev_last_time_checked is the last time the previous location was checked in Unix time
	Prev_last_time_checked int64
}
