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

package GPTCommunicator

import (
	"Utils"
	"bufio"
	"io"
	"log"
	"os/exec"
)

type _OllamaRequest struct {
	Model string `json:"model"`
	Prompt string `json:"prompt"`
	Suffix string `json:"suffix"`
	Images []string `json:"images"`

	Format string `json:"format"`
	Options _ModelfileParams `json:"options"`
	System string `json:"system"`
	Template string `json:"template"`
	Context []int `json:"context"`
	Stream bool `json:"stream"`
	Raw bool `json:"raw"`
	Keep_alive string `json:"keep_alive"`
}

type _ModelfileParams struct {
	Num_ctx int `json:"num_ctx"`
	Temperature float32 `json:"temperature"`
}

type _OllamaResponse struct {
	// Adjust field names and types based on the expected JSON structure
	Model string `json:"model"`
	Created_at string `json:"created_at"`
	Response string `json:"response"`
	Done bool `json:"done"`
	Total_duration int `json:"total_duration"`
	Load_duration int `json:"load_duration"`
	Prompt_eval_count int `json:"prompt_eval_count"`
	Prompt_eval_duration int `json:"prompt_eval_duration"`
	Eval_count int `json:"eval_count"`
	Eval_duration int `json:"eval_duration"`
	Context []int `json:"context"`
}



/*// Directory for processing text files
const _TO_PROCESS_REL_FOLDER string = "to_process"

// Command prefixes for Wolfram Alpha and Wikipedia
const ASK_WOLFRAM_ALPHA string = "/askWolframAlpha "
const SEARCH_WIKIPEDIA string = "/searchWikipedia "


var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
	modGenInfo_GL  *ModsFileInfo.Mod7GenInfo
	modUserInfo_GL *ModsFileInfo.Mod7UserInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_7
		modUserInfo_GL = &Utils.User_settings_GL.GPTCommunicator


		// Set initial module state
		modGenInfo_GL.State = ModsFileInfo.MOD_7_STATE_STARTING

		// Stop any running Ollama instances to start fresh
		restartOllama()

		time.Sleep(1 * time.Second)

		var system_info string = modUserInfo_GL.System_info
		system_info = strings.Replace(system_info, "3234_WEEKDAY", time.Now().Weekday().String(), -1)
		system_info = strings.Replace(system_info, "3234_DAY", strconv.Itoa(time.Now().Day()), -1)
		system_info = strings.Replace(system_info, "3234_MONTH", time.Now().Month().String(), -1)
		system_info = strings.Replace(system_info, "3234_YEAR", strconv.Itoa(time.Now().Year()), -1)

		// Load visor introduction text
		var visor_intro string = *moduleInfo_GL.ModDirsInfo.ProgramData.Add2(false, "visor_intro.txt").ReadTextFile()
		visor_intro = strings.Replace(visor_intro, "\n", " ", -1)
		visor_intro = strings.Replace(visor_intro, "\"", "\\\"", -1)
		visor_intro = strings.Replace(visor_intro, "3234_NICK", modUserInfo_GL.User_nickname, -1)

		// Initialize memory string
		var memories string = strings.Join(modGenInfo_GL.Memories, ". ")
		memories = strings.Replace(memories, "\"", "\\\"", -1)

		var system_prompt string =
			system_info + " \n\n" +
			"Memories stored about the user: " + memories + "\n\n" +
			"About you:" + visor_intro
		if system_prompt == "" {}

		// Declare and assign context sizes
		//var smart_ctx_size int = 12288
		var dumb_ctx_size int = 4096



		var request_data _OllamaRequest = _OllamaRequest{
			Model:      "llama3.2:latest",
			Prompt:     "What do you remember about me?",
			Suffix:     "",
			Images:     nil,
			Format:     "",
			Options:    _ModelfileParams{
				Num_ctx:     dumb_ctx_size,
				Temperature: 0.8,
			},
			System:     system_prompt,
			Template:   "",
			Context:    nil,
			Stream:     true,
			Raw:        false,
			Keep_alive: "5m",
		}

		jsonData, err := json.Marshal(request_data)
		if err != nil {
			fmt.Printf("Error marshalling JSON: %v\n", err)
			os.Exit(1)
		}

		// Create the POST request
		req, err := http.NewRequest("POST", "http://127.0.0.1:11434/api/generate", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		// Use a JSON decoder to handle the streamed response
		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var response _OllamaResponse
			if err := decoder.Decode(&response); err != nil {
				fmt.Printf("Error decoding JSON: %v\n", err)
				break
			}

			// Process each JSON object as it arrives
			fmt.Printf("Response text: %v\n", response)
		}



		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1000000000) {
				return
			}
		}
	}
}*/

func startOllama(instance_type string) (*bufio.Writer, io.ReadCloser, io.ReadCloser) {
	cmd := exec.Command(Utils.GetShellSHELL("", ""))
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Printf("Error starting %s a shell session for an LLM: %v", instance_type, err)

		return nil, nil, nil
	}

	writer := bufio.NewWriter(stdin)

	return writer, stdout, stderr
}

/*
stopOllama stop the Ollama service.
*/
func stopOllama() {
	_, _ = Utils.ExecCmdSHELL([]string{"sudo systemctl stop ollama.service"})
}

/*
restartOllama restarts the Ollama service.
*/
func restartOllama() {
	_, _ = Utils.ExecCmdSHELL([]string{"sudo systemctl restart ollama.service"})
}
