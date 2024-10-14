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

package UtilsSWA

import (
	"strings"
)

/*
IsRootAvailableROOT checks if the device can run root commands.

-----------------------------------------------------------

– Returns:
  - true if the device can run root commands, false otherwise.
*/
func IsRootAvailableROOT() bool {
	// The original implementation was gotten from
	// http://muzikant-android.blogspot.com/2011/02/how-to-get-root-access-and-execute.html, was then adapted to the
	// app and finally changed to call the executeShellCmd() function.
	// So the way it checks root is supposed to be the same (hopefully it's the same and I didn't mess up). Only the
	// shell part is now removed because it is on the mentioned function.
	// EDIT: it's no longer going according to the mentioned webpage. Now it just checks yes or no, not yes, no
	// because denied, or no because not rooted.

	// Keep attempt_su false and input the whole command or infinite recursion will pay a visit.
	cmd_output, err := ExecCmdSHELL(false, "su -c id")
	if err != nil {
		return false
	}

	// Root denied, could be error 13 of permission denied, for example - happens with Magisk. With SuperSU,
	// error code 1 is returned. Which means, don't check for specific error codes except the file not found
	// one (127), which for sure indicates there is no su binary available - how could it return [didn't get to
	// finish... It could return 127 for no permission... Nothing stops it from not returning 127.]
	// EDIT: I don't think there is a reliable way to check if root is available or not on the device (I don't
	// mean denied - I mean, device not rooted at all) without checking the default binary locations. I don't
	// know if that's reliable either, and so far no need to know if the user decided not to give root
	// permissions while being able to, so whatever.

	// All that is needed to know if there is root access is to know if the User ID is 0, which means root. If it's
	// not root, then permission to access it was either denied or the su binary is non-existent (device not
	// rooted).
	if strings.Contains(string(GetStdoutSHELL(cmd_output)), "uid=0") {
		return true
	} else {
		return false
	}
}
