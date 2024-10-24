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
	"os/exec"
	"runtime"
)

/*
StartConAppPROCESSES starts a new separate and independent console process with the given path, with hidden window.

-----------------------------------------------------------

– Params:
  - path – the full path of the program to start
  - arg – an optional argument to pass to the program (ONLY FOR WINDOWS)

– Returns:
  - true if the process was started correctly, false otherwise
*/
func StartConAppPROCESSES(path GPath, arg string) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd.exe", "/C", "start", "conhost.exe", path.GPathToStringConversion(), arg)
		err := cmd.Run()
		if err != nil {
			return false
		}
	} else {
		cmd := exec.Command(GetShellSHELL("", ""), "-c", "nohup " + path.GPathToStringConversion() + " </dev/null >/dev/null 2>&1 &")
		err := cmd.Run()
		if err != nil {
			return false
		}
	}

	return true
}

/*
KillAllPROCESSES kills all the processes with the given name.

-----------------------------------------------------------

– Params:
  - prog_name – the name of the process to kill
 */
func KillAllPROCESSES(prog_name string) {
	if runtime.GOOS == "windows" {
		_, _ = ExecCmdSHELL([]string{"taskkill /F /IM " + prog_name + "{{EXE}}"})
	} else {
		_, _ = ExecCmdSHELL([]string{"killall " + prog_name + "{{EXE}}"})
	}
}
