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

package TasksExecutor

import (
	"Speech"
	"SpeechQueue/SpeechQueue"
	"TEHelper/TEHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"log"
)

// Tasks Executor //

var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		go func() {
			for {
				var task *ModsFileInfo.Task = TEHelper.CheckDueTasks()
				if task == nil {
					break
				}

				log.Println("Task! -->", task.Id)

				if task.Message != "" {
					Speech.QueueSpeech(task.Message, SpeechQueue.PRIORITY_MEDIUM, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
				}

				if task.Command != "" {
					Utils.SendToModChannel(Utils.NUM_MOD_CmdsExecutor, "SentenceInternal", task.Command)
				}
			}
		}()

		for {
			TEHelper.UpdateUserLocation(&Utils.Gen_settings_GL.MOD_12.User_location)

			if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
				TEHelper.StopChecker()

				return
			}
		}
	}
}
