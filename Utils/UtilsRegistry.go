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

// This is a copy from UtilsSWA/UtilsRegistry.go so that I can add the Registry to the Generated Settings at the same
// time that I can compile the Main Libraries with the Registry utilities included. Any more decent way welcomed...

// Value represents a value in the registry
type Value struct {
	// Key is the Key of the value
	Key string
	// Pretty_name is the pretty name of the value
	Pretty_name string
	// Description is the Description of the value
	Description string
	// Type_ is the type of the value
	Type_ string

	// Prev_data is the previous data of the value
	Prev_data string
	// Time_updated_prev is the time the previous data was updated in milliseconds
	Time_updated_prev int64
	// Curr_data is the current data of the value
	Curr_data string
	// Time_updated_curr is the time the data was updated in milliseconds
	Time_updated_curr int64
}
