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

package Utils

import (
	"bufio"
	"errors"
	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/shirou/gopsutil/v4/process"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// _BIN_REL_DIR is the relative path to the binaries' directory from PersonalConsts._VISOR_DIR.
	_BIN_REL_DIR string = "bin"
	// _DATA_REL_DIR is the relative path to the data directory from PersonalConsts._VISOR_DIR.
	_DATA_REL_DIR string = "data"
	// _TEMP_FOLDER is the relative path to the temporary folder from PersonalConsts._VISOR_DIR.
	_TEMP_FOLDER string = _DATA_REL_DIR + "/Temp"
	// _USER_DATA_REL_DIR is the relative path to the user data directory from PersonalConsts._VISOR_DIR.
	_USER_DATA_REL_DIR string = _DATA_REL_DIR + "/UserData"
	// _PROGRAM_DATA_REL_DIR is the relative path to the program data directory from PersonalConsts._VISOR_DIR.
	_PROGRAM_DATA_REL_DIR string = _DATA_REL_DIR + "/ProgramData"
	// _WEBSITE_FILES_REL_DIR is the relative path to the website files directory from PersonalConsts._WEBSITE_DIR.
	_WEBSITE_FILES_REL_DIR string = "files_EOG"
)

// _MOD_FOLDER_PREFFIX is the preffix of the modules' folders.
const _MOD_FOLDER_PREFFIX string = "MOD_"

// _MOD_GEN_ERROR_CODE is the exit code of a module when a general error occurs.
const _MOD_GEN_ERROR_CODE int = 3234

const (
	NUM_MOD_VISOR           int = iota // This is a special one. Includes both the client and the server version main apps
	NUM_MOD_ModManager
	NUM_MOD_SMARTChecker
	NUM_MOD_Speech
	NUM_MOD_RssFeedNotifier
	NUM_MOD_EmailSender
	NUM_MOD_OnlineInfoChk
	NUM_MOD_GPTCommunicator
	NUM_MOD_WebsiteBackend
	NUM_MOD_TasksExecutor
	NUM_MOD_SystemChecker
	NUM_MOD_SpeechRecognition
	NUM_MOD_UserLocator
	NUM_MOD_CmdsExecutor

	MODS_ARRAY_SIZE
)
// MOD_NUMS_NAMES is a map of the numbers of the modules and their names. Use with the NUM_MOD_ constants.
var MOD_NUMS_NAMES map[int]string = map[int]string{
	NUM_MOD_VISOR:             "V.I.S.O.R.",
	NUM_MOD_ModManager:        "Modules Manager",
	NUM_MOD_SMARTChecker:      "S.M.A.R.T. Checker",
	NUM_MOD_Speech:            "Speech",
	NUM_MOD_RssFeedNotifier:   "RSS Feed Notifier",
	NUM_MOD_EmailSender:       "Email Sender",
	NUM_MOD_OnlineInfoChk:     "Online Information Checker",
	NUM_MOD_GPTCommunicator:   "GPT Communicator",
	NUM_MOD_WebsiteBackend:    "Website Backend",
	NUM_MOD_TasksExecutor:     "Tasks Executor",
	NUM_MOD_SystemChecker:     "System Checker",
	NUM_MOD_SpeechRecognition: "Speech Recognition",
	NUM_MOD_UserLocator:       "User Locator",
	NUM_MOD_CmdsExecutor:      "Commands Executor",
}

const (
	MOD_CLIENT int = 1 << 0
	MOD_SERVER int = 1 << 1
	MOD_BOTH   int = MOD_CLIENT | MOD_SERVER
)
// MOD_NUMS_SUPPORT is a map of the numbers of the modules and if they are supported on the server version, client
// version, or both.
var MOD_NUMS_SUPPORT map[int]int = map[int]int{
	NUM_MOD_VISOR:             MOD_BOTH,
	NUM_MOD_ModManager:        MOD_BOTH,
	NUM_MOD_SMARTChecker:      MOD_SERVER,
	NUM_MOD_Speech:            MOD_CLIENT,
	NUM_MOD_RssFeedNotifier:   MOD_SERVER,
	NUM_MOD_EmailSender:       MOD_SERVER,
	NUM_MOD_OnlineInfoChk:     MOD_SERVER,
	NUM_MOD_GPTCommunicator:   MOD_SERVER,
	NUM_MOD_WebsiteBackend:    MOD_SERVER,
	NUM_MOD_TasksExecutor:     MOD_CLIENT,
	NUM_MOD_SystemChecker:     MOD_CLIENT,
	NUM_MOD_SpeechRecognition: MOD_CLIENT,
	NUM_MOD_UserLocator:       MOD_CLIENT,
	NUM_MOD_CmdsExecutor:      MOD_CLIENT,
}

// _LOOP_TIME_S is the number of seconds to wait for the next timestamp to be registered by a module (must be more than
// a second higher than the actual time, for some reason).
const _LOOP_TIME_S int64 = 5

type _ModDirsInfo struct {
	// ProgramData is the path to the directory of the program data files.
	ProgramData GPath
	// UserData is the path to the directory of the private user data files.
	UserData GPath
	// Temp is the path to the directory of the private temporary files of the module.
	Temp GPath
}

type ModuleInfo struct {
	// Name is the name of the module.
	Name string
	// Num is the number of the module.
	Num int
	// ModDirsInfo is the information about the directories of the module.
	ModDirsInfo _ModDirsInfo
}

type Module struct {
	// Num is the number of the module.
	Num int
	// Name is the name of the module.
	Name string
	// Stop is set to true if the module should stop.
	Stop    bool
	// Stopped is set to true if the module has stopped.
	Stopped bool
	// Enabled is set to true if the module is enabled.
	Enabled bool
}

/*
RealMain is the type of the realMain() function of a module.

realMain is the function that does the actual work of a module (it's equivalent to what main() would normally be).

-----------------------------------------------------------

– Params:
  - module_stop – a pointer to a boolean that is set to true if the module should stop
  - moduleInfo_any – the ModuleInfo struct of the module
*/
type RealMain func(module_stop *bool, moduleInfo_any any)

/*
ModStartup does the startup routine for a module and executes its realMain() function, catching any fatal errors and
sending an email with them.

Call this as the ONLY thing in the Start() function of a module.

-----------------------------------------------------------

– Params:
  - realMain – a pointer to the realMain() function of the module
  - module – a pointer to the Module struct of the module
*/
func ModStartup(realMain RealMain, module *Module) {
	ModStartup2(realMain, module, false)
}
/*
ModStartup2 is the main function for ModStartup. Read everything there, except one different parameter.

-----------------------------------------------------------

– Params:
  - server – true if the version running is the server version, false if it's the client version
 */
func ModStartup2(realMain RealMain, module *Module, server bool) {
	// Module startup routine //

	var mod_num = module.Num
	var mod_name = module.Name

	if mod_num == NUM_MOD_VISOR {
		printStartupSequenceMODULES(mod_name)

		VISOR_server_GL = server

		if err := loadDeviceSettings(); err != nil {
			log.Println("CRITICAL ERROR: Error obtaining device settings - aborting")
			panic(err)
		}

		if err := loadGenSettings(server); err != nil {
			log.Println("warning: Error obtaining generated settings - aborting")
			log.Println(err)

			log.Println("Overwrite settings with empty file? Press ENTER to overwrite, or Ctrl+C to abort.")
			_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		}

		go func() {
			// Keep reloading the device/user settings and saving the generated settings global variables in case it's
			// MOD_0 that's running.
			for {
				if module.Stop {
					break
				}

				// Always reload the device settings
				err := loadDeviceSettings()
				if err != nil {
					module.Stop = true

					log.Println("warning: Error obtaining device settings - aborting")

					break
				}

				saveGenSettings(server)

				time.Sleep(5 * time.Second)
			}
		}()
	}

	if !IsModSupportedMODULES(mod_num) {
		panic(errors.New("module " + strconv.Itoa(mod_num) + " is not supported on this system"))
	}

	var moduleInfo ModuleInfo = ModuleInfo{
		Name:       mod_name,
		Num:        mod_num,
		ModDirsInfo: _ModDirsInfo{
			ProgramData: getProgramDataDirMODULES(mod_num),
			UserData:    GetUserDataDirMODULES(mod_num),
			Temp:        getModTempDirMODULES(mod_num),
		},
	}

	var errs bool = false
	var to_do func()

	if moduleInfo.signalledToStop() {
		log.Println("Module " + strconv.Itoa(mod_num) + " was signalled to stop before starting. Exiting...")

		goto end
	}

	// Start the loopSleep() routine asynchronously
	go func() {
		for {
			if moduleInfo.loopSleep() {
				module.Stop = true

				break
			}
		}
	}()

	to_do = func() {
		module.Stopped = false

		Tcef.Tcef{
			Try: func() {
				// Execute realMain()
				realMain(&module.Stop, moduleInfo)
			},
			Catch: func(e Tcef.Exception) {
				errs = true

				var str_error string = GetFullErrorMsgGENERAL(e)

				// Print the error and send an email with it
				log.Println(str_error)
				if err := SendModErrorEmailMODULES(mod_num, str_error); nil != err {
					log.Println("Error sending email with error:\n" + GetFullErrorMsgGENERAL(err) + "\n-----\n" + str_error)
				}
			},
		}.Do()

		module.Stopped = true
	}

	if mod_num == NUM_MOD_VISOR {
		// Don't run in another thread if it's the main program - it must be run on the main thread.

		if isVISORRunningMODULES(server) {
			log.Println("Module " + strconv.Itoa(mod_num) + " is already running. Exiting...")

			goto end
		}

		InitializeCommsChannels()

		moduleInfo.updateVISORRunInfo(server)

		to_do()
	} else {
		go func() {
			to_do()
		}()
	}

	end:

	if mod_num == NUM_MOD_VISOR {
		printShutdownSequenceMODULES(errs, moduleInfo.Name, moduleInfo.Num)

		// Delete the PID file
		var suffix = "_Client"
		if server {
			suffix = "_Server"
		}
		_ = GetUserDataDirMODULES(mod_num).Add2(false, "PID=" + strconv.Itoa(os.Getpid()) + suffix).Remove()

		if errs {
			os.Exit(_MOD_GEN_ERROR_CODE)
		}
	}
}

/*
GetModNameMODULES gets the name of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the name of the module or an empty string if the module number is invalid
*/
func GetModNameMODULES(mod_num int) string {
	if mod_name, ok := MOD_NUMS_NAMES[mod_num]; ok {
		return mod_name
	}

	return "INVALID MODULE NUMBER"
}

/*
SendModErrorEmailMODULES directly sends an email to the developer with the error message.

This function does *not* use any modules to do anything. Only utility functions. So it can be used from any
module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module from which the error occurred
  - error – the error message

– Returns:
  - nil if the email was sent successfully, otherwise an error
*/
func SendModErrorEmailMODULES(mod_num int, err_str string) error {
	var things_replace map[string]string = map[string]string{
		MODEL_INFO_MSG_BODY_EMAIL : err_str,
		MODEL_INFO_DATE_TIME_EMAIL: GetDateTimeStrTIMEDATE(-1),
	}
	var email_info = GetModelFileEMAIL(MODEL_FILE_INFO, things_replace)
	email_info.Subject = "Error in module: " + GetModNameMODULES(mod_num)

	if mod_num == NUM_MOD_EmailSender {
		// Send the email directly
		message_eml, mail_to, success := prepareEmlEMAIL(email_info)
		if !success {
			return errors.New("error preparing email")
		}

		return SendEmailEMAIL(message_eml, mail_to, true)
	} else {
		// Queue the email
		return QueueEmailEMAIL(email_info)
	}
}

/*
LoopSleep sleeps for _LOOP_TIME_S seconds and checks if the module was signalled to stop.

-----------------------------------------------------------

– Returns:
  - true if the module should stop, false otherwise
*/
func (moduleInfo *ModuleInfo) loopSleep() bool {
	if moduleInfo.signalledToStop() {
		return true
	}

	time.Sleep(time.Duration(_LOOP_TIME_S) * time.Second)

	return false
}

/*
signalledToStop checks if the module was signalled to stop.

-----------------------------------------------------------

– Returns:
  - true if the module was signalled to stop, false otherwise
*/
func (moduleInfo *ModuleInfo) signalledToStop() bool {
	var stop_file_1_path GPath = moduleInfo.ModDirsInfo.UserData.Add2(false, "STOP")
	var stop_file_2_path GPath = moduleInfo.ModDirsInfo.UserData.Add2(false, "STOP_p")
	var stop_file_3_path GPath = getVISORDirFILESDIRS().Add2(false, _USER_DATA_REL_DIR, "STOP")
	if stop_file_1_path.Exists() {
		err := stop_file_1_path.Remove()
		if nil != err {
			panic(err)
		}

		return true
	}
	if stop_file_2_path.Exists() || stop_file_3_path.Exists() {
		return true
	}

	return false
}

/*
printStartupSequenceMODULES prints the startup sequence of a module.

-----------------------------------------------------------

– Params:
  - mod_name – the name of the module
*/
func printStartupSequenceMODULES(mod_name string) {
	log.Println("//------------------------------------------\\\\")
	log.Println("--- " + mod_name + " ---")
	log.Println("V.I.S.O.R. Systems")
	log.Println("------------------")
	log.Println()
}

/*
printShutdownSequenceMODULES prints the shutdown sequence of a module.

-----------------------------------------------------------

– Params:
  - errors – true if the module is exiting with errors, false otherwise
  - mod_name – the name of the module
  - mod_num – the number of the module
*/
func printShutdownSequenceMODULES(errors bool, mod_name string, mod_num int) {
	log.Println()
	log.Println("---------")
	if errors {
		log.Println("Exiting with ERRORS the module \"" + mod_name + "\" (number " + strconv.Itoa(mod_num) + ")...")
	} else {
		log.Println("Exiting normally the module \"" + mod_name + "\" (number " + strconv.Itoa(mod_num) + ")...")
	}
	log.Println("\\\\------------------------------------------//")
}

/*
getProgramDataDirMODULES gets the full path to the program data directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the program data directory of the module
*/
func getProgramDataDirMODULES(mod_num int) GPath {
	return getVISORDirFILESDIRS().Add2(true, _PROGRAM_DATA_REL_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num))
}

/*
GetUserDataDirMODULES gets the full path to the private user data directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the private data directory of the module
*/
func GetUserDataDirMODULES(mod_num int) GPath {
	return getVISORDirFILESDIRS().Add2(true, _USER_DATA_REL_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num))
}

/*
getModTempDirMODULES gets the full path to the private temporary directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the private temporary directory of the module
*/
func getModTempDirMODULES(mod_num int) GPath {
	return getVISORDirFILESDIRS().Add2(true, _TEMP_FOLDER, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num))
}

/*
updateVISORRunInfo updates the information about the running of VISOR.

-----------------------------------------------------------

– Returns:
  - the path to the file containing the information about the running of the module
*/
func (moduleInfo *ModuleInfo) updateVISORRunInfo(server bool) {
	files, _ := os.ReadDir(GetUserDataDirMODULES(NUM_MOD_VISOR).GPathToStringConversion())

	var curr_pid string = strconv.Itoa(os.Getpid())
	var file_exists bool = false

	var suffix = "_Client"
	if server {
		suffix = "_Server"
	}

	// Remove all the old info files
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "PID=") && strings.HasSuffix(file.Name(), suffix) {
			var pid_str string = strings.Split(file.Name(), "=")[1]
			pid_str = strings.Split(pid_str, "_")[0]
			if pid_str != curr_pid {
				_ = moduleInfo.ModDirsInfo.UserData.Add2(false, file.Name()).Remove()
			} else {
				file_exists = true
			}
		}
	}

	if !file_exists {
		var new_info_file GPath = GetUserDataDirMODULES(NUM_MOD_VISOR).Add2(false, "PID=" + curr_pid + suffix)
		err := new_info_file.Create(true)
		if nil != err {
			panic(err)
		}
	}
}

/*
isVISORRunningMODULES checks if VISOR is already running.

-----------------------------------------------------------

– Params:
  - server – true if the version running is the server version, false if it's the client version

– Returns:
  - true if the module is running, false otherwise
*/
func isVISORRunningMODULES(server bool) bool {
	var curr_pid int = os.Getpid()
	files, err := os.ReadDir(GetUserDataDirMODULES(NUM_MOD_VISOR).GPathToStringConversion())
	if nil != err {
		return false
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "PID=") {
			if server && !strings.HasSuffix(file.Name(), "_Server") {
				continue
			} else if !server && !strings.HasSuffix(file.Name(), "_Client") {
				continue
			}

			var file_path GPath = GetUserDataDirMODULES(NUM_MOD_VISOR).Add2(false, file.Name())

			// File name example: PID=1243_Server
			var pid_str string = strings.Split(file.Name(), "=")[1]
			pid_str = strings.Split(pid_str, "_")[0]

			var pid int
			if pid, err = strconv.Atoi(pid_str); nil != err {
				_ = file_path.Remove()

				continue
			}

			id_pid_running, _ := process.PidExists(int32(pid))
			if pid != curr_pid && id_pid_running {
				return true
			}
		}
	}

	return false
}

/*
IsModSupportedMODULES checks if a module is supported on the current machine.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - true if the module is supported, false otherwise
 */
func IsModSupportedMODULES(mod_num int) bool {
	switch mod_num {
		case NUM_MOD_VISOR:
			return true
		case NUM_MOD_ModManager:
			return true
		case NUM_MOD_SMARTChecker:
			// Check if the command "smartctl" is available
			output, err := ExecCmdSHELL([]string{"smartctl{{EXE}} --version"})
			if err != nil {
				return false
			}

			return output.Exit_code == 0
		case NUM_MOD_Speech:
			return runtime.GOOS == "windows"
		case NUM_MOD_RssFeedNotifier:
			return true
		case NUM_MOD_EmailSender:
			// Check if the command "curl" is available
			output, err := ExecCmdSHELL([]string{"curl{{EXE}} --version"})
			if err != nil {
				return false
			}

			return output.Exit_code == 0
		case NUM_MOD_OnlineInfoChk:
			// Check if the command "chromedriver" is available
			output, err := ExecCmdSHELL([]string{"chromedriver{{EXE}} --version"})
			if err != nil {
				return false
			}

			return output.Exit_code == 0
		case NUM_MOD_GPTCommunicator:
			// Check if the command "llama-cli" is available
			output, err := ExecCmdSHELL([]string{"llama-cli{{EXE}} --version"})
			if err != nil {
				return false
			}

			return output.Exit_code == 0
		case NUM_MOD_WebsiteBackend:
			return true
		case NUM_MOD_TasksExecutor:
			return true
		case NUM_MOD_SystemChecker:
			return runtime.GOOS == "windows"
		case NUM_MOD_SpeechRecognition:
			return runtime.GOOS == "windows"
		case NUM_MOD_UserLocator:
			return true
		case NUM_MOD_CmdsExecutor:
			return true
		default:
			return false
	}
}

/*
SignalModulesStopMODULES signals all the modules to stop and waits for them to stop, also closing all communication
channels.

-----------------------------------------------------------

– Params:
  - modules – the list of modules
*/
func SignalModulesStopMODULES(modules []Module) {
	// Stop the modules gracefully before forcing an exit and wait for them to stop
	for {
		// Begin with the Manager (i := 1). VISOR doesn't count - of course it's running, else we wouldn't be here.

		for i := 1; i < MODS_ARRAY_SIZE; i++ {
			modules[i].Stop = true
		}

		var all_stopped bool = true
		for i := 1; i < MODS_ARRAY_SIZE; i++ {
			if !modules[i].Stopped {
				all_stopped = false

				break
			}
		}

		if all_stopped {
			break
		}

		time.Sleep(1 * time.Second)
	}

	modules[NUM_MOD_VISOR].Stop = true

	// Give time for threads to stop
	time.Sleep(1 * time.Second)
}
