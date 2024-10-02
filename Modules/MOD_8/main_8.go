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

package MOD_8

import (
	"ULComm/ULComm"
	"Utils"
	"context"
	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/gorilla/websocket"
	"github.com/yousifnimah/Cryptx/CRC16"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Website Backend //

const MAX_CHANNELS int = 100

const PONG_WAIT = 60 * time.Second // Allow X time before considering the client unreachable.
const PING_PERIOD = 30 * time.Second // Must be less than PONG_WAIT.

var channels_GL [MAX_CHANNELS]chan []byte = [MAX_CHANNELS]chan []byte{}
var used_channels_GL [MAX_CHANNELS]bool = [MAX_CHANNELS]bool{}

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		go func() {
			for {
				var comms_map map[string]any = <- Utils.ModsCommsChannels_GL[Utils.NUM_MOD_WebsiteBackend]
				if comms_map == nil {
					return
				}

				var message []byte = comms_map["Message"].([]byte)
				for i := 0; i < MAX_CHANNELS; i++ {
					if used_channels_GL[i] {
						channels_GL[i] <- message
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
			err := srv.ListenAndServeTLS(Utils.User_settings_GL.MOD_8.Cert_file, Utils.User_settings_GL.MOD_8.Key_file)
			if err != nil {
				log.Println("ListenAndServeTLS error:", err)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1000000000) {
				for i := 0; i < MAX_CHANNELS; i++ {
					if used_channels_GL[i] {
						close(channels_GL[i])
					}
				}

				if err := srv.Shutdown(ctx); err == nil {
					log.Println("Server stopped gracefully")
				} else {
					log.Println("Server shutdown error:", err)
				}

				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func webSocketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocketsHandler called")

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	conn.SetPongHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(PONG_WAIT))

		return nil
	})

	ticker := time.NewTicker(PING_PERIOD)
	defer ticker.Stop()

	var channel_num int = registerChannel()
	if channel_num == -1 {
		log.Println("No available channels")

		return
	}

	go func() {
		for {
			select {
				case <- ticker.C:
					if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
						log.Println("Ping error:", err)

						return
					}
				case message := <- channels_GL[channel_num]:
					if message == nil {
						return
					}

					if err := conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
						log.Println("Write error:", err)

						return
					} else {
						log.Printf("Message sent 2. Length: %d; CRC16: %d; Content: %s", len(message),
							CRC16.Result(message, "CCIT_ZERO"), message[:strings.Index(string(message), "|")])
					}
			}
		}
	}()

	for {
		// Read message from WebSocket
		message_type, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)

			break
		}

		if message_type != websocket.BinaryMessage {
			continue
		}
		var message_str string = string(message)
		var index_bar int = strings.Index(message_str, "|")
		var index_2nd_bar int = strings.Index(message_str[index_bar + 1:], "|")
		if index_bar == -1 || index_2nd_bar == -1 {
			continue
		}

		// Print received message
		log.Printf("Received: %s", message[:index_bar + index_2nd_bar + 2])

		var message_parts []string = strings.Split(message_str, "|")
		if len(message_parts) < 3 {
			continue
		}
		var msg_to string = message_parts[0]
		var type_ string = message_parts[1]
		var bytes []byte = message[index_bar + index_2nd_bar + 2:]

		var partial_resp []byte = handleMessage(type_, bytes)
		if msg_to != "NONE" {
			var response []byte = []byte(msg_to + "|")
			response = append(response, partial_resp...)

			if err := conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
				log.Println("Write error:", err)

				break
			} else {
				log.Printf("Message sent 1. Length: %d; CRC16: %d; Content: %s", len(response),
					CRC16.Result(response, "CCIT_ZERO"), response[:strings.Index(string(response), "|")])
			}
		} else {
			log.Println("Giving no response")
		}
	}

	unregisterChannel(channel_num)

	log.Println("WebSocketsHandler ended")
}

func handleMessage(type_ string, bytes []byte) []byte {
	switch type_ {
		case "Echo":
			// Example: "Hello world!"
			return bytes
		case "GPT":
			// Example: a compressed string
			if err := Utils.GetUserDataDirMODULES(Utils.NUM_MOD_GPTCommunicator).Add2(false, "to_process", "test.txt").
				WriteTextFile(Utils.DecompressString(bytes), false); err == nil {
					log.Println("File written")

					return []byte("OK")
			} else {
				log.Println("Error writing file:", err)
			}
		case "Email":
			// Example: "email_to@gmail.com|" + a compressed EML file
			if err := Utils.QueueEmailEMAIL(Utils.EmailInfo{
				Mail_to: strings.Split(string(bytes), "|")[0],
				Eml:     Utils.DecompressString(bytes[strings.Index(string(bytes), "|") + 1:]),
			}); err == nil {
				log.Println("Email queued")

				return []byte("OK")
			} else {
				log.Println("Error queuing email:", err)
			}
		case "UserLocator":
			// Example: "Device_Id|" + a compressed JSON file, or "Device_Id|time=" + timestamp in ms
			var bytes_split []string = strings.Split(string(bytes), "|")
			var device_id string = bytes_split[0]
			// If the message is a timestamp, just update the last communication time of the device.
			if strings.HasPrefix(bytes_split[1], "time=") {
				timestamp, err := strconv.ParseInt(strings.Split(bytes_split[1], "time=")[1], 10, 64)
				if err != nil {
					log.Println("Error parsing timestamp:", err)

					break
				}

				var file Utils.GPath = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_UserLocator).Add2(false, "devices", device_id + ".json")
				var p_json *string = file.ReadTextFile()
				if p_json == nil {
					log.Println("Error reading file")

					break
				}

				var device_info ULComm.DeviceInfo
				if err = Utils.FromJsonGENERAL([]byte(*p_json), &device_info); err != nil {
					log.Println("Error parsing JSON:", err)

					break
				}

				device_info.Last_comm = timestamp

				if err = file.WriteTextFile(*Utils.ToJsonGENERAL(device_info), false); err == nil {
					log.Println("File written")

					return []byte("OK")
				} else {
					log.Println("Error writing file:", err)
				}
			} else {
				// If the message is a JSON file, update the device's information.
				var json string = Utils.DecompressString(bytes[strings.Index(string(bytes), "|") + 1:])
				if err := Utils.GetUserDataDirMODULES(Utils.NUM_MOD_UserLocator).
					Add2(false, "devices", device_id + ".json").WriteTextFile(json, false); err == nil {
						log.Println("File written")

						return []byte("OK")
				} else {
					log.Println("Error writing file:", err)
				}
			}
		case "File":
			// Example: "true|partial_path" or "false|partial_path", where true means to get the CRC16 checksum and
			// false means to get the file contents
			var bytes_split []string = strings.Split(string(bytes), "|")
			var get_crc16 bool = bytes_split[0] == "true"
			var partial_path string = bytes_split[1]

			var p_file_contents *string = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, partial_path).ReadTextFile()
			if p_file_contents == nil {
				log.Println("Error file not found:", partial_path)

				break
			}

			log.Println("File read")
			if get_crc16 {
				var crc16 uint16 = CRC16.Result([]byte(*p_file_contents), "CCIT_ZERO")
				var crc16_bytes []byte = make([]byte, 2)
				crc16_bytes[0] = byte(crc16 >> 8)
				crc16_bytes[1] = byte(crc16)
				return crc16_bytes
			} else {
				return Utils.CompressString(*p_file_contents)
			}
	}

	return []byte("ERROR")
}

func registerChannel() int {
	for i := 0; i < MAX_CHANNELS; i++ {
		if !used_channels_GL[i] {
			channels_GL[i] = make(chan []byte)
			used_channels_GL[i] = true

			return i
		}
	}

	return -1
}

func unregisterChannel(channel_num int) {
	if channel_num >= 0 && channel_num < MAX_CHANNELS {
		used_channels_GL[channel_num] = false
		channels_GL[channel_num] = nil
	}
}
