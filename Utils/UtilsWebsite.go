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
	"net/http"
	"net/url"
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
	// Text3 is the third text (optional)
	Text3 string
}

/*
GetPageContentsWEBSITE gets the page contents from the given VISOR's website page.

This function can be used in libraries (requests the website URL and the password instead of getting them from
PersonalConsts_GL).

-----------------------------------------------------------

– Params:
  - partial_url – the partial URL of the page to get the contents from. Example: files_EOG/gpt_text.txt to get from,
    for example, https://www.visor.com/files_EOG/gpt_text.txt

– Returns:
  - the page contents or nil if an error occurred
*/
func GetPageContentsWEBSITE(partial_url string) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", PersonalConsts_GL.WEBSITE_URL + partial_url, nil)
	if err != nil {
		return nil
	}
	req.SetBasicAuth("VISOR", PersonalConsts_GL.WEBSITE_PW)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	bodyText, _ := io.ReadAll(resp.Body)

	return bodyText
}

/*
SubmitFormWEBSITE sends a form to the given VISOR's website.

-----------------------------------------------------------

– Params:
  - form – the form to send
  - website – the website URL
  - passwd – the password to access the website

– Returns:
  - true if the form was submitted successfully, false otherwise
*/
func SubmitFormWEBSITE(form WebsiteForm) error {
	formData := url.Values{
		"type": {form.Type},
		"text1":  {form.Text1},
		"text2":  {form.Text2},
		"text3":  {form.Text3},
	}

	// Convert form data to a format suitable for HTTP requests
	formDataEncoded := formData.Encode()

	// Create a new POST request with the form data
	req, err := http.NewRequest("POST", PersonalConsts_GL.WEBSITE_URL + "submit-form", bytes.NewBufferString(formDataEncoded))
	if err != nil {
		return err
	}
	req.SetBasicAuth("VISOR", PersonalConsts_GL.WEBSITE_PW)

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create an HTTP client and send the request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return errors.New("response status code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}
