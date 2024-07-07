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
InitPersonalConsts initializes the personal constants.

-----------------------------------------------------------

– Params:
  - device_id – the device ID
  - website_url – the URL of VISOR's website
  - website_pw – the password of VISOR's website
 */
func InitPersonalConsts(device_id string, website_url string, website_pw string) {
	Utils.User_settings_GL.PersonalConsts.Device_ID = device_id
	Utils.User_settings_GL.PersonalConsts.Website_url = website_url
	Utils.User_settings_GL.PersonalConsts.Website_pw = website_pw
}
