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

package UtilsSWA

import (
	"strings"

	"Utils"
)

/*
PathFILESDIRS returns the path and whether it describes a directory or not.

-----------------------------------------------------------

– Params:
  - paths_list – the paths list with subpaths separated by the null character

– Returns:
  - the path and whether it describes a directory or not separated by "|"
*/
func PathFILESDIRS(describes_dir bool, paths_list string) string {
	var array1 []string = strings.Split(paths_list, "\000")
	var array2 []any = make([]any, len(array1))
	for i := 0; i < len(array1); i++ {
		array2[i] = array1[i]
	}

	var gPath Utils.GPath = Utils.PathFILESDIRS(describes_dir, "/", array2...)
	var describes_dir_str string = "false"
	if gPath.DescribesDir() {
		describes_dir_str = "true"
	}

	return gPath.GPathToStringConversion() + "|" + describes_dir_str
}
