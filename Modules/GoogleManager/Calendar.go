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

package GoogleManager

import (
	"Utils"
	"Utils/ModsFileInfo"
	"context"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func storeCalendarsEvents(client *http.Client) bool {
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		Utils.LogLnError("Unable to retrieve Calendar client:", err)

		return false
	}

	// Get the list of all calendars
	calendarList, err := service.CalendarList.List().Do()
	if err != nil {
		Utils.LogLnError("Unable to retrieve Calendar list:", err)

		return false
	}

	// Store the calendars
	if getModGenSettings().Calendars == nil {
		getModGenSettings().Calendars = make(map[string]ModsFileInfo.GCalendar)
	}
	for _, calendarListEntry := range calendarList.Items {
		var summary string = calendarListEntry.SummaryOverride
		if summary == "" {
			summary = calendarListEntry.Summary
		}
		var cal_enabled bool = true
		if cal, ok := getModGenSettings().Calendars[calendarListEntry.Id]; ok {
			// If the calendar is already on the list, use its enabled state (else, default to true)
			cal_enabled = cal.Enabled
		}
		getModGenSettings().Calendars[calendarListEntry.Id] = ModsFileInfo.GCalendar{
			Title:   summary, // Don't check if the calendar is new or not - this field must be always updated
			Enabled: cal_enabled,
		}
	}

	// Calculate the start of the current week (Monday)
	var now time.Time = time.Now()
	var weekday int = int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Make Sunday 7 instead of 0 for easier calculation
	}
	start_of_Week := now.AddDate(0, 0, -weekday+1) // Go back to Monday

	// Calculate the end of the next week (Sunday)
	end_of_next_week := start_of_Week.AddDate(0, 0, 13)

	// Set time range for events
	//log.Println("Getting all events for this week and next week")

	// Reset the events map every time we update the events
	var events_final []ModsFileInfo.GEvent = nil

	// Iterate over each calendar and retrieve events
	for _, calendarListEntry := range calendarList.Items {
		if !getModGenSettings().Calendars[calendarListEntry.Id].Enabled {
			continue
		}

		//log.Printf("Calendar: %s\n", calendarListEntry.Summary)

		events, err := service.Events.List(calendarListEntry.Id).
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(start_of_Week.Format(time.RFC3339)).
			TimeMax(end_of_next_week.Format(time.RFC3339)).
			MaxResults(9999).
			OrderBy("startTime").
			Do()
		if err != nil {
			//log.Printf("Unable to retrieve events for calendar %s: %v\n", calendarListEntry.Summary, err)

			continue
		}

		// Display the events
		if len(events.Items) == 0 {
			//log.Println("No upcoming events found.")
		} else {
			for _, item := range events.Items {
				var start_date string = item.Start.DateTime
				if start_date == "" {
					temp_time, _ := time.Parse("2006-01-02", item.Start.Date)
					start_date = temp_time.Format(time.RFC3339)
				}
				var end_date string = item.End.DateTime
				if end_date == "" {
					temp_time, _ := time.Parse("2006-01-02", item.End.Date)
					end_date = temp_time.Format(time.RFC3339)
				}
				//log.Printf("%s<->%s - %s\n", start_date, end_date, item.Summary)

				start_date_parsed, _ := time.Parse(time.RFC3339, start_date)
				end_date_parsed, _ := time.Parse(time.RFC3339, end_date)

				//log.Println(end_date_parsed.Sub(start_date_parsed))

				var duration_min int64 = int64(end_date_parsed.Sub(start_date_parsed).Minutes())

				// Store the event
				events_final = append(events_final, ModsFileInfo.GEvent{
					Id:           item.Id,
					Calendar_id:  calendarListEntry.Id,
					Summary:      item.Summary,
					Location:     item.Location,
					Description:  item.Description,
					Start_time_s: start_date_parsed.Unix(),
					Duration_min: duration_min,
				})
			}
		}

		getModGenSettings().Events = events_final
	}

	return true
}

func addEvent(event ModsFileInfo.GEvent, client *http.Client) bool {
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		Utils.LogLnError("Unable to retrieve Calendar client:", err)

		return false
	}

	var startTime time.Time = time.Unix(event.Start_time_s, 0)
	var endTime time.Time = startTime.Add(time.Duration(event.Duration_min) * time.Minute)

	var google_event *calendar.Event = &calendar.Event{
		Summary:     event.Summary,
		Location:    event.Location,
		Description: event.Description,
		Start: &calendar.EventDateTime{
			DateTime: startTime.Format(time.RFC3339),
			TimeZone: "UTC", // TODO: make timezone configurable
		},
		End: &calendar.EventDateTime{
			DateTime: endTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	// Insert the event into the primary calendar
	createdEvent, err := service.Events.Insert("primary", google_event).Do()
	if err != nil {
		Utils.LogLnError("Unable to create event:", err)

		return false
	}

	Utils.LogLnInfo("Event created: " + createdEvent.HtmlLink)

	return true
}
