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

package GPTComm

import (
	"Utils"
	"strconv"
	"strings"
	"time"
)

type File struct {
	Is_image bool
	Size int
	Contents []byte
}

var files_to_send_GL []File = nil

const MODEL_TYPE_TEXT string = "T"
const MODEL_TYPE_VISION string = "T+V"

const ROLE_USER string = "USER"
const ROLE_TOOL string = "TOOL"
const SESSION_TYPE_NEW string = "NEW"
const SESSION_TYPE_TEMP string = "TEMP"
const SESSION_TYPE_ACTIVE string = "ACTIVE"
/*
SendText sends the given text to the LLM model along with the selected files added by AddFileToSend().

After this function returns, all files have been removed from the list of files to send.

-----------------------------------------------------------

– Params:
  - text – the text to send or an empty string to just get the return value
  - session_type – one of the SESSION_TYPE_-started constants or a session ID (ignored in case `text` is empty)
  - role – the role of the message (one of the ROLE_-started constants)
  - more_coming – whether more messages are coming or not and the LLM should wait for them before replying

– Returns:
  - the state of the GPT Communicator module or -1 in case of errors
*/
func SendText(text string, session_type string, role string, more_coming bool) int32 {
	var message []byte = []byte("GPT|process|")
	if text != "" {
		var curr_location string = Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_12.User_location.Curr_location
		var date_time string = time.Now().Weekday().String() + " " + time.Now().Format("2006-01-02 15:04")

		var new_text string = text
		if role == ROLE_USER {
			new_text = "[current user location: " + curr_location + " | date/time: " + date_time + "]" + text
		}

		var images []File = nil
		var audios []File = nil
		for _, file := range files_to_send_GL {
			if file.Is_image {
				images = append(images, file)
			} else {
				audios = append(audios, file)
			}
		}

		///////////////////////////////////////////
		var metadata string = "["

		metadata += strconv.Itoa(len(images)) + "|" + strconv.Itoa(len(audios)) + "|"
		for _, file := range images {
			metadata += strconv.Itoa(file.Size) + "|"
		}
		for _, file := range audios {
			metadata += strconv.Itoa(file.Size) + "|"
		}

		metadata += Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id + "|" + session_type + "|" + role + "|" +
			strconv.FormatBool(more_coming)

		metadata += "]"
		///////////////////////////////////////////

		message = append(message, metadata + new_text + "\x00"...)

		for _, file := range images {
			message = append(message, file.Contents...)
		}
		for _, file := range audios {
			message = append(message, file.Contents...)
		}
	}

	var successful_queue bool = Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, 1, message)
	files_to_send_GL = nil
	if !successful_queue {
		return -1
	}

	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 1, 10)
	if comms_map == nil {
		return -1
	}

	var response []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	ret, _ := strconv.Atoi(string(response))

	return int32(ret)
}

/*
AddFileToSend adds a file to the list of files to send.

-----------------------------------------------------------

– Params:
  - is_image – whether the file is an image or not (an audio then)
  - contents – the raw contents of the file
 */
func AddFileToSend(is_image bool, contents []byte) {
	files_to_send_GL = append(files_to_send_GL, File{
		Is_image: is_image,
		Size:     len(contents),
		Contents: contents,
	})
}

/*
GetModuleState gets the state of the GPT Communicator module.

DON'T delete the function. It's useful for use in a thread other than the one used to send text.
*/
func GetModuleState() int32 {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, 5, []byte("GPT|process|")) {
		return -1
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 5, 10)
	if comms_map == nil {
		return -1
	}

	var response []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	ret, _ := strconv.Atoi(string(response))

	return int32(ret)
}

/*
GetMemories gets the memories from the GPT.

-----------------------------------------------------------

– Returns:
  - the memories separated by new lines
 */
func GetMemories() string {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, 3, []byte("G_S|true|GPTMem")) {
		return ""
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 3, 10)
	if comms_map == nil {
		return ""
	}

	var json_bytes []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	var memories []string
	if err := Utils.FromJsonGENERAL(json_bytes, &memories); err != nil {
		return ""
	}

	return strings.Join(memories, "\n")
}

/*
SetMemories sets the memories in the GPT.

-----------------------------------------------------------

– Params:
  - json – the memories separated by new lines
 */
func SetMemories(memories_str string) bool {
	var memories []string = strings.Split(memories_str, "\n")

	var message []byte = []byte("S_S|GPTMem|")
	message = append(message, *Utils.ToJsonGENERAL(memories)...)

	return Utils.QueueNoResponseMessageSERVER(message)
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
  - the entry or an empty entry with time = -1 if it doesn't exist or there's no Internet connection
*/
func getEntry(time int64, num int) *_Entry {
	if !Utils.IsCommunicatorConnectedSERVER() {
		return &_Entry{
			device_id: "",
			text:      "",
			time_ms:   -1,
		}
	}

	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GPTComm, 2, []byte("File|false|gpt_text.txt")) {
		return &_Entry{
			device_id: "",
			text:      "",
			time_ms:   -1,
		}
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GPTComm, 2, 10)
	if comms_map == nil {
		return &_Entry{
			device_id: "",
			text:      "",
			time_ms:   -1,
		}
	}

	var file_contents string = string(comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte))
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

	time_int, err := strconv.ParseInt(time_str, 10, 64)
	if err != nil {
		return -1
	}

	return time_int
}

/*
getTextFromEntry gets the text from the entry.

-----------------------------------------------------------

– Params:
  - entry – the entry

– Returns:
  - the text
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
