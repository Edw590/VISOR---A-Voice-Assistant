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

// GEN_ERROR is the error message for general errors
const GEN_ERROR string = "3234_ERROR"

var website_url_GL string = ""
var website_pw_GL string = ""

/*
SetWebsiteInfo sets VISOR's website url and password.

-----------------------------------------------------------

– Params:
  - url – the website url
  - pw – the website password
*/
func SetWebsiteInfo(url string, pw string) {
	website_url_GL = url
	website_pw_GL = pw
}

/*
GetNumEntries gets the number of entries in the text file.

-----------------------------------------------------------

– Returns:
  - the number of entries
 */
func GetNumEntries() int {
	var page_contents string = string(Utils.GetPageContentsWEBSITE(website_url_GL + "files_EOG/gpt_text.txt", website_pw_GL))
	var entries []string = strings.Split(page_contents, "[3234_START:")

	for i := 0; i < len(entries); i++ {
		if entries[i] == "" {
			Utils.DelElemSLICES(&entries, i)
		}
	}

	return len(entries)
}

/*
GetEntry gets the entry at the given number or time.

If -1 is provided on both parameters, it will return the last entry. The time parameter is prioritized over the number
parameter.

-----------------------------------------------------------

– Params:
  - time – the time of the entry or -1 if the entry is to be found by number
  - num – the number of the entry or negative to count from the end

– Returns:
  - the entry or nil if an error occurred
 */
func GetEntry(time int64, num int32) *Entry {
	var page_contents string = string(Utils.GetPageContentsWEBSITE(website_url_GL + "files_EOG/gpt_text.txt", website_pw_GL))
	if page_contents == "" {
		return nil
	}
	var entries []string = strings.Split(page_contents, "[3234_START:")

	if time != -1 {
		for _, entry := range entries {
			if getTimeFromEntry(entry) == time {
				return &Entry{
					text: getTextFromEntry(entry),
					time: getTimeFromEntry(entry),
				}
			}
		}
	} else {
		if len(entries) == 0 {
			return nil
		}

		if num < 0 {
			num = int32(len(entries)) + num
		} else if num >= int32(len(entries)) {
			num = int32(len(entries)) - 1
		} else {
			// Do nothing
		}

		var entry string = entries[num]

		return &Entry{
			text: getTextFromEntry(entry),
			time: getTimeFromEntry(entry),
		}
	}

	return nil
}

func SendText(text string) error {
	return Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
		Name: "GPT",
		Text1: text,
	})
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
	// It comes like: "[3234_START:time]text[3234_END]"
	// Filtered gets it like: "time|text[3234_END]"
	entry = entry[:strings.Index(entry, "]")] + "|" + entry[strings.Index(entry, "]") + 1:]
	var parts []string = strings.Split(entry, "|")

	if len(parts) < 2 {
		return -1
	}

	time, err := strconv.ParseInt(parts[0], 10, 64)
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
	entry = entry[:strings.Index(entry, "]")] + "|" + entry[strings.Index(entry, "]") + 1:]
	var parts []string = strings.Split(entry, "|")

	if len(parts) < 2 {
		return GEN_ERROR
	}

	// Remove all after END_ENTRY if it exists
	var end_entry_exists bool = strings.Contains(parts[1], END_ENTRY)

	parts = strings.Split(parts[1], END_ENTRY)

	if end_entry_exists {
		parts[0] += END_ENTRY
	}

	return parts[0]
}
