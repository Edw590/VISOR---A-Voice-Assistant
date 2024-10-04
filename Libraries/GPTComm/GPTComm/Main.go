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

package GPTComm

import (
	"Utils"
	"strconv"
	"strings"
)

/*
SendText sends the given text to the LLM model.

-----------------------------------------------------------

– Params:
  - text – the text to send
*/
func SendText(text string) {
	var message []byte = []byte("GPT|")
	message = append(message, Utils.CompressString("[" + Utils.User_settings_GL.PersonalConsts.Device_ID + "]" + text)...)
	Utils.QueueNoResponseMessageSERVER(message)
}

/*
getEntry gets the entry at the given number or time.

If -1 is provided on both parameters, it will return the last entry. The time parameter is prioritized over the number
parameter.

-----------------------------------------------------------

– Params:
  - time – the time of the entry or -1 if the entry is to be found by number
  - num – the number of the entry or negative to count from the end

– Returns:
  - the entry or an empty entry with time = -1 if it doesn't exist
*/
func getEntry(time int64, num int) *_Entry {
	Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, []byte("File|false|gpt_text.txt"))
	var comms_map map[string]any = <- Utils.LibsCommsChannels_GL[Utils.NUM_LIB_GPTComm]
	if comms_map == nil {
		return &_Entry{
			device_id: "",
			text:      "",
			time_ms:   -1,
		}
	}

	var file_contents string = Utils.DecompressString(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte))
	if file_contents == "" {
		return &_Entry{
			device_id: "",
			text:      "",
			time_ms:   -1,
		}
	}
	var entries []string = strings.Split(file_contents, "[3234_START:")

	if time != -1 {
		for _, entry := range entries {
			if entry == "" {
				continue
			}

			if getTimeFromEntry(entry) == time {
				return &_Entry{
					device_id: getDeviceIdFromEntry(entry),
					text:      getTextFromEntry(entry),
					time_ms:   getTimeFromEntry(entry),
				}
			}
		}
	} else {
		if len(entries) == 0 {
			return &_Entry{
				device_id: "",
				text:      "",
				time_ms:   -1,
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
			return &_Entry{
				device_id: getDeviceIdFromEntry(entry),
				text:      getTextFromEntry(entry),
				time_ms:   getTimeFromEntry(entry),
			}
		}
	}

	return &_Entry{
		device_id: "",
		text:      "",
		time_ms:   -1,
	}
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
