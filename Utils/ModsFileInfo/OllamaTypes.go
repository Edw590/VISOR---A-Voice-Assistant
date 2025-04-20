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

package ModsFileInfo

// ----- Prepared for Ollama 0.6.2 -----

// Taken from: https://github.com/ollama/ollama/blob/main/api/types.go.

type OllamaChatRequest struct {
	Model string `json:"model"`
	Messages []OllamaMessage `json:"messages"`

	Format  string `json:"format,omitempty"`
	Options OllamaOptions `json:"options"`
	Stream bool `json:"stream,omitempty"`
	Keep_alive string `json:"keep_alive,omitempty"` // This must be a string

	Tools OllamaTools `json:"tools,omitempty"`
}

type OllamaOptions struct {
	Num_keep int `json:"num_keep"`
	Num_ctx int32 `json:"num_ctx"`
	Temperature float32 `json:"temperature"`
}

type OllamaChatResponse struct {
	Model string `json:"model"`
	Message OllamaMessage `json:"message"`
}

type OllamaMessage struct {
	Role string `json:"role"`
	Content string `json:"content"`
	Images []ImageData `json:"images,omitempty"`
	Tool_calls []OllamaToolCall `json:"tool_calls,omitempty"`
	Timestamp_s int64
}

type ImageData []byte

type OllamaTools []OllamaTool

type OllamaToolCall struct {
	Function OllamaToolCallFunction `json:"function"`
}

type OllamaToolCallFunction struct {
	Index     int                             `json:"index,omitempty"`
	Name      string                          `json:"name"`
	Arguments OllamaToolCallFunctionArguments `json:"arguments"`
}

type OllamaToolCallFunctionArguments map[string]any

type OllamaTool struct {
	Type     string       `json:"type"`
	Function OllamaToolFunction `json:"function"`
}

type OllamaToolFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties map[string]struct {
			Type        string   `json:"type"`
			Description string   `json:"description"`
			Enum        []string `json:"enum,omitempty"`
		} `json:"properties"`
	} `json:"parameters"`
}