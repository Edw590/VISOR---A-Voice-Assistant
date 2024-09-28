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
	"Utils"
	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/yousifnimah/Cryptx/CRC16"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Website Backend //

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		var srv *http.Server = nil
		go func() {
			Tcef.Tcef{
				Try: func() {
					// Try to register. If it's already registered, ignore the panic.
					http.HandleFunc("/submit-form", formHandler)
					http.HandleFunc("/file/", handleGetRequest)
				},
			}.Do()

			//log.Println("Server running on port 8080")
			srv = &http.Server{Addr: ":8080"}
			_ = srv.ListenAndServe()
		}()

		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1000000000) {
				_ = srv.Shutdown(nil)

				return
			}
		}
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request metho: " + r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Println("ParseForm() err: " + err.Error())
		http.Error(w, "ParseForm() err: " + err.Error(), http.StatusInternalServerError)

		return
	}

	var type_ string = r.FormValue("type")
	var text1 string = r.FormValue("text1")
	//var text2 string = r.FormValue("text2")
	file, file_header, err := r.FormFile("file")

	//log.Println("Form:", r)
	//log.Println("Type: " + type_)
	//log.Println("Text1: " + text1)
	//log.Println("Text2: " + text2)

	var file_bytes []byte = nil
	if err == nil {
		file_bytes = make([]byte, file_header.Size)
		_, _ = file.Read(file_bytes)
	}
	//log.Println("File:", file_bytes)

	switch type_ {
		case "GPT":
			log.Println("GPT")
			// File: the text to process, compressed
			// Returns: nothing
			_ = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_GPTCommunicator).Add2(false, "to_process", "test.txt").
				WriteTextFile(Utils.DecompressString(file_bytes), false)
		case "Email":
			log.Println("Email")
			// Text1: the email address to send to
			// File: the EML file to send, compressed
			// Returns: nothing
			_ = Utils.QueueEmailEMAIL(Utils.EmailInfo{
				Mail_to: text1,
				Eml:     Utils.DecompressString(file_bytes),
			})
		case "UserLocator":
			// Text1: the device ID
			// File: the JSON data to write, compressed
			// Returns: nothing
			log.Println("UserLocator")
			_ = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_UserLocator).Add2(false, "devices", text1 + ".json").
				WriteTextFile(Utils.DecompressString(file_bytes), false)
		default:
			// Do nothing
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	url_str, _ := url.Parse(r.URL.String())
	params, _ := url.ParseQuery(url_str.RawQuery)

	var get_crc16 bool = params.Get("crc") == "true"

	// Remove "/file/" and any parameters from the URL
	var file_path string = r.URL.Path
	file_path = strings.TrimPrefix(file_path, "/file/")
	file_path = strings.Split(file_path, "?")[0]
	if strings.HasSuffix(file_path, "/") {
		http.Error(w, "Not a file", http.StatusNotFound)

		return
	}
	var p_file_contents *string = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, file_path).ReadTextFile()
	if p_file_contents == nil {
		http.Error(w, "File not found", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	if get_crc16 {
		var crc16 uint16 = CRC16.Result([]byte(*p_file_contents), "CCIT_ZERO")
		var crc16_bytes []byte = make([]byte, 2)
		crc16_bytes[0] = byte(crc16 >> 8)
		crc16_bytes[1] = byte(crc16)
		_, _ = w.Write(crc16_bytes)
	} else {
		_, _ = w.Write([]byte(*p_file_contents))
	}
}
