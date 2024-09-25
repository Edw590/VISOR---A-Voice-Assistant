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
	"crypto/md5"
	Tcef "github.com/Edw590/TryCatch-go"
	"log"
	"net/http"
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "ParseForm() err: " + err.Error(), http.StatusInternalServerError)

		return
	}

	var type_ string = r.FormValue("type")
	var text1 string = r.FormValue("text1")
	var text2 string = r.FormValue("text2")
	//var text3 string = r.FormValue("text3")

	switch type_ {
		case "GPT":
			log.Println("GPT")
			// Text1 is the text to process
			_ = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_GPTCommunicator).Add2(false, "to_process", "test.txt").
				WriteTextFile(text1, false)
		case "Email":
			log.Println("Email")
			// Text1 is the email address to send to
			// Text2 is the EML file to send
			_ = Utils.QueueEmailEMAIL(Utils.EmailInfo{
				Mail_to: text1,
				Eml: text2,
			})
		case "UserLocator":
			// Text1 is the device ID
			// Text2 is the JSON data to write
			log.Println("UserLocator")
			_ = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_UserLocator).Add2(false, "devices", text1 + ".json").
				WriteTextFile(text2, false)
		case "GET":
			// Text1 is true if it's to get a file, false if it's to get its MD5 hash
			// Text2 is the file path
			var file_bytes []byte = Utils.GetWebsiteFilesDirFILESDIRS().Add2(false, text2).ReadFile()
			if text1 == "true" {
				_, _ = w.Write(file_bytes)
			} else {
				var hash [16]byte = md5.Sum(file_bytes)
				_, _ = w.Write(hash[:])
			}
		default:
			// Do nothing
	}
}
