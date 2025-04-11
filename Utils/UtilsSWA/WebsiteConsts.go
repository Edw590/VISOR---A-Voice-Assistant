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

package UtilsSWA

import "Utils"

/*
InitWebsiteConsts initializes the website constants.

Useful for before the User Settings are downloaded - this sets the constants needed to access the server and download
the User Settings file.

-----------------------------------------------------------

– Params:
  - website_domain – the domain of VISOR's website
  - website_pw – the password of VISOR's website
 */
func InitWebsiteConsts(website_domain string, website_pw string) {
	Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Website_domain = website_domain
	Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Website_pw = website_pw
}
