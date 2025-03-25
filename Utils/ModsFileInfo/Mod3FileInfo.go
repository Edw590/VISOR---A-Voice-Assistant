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

package ModsFileInfo

// Mod3UserInfo is the format of the custom information file about this specific module.
type Mod3UserInfo struct {
	// Disks_info is the information about the disks
	Disks_info []DiskInfo
}

type DiskInfo struct {
	// Id is the disk serial number
	Id string
	// Enabled is whether the disk is enabled
	Enabled bool
	// Label is the disk label
	Label string
	// Is_HDD is true if the disk is an HDD, false otherwise
	Is_HDD bool
}

// Mod3GenInfo is the format of the custom generated information about this specific module.
type Mod3GenInfo struct {
	// Disks_info is the information about the disks
	Disks_info []DiskInfo2
}

type DiskInfo2 struct {
	// Id is the disk serial number
	Id string
	// Last_short_test_s is the timestamp of the last short test in seconds
	Last_short_test_s int64
	// Last_long_test_s is the timestamp of the last long test in seconds
	Last_long_test_s int64
}
