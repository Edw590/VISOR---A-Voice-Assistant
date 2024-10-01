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
)

const COMMS_MAP_SRV_KEY string = "SrvComm"

var srvComm_gen_ch_in chan []byte = make(chan []byte)
var srvComm_gen_ch_out chan []byte = make(chan []byte, 1000)
var srvComm_stop bool = false

/*
StartCommunicatorSERVER starts the communicator.

This function does not return until the communicator is stopped.
*/
func StartCommunicatorSERVER() {
	// Define the WebSocket server address
	u := url.URL{Scheme: "wss", Host: User_settings_GL.PersonalConsts.Website_domain, Path: "/ws"}
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
		},
	}

	// Establish WebSocket connection
	c, r, err := dialer.Dial(u.String(), headers)
	if err != nil {
		log.Println("Response:", r)
		log.Println("Dial error:", err)

		return
	}
	defer c.Close()

	go func() {
		for {
			message_type, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				srvComm_stop = true

				return
			}
			if message_type != websocket.BinaryMessage {
				continue
			}

			//log.Printf("Received message: %s", string(message))

			var msg_to string = strings.Split(string(message), "|")[0]
			var index_bar int = strings.Index(string(message), "|")
			var truncated_msg []byte = message[index_bar + 1:]
			if msg_to == "GEN" {
				srvComm_gen_ch_in <- truncated_msg

				continue
			}
			var msg_to_split []string = strings.Split(msg_to, "_")
			var to_mod bool = msg_to_split[0] == "MOD"
			num, err := strconv.Atoi(msg_to_split[1])
			if err != nil {
				log.Println("Error converting module number:", err)

				continue
			}

			if to_mod {
				ModsCommsChannels_GL[num] <- map[string]any{COMMS_MAP_SRV_KEY: truncated_msg}
			} else {
				LibsCommsChannels_GL[num] <- map[string]any{COMMS_MAP_SRV_KEY: truncated_msg}
			}
		}
	}()

	go func() {
		for {
			var message []byte = <-srvComm_gen_ch_out

			err := c.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				log.Println("Write error:", err)
				srvComm_stop = true

				return
			}
			//log.Printf("Sent message: %s", message)
		}
	}()

	log.Println("Communicator started")

	for {
		if WaitWithStopTIMEDATE(&srvComm_stop, 1000000000) {
			_ = c.Close()

			return
		}
	}
}

/*
GetGeneralMessageSERVER gets a general message from the server.

The message is sent by QueueGeneralMessageSERVER().

If no message is available, the function will wait until a message is received.
*/
func GetGeneralMessageSERVER() []byte {
	return <-srvComm_gen_ch_in
}

/*
QueueGeneralMessageSERVER queues a general message to be sent to the server.

It is received by GetGeneralMessageSERVER().

-----------------------------------------------------------

– Params:
  - message – the message to be sent
*/
func QueueGeneralMessageSERVER(message []byte) {
	var new_msg []byte = append([]byte("GEN|"), message...)
	srvComm_gen_ch_out <- new_msg
}

/*
QueueMessageSERVER queues a message to be sent to the server.

-----------------------------------------------------------

– Params:
  - is_mod – true if this function was called from a module, false if it was called from a library
  - mod_num – the number of the module or library that called this function
  - message – the message to be sent
*/
func QueueMessageSERVER(is_mod bool, num int, message []byte) {
	var mod_lib string = "MOD"
	if !is_mod {
		mod_lib = "LIB"
	}
	var message_str string = mod_lib + "_" + strconv.Itoa(num) + "|"
	var new_msg []byte = append([]byte(message_str), message...)
	srvComm_gen_ch_out <- new_msg
}

/*
QueueNoResponseMessageSERVER queues a message to be sent to the server without expecting a response.

-----------------------------------------------------------

– Params:
  - message – the message to be sent
*/
func QueueNoResponseMessageSERVER(message []byte) {
	var new_msg []byte = append([]byte("NONE|"), message...)
	srvComm_gen_ch_out <- new_msg
}

/*
StopCommunicatorSERVER stops the communicator.
*/
func StopCommunicatorSERVER() {
	srvComm_stop = true
}
