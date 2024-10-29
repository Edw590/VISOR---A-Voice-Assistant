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

package TEHelper

import (
	"Utils/ModsFileInfo"
	"strconv"
)

/*
GetTasks returns all tasks.

-----------------------------------------------------------

– Returns:
  - all tasks
*/
func GetTasks() []ModsFileInfo.Task {
	return modUserInfo_GL.Tasks
}

/*
GetIdsList returns a list of all tasks' IDs.

-----------------------------------------------------------

– Returns:
  - a list of all tasks' IDs separated by "|"
*/
func GetIdsList() string {
	var ids string
	for _, task := range modUserInfo_GL.Tasks {
		ids += strconv.Itoa(task.Id) + "|"
	}

	return ids
}

/*
GetTaskById returns a task by its ID.

-----------------------------------------------------------

– Params:
  - id – the task ID

– Returns:
  - the task or nil if the task was not found
*/
func GetTaskById(id int) *ModsFileInfo.Task {
	for _, task := range modUserInfo_GL.Tasks {
		if task.Id == id {
			return &task
		}
	}

	return nil
}
