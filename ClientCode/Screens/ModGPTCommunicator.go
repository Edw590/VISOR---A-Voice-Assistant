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

package Screens

import (
	"GPTComm/GPTComm"
	"Speech"
	"SpeechQueue/SpeechQueue"
	"Utils"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"time"
)

func ModGPTCommunicator() fyne.CanvasObject {
	Current_screen_GL = ID_MOD_GPT_COMM

	return container.NewAppTabs(
		container.NewTabItem("Main", gptCommunicatorCreateCommunicatorTab()),
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
	var label_info *widget.Label = widget.NewLabel("List of memories stored for the smart LLM, one per line " +
		"(maximize the window):")

	var memories_text *widget.Entry = widget.NewMultiLineEntry()
	memories_text.SetPlaceHolder("Stored memories on the smart LLM")
	memories_text.Wrapping = fyne.TextWrapWord
	memories_text.SetMinRowsVisible(100)
	memories_text.SetText(GPTComm.GetMemories())

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
		"--> (Ask what you have for today/tomorrow/this week/next week - Google Calendar events. After asking and if " +
			"you asked with the server connected (meaning will be the LLM answering), you can then talk normally " +
			"about your tasks with VISOR)\n",
	)
	label_list_commands.Wrapping = fyne.TextWrapWord

	return createMainContentScrollUTILS(
		label_example_command,
		label_list_commands,
	)
}

func gptCommunicatorCreateSettingsTab() *container.Scroll {
	var label_supported_models *widget.Label = widget.NewLabel("Supported models: Llama3.1 and 3.2 only")

	var entry_smart_model_loc *widget.Entry = widget.NewEntry()
	entry_smart_model_loc.Validator = validation.NewRegexp(`^.*\.gguf$`, "The model location must end with .gguf")
	entry_smart_model_loc.SetPlaceHolder("GGUF location for the smart LLM (used for normal conversation)")
	entry_smart_model_loc.SetText(Utils.User_settings_GL.GPTCommunicator.Model_smart_loc)

	var entry_dumb_model_loc *widget.Entry = widget.NewEntry()
	entry_dumb_model_loc.Validator = validation.NewRegexp(`^.*\.gguf$`, "The model location must end with .gguf")
	entry_dumb_model_loc.SetPlaceHolder("GGUF location for the dumb LLM (used to summarize things)")
	entry_dumb_model_loc.SetText(Utils.User_settings_GL.GPTCommunicator.Model_dumb_loc)

	var entry_system_info *widget.Entry = widget.NewEntry()
	entry_system_info.SetPlaceHolder("LLM system information")
	entry_system_info.SetText(Utils.User_settings_GL.GPTCommunicator.System_info)

	var entry_user_nickname *widget.Entry = widget.NewEntry()
	entry_user_nickname.SetPlaceHolder("User nickname")
	entry_user_nickname.SetText(Utils.User_settings_GL.GPTCommunicator.User_nickname)

	var btn_save *widget.Button = widget.NewButton("Save", func() {
		Utils.User_settings_GL.GPTCommunicator.Model_smart_loc = entry_smart_model_loc.Text
		Utils.User_settings_GL.GPTCommunicator.Model_dumb_loc = entry_dumb_model_loc.Text
		Utils.User_settings_GL.GPTCommunicator.System_info = entry_system_info.Text
		Utils.User_settings_GL.GPTCommunicator.User_nickname = entry_user_nickname.Text
	})
	btn_save.Importance = widget.SuccessImportance

	return createMainContentScrollUTILS(
		label_supported_models,
		entry_smart_model_loc,
		entry_dumb_model_loc,
		entry_system_info,
		entry_user_nickname,
		btn_save,
	)
}

func gptCommunicatorCreateCommunicatorTab() *container.Scroll {
	var text_to_send *widget.Entry = widget.NewMultiLineEntry()
	text_to_send.Wrapping = fyne.TextWrapWord
	text_to_send.SetMinRowsVisible(6) // 6 lines, like ChatGPT has
	text_to_send.SetPlaceHolder(
		"Text to send to VISOR (commands, or normal text to the smart LLM)\n" +
		"- /clear to clear the context by restarting the LLMs\n" +
		"- /stop to stop the LLM while it's generating text by restarting both\n" +
		"- /mem to memorize the conversation\n" +
		"- /memmem to summarize the list of memories (use with caution)\n",
	)

	var btn_send_text *widget.Button = widget.NewButton("Send text", func() {
		Utils.SendToModChannel(Utils.NUM_MOD_CmdsExecutor, "Sentence", text_to_send.Text)
	})

	var btn_send_text_gpt_smart *widget.Button = widget.NewButton("Send text directly to the LLM (smart)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Not connected to the server."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		if !GPTComm.SendText(text_to_send.Text, true) {
			Speech.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})

	var btn_send_text_gpt_dumb *widget.Button = widget.NewButton("Send text directly to the LLM (dumb)", func() {
		if !Utils.IsCommunicatorConnectedSERVER() {
			var speak string = "GPT unavailable. Not connected to the server."
			Speech.QueueSpeech(speak, SpeechQueue.PRIORITY_USER_ACTION, SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)

			return
		}

		if !GPTComm.SendText(text_to_send.Text, false) {
			Speech.QueueSpeech("Sorry, the GPT is busy at the moment. Text on hold.", SpeechQueue.PRIORITY_USER_ACTION,
				SpeechQueue.MODE1_ALWAYS_NOTIFY, "", 0)
		}
	})

	var response_text *widget.Entry = widget.NewMultiLineEntry()
	response_text.SetPlaceHolder("Response from the smart LLM")
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
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(
		text_to_send,
		btn_send_text,
		btn_send_text_gpt_smart,
		btn_send_text_gpt_dumb,
		response_text,
	)
}
