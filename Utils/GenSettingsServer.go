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

//go:build server

package Utils

import "Utils/ModsFileInfo"

type GenSettings struct {
	Device_settings ModsFileInfo.DeviceSettings
	MOD_3           ModsFileInfo.Mod3GenInfo
	MOD_4           ModsFileInfo.Mod4GenInfo
	MOD_5           ModsFileInfo.Mod5GenInfo
	MOD_6           ModsFileInfo.Mod6GenInfo
	MOD_7           ModsFileInfo.Mod7GenInfo
	Registry        []*Value
}
