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

package main

import (
	"fmt"

	"Utils"
)

func main() {
	fmt.Println("Hello, World!")

	/*password1 := []byte("this is a test")
	password2 := []byte("this is one other test")
	message := utf7.UTF7EncodeBytes([]byte("this is another test ´1ºªá¨nñë€§«"))
	associated_data := []byte("Test 44")

	fmt.Println(string(password1))
	fmt.Println(string(password2))
	tmp, _ := utf7.UTF7DecodeBytes(message)
	fmt.Println(string(tmp))
	fmt.Println(string(associated_data))

	bytes:= UtilsSWA.EncryptBytesCRYPTOENDECRYPT(password1, password2, message, associated_data)
	fmt.Println(string(bytes))

	bytes = UtilsSWA.DecryptBytesCRYPTOENDECRYPT(password1, password2, bytes, associated_data)
	tmp, _ = utf7.UTF7DecodeBytes(bytes)
	fmt.Println(string(tmp))

	fmt.Println(UtilsSWA.BytesToHexDATACONV(tmp))
	fmt.Println(UtilsSWA.BytesToOctalDATACONV(tmp))*/

	/*var commands_list string = "bash" + UtilsSWA.CMD_SEP + "ps -p $$"
	output, err := UtilsSWA.ExecCmdSHELL(commands_list, false)
	if err != nil {
		fmt.Println(err)
		fmt.Println("-----------")
		fmt.Println(UtilsSWA.GetExitCodeSHELL(output))
	} else {
		fmt.Println(err)
		fmt.Println("-----------")
		fmt.Println(UtilsSWA.GetExitCodeSHELL(output))
		fmt.Println("-----------")
		fmt.Println(string(UtilsSWA.GetStdoutSHELL(output)))
		fmt.Println("-----------")
		fmt.Println(string(UtilsSWA.GetStderrSHELL(output)))
	}*/

	var commands_list []string = []string{
		"ps -p $$",
		"bash",
		"ps -p $$",
	}
	output, err := Utils.ExecCmdSHELL(commands_list)
	if err != nil {
		fmt.Println(err)
		fmt.Println("-----------")
		fmt.Println(output.Exit_code)
	} else {
		fmt.Println(err)
		fmt.Println("-----------")
		fmt.Println(output.Exit_code)
		fmt.Println("-----------")
		fmt.Println(output.Stdout_str)
		fmt.Println("-----------")
		fmt.Println(output.Stderr_str)
	}
}
