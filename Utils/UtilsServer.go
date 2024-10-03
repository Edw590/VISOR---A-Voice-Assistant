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

package Utils

import (
	"crypto/tls"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const COMMS_MAP_SRV_KEY string = "SrvComm"

var srvComm_gen_ch_in_GL chan []byte
var srvComm_gen_ch_out_GL chan []byte
var srvComm_stop_GL bool = false
var srvComm_started_GL bool = false

/*
StartCommunicatorSERVER starts the communicator.

This function does not return until the communicator is stopped, or returns in case the communicator is already started.

-----------------------------------------------------------

– Returns:
  - bool – true if the communicator was started or was already started, false if it an error occurred and it did not
	start
*/
func StartCommunicatorSERVER() bool {
	if srvComm_started_GL {
		return true
	}
	srvComm_started_GL = true

	srvComm_stop_GL = false
	srvComm_gen_ch_in_GL = make(chan []byte)
	srvComm_gen_ch_out_GL = make(chan []byte, 1000)
	var routines_working [2]bool

	// Define the WebSocket server address
	u := url.URL{Scheme: "wss", Host: User_settings_GL.PersonalConsts.Website_domain + ":3234", Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	// Create Basic Auth credentials (username:password)
	username := "VISOR"
	password := User_settings_GL.PersonalConsts.Website_pw
	auth := username + ":" + password
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	// Define request headers including Authorization (Basic Auth)
	headers := http.Header{}
	headers.Set("Authorization", authHeader)

	// Create a custom WebSocket dialer
	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Disable certificate verification
			ClientSessionCache: tls.NewLRUClientSessionCache(32), // Use an LRU session cache for resumption
		},
	}

	// Establish WebSocket connection
	conn, _, err := dialer.Dial(u.String(), headers)
	if err != nil {
		//log.Println("Response:", r)
		//log.Println("Dial error:", err)

		srvComm_started_GL = false

		return false
	}
	defer conn.Close()

	go func() {
		routines_working[0] = true
		for {
			message_type, message, err := conn.ReadMessage()
			if err != nil {
				//log.Println("Read error:", err)
				srvComm_stop_GL = true

				break
			}
			if message_type != websocket.BinaryMessage {
				continue
			}

			//log.Printf("Received message: %s", string(message))

			var msg_to string = strings.Split(string(message), "|")[0]
			var index_bar int = strings.Index(string(message), "|")
			var truncated_msg []byte = message[index_bar + 1:]
			if msg_to == "G" {
				srvComm_gen_ch_in_GL <- truncated_msg

				continue
			}
			var msg_to_split []string = strings.Split(msg_to, "_")
			var to_mod bool = msg_to_split[0] == "M"
			num, err := strconv.Atoi(msg_to_split[1])
			if err != nil {
				//log.Println("Error converting module number:", err)

				continue
			}

			if to_mod {
				ModsCommsChannels_GL[num] <- map[string]any{COMMS_MAP_SRV_KEY: truncated_msg}
			} else {
				LibsCommsChannels_GL[num] <- map[string]any{COMMS_MAP_SRV_KEY: truncated_msg}
			}
		}
		routines_working[0] = false
	}()

	go func() {
		var first_message bool = true
		routines_working[1] = true
		for {
			var message []byte
			if first_message {
				message = []byte(User_settings_GL.PersonalConsts.Device_ID)
				first_message = false
			} else {
				message = <- srvComm_gen_ch_out_GL
				if message == nil {
					break
				}
			}

			err := conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				//log.Println("Write error:", err)
				srvComm_stop_GL = true

				break
			}
			//log.Printf("Sent message: %s", message)
		}
		routines_working[1] = false
	}()

	log.Println("Communicator started")

	for {
		if WaitWithStopTIMEDATE(&srvComm_stop_GL, 1000000000) {
			close(srvComm_gen_ch_in_GL)
			close(srvComm_gen_ch_out_GL)
			_ = conn.Close()
			for {
				if !routines_working[0] && !routines_working[1] {
					log.Println("Communicator stopped")

					srvComm_started_GL = false

					return true
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

/*
GetGeneralMessageSERVER gets a general message from the server.

The message is sent by QueueGeneralMessageSERVER().

If no message is available, the function will wait until a message is received.
*/
func GetGeneralMessageSERVER() []byte {
	if !srvComm_started_GL {
		return nil
	}

	return <- srvComm_gen_ch_in_GL
}

/*
QueueGeneralMessageSERVER queues a general message to be sent to the server.

It is received by GetGeneralMessageSERVER().

-----------------------------------------------------------

– Params:
  - message – the message to be sent
*/
func QueueGeneralMessageSERVER(message []byte) {
	if !srvComm_started_GL {
		return
	}

	var new_msg []byte = append([]byte("G|"), message...)
	srvComm_gen_ch_out_GL <- new_msg
}

/*
QueueMessageSERVER queues a message to be sent to the server.

-----------------------------------------------------------

– Params:
  - is_mod – true if this function was called from a module, false if it was called from a library
  - num – the number of the module or library that called this function
  - message – the message to be sent
*/
func QueueMessageSERVER(is_mod bool, num int, message []byte) {
	if !srvComm_started_GL {
		return
	}

	var mod_lib string = "M"
	if !is_mod {
		mod_lib = "L"
	}
	var message_str string = mod_lib + "_" + strconv.Itoa(num) + "|"
	var new_msg []byte = append([]byte(message_str), message...)
	srvComm_gen_ch_out_GL <- new_msg
}

/*
QueueNoResponseMessageSERVER queues a message to be sent to the server without expecting a response.

-----------------------------------------------------------

– Params:
  - message – the message to be sent
*/
func QueueNoResponseMessageSERVER(message []byte) {
	if !srvComm_started_GL {
		return
	}

	var new_msg []byte = append([]byte("N|"), message...)
	srvComm_gen_ch_out_GL <- new_msg
}

/*
StopCommunicatorSERVER stops the communicator.
*/
func StopCommunicatorSERVER() {
	srvComm_stop_GL = true
}
