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

const COMMUNICATOR_ABOUT string =
`This screen is about VISOR's communicator. You can use the communicator to send and receive messages to/from him.

You can send commands which will be processed locally or you can send anything else you want and that will be processed by the LLM (Large Language Model) on the server - if you set it up and is connected, or else VISOR will warn that the GPT model it's not available.

If you send a non-command sentence, it can take a bit to start receiving the response. Depends on how slow the model is on your computer (choose one that makes it so that speaking is not interrupted to wait for more sentences).

The memories are "infinite". You can add as many as you want and VISOR will use them all. The more you add, the more he knows about you and better the responses are. Though, the more you add, the more time he'll take to initialize.

Also currently you must explicitly tell him to memorize the conversation. The command to summarize the memories must be used with caution as the LLM may remove important parts of memories because it thinks they're not important. So, use it only when you're sure you want to summarize the memories and back them up first.`

const RSS_ABOUT string =
`The feeds are checked every 2 minutes. When you add a new one, you'll receive no immediate notifications, but you'll receive the next updates.

Currently I've only tested this with YouTube channels and playlists, and Stack Exchange feeds. It should work for other sites though. Just add and see if you get updates. If you don't, tell me so I can add support for it.`

const TASKS_ABOUT string =
`This is where you can configure VISOR to do tasks for you when you need him to.

Notes about the location trigger:
- You can input for example home_* and that means trigger when you're at any location that starts with "home_"
- You can also input for example +university which means trigger when you *arrive* at the university (or -university to trigger when you *leave* the university)
- And of course you can mix both cases (e.g. +home_*, which means trigger when you arrive at any location started with "home_")

About the programmable condition:
- It's written in Go and each time you write a boolean expression, you have to write it inside "bool(expression here)". Example: bool(sound_volume < 80). If you don't do this, it won't work (limitation of the library I used to process the conditions).
- The variables you can use are all the ones in the Registry, but written in lowercase (example: SOUND_VOLUME is written as sound_volume). Go there and grab any you want, including the manual ones if you need them.`

const LOCATOR_ABOUT string =
`Here you can add locations to be used for Tasks. Currently only Wi-Fi and Bluetooth location detectors are supported. Later GPS should be supported too.`

const ONLINE_INFO_CHK_ABOUT string =
`This module still needs development, but currently it checks the weather and news from where you tell it to.

You can then request the weather and the news in the Communicator screen, and VISOR will (currently) list *all* the locations, even if you want only one of them, so have that in mind.`

const SYS_CHK_ABOUT string =
`This module keeps checking the state of the system so that other modules can use this information (for example for the Tasks).`

const SMART_ABOUT string =
`This module starts a full S.M.A.R.T. scan of the disks you list here (NOTE: this is not for the client's disks - it's for the server ones), every first day of the month and sends an email with the complete report after the scan is complete.`

const REGISTRY_ABOUT string =
`The Registry is where you can see all the "public" variables VISOR is keeping track of (there's internal ones scattered across the program, but these are the important ones accessible to the user).

There are the automatic ones set by VISOR, and the manual ones you can set yourself through the various settings tabs of the app.`