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

package EmailSender

import (
	"Utils/ModsFileInfo"
	"errors"
	"os"
	"strings"
	"time"

	"Utils"
)

// _MAX_EMAILS_HOUR is the maximum number of emails that can be sent per hour according to Google, which are 20. But
// we'll go with 15 to be safe about emails sent without this module's help (error emails).
const _MAX_EMAILS_HOUR = 15

const _TIME_SLEEP_S int = 5

type emailSent struct {
	email  string
	time_s int64
}

var modDirsInfo_GL Utils.ModDirsInfo
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	var to_send_dir Utils.GPath = modDirsInfo_GL.UserData.Add2(true, Utils.TO_SEND_REL_FOLDER)
	if !to_send_dir.Exists() {
		_ = to_send_dir.Create(false)
	}

	//log.Println("Checking for emails to send in \"" + to_send_dir.GPathToStringConversion() + "\"...")

	for {
		var files_to_send []Utils.FileInfo = to_send_dir.GetFileList()
		if files_to_send == nil {
			panic(errors.New("couldn't read directory \"" + to_send_dir.GPathToStringConversion() + "\""))
		}

		var last_file_sent emailSent
		for len(files_to_send) > 0 {
			file_to_send, idx_to_remove := Utils.GetOldestFileFILESDIRS(files_to_send)
			if *file_to_send.GPath.ReadTextFile() == last_file_sent.email && time.Now().Unix() - last_file_sent.time_s < 60 {
				// Don't send the same email twice or more in a row.
				Utils.DelElemSLICES(&files_to_send, idx_to_remove)

				continue
			}

			// ... and send it.
			var mail_to string = strings.TrimSuffix(file_to_send.Name, ".eml")
			mail_to = mail_to[Utils.RAND_STR_LEN:]

			//log.Println("--------------------")
			//log.Println("Sending email file " + file_to_send.Name + " to " + mail_to + "...")

			if !reachedMaxEmailsHour() {
				if err := Utils.SendEmailEMAIL(*file_to_send.GPath.ReadTextFile(), mail_to, false); err == nil {
					if time.Now().Hour() != getModGenSettings().Hour {
						getModGenSettings().Hour = time.Now().Hour()
						getModGenSettings().Num_emails_hour = 0
					}
					getModGenSettings().Num_emails_hour++
					//log.Println("Email sent successfully.")

					last_file_sent.email = *file_to_send.GPath.ReadTextFile()
					last_file_sent.time_s = time.Now().Unix()

					// Remove the file
					Utils.DelElemSLICES(&files_to_send, idx_to_remove)
					if os.Remove(file_to_send.GPath.GPathToStringConversion()) == nil {
						//log.Println("File deleted successfully.")
					} else {
						//log.Println("Error deleting file.")
					}
				} else {
					//log.Println("Error sending email with error\n" + Utils.GetFullErrorMsgGENERAL(err))

					panic(err)
				}
			} else {
				//log.Println("The maximum number of emails per hour has been reached (" +
				//	strconv.Itoa(_MAX_EMAILS_HOUR) + "). Waiting for the next hour...")

				goto end_loop
			}

			// No mega fast email spamming - don't want the account blocked.
			if Utils.WaitWithStopDATETIME(module_stop, 1) {
				return
			}
		}

		end_loop:

		if Utils.WaitWithStopDATETIME(module_stop, _TIME_SLEEP_S) {
			return
		}
	}
}

/*
reachedMaxEmailsHour returns true if the maximum number of emails per hour has been reached.

-----------------------------------------------------------

– Returns:
  - true if the maximum number of emails per hour has been reached, false otherwise.
 */
func reachedMaxEmailsHour() bool {
	return getModGenSettings().Num_emails_hour >= _MAX_EMAILS_HOUR &&
		time.Now().Hour() == getModGenSettings().Hour
}

func getModGenSettings() *ModsFileInfo.Mod5GenInfo {
	return &Utils.GetGenSettings().MOD_5
}
