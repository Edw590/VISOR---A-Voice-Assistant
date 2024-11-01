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

//go:build windows

package Utils

import (
	"github.com/itchyny/volume-go"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
	"os/exec"
	"strconv"
	"syscall"
)

/*
RunningAsAdminPROCESSES checks if the program is running as administrator/root.

-----------------------------------------------------------

– Returns:
  - true if the program is running as admin, false otherwise
*/
func RunningAsAdminPROCESSES() bool {
	// Took from https://github.com/golang/go/issues/28804

	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,      // Authority
		2,                                   // Revision
		windows.SECURITY_BUILTIN_DOMAIN_RID, // Domain RID
		windows.DOMAIN_ALIAS_RID_ADMINS,     // Administrator RID
		0, 0, 0, 0, 0, 0,                    // Sub-authority values
		&sid)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		return false
	}

	return member
}

/*
HideConsoleWindowPROCESSES hides the console window of the program.

Notice: on Windows only works if the program is started with conhost.exe (always is except when it's started by the
new Windows Terminal). So use StartConAppPROCESSES() to start the program with conhost.exe.
*/
func HideConsoleWindowPROCESSES() {
	// Hide the console window using the SW_HIDE flag
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}

/*
GenerateCtrlCPROCESSES generates a Ctrl+C event to a process group.

This function has NOT been tested!

-----------------------------------------------------------

– Params:
  - process_group_id – the process group ID

– Returns:
  - an error if the event couldn't be generated, nil otherwise
*/
func GenerateCtrlCPROCESSES(cmd *exec.Cmd, process_group_id uint32) error {
	if cmd.SysProcAttr == nil {
		// Set the process to create a new process group
		// WARNING: this might need to be done before calling cmd.Start()
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	}

	var kernel32 = syscall.NewLazyDLL("kernel32.dll")
	var procGenerateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")

	// Call GenerateConsoleCtrlEvent with CTRL_C_EVENT
	r, _, err := procGenerateConsoleCtrlEvent.Call(
		syscall.CTRL_C_EVENT,      // Ctrl+C event
		uintptr(process_group_id)) // Process group ID
	if r == 0 {
		return err
	}

	return nil
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
	return toggleNetworkInterfaceCONNECTIVITY("Wi-Fi", enable)
}

/*
ToggleEthernetCONNECTIVITY toggles the Ethernet interface.

-----------------------------------------------------------

– Params:
  - enable – true to enable the Ethernet interface, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
 */
func ToggleEthernetCONNECTIVITY(enable bool) bool {
	return toggleNetworkInterfaceCONNECTIVITY("Ethernet", enable)
}

/*
ToggleNetworkingCONNECTIVITY toggles the Wi-Fi and Ethernet interfaces.

-----------------------------------------------------------

– Params:
  - enable – true to enable the interfaces, false to disable them

– Returns:
  - true if the operation was successful for both, false otherwise
*/
func ToggleNetworkingCONNECTIVITY(enable bool) bool {
	return ToggleWifiCONNECTIVITY(enable) && ToggleEthernetCONNECTIVITY(enable)
}

/*
toggleNetworkInterfaceCONNECTIVITY toggles a network interface.

-----------------------------------------------------------

– Params:
  - interface_name – the name of the interface
  - enable – true to enable the interface, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
 */
func toggleNetworkInterfaceCONNECTIVITY(interface_name string, enable bool) bool {
	var en_dis string = "disabled"
	if enable {
		en_dis = "enabled"
	}
	cmd_output, err := ExecCmdSHELL([]string{"netsh interface set interface " + interface_name + " " + en_dis})
	if err != nil {
		return false
	}

	return cmd_output.Exit_code == 0
}

/*
SetVolumeVOLUME sets the system volume.

-----------------------------------------------------------

– Params:
  - vol – the volume to set (0-100)

– Returns:
  - true if the operation was successful, false otherwise
 */
func SetVolumeVOLUME(vol int) bool {
	err := volume.SetVolume(vol)
	if err != nil {
		var vol_nircmd int = vol * 65535 / 100
		return CheckCmdOutput(ExecCmdSHELL([]string{"nircmdc{{EXE}} setsysvolume " + strconv.Itoa(vol_nircmd)}))
	}

	return true
}

/*
SetMutedVOLUME mutes or unmutes the system volume.

-----------------------------------------------------------

– Params:
  - mute – true to mute the volume, false to unmute it

– Returns:
  - true if the operation was successful, false otherwise
 */
func SetMutedVOLUME(mute bool) bool {
	var err error
	if mute {
		err = volume.Mute()
		if err != nil {
			return CheckCmdOutput(ExecCmdSHELL([]string{"nircmdc{{EXE}} mutesysvolume 1"}))
		}
	} else {
		err = volume.Unmute()
		if err != nil {
			return CheckCmdOutput(ExecCmdSHELL([]string{"nircmdc{{EXE}} mutesysvolume 0"}))
		}
	}

	return true
}

/*
GetOSVersion gets the OS version.

-----------------------------------------------------------

– Returns:
  - the major version
  - the minor version
 */
func GetOSVersion() (int, int, int) {
	maj, min, patch := windows.RtlGetNtVersionNumbers()

	return int(maj), int(min), int(patch)
}
