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

package SettingsSync

import (
	"Utils"
	"strings"
)

/*
GetTempLocsOIC returns the locations to get the weather from, from the user settings.

-----------------------------------------------------------

– Returns:
  - the locations separated by "\n"
 */
func GetTempLocsOIC() string {
	return strings.Join(Utils.GetUserSettings().OnlineInfoChk.Temp_locs, "\n")
}

/*
SetTempLocsOIC sets the locations to get the weather from, in the user settings.

-----------------------------------------------------------

– Params:
  - locs – the locations separated by "\n"
 */
func SetTempLocsOIC(locs string) {
	Utils.GetUserSettings().OnlineInfoChk.Temp_locs = strings.Split(locs, "\n")
}

/*
GetNewsLocsOIC returns the locations to get the news from, from the user settings.

-----------------------------------------------------------

– Returns:
  - the locations separated by "\n"
 */
func GetNewsLocsOIC() string {
	return strings.Join(Utils.GetUserSettings().OnlineInfoChk.News_locs, "\n")
}

/*
SetNewsLocsOIC sets the locations to get the news from, in the user settings.

-----------------------------------------------------------

– Params:
  - locs – the locations separated by "\n"
 */
func SetNewsLocsOIC(locs string) {
	Utils.GetUserSettings().OnlineInfoChk.News_locs = strings.Split(locs, "\n")
}
