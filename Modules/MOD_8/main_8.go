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
	"fmt"
	"log"
	"net/http"
)

// Modules Manager //

var srv_GL *http.Server = nil

type _MGIModSpecInfo any
var (
	realMain Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo[_MGIModSpecInfo]
)
func Start(module *Utils.Module) {Utils.ModStartup[_MGIModSpecInfo](Utils.NUM_MOD_WebsiteBackend, realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo[_MGIModSpecInfo])

		go func() {
			if srv_GL == nil {
				http.HandleFunc("/submit-form", formHandler)

				log.Println("Server running on port 8080")
				srv_GL = &http.Server{Addr: ":8080"}
				_ = srv_GL.ListenAndServe()
				srv_GL = nil
			}
		}()

		for {
			if Utils.WaitWithStop(module_stop, 1000000000) {
				return
			}
		}
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "ParseForm() err: " + err.Error(), http.StatusInternalServerError)

			return
		}

		var type_ string = r.FormValue("type")
		var text1 string = r.FormValue("text1")
		var text2 string = r.FormValue("text2")
		var text3 string = r.FormValue("text3")
		// Process form data here
		_, _ = fmt.Fprintf(w, "Received type: %s\ntext1: %s\n text2: %s\n text3: %s\n", type_, text1, text2, text3)

		if type_ == "GPT" {
			_ = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_GPTCommunicator).Add2(false, "to_process", "test.txt").
				WriteTextFile(text1, false)
		} else if type_ == "Email" {
			_ = Utils.QueueEmailEMAIL(Utils.EmailInfo{
				Mail_to:    text1,
				Eml:        text2,
			})
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
