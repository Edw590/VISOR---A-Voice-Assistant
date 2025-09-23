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
	"strings"

	"github.com/fatih/color"
)

const _MAX_FILE_NAME_LEN int = 6
const _MAX_LINE_NUM_LEN int = 4

var log_level_GL int = 90

const (
	LOG_LEVEL_ERROR int = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
)

/*
SetLogLevel sets the global log level.

-----------------------------------------------------------

– Params:
  - log_level – the log level to set (one of the LOG_LEVEL_* constants)
 */
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
	_, file, line, _ := runtime.Caller(1)
	logPrintln(LOG_LEVEL_ERROR, file, line, args...)
}

/*
LogfError logs an error message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfError(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintf(LOG_LEVEL_ERROR, file, line, format, args...)
}

/*
LogLnWarning logs a warning message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnWarning(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintln(LOG_LEVEL_WARNING, file, line, args...)
}

/*
LogfWarning logs a warning message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfWarning(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintf(LOG_LEVEL_WARNING, file, line, format, args...)
}

/*
LogLnInfo logs an info message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnInfo(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintln(LOG_LEVEL_INFO, file, line, args...)
}

/*
LogfInfo logs an info message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfInfo(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintf(LOG_LEVEL_INFO, file, line, format, args...)
}

/*
LogLnDebug logs a debug message.

-----------------------------------------------------------

– Params:
  - args – the arguments to log
*/
func LogLnDebug(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintln(LOG_LEVEL_DEBUG, file, line, args...)
}

/*
LogfDebug logs a debug message with a format.

-----------------------------------------------------------

– Params:
  - format – the format string
  - args – the arguments to log
 */
func LogfDebug(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	logPrintf(LOG_LEVEL_DEBUG, file, line, format, args...)
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
func logPrintln(log_level int, file string, line int, args ...any) {
	if log_level > log_level_GL {
		return
	}

	var caps_file_name string = GetInitialsOfFileNameLOGGING(file)
	var line_str string = FormatLineNumberLOGGING(line)
	var middle_str string = caps_file_name + ":" + line_str + "|>"

	switch log_level {
		case LOG_LEVEL_ERROR:
			color.Set(color.FgHiRed)
			log.Println("-- E:" + middle_str, fmt.Sprint(args...))
			color.Unset()
		case LOG_LEVEL_WARNING:
			color.Set(color.FgYellow)
			log.Println("-- W:" + middle_str, fmt.Sprint(args...))
			color.Unset()
		case LOG_LEVEL_INFO:
			color.Set(color.FgCyan)
			log.Println("-- I:" + middle_str, fmt.Sprint(args...))
			color.Unset()
		case LOG_LEVEL_DEBUG:
			log.Println("-- D:" + middle_str, fmt.Sprint(args...))
		default:
			// Won't get here
	}
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
func logPrintf(log_level int, file string, line int, format string, args ...any) {
	if log_level > log_level_GL {
		return
	}

	var caps_file_name string = GetInitialsOfFileNameLOGGING(file)
	var line_str string = FormatLineNumberLOGGING(line)
	var middle_str string = caps_file_name + ":" + line_str + "|> "

	switch log_level {
		case LOG_LEVEL_ERROR:
			color.Set(color.FgHiRed)
			log.Printf("-- E:" + middle_str + format, args...)
			color.Unset()
		case LOG_LEVEL_WARNING:
			color.Set(color.FgYellow)
			log.Printf("-- W:" + middle_str + format, args...)
			color.Unset()
		case LOG_LEVEL_INFO:
			color.Set(color.FgCyan)
			log.Printf("-- I:" + middle_str + format, args...)
			color.Unset()
		case LOG_LEVEL_DEBUG:
			log.Printf("-- D:" + middle_str + format, args...)
		default:
			// Won't get here
	}
}

/*
FormatLineNumberLOGGING formats a line number to a fixed length string.

-----------------------------------------------------------

– Params:
  - line – the line number to format

– Returns:
  - the formatted line number string, up to _MAX_LINE_NUM_LEN characters
*/
func FormatLineNumberLOGGING(line int) string {
	var line_str string = fmt.Sprintf("%d", line)
	// Up to _MAX_LINE_NUM_LEN characters
	if len(line_str) > _MAX_LINE_NUM_LEN {
		line_str = line_str[:_MAX_LINE_NUM_LEN-2] + ".."
	} else {
		// But always _MAX_LINE_NUM_LEN characters
		for len(line_str) < _MAX_LINE_NUM_LEN {
			line_str += " "
		}
	}

	return line_str
}

/*
GetInitialsOfFileNameLOGGING returns the initials of a file name.

-----------------------------------------------------------

– Params:
  - file_name – the full path of the file

– Returns:
  - the initials of the file name, up to _MAX_FILE_NAME_LEN characters
*/
func GetInitialsOfFileNameLOGGING(file_name string) string {
	// Note: this function is prepared for Pascal Case only

	// "file_name" comes as the full path, we only want the file name
	var just_file string = file_name
	if strings.Contains(file_name, "/") {
		just_file = file_name[strings.LastIndex(file_name, "/")+1:]
	}

	// Get all the first letters of each word of the file name
	// e.g. "ModGPTCommunicator1.go" -> "MGC1"
	var letters string = ""
	var prev_was_upper bool = false
	var next_is_upper bool = false
	for n, char := range just_file {
		// If we reached the file extension, stop
		if char == '.' {
			break
		}

		if n + 1 < len(just_file) {
			next_is_upper = just_file[n + 1] >= 'A' && just_file[n + 1] <= 'Z'
		}
		if n - 1 >= 0 {
			prev_was_upper = just_file[n - 1] >= 'A' && just_file[n - 1] <= 'Z'
		}

		if prev_was_upper && next_is_upper {
			continue
		}

		if (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			letters += string(char)
		}
	}

	// Up to _MAX_FILE_NAME_LEN characters
	if len(letters) > _MAX_FILE_NAME_LEN {
		letters = letters[:_MAX_FILE_NAME_LEN-2] + ".."
	} else {
		// But always _MAX_FILE_NAME_LEN characters
		for len(letters) < _MAX_FILE_NAME_LEN {
			letters = " " + letters
		}
	}

	return letters
}
