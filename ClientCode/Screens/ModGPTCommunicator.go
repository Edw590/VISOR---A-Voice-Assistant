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

package Screens

import (
	"GPTComm/GPTComm"
	"Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"Utils/ModsFileInfo"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sort"
	"strconv"
	"strings"
	"time"
)

type _SessionInfo struct {
	id string
	session ModsFileInfo.Session
}

func ModGPTCommunicator() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_GPT_COMM

	return container.NewAppTabs(
		container.NewTabItem("Main", gptCommunicatorCreateMainTab()),
		container.NewTabItem("Chats", gptCommunicatorCreateSessionsTab()),
		container.NewTabItem("List of commands", gptCommunicatorCreateListCommandsTab()),
		container.NewTabItem("Memories", gptCommunicatorCreateMemoriesTab()),
		container.NewTabItem("Settings", gptCommunicatorCreateSettingsTab()),
		container.NewTabItem("About", gptCommunicatorCreateAboutTab()),
	)
}

func gptCommunicatorCreateAboutTab() *container.Scroll {
	var label_info *widget.Label = widget.NewLabel(COMMUNICATOR_ABOUT)
	label_info.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(label_info)
}

func gptCommunicatorCreateMemoriesTab() *container.Scroll {
	var memories string = "[Not connected to the server to get the memories]"
	var num_memories int = 0
	if Utils.IsCommunicatorConnectedSERVER() {
		 memories = GPTComm.GetMemories()
		if memories != "" {
			num_memories = strings.Count(memories, "\n") + 1
		}
	}

	var label_info *widget.Label = widget.NewLabel("List of memories stored, one per line (maximize the window). " +
		"Number of memories: " + strconv.Itoa(num_memories) + ".")
	label_info.Wrapping = fyne.TextWrapWord

	var memories_text *widget.Entry = widget.NewMultiLineEntry()
	memories_text.SetPlaceHolder("Stored memories")
	memories_text.Wrapping = fyne.TextWrapWord
	memories_text.SetMinRowsVisible(100)
	memories_text.SetText(memories)

	var btn_save *widget.Button = widget.NewButton("Save memories", func() {
		GPTComm.SetMemories(memories_text.Text)
	})
	btn_save.Importance = widget.SuccessImportance

	return createMainContentScrollUTILS(
		label_info,
		btn_save,
		memories_text,
	)
}

func gptCommunicatorCreateListCommandsTab() *container.Scroll {
	var label_example_command *widget.Label = widget.NewLabel("Example of a complex command VISOR understands " +
		"(always without punctuation - must not be present): \"turn it on. turn on the wifi, and... and the airplane " +
		"mode, get it it on. no, don't turn it on. turn off airplane mode and also the wifi, please.\"")
	label_example_command.Wrapping = fyne.TextWrapWord

	var label_list_commands *widget.Label = widget.NewLabel(
		"List of all commands and variations available (optional words in [...] and generic descriptions in (...):\n\n" +
		"(Note: there is more than one way to say a command, with synonyms and random words in between (\"switch on " +
			"the phone's wifi\", \"what's the current time\", \"terminate the phone call\").)\n\n" +
		"--> (Ask for the time)\n" +
		"--> (Ask for the date)\n" +
		"--> Turn on/off Wi-Fi\n" +
		"--> (Ask for the battery percentage/status/level(s))\n" +
		"--> (Ask for the weather)\n" +
		"--> (Ask for the news)\n" +
		"--> Turn on/off Ethernet\n" +
		"--> Turn on/off networking\n" +
		"--> (Ask what you have for today/tomorrow/this week/next week - Google Calendar events and tasks (tasks " +
			"are only retrieved if you ask for today or tomorrow). After asking and if you asked with the server " +
			"connected (meaning will be the LLM answering), you can then talk normally about your events and tasks " +
			"with VISOR)\n",
	)
	label_list_commands.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(
		label_example_command,
		label_list_commands,
	)
}

func gptCommunicatorCreateSettingsTab() *container.Scroll {
	var server_uri *widget.Entry = widget.NewEntry()
	server_uri.SetPlaceHolder("GPT Server uri (example localhost:11434)")
	server_uri.SetText(Utils.GetUserSettings().GPTCommunicator.Server_uri)

	var entry_model_name *widget.Entry = widget.NewEntry()
	entry_model_name.SetPlaceHolder("GPT model name (example: llama3.2)")
	entry_model_name.SetText(Utils.GetUserSettings().GPTCommunicator.Model_name)

	//var checkbox_model_has_tool_role *widget.Check = widget.NewCheck("Is the tool role available for the model?", nil)

	var entry_ctx_size *widget.Entry = widget.NewEntry()
	entry_ctx_size.SetPlaceHolder("GPT context size (example: 4096)")
	entry_ctx_size.SetText(strconv.Itoa(int(Utils.GetUserSettings().GPTCommunicator.Context_size)))
	entry_ctx_size.Validator = validation.NewRegexp(`^(\d+)?$`, "Context size must be numberic")

	var entry_temperature *widget.Entry = widget.NewEntry()
	entry_temperature.SetPlaceHolder("GPT temperature (example: 0.8)")
	entry_temperature.SetText(strconv.FormatFloat(float64(Utils.GetUserSettings().GPTCommunicator.Temperature), 'f', -1, 32))
	entry_temperature.Validator = func(s string) error {
		value, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return errors.New("temperature must be a decimal number")
		}
		if value < 0 || value > 1 {
			return errors.New("temperature must be between 0 and 1")
		}

		return nil
	}

	var entry_system_info *widget.Entry = widget.NewMultiLineEntry()
	entry_system_info.SetPlaceHolder("LLM system information (remove any current date/time - that's automatic)")
	entry_system_info.SetMinRowsVisible(3)
	entry_system_info.SetText(Utils.GetUserSettings().GPTCommunicator.System_info)

	var entry_user_nickname *widget.Entry = widget.NewEntry()
	entry_user_nickname.SetPlaceHolder("User nickname (Sir, for example)")
	entry_user_nickname.SetText(Utils.GetUserSettings().GPTCommunicator.User_nickname)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.GetUserSettings().GPTCommunicator.Server_uri = server_uri.Text
		Utils.GetUserSettings().GPTCommunicator.Model_name = entry_model_name.Text
		//Utils.GetUserSettings().GPTCommunicator.Model_has_tool_role = checkbox_model_has_tool_role.Checked
		value1, _ := strconv.ParseInt(entry_ctx_size.Text, 10, 32)
		Utils.GetUserSettings().GPTCommunicator.Context_size = int32(value1)
		value2, _ := strconv.ParseFloat(entry_temperature.Text, 32)
		Utils.GetUserSettings().GPTCommunicator.Temperature = float32(value2)
		Utils.GetUserSettings().GPTCommunicator.System_info = entry_system_info.Text
		Utils.GetUserSettings().GPTCommunicator.User_nickname = entry_user_nickname.Text
	})
	btn_save.Importance = widget.SuccessImportance

	return createMainContentScrollUTILS(
		server_uri,
		entry_model_name,
		//checkbox_model_has_tool_role,
		entry_ctx_size,
		entry_temperature,
		entry_system_info,
		entry_user_nickname,
		btn_save,
	)
}

func gptCommunicatorCreateSessionsTab() *container.Scroll {
	if !Utils.IsCommunicatorConnectedSERVER() {
		return createMainContentScrollUTILS(widget.NewLabel("[Not connected to the server to get the chats]"))
	}

	var session_ids_str string = GPTComm.GetSessionIdsList()
	if session_ids_str == "" {
		return createMainContentScrollUTILS()
	}

	var sessions_info []_SessionInfo = nil
	var session_ids []string = strings.Split(session_ids_str, "|")
	for _, session_id := range session_ids {
		var session_info _SessionInfo
		session_info.id = session_id
		session_info.session.Name = GPTComm.GetSessionName(session_id)
		session_info.session.Created_time_s = GPTComm.GetSessionCreatedTime(session_id)
		sessions_info = append(sessions_info, session_info)
	}

	sort.SliceStable(sessions_info, func(i, j int) bool {
		return sessions_info[i].session.Created_time_s > sessions_info[j].session.Created_time_s
	})

	var entries_map map[string]*widget.Entry = make(map[string]*widget.Entry)
	var accordion *widget.Accordion = widget.NewAccordion()
	for _, session_info := range sessions_info {
		if session_info.id == "temp" || session_info.id == "dumb" {
			continue
		}

		var title string = session_info.session.Name + " - " +
			Utils.GetDateTimeStrTIMEDATE(GPTComm.GetSessionCreatedTime(session_info.id) * 1000)

		accordion.Append(widget.NewAccordionItem(trimAccordionTitleUTILS(title),
			createSessionView(entries_map, session_info)))
	}

	go func() {
		for {
			if Current_screen_GL == ID_MOD_GPT_COMM {
				if Utils.IsCommunicatorConnectedSERVER() {
					session_ids_str = GPTComm.GetSessionIdsList()
					if session_ids_str != "" {
						session_ids = strings.Split(session_ids_str, "|")
						for _, session_id := range session_ids {
							if session_id == "temp" || session_id == "dumb" {
								continue
							}

							var session_history_str string = GPTComm.GetSessionHistory(session_id)
							if session_history_str == "" {
								continue
							}

							var session_history []string = strings.Split(session_history_str, "\000")
							var msg_content_str string = ""
							for _, message := range session_history {
								var message_parts_pipe []string = strings.Split(message, "|")
								var index_first_pipe int = strings.Index(message, "|")
								var message_parts_slash []string = strings.Split(message_parts_pipe[0], "/")

								var msg_role = message_parts_slash[0]
								switch msg_role {
									case "system":
										continue
									case "assistant":
										msg_role = "VISOR"
									case "user":
										msg_role = "YOU"
								}

								if len(message_parts_pipe) < 2 || message_parts_pipe[1] == "" {
									// Means no message (so maybe was a "SYSTEM TASK" message - ignore those)
									continue
								}

								var msg_timestamp_s, _ = strconv.ParseInt(message_parts_slash[1], 10, 64)
								var msg_content = message[index_first_pipe+1:]

								msg_content_str +=
									"-----------------------------------------------------------------------\n" +
										"|" + strings.ToUpper(msg_role) + "| on " +
										Utils.GetDateTimeStrTIMEDATE(msg_timestamp_s*1000) + ":\n" + msg_content + "\n\n"
							}
							if len(msg_content_str) > 2 {
								msg_content_str = msg_content_str[:len(msg_content_str)-2]
							}

							if entries_map[session_id] != nil && entries_map[session_id].Text != msg_content_str {
								entries_map[session_id].SetText(msg_content_str)
							}
						}
					}
				}
			} else {
				break
			}

			time.Sleep(5 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(accordion)
}

func createSessionView(entries_map map[string]*widget.Entry, session_info _SessionInfo) *fyne.Container {
	var label_date *widget.Label = widget.NewLabel("Created on " +
		Utils.GetDateTimeStrTIMEDATE(GPTComm.GetSessionCreatedTime(session_info.id) * 1000))

	var entry_name *widget.Entry = widget.NewEntry()
	entry_name.SetPlaceHolder("Chat name")
	entry_name.SetText(session_info.session.Name)

	var btn_save *widget.Button = widget.NewButton("Save name", func() {
		if entry_name.Text == "" {
			dialog.ShowError(errors.New("the chat name must not be empty"), Current_window_GL)

			return
		}

		GPTComm.SetSessionName(session_info.id, entry_name.Text)

		reloadScreen()
	})
	btn_save.Importance = widget.SuccessImportance

	var btn_delete *widget.Button = widget.NewButton("Delete chat", func() {
		createConfirmationDialogUTILS("Are you sure you want to delete this chat?", func(confirmed bool) {
			if confirmed {
				GPTComm.DeleteSession(session_info.id)

				reloadScreen()
			}
		})
	})
	btn_delete.Importance = widget.DangerImportance

	var entry_history *widget.Entry = widget.NewMultiLineEntry()
	entry_history.Wrapping = fyne.TextWrapWord
	entry_history.SetMinRowsVisible(45)
	entries_map[session_info.id] = entry_history

	return container.NewVBox(
		label_date,
		entry_name,
		container.New(layout.NewGridLayout(2), btn_save, btn_delete),
		entry_history,
	)
}

func gptCommunicatorCreateMainTab() *container.Scroll {
	var label_gpt_comm_state *widget.Label = widget.NewLabel("GPT state: error")

	var text_to_send *widget.Entry = widget.NewMultiLineEntry()
	text_to_send.Wrapping = fyne.TextWrapWord
	text_to_send.SetMinRowsVisible(6) // 6 lines, like ChatGPT has
	text_to_send.SetPlaceHolder("Text to send to VISOR (without punctuation for command detection)\n" +
		"- /stop to stop the LLM while it's generating text")

	var btn_send_text *widget.Button = widget.NewButton("Send text", func() {
		Utils.SendToModChannel(Utils.NUM_MOD_CmdsExecutor, 0, "Sentence", text_to_send.Text)
	})

	var btn_send_text_gpt_smart *widget.Button = widget.NewButton("Send text directly to the LLM (new chat)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Not connected to the server."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		var speak string = ""
		switch GPTComm.SendText(text_to_send.Text, GPTComm.SESSION_TYPE_NEW, GPTComm.ROLE_USER, false) {
			case ModsFileInfo.MOD_7_STATE_STOPPED:
				speak = "The GPT is stopped. Text on hold."
			case ModsFileInfo.MOD_7_STATE_STARTING:
				speak = "The GPT is starting up. Text on hold."
			case ModsFileInfo.MOD_7_STATE_BUSY:
				speak = "The GPT is busy. Text on hold."
		}
		if speak != "" {
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})

	var btn_send_text_gpt_dumb *widget.Button = widget.NewButton("Send text directly to the LLM (temp chat)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Not connected to the server."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		var speak string = ""
		switch GPTComm.SendText(text_to_send.Text, GPTComm.SESSION_TYPE_TEMP, GPTComm.ROLE_USER, false) {
			case ModsFileInfo.MOD_7_STATE_STOPPED:
				speak = "The GPT is stopped. Text on hold."
			case ModsFileInfo.MOD_7_STATE_STARTING:
				speak = "The GPT is starting up. Text on hold."
			case ModsFileInfo.MOD_7_STATE_BUSY:
				speak = "The GPT is busy. Text on hold."
		}
		if speak != "" {
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})

	var response_text *widget.Entry = widget.NewMultiLineEntry()
	response_text.SetPlaceHolder("Response")
	response_text.Wrapping = fyne.TextWrapWord
	response_text.SetMinRowsVisible(100)

	go func() {
		var old_text string = ""
		for {
			if Current_screen_GL == ID_MOD_GPT_COMM {
				var new_text string = GPTComm.GetLastText()
				if new_text != old_text {
					old_text = new_text
					response_text.SetText(new_text)
				}

				var gpt_state string = "[Not connected to the server to get the GPT state]"
				if Utils.IsCommunicatorConnectedSERVER() {
					switch GPTComm.GetModuleState() {
						case ModsFileInfo.MOD_7_STATE_STOPPED:
							gpt_state = "stopped"
						case ModsFileInfo.MOD_7_STATE_STARTING:
							gpt_state = "starting"
						case ModsFileInfo.MOD_7_STATE_READY:
							gpt_state = "ready"
						case ModsFileInfo.MOD_7_STATE_BUSY:
							gpt_state = "busy"
						default:
							gpt_state = "invalid"
					}
				}
				label_gpt_comm_state.SetText("GPT state: " + gpt_state)
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(
		label_gpt_comm_state,
		text_to_send,
		btn_send_text,
		btn_send_text_gpt_smart,
		btn_send_text_gpt_dumb,
		response_text,
	)
}
