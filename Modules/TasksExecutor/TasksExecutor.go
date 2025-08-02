/*******************************************************************************
 * Copyright 2023-2025 The V.I.S.O.R. authors
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

package TasksExecutor

import (
	"Speech"
	"SpeechQueue"
	"TEHelper"
	"Utils"
	"Utils/ModsFileInfo"
)

var (
	modDirsInfo_GL  Utils.ModDirsInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(main, module)}
func main(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	go func() {
		for {
			var task *ModsFileInfo.Task = TEHelper.CheckDueTasks()
			if task == nil {
				break
			}

			Utils.LogLnDebug(task.Id)

			if task.Message != "" {
				Speech.QueueSpeech(task.Message, SpeechQueue.PRIORITY_MEDIUM, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
			}

			if task.Command != "" {
				Utils.SendToModChannel(Utils.NUM_MOD_CmdsExecutor, 0, "SentenceInternal", task.Command)
			}
		}
	}()

	for {
		if Utils.WaitWithStopDATETIME(module_stop, 1000000000) {
			TEHelper.StopChecker()

			return
		}
	}
}
