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

import (
	"Utils"
)

/*
GetDateTimeStrDATETIME gets the current time and date in the format DATE_TIME_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the time in seconds or -1 for the current time

– Returns:
  - the current time and date in the default format
 */
func GetDateTimeStrDATETIME(s int64) string {
	return Utils.GetDateTimeStrTIMEDATE(s)
}

/*
GetDateStrDATETIME gets the current date in the format DATE_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the date in seconds or -1 for the current date

– Returns:
  - the current date in the default format
 */
func GetDateStrDATETIME(s int64) string {
	return Utils.GetDateStrTIMEDATE(s)
}

/*
GetTimeStrDATETIME gets the current time in the format TIME_FORMAT.

-----------------------------------------------------------

– Params:
  - s – the time in seconds or -1 for the current time

– Returns:
  - the current time in the default format
 */
func GetTimeStrDATETIME(s int64) string {
	return Utils.GetTimeStrTIMEDATE(s)
}
