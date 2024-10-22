/*******************************************************************************
 * Copyright 2023-2024 The V.I.S.O.R. authors
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
	"bytes"
	"errors"
	"mime/quotedprintable"
	"os"
	"strings"
)

// EmailInfo is the info needed to send an email through QueueEmail().
type EmailInfo struct {
	// Sender name (can be anything)
	Sender string
	// Mail_to is the email address to send the email to.
	Mail_to string
	// Subject of the email.
	Subject string
	// Html is the HTML body of the email.
	Html    string
	// Multipart is the list of multipart items to attach to the email aside from the main HTML.
	Multiparts []Multipart
	// Eml is the EML file of the email. If it's not empty, the email will be sent with this file and all the fields
	// except Mail_to will be ignored.
	Eml     string
}

// Multipart is an item to attach to an email as described in RFC 1521.
type Multipart struct {
	Content_type              string
	Content_transfer_encoding string
	Content_id                string
	Body                      string
}

const (
	MODEL_INFO_MSG_BODY_EMAIL  string = "|3234_MSG_BODY|"
	MODEL_INFO_DATE_TIME_EMAIL string = "|3234_DATE_TIME|"

	MODEL_RSS_ENTRY_TITLE_EMAIL       string = "|3234_ENTRY_TITLE|"
	MODEL_RSS_ENTRY_AUTHOR_EMAIL      string = "|3234_ENTRY_AUTHOR|"
	MODEL_RSS_ENTRY_DESCRIPTION_EMAIL string = "|3234_ENTRY_DESCRIPTION|"
	MODEL_RSS_ENTRY_URL_EMAIL         string = "|3234_ENTRY_URL|"
	MODEL_RSS_ENTRY_PUB_DATE_EMAIL    string = "|3234_ENTRY_PUB_DATE|"
	MODEL_RSS_ENTRY_UPD_DATE_EMAIL    string = "|3234_ENTRY_UPD_DATE|"

	MODEL_YT_VIDEO_HTML_TITLE_EMAIL        string = "|3234_HTML_TITLE|"
	MODEL_YT_VIDEO_CHANNEL_NAME_EMAIL      string = "|3234_CHANNEL_NAME|"
	MODEL_YT_VIDEO_CHANNEL_CODE_EMAIL      string = "|3234_CHANNEL_CODE|"
	MODEL_YT_VIDEO_CHANNEL_IMAGE_EMAIL     string = "|3234_CHANNEL_IMAGE|"
	MODEL_YT_VIDEO_VIDEO_TITLE_EMAIL       string = "|3234_VIDEO_TITLE|"
	MODEL_YT_VIDEO_VIDEO_DESCRIPTION_EMAIL string = "|3234_VIDEO_DESCRIPTION|"
	MODEL_YT_VIDEO_VIDEO_CODE_EMAIL        string = "|3234_VIDEO_CODE|"
	MODEL_YT_VIDEO_VIDEO_IMAGE_EMAIL       string = "|3234_VIDEO_IMAGE|"
	MODEL_YT_VIDEO_VIDEO_TIME_EMAIL        string = "|3234_VIDEO_TIME|"
	MODEL_YT_VIDEO_VIDEO_TIME_COLOR_EMAIL  string = "|3234_VIDEO_TIME_COLOR|"
	MODEL_YT_VIDEO_PLAYLIST_CODE_EMAIL     string = "|3234_PLAYLIST_CODE|"
	MODEL_YT_VIDEO_SUBSCRIPTION_LINK_EMAIL string = "|3234_SUBSCRIPTION_LINK|"
	MODEL_YT_VIDEO_SUBSCRIPTION_NAME_EMAIL string = "|3234_SUBSCRIPTION_NAME|"

	MODEL_DISKS_SMART_DISK_LABEL_EMAIL        string = "|3234_DISK_LABEL|"
	MODEL_DISKS_SMART_DISK_SERIAL_EMAIL       string = "|3234_DISK_SERIAL|"
	MODEL_DISKS_SMART_DISK_PARTITION_EMAIL    string = "|3234_DISK_PARTITION|"
	MODEL_DISKS_SMART_DISKS_SMART_HTML_EMAIL  string = "|3234_DISKS_SMART_HTML|"

)

const RAND_STR_LEN int = 10

const TO_SEND_REL_FOLDER string = "to_send"
const _EMAIL_MODELS_FOLDER string = "email_models"

const _TEMP_EML_FILE string = "msg_temp.eml"

const MODEL_FILE_INFO string = "model_email_info.html"
const MODEL_FILE_RSS string = "model_email_rss.html"
const MODEL_FILE_YT_VIDEO string = "model_email_video_YouTube.html"
const MODEL_FILE_DISKS_SMART string = "model_email_disks_smart.html"
const _MODEL_FILE_MESSAGE_EML string = "model_message.eml"
/*
GetModelFileEMAIL returns the contents of an email model file.

-----------------------------------------------------------

– Params:
  - file_name – the name of the file
  - things_replace – the map of things to replace in the file

– Returns:
  - an instance of EmailInfo with the EmailInfo.Sender, EmailInfo.Mail_to and EmailInfo.Html filled and ready
*/
func GetModelFileEMAIL(file_name string, things_replace map[string]string) EmailInfo {
	var sender string
	switch file_name {
		case MODEL_FILE_INFO:
			sender = "VISOR - Info"
		case MODEL_FILE_RSS:
			sender = "VISOR - RSS"
		case MODEL_FILE_YT_VIDEO:
			sender = "YouTube"
		case MODEL_FILE_DISKS_SMART:
			sender = "VISOR - S.M.A.R.T."
		default:
			sender = "VISOR"
	}

	var msg_html string = *getProgramDataDirMODULES(NUM_MOD_EmailSender).Add2(false, _EMAIL_MODELS_FOLDER, file_name).
		ReadTextFile()
	for key, value := range things_replace {
		msg_html = strings.ReplaceAll(msg_html, key, value)
	}

	return EmailInfo{
		Sender:     sender,
		Mail_to:    User_settings_GL.PersonalConsts.User_email_addr,
		Subject:    "",
		Html:       msg_html,
		Multiparts: nil,
	}
}

/*
QueueEmailEMAIL queues an email to be sent by the Email Sender module.

-----------------------------------------------------------

– Params:
  - emailInfo – the email info

– Returns:
  - nil if the email was queued successfully, otherwise an error
*/
func QueueEmailEMAIL(emailInfo EmailInfo) error {
	var message_eml string
	if emailInfo.Eml != "" {
		message_eml = emailInfo.Eml
	} else {
		eml, _, success := prepareEmlEMAIL(emailInfo)
		if !success {
			return errors.New("error preparing the EML file")
		}
		message_eml = eml
	}

	if VISOR_server_GL {
		// Keep trying to create a file with a unique name.
		var file_name string = ""
		var to_send_dir GPath = GetUserDataDirMODULES(NUM_MOD_EmailSender).Add2(true, TO_SEND_REL_FOLDER)
		for {
			var rand_string string = RandStringGENERAL(RAND_STR_LEN)
			_, err := os.ReadFile(to_send_dir.Add2(false, rand_string+emailInfo.Mail_to + ".eml").
				GPathToStringConversion())
			if nil != err {
				// If the file doesn't exist, choose that name.
				file_name = rand_string + emailInfo.Mail_to + ".eml"

				return to_send_dir.Add2(false, file_name).WriteTextFile(message_eml, false)
			}
		}
	} else {
		var message []byte = []byte("Email|" + emailInfo.Mail_to + "|")
		message = append(message, CompressString(message_eml)...)
		QueueNoResponseMessageSERVER(message)

		return nil
	}
}

/*
SendEmailEMAIL sends an email with the given message and receiver.

***DO NOT USE OUTSIDE THE EMAIL SENDER OR WEBSITE BACKEND MODULES***

-----------------------------------------------------------

– Params:
  - message_eml – the complete message to be sent in EML format
  - mail_to – the receiver of the email
  - emergency_email – true if the email is an emergency email and so will make this function halt until the connection
	is made, false otherwise

– Returns:
  - nil if the email was sent successfully, otherwise an error
*/
func SendEmailEMAIL(message_eml string, mail_to string, emergency_email bool) error {
	if err := getModTempDirMODULES(NUM_MOD_EmailSender).Add2(false, _TEMP_EML_FILE).WriteTextFile(message_eml, false); nil != err {
		return err
	}
	_, err := ExecCmdSHELL([]string{getCurlStringEMAIL(mail_to, emergency_email)})

	return err
}

/*
ToQuotedPrintableEMAIL converts a string to a quoted printable string.

-----------------------------------------------------------

– Params:
  - str – the string to convert

– Returns:
  - the quoted printable string or nil if an error occurs
*/
func ToQuotedPrintableEMAIL(str string) *string {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(str))
	if nil != err {
		return nil
	}
	err = w.Close()
	if nil != err {
		return nil
	}
	ret := ac.String()

	return &ret
}

/*
prepareEmlEMAIL prepares the EML file of the email.

-----------------------------------------------------------

– Params:
  - emailInfo – the email info
  - multiparts – the list of multipart items to attach to the email aside from the main HTML or nil to ignore

– Returns:
  - the email EML file to be sent
  - the receiver of the email
  - nil if the email was queued successfully, otherwise an error
*/
func prepareEmlEMAIL(emailInfo EmailInfo) (string, string, bool) {
	var p_message_eml *string = getProgramDataDirMODULES(NUM_MOD_EmailSender).Add2(false, _EMAIL_MODELS_FOLDER,
		_MODEL_FILE_MESSAGE_EML).ReadTextFile()
	if p_message_eml == nil {
		return "", "", false
	}
	var message_eml string = *p_message_eml

	emailInfo.Html = strings.ReplaceAll(emailInfo.Html, "|3234_EML_SUBJECT|", emailInfo.Subject)
	emailInfo.Html = strings.ReplaceAll(emailInfo.Html, "|3234_EML_SENDER_NAME|", emailInfo.Sender)

	message_eml = strings.ReplaceAll(message_eml, "|3234_EML_HTML|", *ToQuotedPrintableEMAIL(emailInfo.Html))
	message_eml = strings.ReplaceAll(message_eml, "|3234_EML_SUBJECT|", emailInfo.Subject)
	message_eml = strings.ReplaceAll(message_eml, "|3234_EML_SENDER_NAME|", emailInfo.Sender)

	var multiparts_str string = ""
	if nil != emailInfo.Multiparts {
		for _, multipart := range emailInfo.Multiparts {
			multiparts_str += "\n--|3234_EML_BOUNDARY|\n" +
						"Content-Type: " + multipart.Content_type + "\n" +
						"Content-Transfer-Encoding: " + multipart.Content_transfer_encoding + "\n" +
						"Content-ID: <" + multipart.Content_id + ">\n" +
						"\n" +
						multipart.Body + "\n\n"
		}
	}
	message_eml = strings.ReplaceAll(message_eml, "|3234_EML_MULTIPARTS|", multiparts_str)

	var msg_boundary string = RandStringGENERAL(25)
	for {
		if !strings.Contains(message_eml, msg_boundary) {
			break
		}
		msg_boundary = RandStringGENERAL(25)
	}
	message_eml = strings.ReplaceAll(message_eml, "|3234_EML_BOUNDARY|", msg_boundary)


	return message_eml, emailInfo.Mail_to, true
}

/*
getCurlStringEMAIL gets the cURL string that sends an email with the default message file path and sender and receiver.

-----------------------------------------------------------

– Params:
  - mail_to – the receiver of the email
  - emergency_email – true if the email is an emergency email and so will make this function halt until the connection
    is made, false otherwise

– Returns:
  - the string ready to be executed by the system
*/
func getCurlStringEMAIL(mail_to string, emergency_email bool) string {
	var timeout string = "10"
	if emergency_email {
		timeout = "100000"
	}

	if User_settings_GL.PersonalConsts.VISOR_email_addr == "" || User_settings_GL.PersonalConsts.VISOR_email_pw == "" {
		return "echo \"No email address or password set\""
	}

	return "curl{{EXE}} --location --connect-timeout " + timeout + " --verbose \"smtp://smtp.gmail.com:587\" --user \"" +
		User_settings_GL.PersonalConsts.VISOR_email_addr + ":" + User_settings_GL.PersonalConsts.VISOR_email_pw +
		"\" --mail-rcpt \"" + mail_to + "\" --upload-file \"" +
		getModTempDirMODULES(NUM_MOD_EmailSender).Add2(false, _TEMP_EML_FILE).GPathToStringConversion() + "\" --ssl-reqd"
}
