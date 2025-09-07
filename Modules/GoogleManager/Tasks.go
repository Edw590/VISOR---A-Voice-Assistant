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

	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func storeTasks(client *http.Client) bool {
	// Create a new Tasks service.
	srv, err := tasks.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		Utils.LogfError("Unable to retrieve Tasks client: %v\n", err)

		return false
	}

	// Retrieve the user's task lists.
	task_lists, err := srv.Tasklists.List().MaxResults(9999).Do()
	if err != nil {
		Utils.LogfError("Unable to retrieve Task lists: %v\n", err)

		return false
	}

	setTokenValid()

	if len(task_lists.Items) == 0 {
		//log.Println("No task lists found.")

		return false
	}

	var tasks_final []ModsFileInfo.GTask = nil

	// Print task list names and retrieve tasks from the primary task list.
	for _, list := range task_lists.Items {
		//log.Printf("Tasks for list: %s\n", list.Title)
		tasks_list, err := srv.Tasks.List(list.Id).MaxResults(9999).Do()
		if err != nil {
			//log.Printf("Unable to retrieve tasks for list %s: %v", list.Title, err)

			continue
		}

		if len(tasks_list.Items) == 0 {
			//log.Println("No tasks found in this list.")
			continue
		}

		// Print task details.
		for _, task := range tasks_list.Items {
			//log.Printf("- %s (Status: %s)\n", task.Title, task.Status)
			//log.Printf("  Notes: %s\n", task.Notes)

			var task_date_time time.Time = time.Unix(0, 0)
			if task.Due != "" {
				task_date_time, _ = time.Parse(time.RFC3339, task.Due)
			}

			tasks_final = append(tasks_final, ModsFileInfo.GTask{
				Id:        task.Id,
				Title:     task.Title,
				Details:   task.Notes,
				Date_s:    task_date_time.Unix(),
				Completed: task.Status == "completed",
			})
		}
	}

	getModGenSettings().Tasks = tasks_final

	return true
}
