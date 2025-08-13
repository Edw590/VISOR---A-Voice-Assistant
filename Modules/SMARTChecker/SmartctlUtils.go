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

package SMARTChecker

import (
	"Utils"
	"strconv"
	"strings"
)

// Critical attributes, as per the Wikipedia page on SMART
var critical_attributes_GL []int = []int{
	1,   // Read Error Rate
	5,   // Reallocated Sectors Count
	10,  // Spin Retry Count
	184, // End-to-End error / IOEDC
	187, // Reported Uncorrectable Errors
	188, // Command Timeout
	196, // Reallocation Event Count
	197, // Current Pending Sector Count
	198, // (Offline) Uncorrectable Sector Count
	201, // Soft Read Error Rate or TA Counter Detected
}

func getHTMLReport(disk_partition string) string {
	stdOutErrCmd, _ := Utils.ExecCmdSHELL([]string{"smartctl{{EXE}} -x " + disk_partition})

	var output []string = strings.Split(stdOutErrCmd.Stdout_str, "\n")

	for i, line := range output {
		// Process the "Vendor Specific SMART Attributes with Thresholds" section
		if strings.Contains(line, "ID#") {
			for j := i + 1; j < len(output); j++ {
				if strings.Contains(output[j], "K auto-keep") {
					// Means the section just ended
					break
				}

				var line_trimmed []string = strings.Fields(output[j])

				id, _ := strconv.Atoi(line_trimmed[0])
				//var attribute_name string = strings.ToLower(line_trimmed[1])
				var flags string = line_trimmed[2]
				//value, _ := strconv.Atoi(line_trimmed[3])
				worst, _ := strconv.Atoi(line_trimmed[4])
				threshold, _ := strconv.Atoi(line_trimmed[5])
				var fail string = line_trimmed[6]
				raw_value, _ := strconv.Atoi(line_trimmed[7])

				var is_critical bool = false
				for _, critical_id := range critical_attributes_GL {
					if id == critical_id {
						is_critical = true

						break
					}
				}

				if strings.Contains(flags, "P") || is_critical {
					// Pre-fail or critical attribute (underline it)
					output[j] = "<u>" + output[j] + "</u>"
				}

				if fail != "-" {
					// Failing now or has failed (red)
					output[j] = "<strong style='color:#FF0000;'>" + output[j] + "</strong>"

					continue
				}

				if is_critical && raw_value > 0 {
					// Critical attribute is above 0 (yellow)
					output[j] = "<strong style='color:#FFA500;'>" + output[j] + "</strong>"

					continue
				}

				if worst <= threshold {
					// Attribute is below or on threshold (red)
					output[j] = "<strong style='color:#FF0000;'>" + output[j] + "</strong>"

					continue
				} else if worst - int(float64(threshold)*0.2) <= threshold {
					// Attribute is close to threshold (orange)
					output[j] = "<strong style='color:#FF7800;'>" + output[j] + "</strong>"

					continue
				} else {
					// Attribute is above threshold, nothing to do
				}
			}
		}


		// Important titles here
		if
		strings.Contains(line, "SMART overall-health self-assessment test result:") ||
			strings.Contains(line, "Vendor Specific SMART Attributes with Thresholds:") ||
			strings.Contains(line, "SMART Error Log Version:") ||
			strings.Contains(line, "SMART Self-test log structure revision number") ||
			strings.Contains(line, "SMART Extended Comprehensive Error Log Version:") ||
			strings.Contains(line, "SMART Extended Self-test Log Version:") ||
			strings.Contains(line, "SCT Temperature History Version:") {
			output[i] = "<strong><em><u>" + line + "</u></em></strong>"
		}

		// Error detection
		if
			strings.Contains(line, "Device Error Count:") ||
			strings.Contains(line, "FAILING_NOW") ||
			strings.Contains(line, "NOW") ||
			strings.Contains(line, "SMART overall-health self-assessment test result: FAILED!") ||
			(strings.Contains(line, "Self-test execution status:") && !strings.Contains(line, "(   0)")) || // Means in case the error code is not 0 (no error)
			strings.Contains(line, "Completed: ") { // Means in case the test completed with errors
			output[i] = "<strong style='color:#FF0000;'>" + line + "</strong>"
		}
	}

	return "<pre>" + strings.Join(output, "\n") + "</pre>"
}

/*
getAllAvailablePartitions gets all the available partitions on the device.

-----------------------------------------------------------

– Returns:
  - a list with all the available partitions on the device or nil if smartmontools is not installed
*/
func getAllAvailablePartitions() []string {
	var partitions_list []string
	stdOutErrCmd, _ := Utils.ExecCmdSHELL([]string{"smartctl{{EXE}} --scan"})

	var output_lines []string = strings.Split(stdOutErrCmd.Stdout_str, "\n")
	output_lines = output_lines[:len(output_lines)-1]
	for _, line := range output_lines {
		// This counts as the disk, since each disk has only the NTFS partition (the one that matters, anyways).
		partitions_list = append(partitions_list, strings.Split(line, " ")[0])
	}

	//log.Println("Partitions: " + strconv.Itoa(len(partitions_list)) + " partition(s).")
	//log.Printf("Partitions: %+v\n", partitions_list)

	return partitions_list
}

/*
getActiveDisks gets the disks that are spinning (means on either Active or Idle state) through smartctl.

-----------------------------------------------------------

– Params:
  - partitions – a list of the partitions of the available disks

– Returns:
  - a list of the partitions associated with the disks that are spinning or nil if none are spinning or if smartmontools
	is not installed
*/
func getActiveDisks(partitions_list []string) []string {
	var which_disks_active []string = nil
	for _, partition := range partitions_list {
		stdOutErrCmd, err := Utils.ExecCmdSHELL([]string{"smartctl{{EXE}} -n standby " + partition})
		if nil != err {
			continue
		}
		var output []string = strings.Split(stdOutErrCmd.Stdout_str, "\n")
		if strings.Contains(output[3], "ACTIVE or IDLE") {
			which_disks_active = append(which_disks_active, partition)
		}
	}

	return which_disks_active
}

/*
getDiskSerialPartitions gets the serial number of the disks associated with the partitions.

-----------------------------------------------------------

– Returns:
  - a map of disks serial numbers to their partition names or nil if smartmontools is not installed
*/
func getDiskSerialPartitions(partitions_list []string) map[string]string {
	var disk_serial_partitions map[string]string = make(map[string]string)
	for _, partition := range partitions_list {
		//log.Println("Partition: " + partition)
		stdOutErrCmd, err := Utils.ExecCmdSHELL([]string{"smartctl{{EXE}} -i " + partition})
		if nil != err {
			continue
		}

		var output []string = strings.Split(stdOutErrCmd.Stdout_str, "\n")
		for _, line := range output {
			if strings.HasPrefix(line, "Serial Number:") {
				var serial string = strings.Split(line, ":")[1]
				serial = strings.Trim(serial, " ")
				disk_serial_partitions[serial] = partition

				//log.Println("Serial: " + serial)

				break
			}
		}
	}

	return disk_serial_partitions
}

/*
initiateTest initiates a test on the disk associated with the specified partition, running in the background.

-----------------------------------------------------------

– Params:
  - long_test – true for long test, false for short test
  - partition – partition of the disk on which the test is to be started

– Returns:
  - the number of minutes the chosen test will take or -1 if it failed
*/
func initiateTest(long_test bool, partition string) int {
	// We need to try some times because sometimes it fails to start the test at first.
	var num_tries int = 3 // 3 tries because I think it's enough most of the times
	for i := 0; i < num_tries; i++ {
		var cmd []string
		if long_test {
			cmd = []string{"smartctl{{EXE}} -t long " + partition}
		} else {
			cmd = []string{"smartctl{{EXE}} -t short " + partition}
		}

		stdOutErrCmd, err := Utils.ExecCmdSHELL(cmd)
		if nil != err {
			continue
		}
		var stdout string = stdOutErrCmd.Stdout_str

		if strings.Contains(stdout, "Please wait ") {
			min, err := strconv.Atoi(strings.Split(stdout[strings.Index(stdout, "Please wait "):], " ")[2])
			if nil != err {
				continue
			}

			return min
		} else { // elif "Can't start self-test without aborting current test " in output: --> Doesn't matter, always try to abort.
			cmd = []string{"smartctl{{EXE}} -X " + partition}
			_, _ = Utils.ExecCmdSHELL(cmd)
		}
	}

	return -1
}

/*
checkDiskInTest checks if the disk associated with the specified partition is in test.

-----------------------------------------------------------

– Params:
  - partition – partition of the disk on which the test is to be checked

– Returns:
  - true if the disk is in test, false otherwise
*/
func checkDiskInTest(partition string) bool {
	stdOutErrCmd, err := Utils.ExecCmdSHELL([]string{"smartctl{{EXE}} -c " + partition})
	if nil != err {
		return false
	}
	var output []string = strings.Split(stdOutErrCmd.Stdout_str, "\n")
	for _, line := range output {
		if strings.Contains(line, "Self-test routine in progress") {
			return true
		}
	}

	return false
}
