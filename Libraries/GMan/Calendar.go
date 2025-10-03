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

package GMan

import (
	"Utils"
	"Utils/ModsFileInfo"
	"strconv"
)

var calendars_GL map[string]ModsFileInfo.GCalendar = nil
var events_GL []ModsFileInfo.GEvent = nil

func getEvents() {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GMan, 0, []byte("G_S|true|GManEvents")) {
		return
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GMan, 0, 10)
	if comms_map == nil {
		return
	}

	var json_bytes []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	if err := Utils.FromJsonGENERAL(json_bytes, &events_GL); err != nil {
		return
	}
}

func getCalendars() {
	if !Utils.QueueMessageSERVER(false, Utils.NUM_LIB_GMan, 0, []byte("G_S|true|GManCals")) {
		return
	}
	var comms_map map[string]any = Utils.GetFromCommsChannel(false, Utils.NUM_LIB_GMan, 0, 10)
	if comms_map == nil {
		return
	}

	var json_bytes []byte = comms_map[Utils.COMMS_MAP_SRV_KEY].([]byte)

	if err := Utils.FromJsonGENERAL(json_bytes, &calendars_GL); err != nil {
		return
	}
}

/*
GetCalendarsIdsList updates and gets a list of all calendar IDs.

-----------------------------------------------------------

– Returns
  - a list of all calendar IDs separated by "|"
 */
func GetCalendarsIdsList() string {
	getCalendars()

	var ids_list string = ""
	for id := range calendars_GL {
		ids_list += id + "|"
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetCalendar returns a calendar by its ID.

-----------------------------------------------------------

– Params:
  - calendar_id – the calendar ID

– Returns:
  - the calendar or nil if the calendar was not found
*/
func GetCalendar(calendar_id string) *ModsFileInfo.GCalendar {
	for id, calendar := range calendars_GL {
		if id == calendar_id {
			return &calendar
		}
	}

	return nil
}

/*
SetCalendarEnabled enables or disables a calendar by its ID.

-----------------------------------------------------------

– Params:
  - calendar_id – the calendar ID
  - enabled – true to enable the calendar, false to disable it

– Returns:
  - true if the operation was successful, false otherwise
 */
func SetCalendarEnabled(calendar_id string, enabled bool) bool {
	var message []byte = []byte("GMan|calendar|" + calendar_id + "|" + strconv.FormatBool(enabled))

	return Utils.QueueNoResponseMessageSERVER(message)
}

/*
GetEventsIdsList updates and gets a list of all events' IDs.

Note: it also updates the calendars list.

-----------------------------------------------------------

– Params:
  - only_active_cals – if true, only events from enabled calendars will be returned

– Returns:
  - a list of all events' IDs separated by "|"
*/
func GetEventsIdsList(only_active_cals bool) string {
	getEvents()
	getCalendars()

	var ids_list string = ""
	for _, event := range events_GL {
		if only_active_cals {
			var calendar *ModsFileInfo.GCalendar = GetCalendar(event.Calendar_id)
			var add_event bool = false
			if calendar == nil {
				// Shouldn't happen, but whatever
				add_event = true
			} else if calendar.Enabled {
				add_event = true
			}

			if add_event {
				ids_list += event.Id + "|"
			}
		}
	}
	if len(ids_list) > 0 {
		ids_list = ids_list[:len(ids_list)-1]
	}

	return ids_list
}

/*
GetEvent returns an event by its ID.

-----------------------------------------------------------

– Params:
  - event_id – the event ID

– Returns:
  - the event or nil if the event was not found
*/
func GetEvent(event_id string) *ModsFileInfo.GEvent {
	for i := range events_GL {
		var event *ModsFileInfo.GEvent = &events_GL[i]
		if event.Id == event_id {
			return event
		}
	}

	return nil
}

/*
AddEvent adds a new event to Google Calendar.

-----------------------------------------------------------

– Params:
  - event – the event to add

– Returns:
  - true if the event was added successfully, false otherwise
 */
func AddEvent(event *ModsFileInfo.GEvent) bool {
	var message []byte = []byte("GMan|event|")
	message = append(message, *Utils.ToJsonGENERAL(*event)...)

	return Utils.QueueNoResponseMessageSERVER(message)
}
