/*******************************************************************************
 * Copyright 2023-2024 Edw590
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

package MOD_13

///////////////////////////////////////////////////////////////////
// Commands list

//const CMD_TOGGLE_FLASHLIGHT string = "1";
const CMD_ASK_TIME string = "2";
const CMD_ASK_DATE string = "3";
//const CMD_TOGGLE_WIFI string = "4";
//const CMD_TOGGLE_MOBILE_DATA string = "5";
//const CMD_TOGGLE_BLUETOOTH string = "6";
//const CMD_ANSWER_CALL string = "7";
//const CMD_END_CALL string = "9";
//const CMD_TOGGLE_SPEAKERS  string = "10";
//const CMD_TOGGLE_AIRPLANE_MODE  string = "11";
const CMD_ASK_BATTERY_PERCENT  string = "12";
//const CMD_POWER_SHUT_DOWN  string = "13";
//const CMD_POWER_REBOOT  string = "14";
//const CMD_TAKE_PHOTO  string = "15";
//const CMD_RECORD_MEDIA  string = "16";
//const CMD_SAY_AGAIN  string = "17";
//const CMD_CALL_CONTACT  string = "18";
//const CMD_TOGGLE_POWER_SAVER_MODE  string = "19";
//const CMD_STOP_RECORD_MEDIA  string = "20";
//const CMD_CONTROL_MEDIA  string = "21";
//const CMD_CONFIRM  string = "22";
//const CMD_REJECT  string = "23";
//const CMD_STOP_LISTENING  string = "24";
//const CMD_START_LISTENING  string = "25";
//const CMD_TELL_WEATHER  string = "26";
//const CMD_TELL_NEWS  string = "27";
//const CMD_GONNA_SLEEP  string = "28";

///////////////////////////////////////////////////////////////////
// Return IDs

const RET_ON string = ".00001"
const RET_OFF string = ".00002"

const RET_14_FAST string = ".00001";
const RET_14_NORMAL string = ".00002";
const RET_14_RECOVERY string = ".00003";
const RET_14_SAFE_MODE string = ".00004";
const RET_14_BOOTLOADER string = ".00005";

const RET_15_REAR string = ".00001";
const RET_15_FRONTAL string = ".00002";

const RET_16_AUDIO_1 string = ".00001";
const RET_16_AUDIO_2 string = ".00003";
const RET_16_VIDEO_1 string = ".00002";
const RET_16_VIDEO_2 string = ".00004";

const RET_20_ANY string = ".00001";
const RET_20_AUDIO string = ".00002";
const RET_20_VIDEO string = ".00003";

const RET_21_PLAY string = ".00001";
const RET_21_PAUSE string = ".00002";
const RET_21_STOP string = ".00003";
const RET_21_NEXT string = ".00004";
const RET_21_PREVIOUS string = ".00005";

///////////////////////////////////////////////////////////////////
// Additional command info

// CMDi_INF1_DO_SOMETHING signals that the referring command requires the assistant to do something.
const CMDi_INF1_DO_SOMETHING = "0"
// CMDi_INF1_ONLY_SPEAK signals that the referring command only requires the assistant to say something (like asking
// what time is it).
const CMDi_INF1_ONLY_SPEAK = "1"
// CMDi_INF1_ASSIST_CMD signals that the referring command is an assistance to another command (like saying "I confirm"
// (the previous command)).
//const CMDi_INF1_ASSIST_CMD = ""

var cmdi_info map[string]string = map[string]string{
	//CMD_TOGGLE_FLASHLIGHT:         CMDi_INF1_DO_SOMETHING,     // 1
	CMD_ASK_TIME:                  CMDi_INF1_ONLY_SPEAK,       // 2
	CMD_ASK_DATE:                  CMDi_INF1_ONLY_SPEAK,       // 3
	//CMD_TOGGLE_WIFI:               CMDi_INF1_DO_SOMETHING,     // 4
	//CMD_TOGGLE_MOBILE_DATA:        CMDi_INF1_DO_SOMETHING,     // 5
	//CMD_TOGGLE_BLUETOOTH:          CMDi_INF1_DO_SOMETHING,     // 6
	//CMD_ANSWER_CALL:               CMDi_INF1_DO_SOMETHING,     // 7
	//CMD_END_CALL:                  CMDi_INF1_DO_SOMETHING,     // 9
	//CMD_TOGGLE_SPEAKERS:           CMDi_INF1_DO_SOMETHING,     // 10
	//CMD_TOGGLE_AIRPLANE_MODE:      CMDi_INF1_DO_SOMETHING,     // 11
	CMD_ASK_BATTERY_PERCENT:       CMDi_INF1_ONLY_SPEAK,       // 12
	//CMD_POWER_SHUT_DOWN:           CMDi_INF1_DO_SOMETHING,     // 13
	//CMD_POWER_REBOOT:              CMDi_INF1_DO_SOMETHING,     // 14
	//CMD_TAKE_PHOTO:                CMDi_INF1_DO_SOMETHING,     // 15
	//CMD_RECORD_MEDIA:              CMDi_INF1_DO_SOMETHING,     // 16
	//CMD_SAY_AGAIN:                 CMDi_INF1_ONLY_SPEAK,       // 17
	//CMD_CALL_CONTACT:              CMDi_INF1_DO_SOMETHING,     // 18
	//CMD_TOGGLE_POWER_SAVER_MODE:   CMDi_INF1_DO_SOMETHING,     // 19
	//CMD_STOP_RECORD_MEDIA:         CMDi_INF1_DO_SOMETHING,     // 20
	//CMD_CONTROL_MEDIA:             CMDi_INF1_DO_SOMETHING,     // 21
	//CMD_CONFIRM:                   CMDi_INF1_ASSIST_CMD,       // 22
	//CMD_REJECT:                    CMDi_INF1_ASSIST_CMD,       // 23
	//CMD_STOP_LISTENING:            CMDi_INF1_DO_SOMETHING,     // 24
	//CMD_START_LISTENING:           CMDi_INF1_DO_SOMETHING,     // 25
	//CMD_TELL_WEATHER:              CMDi_INF1_ONLY_SPEAK,       // 26
	//CMD_TELL_NEWS:                 CMDi_INF1_ONLY_SPEAK,       // 27
	//CMD_GONNA_SLEEP:               CMDi_INF1_ONLY_SPEAK,       // 28
}
