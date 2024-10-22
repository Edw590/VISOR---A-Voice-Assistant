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

package main

import (
	"ModulesManager"
	"SettingsSync/SettingsSync"
	"Utils"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"VISOR_Client/Logo"
	"VISOR_Client/Screens"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var my_app_GL fyne.App = nil
var my_window_GL fyne.Window = nil

var (
	realMain      Utils.RealMain = nil
	moduleInfo_GL Utils.ModuleInfo
)

func main() {
	var module Utils.Module = Utils.Module{
		Num:     Utils.NUM_MOD_VISOR,
		Name:    Utils.GetModNameMODULES(Utils.NUM_MOD_VISOR),
		Stop:    false,
		Stopped: false,
		Enabled: true,
	}
	Utils.ModStartup2(realMain, &module, false)
}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		moduleInfo_GL = moduleInfo_any.(Utils.ModuleInfo)

		//////////////////////////////////////////
		// Get the user settings

		go func() {
			for {
				Utils.StartCommunicatorSERVER()

				time.Sleep(1 * time.Second)
			}
		}()

		var user_settings_json string = ""
		var p_user_settings_json *string = Utils.GetBinDirFILESDIRS().Add2(true, Utils.USER_SETTINGS_FILE).ReadTextFile()
		if p_user_settings_json != nil {
			user_settings_json = *p_user_settings_json
		}
		if err := SettingsSync.LoadUserSettings(user_settings_json); err != nil {
			log.Println("Failed to load user settings. Attempting to retrieve them from the server...")
			log.Println("Please enter VISOR's website domain")
			var website_domain string
			_, _ = fmt.Scanln(&website_domain)

			log.Println("Please enter VISOR's website password")
			var website_password string
			_, _ = fmt.Scanln(&website_password)

			Utils.User_settings_GL.General.Website_domain = website_domain
			Utils.User_settings_GL.General.Website_pw = website_password

			UtilsSWA.WaitForNetwork(10)

			// Load or sync the user settings
			if !SettingsSync.SyncUserSettings(false) {
				log.Println("Failed to obtain user settings. Exiting...")

				return
			}
		}

		//////////////////////////////////////////
		// Prepare to hide the window

		if !isOpenGLSupport() {
			log.Println("Required OpenGL version not supported. Exiting...")

			return
		}

		// All mainly alright, let's hide the terminal window
		if runtime.GOOS == "windows" {
			if !Utils.WasArgUsedGENERAL(os.Args, "--conhost") {
				// Restart the process with conhost.exe on Windows to be able to actually hide the window
				if Utils.StartConAppPROCESSES(Utils.GetBinDirFILESDIRS().Add2(true, filepath.Base(os.Args[0])), "--conhost") {
					return
				}
			}
		}
		Utils.HideConsoleWindowPROCESSES()

		//////////////////////////////////////////
		// No terminal window from here on

		// Keep syncing the user settings with the server.
		go func() {
			SettingsSync.SyncUserSettings(true)
		}()

		ClientRegKeys.RegisterValues()

		var modules []Utils.Module
		for i := 0; i < Utils.MODS_ARRAY_SIZE; i++ {
			modules = append(modules, Utils.Module{
				Num:     i,
				Name:    Utils.GetModNameMODULES(i),
				Stop:    true,
				Stopped: true,
				Enabled: true,
			})
		}
		modules[Utils.NUM_MOD_VISOR].Stop = false
		modules[Utils.NUM_MOD_VISOR].Stopped = false
		// The Manager needs to be started first. It'll handle the others.
		modules[Utils.NUM_MOD_ModManager].Stop = false

		ModulesManager.Start(modules)

		// Create a new application
		my_app_GL = app.NewWithID("com.edw590.visor_c")
		my_app_GL.SetIcon(Logo.LogoBlackGmail)
		my_window_GL = my_app_GL.NewWindow("V.I.S.O.R.")

		processCommsChannel()

		// Create the content area with a label to display different screens
		var content_label *widget.Label = widget.NewLabel("Welcome!")
		var content_container *fyne.Container = container.NewVBox(content_label)

		// Create the navigation bar
		var nav_bar *fyne.Container = container.NewVBox(
			widget.NewButton("Home", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.Home()}
				content_container.Refresh()
			}),
			widget.NewButton("Dev Mode", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.DevMode(my_window_GL)}
				content_container.Refresh()
			}),
			widget.NewButton("Communicator", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.Communicator()}
				content_container.Refresh()
			}),
			widget.NewButton("Modules Status", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.ModulesStatus(modules)}
				content_container.Refresh()
			}),
			widget.NewButton("Calendar", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.Calendar()}
				content_container.Refresh()
			}),
			widget.NewButton("Global values", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.GlobalValues()}
				content_container.Refresh()
			}),
			widget.NewButton("System State", func() {
				content_container.Objects = []fyne.CanvasObject{Screens.SystemState()}
				content_container.Refresh()
			}),
		)

		// Create a split container to hold the navigation bar and the content
		var split *container.Split = container.NewHSplit(nav_bar, content_container)
		split.SetOffset(0.2) // Set the split ratio (20% for nav, 80% for content)

		// Set the content of the window
		my_window_GL.SetContent(split)

		var prev_screen fyne.CanvasObject = nil
		// Add system tray functionality
		if desk, ok := my_app_GL.(desktop.App); ok {
			var icon *fyne.StaticResource = Logo.LogoBlackGmail
			var menu *fyne.Menu = fyne.NewMenu("Tray",
				fyne.NewMenuItem("Show", func() {
					// Hide too because in case the window is shown but behind other apps, it won't show. So hiding and
					// showing does it. Maybe this happens because RequestFocus doesn't always work? Who knows. But this
					// fixes whatever the problem is.
					my_window_GL.Hide()
					my_window_GL.Show()
					my_window_GL.RequestFocus()

					// Restore the previous screen state
					Screens.Current_screen_GL = prev_screen
				}),
				fyne.NewMenuItem("Quit (USE THIS ONE)", func() {
					quitApp(modules)
				}),
			)
			desk.SetSystemTrayMenu(menu)
			desk.SetSystemTrayIcon(icon)
		}

		// Minimize to tray on close
		my_window_GL.SetCloseIntercept(func() {
			// Store the previous screen before hiding
			prev_screen = Screens.Current_screen_GL
			Screens.Current_screen_GL = nil
			my_window_GL.Hide()

			// Create and send one-time notification
			var notification_title string = "V.I.S.O.R. minimized"
			var notification_text string = "I'm still running in the background. To quit, use the system tray menu."
			notification := fyne.NewNotification(notification_title, notification_text)
			my_app_GL.SendNotification(notification)
		})

		go func() {
			for {
				if *module_stop {
					SettingsSync.StopUserSettingsSyncer()
					Utils.StopCommunicatorSERVER()
					quitApp(modules)

					return
				}

				Utils.SaveUserSettings()

				time.Sleep(5 * time.Second)
			}
		}()

		// Show and run the application
		my_window_GL.Resize(fyne.NewSize(640, 480))
		my_window_GL.ShowAndRun()
	}
}

/*
processCommsChannel processes in a different thread the communications channel.
*/
func processCommsChannel() {
	go func() {
		for {
			var comms_map map[string]any = <- Utils.ModsCommsChannels_GL[Utils.NUM_MOD_VISOR]
			if comms_map == nil {
				return
			}
			map_value, ok := comms_map["Notification"]
			if ok {
				var notif_info []string = map_value.([]string)
				notification := fyne.NewNotification(notif_info[0], notif_info[1])
				my_app_GL.SendNotification(notification)

				time.Sleep(5 * time.Second)
			} else {
				map_value, ok = comms_map["ShowApp"]
				if ok {
					showWindow()
				}
			}
		}
	}()
}

func isOpenGLSupport() bool {
	err := glfw.Init()
	if err != nil {
		return false
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.Visible, glfw.False)
	window, err := glfw.CreateWindow(1, 1, "", nil, nil)
	if err != nil {
		return false
	}

	defer window.Destroy()

	return true
}

func showWindow() {
	// Hide too because in case the window is shown but behind other apps, it won't show. So hiding and
	// showing does it. Maybe this happens because RequestFocus doesn't always work? Who knows. But this
	// fixes whatever the problem is.
	my_window_GL.Hide()
	my_window_GL.Show()
	my_window_GL.RequestFocus()
}

func quitApp(modules []Utils.Module) {
	Utils.CloseCommsChannels()
	Utils.SignalModulesStopMODULES(modules)

	my_app_GL.Quit()
}
