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

import (
	"Utils"
)

/*
GetDateTimeStrTIMEDATE gets the current time and date in the format DATE_TIME_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current time and date in the default format
 */
func GetDateTimeStrTIMEDATE(millis int64) string {
	return Utils.GetDateTimeStrTIMEDATE(millis)
}

/*
GetDateStrTIMEDATE gets the current date in the format DATE_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current time in the default format
 */
func GetDateStrTIMEDATE(millis int64) string {
	return Utils.GetDateStrTIMEDATE(millis)
}

/*
GetTimeStrTIMEDATE gets the current time in the format TIME_FORMAT.

-----------------------------------------------------------

– Returns:
  - the current time in the default format
 */
func GetTimeStrTIMEDATE(millis int64) string {
	return Utils.GetTimeStrTIMEDATE(millis)
}
