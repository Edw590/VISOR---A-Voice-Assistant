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

package SMARTChecker

import (
	"Utils/ModsFileInfo"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"Utils"
)

// S.M.A.R.T. Checker //

const _LONG_TEST_EACH_S int64 = 30*24*60*60 // 30 days
const _SHORT_TEST_EACH_S int64 = 5*60*60*99999999999 // oo hours (test disabled)//5*60*60 // 5 hours

const _TIME_SLEEP_S int = 60

var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo
	modGenInfo_GL   *ModsFileInfo.Mod3GenInfo
	modUserInfo_GL  *ModsFileInfo.Mod3UserInfo
)
func Start(module *Utils.Module) {Utils.ModStartup(realMain, module)}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)
		modGenInfo_GL = &Utils.Gen_settings_GL.MOD_3
		modUserInfo_GL = &Utils.User_settings_GL.SMARTChecker

		if !Utils.RunningAsAdminPROCESSES() {
			panic(errors.New("this program must be run as administrator/root"))
		}

		const NO_TEST = -1
		const SHORT_TEST = 0
		const LONG_TEST = 1

		var no_test bool = false
		for _, arg := range os.Args {
			if arg == "--notest" {
				no_test = true

				//log.Println("No tests to execute. Will only send the reports.")
			}
		}

		for {
			var disks_to_chk map[string]*ModsFileInfo.DiskInfo = modUserInfo_GL.Disks_info
			if len(disks_to_chk) == 0 {
				//log.Println("No disks to check.")

				goto end_loop
			}

			//log.Println("Checking listed disks...")
			//log.Println()

			for disk_serial, disk_user_info := range disks_to_chk {
				//log.Println("------------------------------------")
				//log.Println("Disk serial: " + disk_serial)
				//log.Println("Disk label: " + disk_user_info.Label)
				//log.Println()
				var partitions_list []string = getAllAvailablePartitions()
				var disksSerialPartitions map[string]string = getDiskSerialPartitions(partitions_list)

				disk_partition, ok := disksSerialPartitions[disk_serial]
				if !ok {
					//log.Println("Disk not found on the system")

					continue
				}

				// Fill the map with spots for the 2 elements that must be in them
				if modGenInfo_GL.Disks_info == nil {
					modGenInfo_GL.Disks_info = make(map[string][]int64, 2)
				}

				// Check which test to execute, or execute none if the time hasn't passed yet.
				var disk_gen_info []int64 = modGenInfo_GL.Disks_info[disk_serial]
				if disk_gen_info == nil {
					disk_gen_info = make([]int64, 2)
					disk_gen_info[SHORT_TEST] = 0
					disk_gen_info[LONG_TEST] = 0
					modGenInfo_GL.Disks_info[disk_serial] = disk_gen_info
				}
				var test_type int = NO_TEST
				if !no_test {
					if time.Now().Unix()-disk_gen_info[LONG_TEST] > _LONG_TEST_EACH_S {
						test_type = LONG_TEST
					} else if time.Now().Unix()-disk_gen_info[SHORT_TEST] > _SHORT_TEST_EACH_S {
						test_type = SHORT_TEST
					} else {
						//log.Println("Time has not passed yet for the tests to be executed.")

						continue
					}

					if disk_user_info.Is_HDD {
						// If enough time passed already, check if the disk is spinning or not.
						if !Utils.ContainsSLICES(getActiveDisks(partitions_list), disk_partition) {
							// If disk is not spinning skip short test or status check, but never skip if it's a long test.
							if test_type != LONG_TEST {
								//log.Println("Disk not spinning, skipping test.")

								continue
							}
						}
					} else {
						// Never spins, so always execute the test.
					}
				}

				if test_type != NO_TEST {
					// Start the test and retrieve the test time
					var test_time int = initiateTest(test_type == LONG_TEST, disk_partition)

					// The total waiting time is the time the test will take in minutes + some time to make sure it's
					// finished.
					var test_type_str string
					if test_type == SHORT_TEST {
						test_type_str = "Short"
					} else {
						test_type_str = "Long"
					}

					if test_time == -1 {
						var msg_body string = "ATTENTION - The " + strings.ToLower(test_type_str) +
							" test could not be executed. An attempt was made to abort any testing being " +
							"executed, but still it was not possible to execute the test. Please execute it " +
							"manually."
						var things_replace = map[string]string{
							Utils.MODEL_INFO_DATE_TIME_EMAIL: Utils.GetDateTimeStrTIMEDATE(-1),
							Utils.MODEL_INFO_MSG_BODY_EMAIL:  msg_body,
						}
						var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_INFO, things_replace)
						email_info.Subject = test_type_str + " test could NOT be started on " + disk_user_info.Label
						Utils.QueueEmailEMAIL(email_info)

						//log.Println(msg_body)

						continue
					}

					var millis_begin int64 = time.Now().UnixMilli()
					var date_time_begin string = Utils.GetDateTimeStrTIMEDATE(millis_begin)

					// Here, this will wait until the test is concluded to report the log (test time + some waiting
					// period to be sure the test is over).
					// Then, this sends an email warning the test has began. That's in case a test takes 4 hours and
					// after 4:30 hours there's no email, for example, then there must have been a problem in the script
					// (not supposed).
					var msg_body string = test_type_str + " test on " + disk_user_info.Label + " (" + disk_partition + ") " +
						"started on " + date_time_begin + ".\n\n" +
						"Test duration : " + strconv.Itoa(test_time) + " minutes.\n\n" +
						"The results will be ready on or before " +
						Utils.GetDateTimeStrTIMEDATE(millis_begin+int64(test_time)*60*1000) + "."
					var things_replace = map[string]string{
						Utils.MODEL_INFO_DATE_TIME_EMAIL: Utils.GetDateTimeStrTIMEDATE(-1),
						Utils.MODEL_INFO_MSG_BODY_EMAIL:  msg_body,
					}
					var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_INFO, things_replace)
					email_info.Subject = test_type_str + " test started on " + disk_user_info.Label
					Utils.QueueEmailEMAIL(email_info)
					//log.Println("Notice email queued")

					/////////////////////////////////////////////////////////
					//log.Println("----------------------------------------------------------------------------")
					//log.Println(msg_body)
					//log.Println("----------------------------------------------------------------------------")
					//log.Println()

					// Wait for the test to finish (it can finish before the supposed time, so this checks every minute).
					for {
						if !checkDiskInTest(disk_partition) {
							break
						}

						if Utils.WaitWithStopTIMEDATE(module_stop, 60) {
							return
						}
					}

					//log.Println("Test finished.")

					// Update the timestamp
					if test_type == SHORT_TEST {
						disk_gen_info[SHORT_TEST] = time.Now().Unix()
					} else {
						disk_gen_info[LONG_TEST] = time.Now().Unix()
					}
				}

				html_report := getHTMLReport(disk_partition)

				var things_replace = map[string]string{
					Utils.MODEL_DISKS_SMART_DISK_LABEL_EMAIL:       disk_user_info.Label,
					Utils.MODEL_DISKS_SMART_DISK_SERIAL_EMAIL:      disk_serial,
					Utils.MODEL_DISKS_SMART_DISK_PARTITION_EMAIL:   disk_partition,
					Utils.MODEL_DISKS_SMART_DISKS_SMART_HTML_EMAIL: html_report,
				}

				var email_info Utils.EmailInfo = Utils.GetModelFileEMAIL(Utils.MODEL_FILE_DISKS_SMART, things_replace)
				email_info.Subject = "S.M.A.R.T. report on " + disk_user_info.Label
				Utils.QueueEmailEMAIL(email_info)
				//log.Println("Report email queued")
			}

			end_loop:

			if no_test {
				// If it's to execute no test, don't loop the program, just send the report and exit.
				return
			}

			if Utils.WaitWithStopTIMEDATE(module_stop, _TIME_SLEEP_S) {
				return
			}
		}
	}
}
