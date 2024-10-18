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

//go:build linux

package Utils

import (
	"os"
	"os/exec"
	"strings"
)

/*
RunningAsAdminPROCESSES checks if the program is running as administrator/root.

-----------------------------------------------------------

– Returns:
  - true if the program is running as admin, false otherwise
*/
func RunningAsAdminPROCESSES() bool {
	stdOutErrCmd, err := ExecCmdSHELL([]string{"id -u"})
	if nil != err {
		return false
	}

	if stdOutErrCmd.Exit_code != 0 {
		return false
	}

	return strings.TrimSpace(stdOutErrCmd.Stdout_str) == "0"
}

/*
HideConsoleWindowPROCESSES does NOTHING.
 */
func HideConsoleWindowPROCESSES() {
	// TODO See if it's needed on Linux too and find a way
}

/*
GenerateCtrlCPROCESSES generates a Ctrl+C event to a process.
*/
func GenerateCtrlCPROCESSES(cmd *exec.Cmd, process_group_id uint32) error {
	return cmd.Process.Signal(os.Interrupt)
}

/*
ToggleWifiCONNECTIVITY toggles the Wi-Fi interface.

-----------------------------------------------------------

– Params:
  - enable – true to enable the Wi-Fi interface, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
*/
func ToggleWifiCONNECTIVITY(enable bool) bool {
	var on_off string = "off"
	if enable {
		on_off = "on"
	}
	cmd_output, err := ExecCmdSHELL([]string{"nmcli radio wifi " + on_off})
	if err != nil {
		return false
	}

	return cmd_output.Exit_code == 0
}

/*
ToggleEthernetCONNECTIVITY toggles the Ethernet interface.

NOT implemented.

-----------------------------------------------------------

– Params:
  - enable – true to enable the Ethernet interface, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
*/
func ToggleEthernetCONNECTIVITY(enable bool) bool {
	return false
}

/*
ToggleNetworkingCONNECTIVITY toggles the networking.

-----------------------------------------------------------

– Params:
  - enable – true to enable the networking, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
 */
func ToggleNetworkingCONNECTIVITY(enable bool) bool {
	var on_off string = "off"
	if enable {
		on_off = "on"
	}
	cmd_output, err := ExecCmdSHELL([]string{"nmcli networking " + on_off})
	if err != nil {
		return false
	}

	return cmd_output.Exit_code == 0
}
