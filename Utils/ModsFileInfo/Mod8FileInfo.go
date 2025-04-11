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

// Mod9GenInfo is the format of the custom generated information about this specific module.
type Mod8GenInfo struct {
	// Active_device_IDs is the list of active device IDs
	Active_device_IDs []string
}

///////////////////////////////////////////////////////////////////////////////

// Mod8UserInfo is the format of the custom information file about this specific module.
type Mod8UserInfo struct {
	// Crt_file is the location of the SSL certificate crt file
	Crt_file string
	// Key_file is the location of the SSL certificate key file
	Key_file string
}
