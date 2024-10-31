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

package ModsFileInfo

// Mod3UserInfo is the format of the custom information file about this specific module.
type Mod3UserInfo struct {
	// Disks_info is the information about the disks. It maps the disk serial number to the disk information struct.
	Disks_info map[string]*DiskInfo
}

type DiskInfo struct {
	// Disk label
	Label string
	// Is the disk an HDD?
	Is_HDD bool
}

// Mod3GenInfo is the format of the custom generated information about this specific module.
type Mod3GenInfo struct {
	// Disks_info is the information about the disks. It maps the disk serial number to an array with the first element
	// being the last short test timestamp and the second element being the last long test timestamp. The timestamps are
	// in seconds.
	Disks_info map[string][]int64
}
