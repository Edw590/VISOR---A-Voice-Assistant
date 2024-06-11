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

package GPT

import (
	"Utils"
	"strconv"
	"strings"
)

type Entry struct {
	// Time is the Unix time in milliseconds
	Time int64
	// Text is the text generated
	Text string
}

const GEN_ERROR string = "3234_ERROR"

var website_url_GL string = ""
var website_pw_GL string = ""

/*
SetWebsiteInfo sets VISOR's website url and password.

-----------------------------------------------------------

– Params:
  - url – the website url
  - pw – the website password
*/
func SetWebsiteInfo(url string, pw string) {
	website_url_GL = url
	website_pw_GL = pw
}

func GetNumEntries() int {
	var page_contents string = string(Utils.GetPageContentsWEBSITE(website_url_GL + "files_EOG/gpt_text.txt", website_pw_GL))
	var entries []string = strings.Split(page_contents, "[3234_START:")

	for i := 0; i < len(entries); i++ {
		if entries[i] == "" {
			Utils.DelElemSLICES(&entries, i)
		}
	}

	return len(entries)
}

func GetEntry(num int) string {
	var page_contents string = string(Utils.GetPageContentsWEBSITE(website_url_GL + "files_EOG/gpt_text.txt", website_pw_GL))
	var entries []string = strings.Split(page_contents, "[3234_START:")

	for i := 0; i < len(entries); i++ {
		if entries[i] == "" {
			Utils.DelElemSLICES(&entries, i)
		}
	}

	if len(entries) == 0 {
		return GEN_ERROR
	}

	if num < 0 {
		num = len(entries) + num
	} else if num >= len(entries) {
		num = len(entries) - 1
	} else {
		// Do nothing
	}

	var entry string = entries[num]

	return entry[:strings.Index(entry, "]")] + "|" + entry[strings.Index(entry, "]") + 1:]
}

func GetTimeFromEntry(entry string) int64 {
	var parts []string = strings.Split(entry, "|")

	if len(parts) < 2 {
		return -1
	}

	time, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return -1
	}

	return time
}

func GetTextFromEntry(entry string) string {
	var parts []string = strings.Split(entry, "|")

	if len(parts) < 2 {
		return GEN_ERROR
	}

	return parts[1]
}

func SendText(text string) error {
	return Utils.SubmitFormWEBSITE(Utils.WebsiteForm{
		Name: "GPT",
		Text1: text,
	})
}
