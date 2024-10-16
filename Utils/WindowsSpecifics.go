//go:build windows

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
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

// == FUNCTION 1: RunningAsAdminPROCESSES ==

/**
 * Checks if the program is running with Administrator (root) privileges.
 *
 * @return true if running as admin, false otherwise
 */
func RunningAsAdminPROCESSES() bool {
	// Step 1: Create a Security Identifier (SID) for the Administrators group
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,      // Authority
		2,                                   // Revision
		windows.SECURITY_BUILTIN_DOMAIN_RID, // Domain RID
		windows.DOMAIN_ALIAS_RID_ADMINS,     // Administrator RID
		0, 0, 0, 0, 0, 0,                    // Sub-authority values
		&sid)
	if err != nil {
		// If SID creation fails, return false
		return false
	}
	defer windows.FreeSid(sid) // Clean up SID memory

	// Step 2: Get current process token (handle)
	// Passing 0 gets the token for the current process
	token := windows.Token(0)

	// Step 3: Check if the token is a member of the Administration Group
	member, err := token.IsMember(sid)
	if err != nil {
		// If membership check fails, return false
		return false
	}

	// Return the membership result
	return member
}

// == FUNCTION 2: HideConsoleWindowPROCESSES ==

/**
 * Hides the console window of the program.
 *
 * Note: Only works on Windows if the program is started with conhost.exe
 * (which is the default, except when using the new Windows Terminal).
 */
func HideConsoleWindowPROCESSES() {
	// Hide the console window using the SW_HIDE flag
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}

// == FUNCTION 3: GenerateCtrlCPROCESSES ==

/**
 * Generates a Ctrl+C event to a process group.
 *
 * @param cmd         The command to generate the event for
 * @param process_group_id The ID of the process group
 * @return error if the event couldn't be generated, nil otherwise
 *
 * Note: This function has NOT been thoroughly tested!
 */

// TODO: Implement GenerateCtrlCPROCESSES usage in the future

/*
func GenerateCtrlCPROCESSES(cmd *exec.Cmd, process_group_id uint32) error {
	// Ensures the command has a SysProAttr set
	if cmd.SysProcAttr == nil {
		// Create a new SysProcAttr to set the creation flags
		cmd.SysProcAttr = &syscall.SysProcAttr{
			// Create a new process group
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	}

	// Load the kernel32.DLL
	var kernel32 = syscall.NewLazyDLL("kernel32.dll")

	// Get the GenerateConsoleCtrlEvent procedure
	var procGenerateConsoleCtrlEvent = kernel32.NewProc("GenerateConsoleCtrlEvent")

	// Call GenerateConsoleCtrlEvent with CTRL_C_EVENT
	r, _, err := procGenerateConsoleCtrlEvent.Call(
		syscall.CTRL_C_EVENT,      // Ctrl+C event
		uintptr(process_group_id)) // Process group ID

	// If the call fails, return the error
	if r == 0 {
		return err
	}

	// If successful, return nil

	return nil
}
*/
