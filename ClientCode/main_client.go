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

package main

import (
	"GPT/GPT"
	MOD_1 "ModManager"
	"OIG/OIG"
	"Registry/Registry"
	"Utils"
	"VISOR_Client/Logo"
	"VISOR_Client/Screens"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var my_app_GL fyne.App = nil
var my_window_GL fyne.Window = nil

type _MGIModSpecInfo any
var (
	realMain        Utils.RealMain = nil
	moduleInfo_GL   Utils.ModuleInfo[_MGIModSpecInfo]
)
func main() {
	var module Utils.Module = Utils.Module{
		Num:     Utils.NUM_MOD_VISOR,
		Name:    Utils.GetModNameMODULES(Utils.NUM_MOD_VISOR),
		Stop:    false,
		Stopped: false,
		Enabled: true,
	}
	Utils.ModStartup2[_MGIModSpecInfo](realMain, &module, false)
}
func init() {realMain =
	func(module_stop *bool, moduleInfo_any any) {
		if !isOpenGLSupport() {
			log.Println("Required OpenGL version not supported. Exiting...")

			return
		}

		// All mainly alright, let's hide the terminal window
		if runtime.GOOS == "windows" {
			if !Utils.WasArgUsedGENERAL(os.Args, "--conhost") {
				// Restart the process with conhost.exe on Windows to be able to actually hide the window
				if Utils.StartConAppPROCESSES(Utils.GetBinDirFILESDIRS().Add2(true, "VISOR.exe"), "--conhost") {
					return
				}
			}
		}
		Utils.HideConsoleWindowPROCESSES()

		//////////////////////////////////////////
		// No terminal window from here on

		GPT.SetWebsiteInfo(Utils.PersonalConsts_GL.WEBSITE_URL, Utils.PersonalConsts_GL.WEBSITE_PW)
		OIG.SetWebsiteInfo(Utils.PersonalConsts_GL.WEBSITE_URL, Utils.PersonalConsts_GL.WEBSITE_PW)

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
		// Just for it to print that VISOR is running
		modules[Utils.NUM_MOD_VISOR].Stop = false
		modules[Utils.NUM_MOD_VISOR].Stopped = false
		// The Manager needs to be started first. It'll handle the others.
		modules[Utils.NUM_MOD_ModManager].Stop = false

		MOD_1.Start(modules)

		// Create a new application
		my_app_GL = app.NewWithID("com.edw590.visor_c")
		my_app_GL.SetIcon(Logo.LogoBlackGmail)
		my_window_GL = my_app_GL.NewWindow("V.I.S.O.R.")

		processNotifications()

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
				content_container.Objects = []fyne.CanvasObject{Screens.DevMode()}
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
				fyne.NewMenuItem("Quit", func() {
					Utils.SignalModulesStopMODULES(modules)

					my_app_GL.Quit()
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
		})

		go func() {
			for {
				if Registry.GetValue(Registry.K_SHOW_APP_SIG).GetBool(true) {
					showWindow()
					Registry.GetValue(Registry.K_SHOW_APP_SIG).SetBool(false, false)
				}

				if Utils.WaitWithStopTIMEDATE(module_stop, 1) {
					return
				}
			}
		}()

		// Show and run the application
		my_window_GL.Resize(fyne.NewSize(640, 480))
		my_window_GL.ShowAndRun()
	}
}

/*
processNotifications processes in a different thread the notifications queued in the notifications folder.
 */
func processNotifications() {
	go func() {
		for {
			var file_list []Utils.FileInfo = Utils.GetUserDataDirMODULES(Utils.NUM_MOD_VISOR).
				Add2(true, Utils.NOTIFS_REL_FOLDER).GetFileList()
			for _, file := range file_list {
				// Display the notification
				notification := fyne.NewNotification(strings.Split(file.Name, "-")[0], *file.GPath.ReadTextFile())
				my_app_GL.SendNotification(notification)

				// Remove the file
				_ = file.GPath.Remove()

				time.Sleep(5 * time.Second)
			}

			time.Sleep(1 * time.Second)
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
