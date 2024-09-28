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
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// WebsiteForm represents a form to be submitted to VISOR's website
type WebsiteForm struct {
	// Type is the form type
	Type string
	// Text1 is the first text
	Text1 string
	// Text2 is the second text (optional)
	Text2 string
	// File is the file bytes to be submitted (optional)
	File []byte
}

/*
GetFileContentsWEBSITE gets the file contents from the given VISOR's website URL.

-----------------------------------------------------------

– Params:
  - partial_path – the partial path of the file to get the contents from. Example: gpt_text.txt to get from
	https://www.visor.com/files_EOG/gpt_text.txt
  - get_file – true if the file contents are to be retrieved, false if the CRC16 checksum of the file is to be retrieved

– Returns:
  - the file contents or the CRC16 checksum, or nil if an error occurred
 */
func GetFileContentsWEBSITE(partial_path string, get_file bool) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var add_to_request string = ""
	if !get_file {
		add_to_request = "?crc=true"
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", User_settings_GL.PersonalConsts.Website_url + "/file/" + partial_path +
		add_to_request, nil)
	if err != nil {
		return nil
	}
	req.SetBasicAuth("VISOR", User_settings_GL.PersonalConsts.Website_pw)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	defer resp.Body.Close()

	body_text, _ := io.ReadAll(resp.Body)

	return body_text
}

/*
SubmitFormWEBSITE sends a form to the given VISOR's webserver and receives its response.

-----------------------------------------------------------

– Params:
  - form – the form to send
  - website – the website URL
  - passwd – the password to access the website

– Returns:
  - true if the form was submitted successfully, false otherwise
*/
func SubmitFormWEBSITE(form WebsiteForm) ([]byte, error) {
	// Create a buffer to hold the multipart form data
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Add the form fields
	if form.Type != "" {
		_ = writer.WriteField("type", form.Type)
	}
	if form.Text1 != "" {
		_ = writer.WriteField("text1", form.Text1)
	}
	if form.Text2 != "" {
		_ = writer.WriteField("text2", form.Text2)
	}
	if form.File != nil {
		part, err := writer.CreateFormFile("file", "file")
		if err != nil {
			return nil, err
		}
		_, err = part.Write(form.File)
		if err != nil {
			return nil, err
		}
	}

	// Close the writer to finalize the form data
	err := writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", User_settings_GL.PersonalConsts.Website_url + "/submit-form", &buffer)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("VISOR", User_settings_GL.PersonalConsts.Website_pw)

	// Set the appropriate headers
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create an HTTP client and send the request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code: " + strconv.Itoa(resp.StatusCode))
	}

	return io.ReadAll(resp.Body)
}

/*
CheckFileChangedWEBSITE checks if the file has changed by comparing the CRC16 checksum of the file with the given CRC16
checksum.

-----------------------------------------------------------

– Params:
  - old_crc16 – the old CRC16 checksum
  - file_path – the file path

– Returns:
  - the new CRC16 checksum if the file has changed, nil otherwise
*/
func CheckFileChangedWEBSITE(old_crc16 []byte, file_path string) []byte {
	var new_crc16 []byte = GetFileContentsWEBSITE(file_path, false)
	if new_crc16 == nil {
		return nil
	}

	if !bytes.Equal(new_crc16, old_crc16) {
		return new_crc16
	}

	return nil
}
