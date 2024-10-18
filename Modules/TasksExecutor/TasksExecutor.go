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
	"TEHelper/TEHelper"
	"Utils"
	"Utils/ModsFileInfo"
	"log"
)

// Tasks Executor //

var (
	realMain       Utils.RealMain = nil
	moduleInfo_GL  Utils.ModuleInfo
	modGenInfo_GL  *ModsFileInfo.Mod9GenInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_9

		var user_location *ModsFileInfo.UserLocation = &Utils.Gen_settings_GL.MOD_12.User_location
		log.Println("User location:", user_location)

		TEHelper.UpdateUserLocation(user_location)

		//TEHelper.LoadLocalTasks()

		go func() {
			for {
				TEHelper.CheckDueTasks()
			}
		}()

		for {
			if Utils.WaitWithStopTIMEDATE(module_stop, 1000000000) {
				TEHelper.StopChecker()

				return
			}
		}
	}
}
