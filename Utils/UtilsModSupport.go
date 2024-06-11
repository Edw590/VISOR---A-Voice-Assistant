/*******************************************************************************
 * Copyright 2023-2023 Edw590
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

import "runtime"

/*
isMOD2Supported checks if the MOD_2 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD2Supported() bool {
	// Check if the command "smartctl" is available
	output, err := ExecCmdSHELL([]string{"smartctl{{EXE}} -h"})
	if err != nil {
		return false
	}

	return output.Exit_code == 0
}

/*
isMOD3Supported checks if the module MOD_4 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD3Supported() bool {
	return runtime.GOOS == "windows"
}

/*
isMOD4Supported checks if the module MOD_4 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD4Supported() bool {
	return true
}

/*
isMOD5Supported checks if the module MOD_5 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD5Supported() bool {
	// Check if the command "curl" is available
	output, err := ExecCmdSHELL([]string{"curl{{EXE}} -h"})
	if err != nil {
		return false
	}

	return output.Exit_code == 0
}

/*
isMOD6Supported checks if the module MOD_6 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD6Supported() bool {
	if runtime.GOOS == "windows" {
		return false
	}

	// Check if the command "/usr/bin/chromedriver" is available
	output, err := ExecCmdSHELL([]string{"/usr/bin/chromedriver{{EXE}} -h"})
	if err != nil {
		return false
	}

	return output.Exit_code == 0
}

/*
isMOD7Supported checks if the module MOD_7 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD7Supported() bool {
	// Check if the command "ollama" is available
	output, err := ExecCmdSHELL([]string{"llamacpp{{EXE}} -h"})
	if err != nil {
		return false
	}

	return output.Exit_code == 0
}

/*
isMOD8Supported checks if the module MOD_8 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD8Supported() bool {
	return true
}

/*
isMOD9Supported checks if the module MOD_9 is supported on the current machine.

-----------------------------------------------------------

– Returns:
  - true if the module is supported, false otherwise.
*/
func isMOD9Supported() bool {
	return true
}
