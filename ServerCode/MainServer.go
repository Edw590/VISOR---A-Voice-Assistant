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

package main

import (
	"ModulesManager"
	"Utils"
	"VISOR_Server/ServerRegKeys"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	flag "github.com/spf13/pflag"
)

var modDirsInfo_GL Utils.ModDirsInfo
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Command line arguments
	var flag_log_level *int = flag.IntP("loglevel", "l", 0, "Log level to use. 0 = ERROR, 1 = WARNING, 2 = INFO, 3 = DEBUG. Default is 0 (ERROR).")
	flag.Bool("status", false, "Keeps printing the status of the modules.")
	flag.Parse()
	Utils.SetLogLevel(*flag_log_level)

	var module Utils.Module = Utils.Module{
		Num:     Utils.NUM_MOD_VISOR,
		Name:    Utils.GetModNameMODULES(Utils.NUM_MOD_VISOR),
		Stop:    false,
		Stopped: false,
		Enabled: true,
	}
	Utils.ModStartup2(realMain, &module, true)
}
func realMain(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	if Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id == "" || Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Type_ == "" ||
			Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Description == "" {
		log.Println("Device settings incomplete. Please enter the missing one(s):")
		if Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id == "" {
			Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Id = Utils.GetInputString("Unique device ID: ")
		}
		if Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Type_ == "" {
			Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Type_ = Utils.GetInputString("Device type (for example " +
				"\"computer\"): ")
		}
		if Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Description == "" {
			Utils.GetGenSettings(Utils.LOCK_UNLOCK).Device_settings.Description = Utils.GetInputString("Device description (for " +
				"example the model, \"Legion Y520\"): ")
		}
	}

	if err := Utils.ReadSettingsFile(true); err != nil {
		Utils.LogLnError("Failed to load user settings. Exiting...")

		return
	}

	if !Utils.RunningAsAdminPROCESSES() {
		Utils.LogLnError("Not running as administrator/root. Exiting...")

		return
	}

	ServerRegKeys.RegisterValues()

	var modules []Utils.Module
	for i := 0; i < Utils.MODS_ARRAY_SIZE; i++ {
		modules = append(modules, Utils.Module{
			Num:     i,
			Name:    Utils.GetModNameMODULES(i),
			Stop:    true,
			Stopped: true,
			Enabled: true,
		})
	}
	modules[Utils.NUM_MOD_VISOR].Stop = false
	modules[Utils.NUM_MOD_VISOR].Stopped = false
	// The Manager needs to be started first. It'll handle the others.
	modules[Utils.NUM_MOD_ModManager].Stop = false

	// Empty the active device IDs list as soon as the server starts. Else devices previously active might be used and
	// stuff sent to them waiting for a response - this is because other modules start before Website Backend which also
	// empties the list.
	Utils.GetGenSettings(Utils.LOCK_UNLOCK).MOD_8.Active_device_IDs = nil

	ModulesManager.Start(modules)

	handleCtrlCGracefully(module_stop)

	var status bool = Utils.WasArgUsedGENERAL(os.Args, "--status")

	for {
		// Wait forever while the other modules do their work
		if status {
			printModulesStatus(modules)
		}

		Utils.WriteSettingsFile(true)

		if Utils.WaitWithStopDATETIME(module_stop, 5) {
			break
		}
	}

	Utils.CloseCommsChannels()
	Utils.SignalModulesStopMODULES(modules)
}

func handleCtrlCGracefully(module_stop *bool) {
	// Copied from https://gist.github.com/jnovack/297cee036f3e5a430aa9444c0ae1b06d
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<- c

		*module_stop = true
	}()
}

func printModulesStatus(modules []Utils.Module) {
	Utils.LogLnDebug("--------------------------------")
	for _, module := range modules {
		Utils.LogLnDebug("--- " + module.Name + " ---")
		Utils.LogLnDebug("- Enabled: " + strconv.FormatBool(module.Enabled) + " | To stop: " +
			strconv.FormatBool(module.Stop) + " | Running: " + strconv.FormatBool(!module.Stopped) + " | Supported: " +
			strconv.FormatBool(Utils.IsModSupportedMODULES(module.Num)))
	}
}
