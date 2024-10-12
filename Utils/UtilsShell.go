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
	"bytes"
	"errors"
	"os/exec"
	"runtime"
	"strings"
)

const GENERIC_ERR int = -230984

// CmdOutput is a struct containing the exit code, stdout and stderr of a command.
type CmdOutput struct {
	// Stdout_str is the stdout of the command as a string, with all line breaks replaced by \n.
	Stdout_str string
	// Stdout is the stdout of the command as a buffer.
	Stdout *bytes.Buffer
	// Stderr_str is the stderr of the command as a string, with all line breaks replaced by \n.
	Stderr_str string
	// Stderr is the stderr of the command as a buffer.
	Stderr *bytes.Buffer
	// Exit_code is the error code returned by the command in case no major error occurred
	Exit_code int
}

/*
ExecCmdSHELL executes a list of commands in a shell and returns the stdout and stderr.

On Windows, the command is executed in powershell.exe; on Linux, it's executed in bash. All elements of the list are
joined using "\n" as the command separator, given to the shell.

ATTENTION: to call any program, add "{{EXE}}" right after the program name in the command string. This will be replaced
by ".exe" on Windows and by "" on Linux. This avoids PowerShell aliases ("curl" is an alias for "Invoke-WebRequest" in
PowerShell but the actual program in Linux, for example - but curl.exe calls the actual program).

-----------------------------------------------------------

– Params:
  - commands_list – the commands to execute

– Returns:
  - the CmdOutput struct containing the stdout, stderr and error code of the command. Note that their string versions
    have all line endings replaced with "\n".
  - the error returned by the command execution, if any. Will be nil in case everything related to the command execution
    went smoothly - CmdOutput.Exit_code can still be non-zero! Will be non-nil if a major error occurred. in which case
    CmdOutput.Exit_code = GENERIC_ERR.
*/
func ExecCmdSHELL(commands_list[] string) (CmdOutput, error) {
	return ExecCmdMainSHELL(commands_list, "", "")
}

/*
ExecCmdMainSHELL is the main function for executing a list of commands in a shell. Check the documentation on
ExecCmdSHELL for more information.

-----------------------------------------------------------

– Params:
  - commands_list – the commansd to execute
  - windows_shell – the Windows shell, or "" to use the default (powershell.exe)
  - linux_shell – the Linux shell, or "" to use the default (bash)

– Returns:
  - the CmdOutput struct containing the stdout, stderr and error code of the command. Note that their string versions
*/
func ExecCmdMainSHELL(commands_list[] string, windows_shell string, linux_shell string) (CmdOutput, error) {
	var shell string = GetShell(windows_shell, linux_shell)

	var commands_str string = ""
	for _, command := range commands_list {
		commands_str += command
		if command != "" && !strings.HasSuffix(command, "\n") {
			commands_str += "\n"
		}
	}
	if runtime.GOOS == "windows" {
		commands_str = strings.Replace(commands_str, "{{EXE}}", ".exe", -1)
	} else {
		commands_str = strings.Replace(commands_str, "{{EXE}}", "", -1)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(shell)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(commands_str)
	err := cmd.Run()

	var stdout_str = strings.ReplaceAll(stdout.String(), "\r\n", "\n")
	stdout_str = strings.ReplaceAll(stdout_str, "\r", "\n")
	var stderr_str = strings.ReplaceAll(stderr.String(), "\r\n", "\n")
	stderr_str = strings.ReplaceAll(stderr_str, "\r", "\n")

	var exit_code int = 0
	if err != nil {
		var exiterr *exec.ExitError
		if errors.As(err, &exiterr) {
			exit_code = exiterr.ExitCode()
			err = nil
		} else {
			exit_code = GENERIC_ERR
		}
	}

	return CmdOutput{
		Stdout_str: stdout_str,
		Stdout:     &stdout,
		Stderr_str: stderr_str,
		Stderr:     &stderr,
		Exit_code:  exit_code,
	}, err
}

/*
GetShell returns the shell to use in the current OS.

-----------------------------------------------------------

– Params:
  - windows_shell – the Windows shell, or "" to use the default (powershell.exe)
  - linux_shell – the Linux shell, or "" to use the default (bash)
 */
func GetShell(windows_shell string, linux_shell string) string {
	if windows_shell == "" {
		windows_shell = "powershell.exe"
	}
	if linux_shell == "" {
		linux_shell = "bash"
	}
	var shell string = ""
	if runtime.GOOS == "windows" {
		shell = windows_shell
	} else {
		shell = linux_shell
	}

	return shell
}
