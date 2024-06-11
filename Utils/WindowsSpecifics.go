//go:build windows

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
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
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
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
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
	win.ShowWindow(win.GetConsoleWindow(), win.SW_HIDE)
}
