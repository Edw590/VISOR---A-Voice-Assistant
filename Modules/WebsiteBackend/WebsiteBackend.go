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

package WebsiteBackend

import (
	"Utils"
	"Utils/ModsFileInfo"
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/gorilla/websocket"
	"github.com/yousifnimah/Cryptx/CRC16"
)

// Website Backend //

const MAX_CLIENTS int = 100

const PONG_WAIT = 120 * time.Second // Allow X time before considering the client unreachable
const PING_PERIOD = 60 * time.Second // Must be less than PONG_WAIT

var channels_GL [MAX_CLIENTS]chan []byte = [MAX_CLIENTS]chan []byte{}
var used_channels_GL [MAX_CLIENTS]bool = [MAX_CLIENTS]bool{}

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_8.Active_device_IDs = nil

	go func() {
		for {
			var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_WebsiteBackend, 0, -1)
			if comms_map == nil {
				return
			}

			var map_value []byte = comms_map["Message"].([]byte)
			if map_value == nil {
				continue
			}

			// Send the message to all clients. Their handler will decide if the message is for them or not.
			for i := 0; i < MAX_CLIENTS; i++ {
				if used_channels_GL[i] {
					// Ignore the panic of writing to a closed channel. Not sure why that happens.
					Tcef.Tcef{
						Try: func() {
							channels_GL[i] <- map_value
						},
					}.Do()
				}
			}
		}
	}()

	var srv *http.Server = nil
	go func() {
		Tcef.Tcef{
			Try: func() {
				// Try to register. If it's already registered, ignore the panic.
				http.HandleFunc("/ws", basicAuth(webSocketsHandler))
			},
		}.Do()

		//log.Println("Server running on port 3234")
		srv = &http.Server{Addr: ":3234"}
		err := srv.ListenAndServeTLS(getModUserInfo().Crt_file, getModUserInfo().Key_file)
		if err != nil {
			Utils.LogLnError(err)

			*module_stop = true
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	for {
		if Utils.WaitWithStopDATETIME(module_stop, -1) {
			_ = srv.Shutdown(ctx)

			for i := 0; i < MAX_CLIENTS; i++ {
				unregisterChannel(i)
			}

			Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_8.Active_device_IDs = nil

			return
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func webSocketsHandler(w http.ResponseWriter, r *http.Request) {
	Utils.LogLnInfo("WebSocketsHandler called")

	var channel_num int = registerChannel()
	if channel_num == -1 {
		Utils.LogLnWarning("No available channels")

		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Utils.LogLnError(err)

		return
	}
	defer conn.Close()

	var device_id string = ""

	_ = conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	conn.SetPongHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(PONG_WAIT))

		return nil
	})

	ticker := time.NewTicker(PING_PERIOD)
	defer ticker.Stop()

	var mutex sync.Mutex

	sendData := func(message_type int, bytes []byte) error {
		mutex.Lock()
		defer mutex.Unlock()

		if err = conn.WriteMessage(message_type, bytes); err != nil {
			return err
		}

		return err
	}

	// Sender
	go func() {
		for {
			select {
				case <- ticker.C:
					if sendData(websocket.PingMessage, nil) != nil {
						Utils.LogLnError(err)

						// If it wasn't possible to ping the client, close the connection.
						_ = conn.Close()

						return
					}
				case message := <- channels_GL[channel_num]:
					if message == nil {
						return
					}

					var message_str string = string(message)
					var msg_device_id string = strings.Split(message_str, "|")[0]
					if msg_device_id != device_id && msg_device_id != "3234_ALL" {
						continue
					}

					var index_bar int = strings.Index(message_str, "|")
					var truncated_msg []byte = message[index_bar + 1:]

					if err = sendData(websocket.BinaryMessage, Utils.CompressBytes(truncated_msg)); err == nil {
						Utils.LogfDebug("Message sent 2. Length: %d; CRC16: %d; To: %s; On: %s\n", len(truncated_msg),
							CRC16.Result(truncated_msg, "CCIT_ZERO"),
							truncated_msg[:strings.Index(string(truncated_msg), "|")], device_id)
					} else {
						Utils.LogLnError(err)
					}
			}
		}
	}()

	// Receiver
	var first_message bool = true
	for {
		// Read message from WebSocket
		message_type, message, err := conn.ReadMessage()
		if err != nil {
			Utils.LogLnError(err)

			break
		}

		if message_type != websocket.BinaryMessage {
			continue
		}
		message = Utils.DecompressBytes(message)

		if first_message {
			first_message = false

			device_id = string(message)

			Utils.LogLnInfo(device_id)

			addActiveDeviceID(device_id)

			continue
		}

		var message_str string = string(message)
		var index_bar int = strings.Index(message_str, "|")
		var index_2nd_bar int = strings.Index(message_str[index_bar + 1:], "|")
		if index_bar == -1 || index_2nd_bar == -1 {
			continue
		}

		// Print received message
		Utils.LogfInfo("Received: %s; From: %s\n", message[:index_bar + index_2nd_bar + 2], device_id)

		var message_parts []string = strings.Split(message_str, "|")
		if len(message_parts) < 3 {
			continue
		}
		var msg_to string = message_parts[0]
		var type_ string = message_parts[1]
		var bytes []byte = message[index_bar + index_2nd_bar + 2:]

		var partial_resp []byte = handleMessage(type_, bytes)
		if msg_to == "N" {
			Utils.LogLnDebug("Returning no response")
		} else {
			var response []byte = []byte(msg_to + "|")
			response = append(response, partial_resp...)
			response = Utils.CompressBytes(response)

			if err = sendData(websocket.BinaryMessage, response); err == nil {
				Utils.LogfDebug("Message sent 1. Length: %d; CRC16: %d; To: %s; On: %s\n", len(response),
					CRC16.Result(response, "CCIT_ZERO"), msg_to, device_id)
			} else {
				Utils.LogLnError(err)
			}
		}
	}

	Utils.LogLnInfo(device_id)

	removeActiveDeviceID(device_id)

	unregisterChannel(channel_num)

	Utils.LogLnInfo("WebSocketsHandler ended")
}

func handleMessage(type_ string, bytes []byte) []byte {
	switch type_ {
		case "Echo":
			// Echo the message.
			// Example: "Hello world!"
			// Returns: the same message
			return bytes
		case "Email":
			// Send an email.
			// Example: "email_to@gmail.com|an EML file"
			// Returns: nothing
			_ = Utils.QueueEmailEMAIL(Utils.EmailInfo{
				Mail_to: strings.Split(string(bytes), "|")[0],
				Eml:     string(bytes[strings.Index(string(bytes), "|") + 1:]),
			})
		case "File":
			// Get a file from the website.
			// Example: "true to get CRC16 checksum, false to get file contents|partial_path"
			// Returns: a CRC16 checksum or a file
			var bytes_split []string = strings.Split(string(bytes), "|")
			var get_crc16 bool = bytes_split[0] == "true"
			var partial_path string = bytes_split[1]

			var p_file_contents *string = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, partial_path).ReadTextFile()
			if p_file_contents == nil {
				Utils.LogLnError(partial_path)

				break
			}

			if get_crc16 {
				return getCRC16([]byte(*p_file_contents))
			} else {
				return []byte(*p_file_contents)
			}
		case "GPT":
			// Send a text to be processed by the GPT model or redirected to the right client.
			// Example: "'process', 'redirect' or 'models'|in case of processing, the data or nothing to just get the
			// return value; in case of redirecting and models, a string"
			// Returns: in case of processing, the GPT Communicator module state; in case of redirecting and models,
			// nothing
			var bytes_str string = string(bytes)
			var action string = strings.Split(bytes_str, "|")[0]

			var data []byte = bytes[strings.Index(bytes_str, "|")+1:]

			switch action {
				case "process":
					var ret []byte = []byte(strconv.Itoa(int(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.State)))

					if len(data) > 0 {
						Utils.SendToModChannel(Utils.NUM_MOD_GPTCommunicator, 2, "Process", data)
					}

					return ret
				case "redirect":
					Utils.SendToModChannel(Utils.NUM_MOD_GPTCommunicator, 0, "Redirect", string(data))
				case "models":
					Utils.SendToModChannel(Utils.NUM_MOD_GPTCommunicator, 1, "Models", string(data))
				default:
					// Nothing
			}
		case "G_S":
			// Get settings.
			// Example: "true to get file contents, false to get CRC16 checksum|one of the strings below"
			// Allowed strings: one of the ones on the switch statement
			// Returns: the user settings in JSON format
			var bytes_split []string = strings.Split(string(bytes), "|")
			var get_json bool = bytes_split[0] == "true"
			var json_origin string = bytes_split[1]
			var settings string = ""
			switch json_origin {
				case "US":
					settings = *Utils.ToJsonGENERAL(*Utils.GetUserSettings(Utils.LOCK_UNLOCK))
				case "Weather":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_6.Weather)
				case "News":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_6.News)
				case "GPTMem":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.Memories)
				case "GPTSessions":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.Sessions)
				case "GManTok":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Token)
				case "GManCals":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Calendars)
				case "GManEvents":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Events)
				case "GManTasks":
					settings = *Utils.ToJsonGENERAL(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Tasks)
				case "GManTokVal":
					settings = strconv.FormatBool(!Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Token_invalid)
				default:
					Utils.LogLnError(json_origin)
			}
			if get_json {
				return []byte(settings)
			} else {
				return getCRC16([]byte(settings))
			}
		case "S_S":
			// Set settings.
			// Example: "one of the strings below|a JSON string"
			// Allowed strings: one of the ones on the switch statement
			// Returns: nothing
			var bytes_split []string = strings.Split(string(bytes), "|")
			var json_dest string = bytes_split[0]
			var settings string = string(bytes[strings.Index(string(bytes), "|")+1:])
			switch json_dest {
				case "US":
					var user_settings Utils.UserSettings
					_ = Utils.FromJsonGENERAL([]byte(settings), &user_settings)
					if user_settings.General.Website_domain != "" &&
							user_settings.General.Website_pw != "" &&
							user_settings.General.User_email_addr != "" {
						// Only accept the new settings if as a start the website information exists (or else the
						// client will be locked out of the server surely), but also if the user email is set (that will
						// mean the settings are not empty).
						*Utils.GetUserSettings(Utils.LOCK_UNLOCK) = user_settings
					}
				case "GPTMem":
					settings = strings.Replace(settings, "\\r", "", -1)
					_ = Utils.FromJsonGENERAL([]byte(settings), &Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.Memories)
				case "GManTok":
					Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Token = settings
				case "GPTSession":
					var instructions []string = strings.Split(settings, "\000")
					var session_id string = instructions[0]
					var action string = instructions[1]
					if action == "delete" {
						delete(Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.Sessions, session_id)
					} else if action == "rename" {
						Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_7.Sessions[session_id].Name = instructions[2]
					}
				default:
					Utils.LogLnError(json_dest)
			}
		case "GMan":
			// Add a calendar event or task, or enable or disable usage of specific calendar IDs.
			// Example: "'event', 'task' or 'calendar'|the JSON of the GEvent or GTask struct; or for the 3rd case, the
			// ID of the calendar|true to enable, false to disable"
			// Returns: nothing
			var bytes_str string = string(bytes)
			var action string = strings.Split(bytes_str, "|")[0]

			var data []byte = bytes[strings.Index(bytes_str, "|")+1:]

			switch action {
				case "event":
					Utils.SendToModChannel(Utils.NUM_MOD_GoogleManager, 0, "Event", data)
				case "task":
					Utils.SendToModChannel(Utils.NUM_MOD_GoogleManager, 0, "Task", data)
				case "calendar":
					var data_str string = string(data)
					var calendar_id string = strings.Split(data_str, "|")[0]
					var enable bool = strings.Split(data_str, "|")[1] == "true"
					var calendars map[string]ModsFileInfo.GCalendar = Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Calendars
					for id := range Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_14.Calendars {
						if id == calendar_id {
							calendars[id] = ModsFileInfo.GCalendar{
								Title:   calendars[id].Title,
								Enabled: enable,
							}

							break
						}
					}
				default:
					// Nothing
			}
	}

	return nil
}

func registerChannel() int {
	for i := 0; i < MAX_CLIENTS; i++ {
		if !used_channels_GL[i] {
			channels_GL[i] = make(chan []byte)
			used_channels_GL[i] = true

			return i
		}
	}

	return -1
}

func unregisterChannel(channel_num int) {
	if channel_num >= 0 && channel_num < MAX_CLIENTS && channels_GL[channel_num] != nil {
		used_channels_GL[channel_num] = false
		if channels_GL[channel_num] != nil {
			close(channels_GL[channel_num])
			channels_GL[channel_num] = nil
		}
	}
}

func getCRC16(bytes []byte) []byte {
	var crc16 uint16 = CRC16.Result(bytes, "CCIT_ZERO")
	var crc16_bytes []byte = make([]byte, 2)
	crc16_bytes[0] = byte(crc16 >> 8)
	crc16_bytes[1] = byte(crc16)

	return crc16_bytes
}

func getModUserInfo() *ModsFileInfo.Mod8UserInfo {
	return &Utils.GetUserSettings(Utils.LOCK_UNLOCK).WebsiteBackend
}

func addActiveDeviceID(device_id string) {
	Utils.GetGenSettings(Utils.ONLY_LOCK).MOD_8.Active_device_IDs =
		append(Utils.GetGenSettings(Utils.ONLY_UNLOCK).MOD_8.Active_device_IDs, device_id)
}

func removeActiveDeviceID(device_id string) {
	for i, id := range Utils.GetGenSettings(Utils.ONLY_LOCK).MOD_8.Active_device_IDs {
		if id == device_id {
			Utils.DelElemSLICES(&Utils.GetGenSettings(Utils.DONT_LOCK_UNLOCK).MOD_8.Active_device_IDs, i)

			break
		}
	}
	Utils.GetGenSettings(Utils.ONLY_UNLOCK)
}
