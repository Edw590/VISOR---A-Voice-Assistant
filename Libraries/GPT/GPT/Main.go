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

package GPT

import (
	"Utils"
	"strconv"
	"strings"
)

/*
GetEntry gets the entry at the given number or time.

If -1 is provided on both parameters, it will return the last entry. The time parameter is prioritized over the number
parameter.

-----------------------------------------------------------

– Params:
  - time – the time of the entry or -1 if the entry is to be found by number
  - num – the number of the entry or negative to count from the end

– Returns:
  - the entry or an empty entry with time = -1 if it doesn't exist
 */
func GetEntry(time int64, num int) *Entry {
	var page_contents string = string(Utils.GetPageContentsWEBSITE("files_EOG/gpt_text.txt"))
	if page_contents == "" {
		return &Entry{
			device_id: "",
			text: "",
			time: -1,
		}
	}
	var entries []string = strings.Split(page_contents, "[3234_START:")

	if time != -1 {
		for _, entry := range entries {
			if entry == "" {
				continue
			}

			if getTimeFromEntry(entry) == time {
				return &Entry{
					device_id: getDeviceIdFromEntry(entry),
					text:      getTextFromEntry(entry),
					time:      getTimeFromEntry(entry),
				}
			}
		}
	} else {
		if len(entries) == 0 {
			return &Entry{
				device_id: "",
				text: "",
				time: -1,
			}
		}

		if num < 0 {
			num = len(entries) + num
		} else if num >= len(entries) {
			num = len(entries) - 1
		} else {
			// Do nothing
		}

		var entry string = entries[num]

		if entry != "" {
			return &Entry{
				device_id: getDeviceIdFromEntry(entry),
				text:      getTextFromEntry(entry),
				time:      getTimeFromEntry(entry),
			}
		}
	}

	return &Entry{
		device_id: "",
		text: "",
		time: -1,
	}
}

func SendText(text string) error {
	_, err := Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
		Type:  "GPT",
		Text1: "[" + Utils.PersonalConsts_GL.DEVICE_ID + "]" + text,
	})

	return err
}

/*
getDeviceIdFromEntry gets the device ID from the entry.

-----------------------------------------------------------

– Params:
  - entry – the entry

– Returns:
  - the device ID
 */
func getDeviceIdFromEntry(entry string) string {
	// It comes like: "time|DEVICE_ID|..."
	entry = entry[strings.Index(entry, "|") + 1:]

	return entry[:strings.Index(entry, "|")]
}

/*
getTimeFromEntry gets the time from the entry.

-----------------------------------------------------------

– Params:
  - entry – the entry

– Returns:
  - the time
*/
func getTimeFromEntry(entry string) int64 {
	// It comes like: "time|..."
	var time_str string = entry[:strings.Index(entry, "|")]

	time, err := strconv.ParseInt(time_str, 10, 64)
	if err != nil {
		return -1
	}

	return time
}

/*
getTextFromEntry gets the text from the entry.

-----------------------------------------------------------

– Params:

*/
func getTextFromEntry(entry string) string {
	// It comes like: "...]text[3234_END]"
	var text string = entry[strings.Index(entry, "]") + 1:]

	// Remove all after END_ENTRY if it exists
	var end_entry_exists bool = strings.Contains(text, END_ENTRY)

	var parts []string = strings.Split(text, END_ENTRY)

	if end_entry_exists {
		parts[0] += END_ENTRY
	}

	return parts[0]
}
