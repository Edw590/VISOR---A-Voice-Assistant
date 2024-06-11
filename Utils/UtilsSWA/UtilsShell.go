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

package UtilsSWA

import (
	"bytes"
	"encoding/binary"
	"strings"

	"Utils"
)

const GENERIC_ERR int32 = int32(Utils.GENERIC_ERR)

// 50 bytes to separate the stdout and stderr in the output of ExecCmdSHELL
var _OUTPUT_SEP []byte = []byte("(K!5pSqW=.h9s60EA'ryI.jS@6SY&uy),qbo4sFWQ_(%@H&(bC")

/*
ExecCmdSHELL executes a list of commands in a shell and returns the stdout and stderr.



-----------------------------------------------------------

– Params:
  - attempt_su – whether to attempt to execute the commands as root (using su -c)
  - command_list – the list of commands to execute, separated by "\n", joined internally by the function

– Returns:
  - first 4 bytes of the output are the exit code in Big Endian; then comes stdout; after that comes _OUTPUT_SEP;
    finally comes stderr
  - the error returned by the command execution, if any. Will be nil in case everything related to the command execution
    went smoothly - CmdOutput.Error_code can still be non-zero! Will be non-nil if a major error occurred. in which case
    CmdOutput.Error_code = GENERIC_ERR.
*/
func ExecCmdSHELL(attempt_su bool, commands_list string) ([]byte, error) {

	// Android needs the full path specified for some reason
	const ANDROID_SH string = "/system/bin/sh"

	// And also Android needs this environment variable set for some reason for some commands to work
	// (https://xdaforums.com/t/running-svc-in-ssh-returns-aborted.4274735/post-85384851 and
	// https://github.com/Magisk-Modules-Repo/ssh/issues/12)
	commands_list = "export ANDROID_DATA=/data\n" + commands_list

	if attempt_su && IsRootAvailable() {
		commands_list = "su\n" + commands_list
	}
	cmd_output, err := Utils.ExecCmdMainSHELL(strings.Split(commands_list, "\n"), "", ANDROID_SH)

	exit_code := make([]byte, 4)
	binary.BigEndian.PutUint32(exit_code, uint32(cmd_output.Exit_code))
	var output []byte = exit_code
	output = append(output, cmd_output.Stdout.Bytes()...)
	output = append(output, _OUTPUT_SEP...)
	output = append(output, cmd_output.Stderr.Bytes()...)

	return output, err
}

func GetExitCodeSHELL(cmd_output []byte) int32 {
	return int32(binary.BigEndian.Uint32(cmd_output[0:4]))
}

func GetStdoutSHELL(cmd_output []byte) []byte {
	return cmd_output[4:bytes.Index(cmd_output, _OUTPUT_SEP)]
}

func GetStderrSHELL(cmd_output []byte) []byte {
	return cmd_output[bytes.Index(cmd_output, _OUTPUT_SEP) + len(_OUTPUT_SEP):]
}
