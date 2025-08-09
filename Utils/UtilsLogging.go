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

package Utils

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var log_level_GL int = 90

const (
	LOG_LEVEL_ERROR int = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

func SetLogLevel(log_level int) {
	if log_level < LOG_LEVEL_ERROR {
		log_level_GL = LOG_LEVEL_ERROR
	} else if log_level > LOG_LEVEL_DEBUG {
		log_level_GL = LOG_LEVEL_DEBUG
	} else {
		log_level_GL = log_level
	}
}

/*
LogLnError logs an error message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnError(args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintln(LOG_LEVEL_ERROR, file, line, fn.Name(), args...)
}

/*
LogfError logs an error message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfError(format string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintf(LOG_LEVEL_ERROR, file, line, fn.Name(), format, args...)
}

/*
LogLnWarning logs a warning message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnWarning(args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintln(LOG_LEVEL_WARNING, file, line, fn.Name(), args...)
}

/*
LogfWarning logs a warning message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfWarning(format string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintf(LOG_LEVEL_WARNING, file, line, fn.Name(), format, args...)
}

/*
LogLnInfo logs an info message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnInfo(args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintln(LOG_LEVEL_INFO, file, line, fn.Name(), args...)
}

/*
LogfInfo logs an info message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfInfo(format string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintf(LOG_LEVEL_INFO, file, line, fn.Name(), format, args...)
}

/*
LogLnDebug logs a debug message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnDebug(args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintln(LOG_LEVEL_DEBUG, file, line, fn.Name(), args...)
}

/*
LogfDebug logs a debug message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfDebug(format string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	logPrintf(LOG_LEVEL_DEBUG, file, line, fn.Name(), format, args...)
}

/*
logPrintln prints a log message with the specified log level.

-----------------------------------------------------------

– Params:
  - log_level – the log level to use
  - file – the file where the log is being printed
  - line – the line number where the log is being printed
  - fn_name – the function name where the log is being printed
  - args – the arguments to log
*/
func logPrintln(log_level int, file string, line int, fn_name string, args ...any) {
	if log_level > log_level_GL {
		return
	}

	// "file" comes as the full path, we only want the file name
	var just_file string = file
	if strings.Contains(file, "/") {
		just_file = file[strings.LastIndex(file, "/") + 1:]
	}
	var prefix string = just_file + ":" + strconv.Itoa(line) + " (" + fn_name + "()):"

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	switch log_level {
		case LOG_LEVEL_ERROR:
			color.Set(color.FgHiRed)
			log.Println(prefix)
			fmt.Println("[E] --", args)
			color.Unset()
		case LOG_LEVEL_WARNING:
			color.Set(color.FgYellow)
			log.Println(prefix)
			fmt.Println("[W] --", args)
			color.Unset()
		case LOG_LEVEL_INFO:
			color.Set(color.FgCyan)
			log.Println(prefix)
			fmt.Println("[I] --", args)
			color.Unset()
		case LOG_LEVEL_DEBUG:
			log.Println(prefix)
			fmt.Println("[D] --", args)
		default:
			// Won't get here
	}

	log.SetFlags(log.LstdFlags)
}

/*
logPrintf prints a log message with the specified log level and format.

-----------------------------------------------------------

– Params:
  - log_level – the log level to use
  - file – the file where the log is being printed
  - line – the line number where the log is being printed
  - fn_name – the function name where the log is being printed
  - format – the format string
  - args – the arguments to log
 */
func logPrintf(log_level int, file string, line int, fn_name string, format string, args ...any) {
	if log_level > log_level_GL {
		return
	}

	// "file" comes as the full path, we only want the file name
	var just_file string = file
	if strings.Contains(file, "/") {
		just_file = file[strings.LastIndex(file, "/") + 1:]
	}
	var prefix string = just_file + ":" + strconv.Itoa(line) + " (" + fn_name + "()):"

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	switch log_level {
		case LOG_LEVEL_ERROR:
			color.Set(color.FgHiRed)
			log.Println(prefix)
			fmt.Printf("[E] -- " + format, args...)
			color.Unset()
		case LOG_LEVEL_WARNING:
			color.Set(color.FgYellow)
			log.Println(prefix)
			fmt.Printf("[W] -- " + format, args...)
			color.Unset()
		case LOG_LEVEL_INFO:
			color.Set(color.FgCyan)
			log.Println(prefix)
			fmt.Printf("[I] -- " + format, args...)
			color.Unset()
		case LOG_LEVEL_DEBUG:
			log.Println(prefix)
			fmt.Printf("[D] -- " + format, args...)
		default:
			// Won't get here
	}

	log.SetFlags(log.LstdFlags)
}
