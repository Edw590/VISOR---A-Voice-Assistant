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

package RRComm

import (
	"Utils"
	"Utils/ModsFileInfo"
)

/*
GetRemindersList returns a list of all reminders.

-----------------------------------------------------------

– Returns:
  - a list of all reminders
 */
func GetRemindersList() *[]ModsFileInfo.Reminder {
	var page_contents []byte = Utils.GetFileContentsWEBSITE("reminders.json", true)

	var reminders []ModsFileInfo.Reminder
	if err := Utils.FromJsonGENERAL(page_contents, &reminders); err != nil {
		return nil
	}

	return &reminders
}

/*
GetIdsList returns a list of all reminders' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all reminders' IDs separated by "|"
 */
func GetIdsList() string {
	var page_contents []byte = Utils.GetFileContentsWEBSITE("reminders.json", true)

	var user_location []ModsFileInfo.Reminder
	if err := Utils.FromJsonGENERAL(page_contents, &user_location); err != nil {
		return ""
	}

	var ids string
	for _, reminder := range user_location {
		ids += reminder.Id + "|"
	}

	return ids
}

/*
GetReminderById returns a reminder by its ID.

-----------------------------------------------------------

– Params:
  - id – the reminder ID

– Returns:
  - the reminder or nil if the reminder was not found
 */
func GetReminderById(id string) *ModsFileInfo.Reminder {
	var p_reminders *[]ModsFileInfo.Reminder = GetRemindersList()
	if p_reminders == nil {
		return nil
	}

	for _, reminder := range *p_reminders {
		if reminder.Id == id {
			return &reminder
		}
	}

	return nil
}
