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

package GoogleManager

import (
	"Utils/ModsFileInfo"
	"context"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

func storeCalendarsEvents(client *http.Client) bool {
	service, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Calendar client: %v\n", err)
		return false
	}

	// Get the list of all calendars
	calendarList, err := service.CalendarList.List().Do()
	if err != nil {
		log.Printf("Unable to retrieve calendar list: %v\n", err)
		return false
	}

	// Calculate the start of the current week (Monday)
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Make Sunday 7 instead of 0 for easier calculation
	}
	startOfWeek := now.AddDate(0, 0, -weekday+1) // Go back to Monday

	// Calculate the end of the next week (Sunday)
	endOfNextWeek := startOfWeek.AddDate(0, 0, 13)

	// Set time range for events
	log.Println("Getting all events for this week and next week")

	// Reset the events map every time we update the events
	modGenInfo_GL.Events = make(map[string][]*ModsFileInfo.Event)

	// Iterate over each calendar and retrieve events
	for _, calendarListEntry := range calendarList.Items {
		log.Printf("Calendar: %s\n", calendarListEntry.Summary)

		events, err := service.Events.List(calendarListEntry.Id).
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(startOfWeek.Format(time.RFC3339)).
			TimeMax(endOfNextWeek.Format(time.RFC3339)).
			MaxResults(999).
			OrderBy("startTime").
			Do()
		if err != nil {
			log.Printf("Unable to retrieve events for calendar %s: %v\n", calendarListEntry.Summary, err)

			continue
		}

		// Display the events
		if len(events.Items) == 0 {
			log.Println("No upcoming events found.")
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
				log.Printf("%s<->%s - %s\n", start_date, end_date, item.Summary)

				start_date_parsed, _ := time.Parse(time.RFC3339, start_date)
				end_date_parsed, _ := time.Parse(time.RFC3339, end_date)

				log.Println(end_date_parsed.Sub(start_date_parsed))

				var duration_min int64 = int64(end_date_parsed.Sub(start_date_parsed).Minutes())

				// Store the event
				modGenInfo_GL.Events[item.Id] = append(modGenInfo_GL.Events[calendarListEntry.Summary],
					&ModsFileInfo.Event{
						Summary:      item.Summary,
						Location:     item.Location,
						Description:  item.Description,
						Start_time:   start_date,
						Duration_min: duration_min,
					},
				)
			}
		}
	}

	return true
}
