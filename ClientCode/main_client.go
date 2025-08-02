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

package main

import (
	"ModulesManager"
	"SettingsSync"
	"Utils"
	"Utils/UtilsSWA"
	"VISOR_Client/ClientRegKeys"
	"VISOR_Client/Logo"
	"VISOR_Client/Screens"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	Tcef "github.com/Edw590/TryCatch-go"
	"github.com/go-gl/glfw/v3.3/glfw"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var my_app_GL fyne.App = nil
var my_window_GL fyne.Window = nil

var quitting_GL bool = false

var modules_GL []Utils.Module = nil

var content_container_GL *fyne.Container = nil

var modDirsInfo_GL Utils.ModDirsInfo
func main() {
	// Command line arguments
	var flag_log_level *int = flag.IntP("loglevel", "l", 0, "Log level to use. 0 = ERROR, 1 = WARNING, 2 = INFO, 3 = " +
		"DEBUG. Default is 0 (ERROR).")
	flag.Bool("nohide", false, "Don't hide the terminal window on startup.")
	flag.Bool("conhost", false, "Set this if the process was started by conhost.exe to be able to hide the terminal " +
		"window on Windows 10 Build 2004 or newer.")
	flag.Parse()
	Utils.SetLogLevel(*flag_log_level)

	var module Utils.Module = Utils.Module{
		Num:     Utils.NUM_MOD_VISOR,
		Name:    Utils.GetModNameMODULES(Utils.NUM_MOD_VISOR),
		Stop:    false,
		Stopped: false,
		Enabled: true,
	}
	Utils.ModStartup2(realMain, &module, false)
}
func realMain(module_stop *bool, moduleInfo_any any) {
	modDirsInfo_GL = moduleInfo_any.(Utils.ModDirsInfo)

	//////////////////////////////////////////
	// Get the user settings

	if err := Utils.ReadSettingsFile(true); err != nil {
		Utils.LogLnInfo("Failed to load user settings. Using empty ones...")
		Utils.LogLnInfo(err)
	}

	//////////////////////////////////////////
	// Prepare to hide the window

	if !isOpenGLSupported() {
		Utils.LogLnError("Required OpenGL version not supported. Exiting...")

		return
	}

	if !Utils.WasArgUsedGENERAL(os.Args, "--nohide") {
		// All mainly alright, let's hide the terminal window
		if runtime.GOOS == "windows" {
			maj, min, patch := Utils.GetOSVersionSYSTEM()
			if maj >= 10 && min >= 0 && patch >= 19041 {
				// Restart the process with conhost.exe on Windows to be able to actually hide the window if we're
				// on Windows 10 Build 2004 or newer (because of the new Windows Terminal).
				if !Utils.WasArgUsedGENERAL(os.Args, "--conhost") {
					if Utils.StartConAppPROCESSES(Utils.GetBinDirFILESDIRS().Add2(true, filepath.Base(os.Args[0])),
							"--conhost") {
						return
					}
				}
			}
		}
		Utils.HideConsoleWindowPROCESSES()
	}

	//////////////////////////////////////////
	// No terminal window from here on

	Utils.StartCommunicatorSERVER()

	// Keep syncing the user settings with the server.
	SettingsSync.SyncUserSettings()

	ClientRegKeys.RegisterValues()

	for i := 0; i < Utils.MODS_ARRAY_SIZE; i++ {
		modules_GL = append(modules_GL, Utils.Module{
			Num:     i,
			Name:    Utils.GetModNameMODULES(i),
			Stop:    true,
			Stopped: true,
			Enabled: true,
		})
	}
	modules_GL[Utils.NUM_MOD_VISOR].Stop = false
	modules_GL[Utils.NUM_MOD_VISOR].Stopped = false
	// The Manager needs to be started first. It'll handle the others.
	modules_GL[Utils.NUM_MOD_ModManager].Stop = false

	ModulesManager.Start(modules_GL)

	// Create a new application
	my_app_GL = app.NewWithID("com.edw590.visor_c")
	my_app_GL.SetIcon(Logo.LogoAdaptiveAllModded)
	my_window_GL = my_app_GL.NewWindow("V.I.S.O.R.")
	my_window_GL.Resize(fyne.NewSize(640, 480))
	Screens.Current_window_GL = my_window_GL

	go processCommsChannel()

	// Create the content area with a label to display different screens
	content_container_GL = container.NewStack()

	// Set the initial screen and lock the app right when it starts.
	Screens.Current_screen_GL = Screens.ID_MOD_GPT_COMM
	lockApp()

	var nav_bar *widget.Tree = &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return tree_index[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := tree_index[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := screens_GL[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: " + uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t)
		},
		OnSelected: func(uid string) {
			showScreen(uid)
		},
	}

	var themes *fyne.Container = container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			my_app_GL.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			my_app_GL.Settings().SetTheme(theme.LightTheme())
		}),
		widget.NewButton("Auto", func() {
			my_app_GL.Settings().SetTheme(theme.DefaultTheme())
		}),
		widget.NewButton("Lock", func() {
			if Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Pin == "" {
				dialog.ShowInformation("No PIN set", "You need to set a PIN in the settings to lock the app.",
					my_window_GL)

				return
			}

			lockApp()
		}),
	)

	var sidebar *fyne.Container = container.NewBorder(nil, themes, nil, nil, nav_bar)

	// Create a split container to hold the sidebar and the content
	var split *container.Split = container.NewHSplit(sidebar, content_container_GL)
	split.SetOffset(0.2) // Set the split ratio (20% for sidebar, 80% for content)

	// Set the content of the window
	my_window_GL.SetContent(split)

	// Add system tray functionality
	if desk, ok := my_app_GL.(desktop.App); ok {
		var icon *fyne.StaticResource = Logo.LogoAdaptiveAllModded
		var menu *fyne.Menu = fyne.NewMenu("Tray",
			fyne.NewMenuItem("Show", func() {
				showWindow()
			}),
			fyne.NewMenuItem("Quit (USE THIS ONE)", func() {
				quitApp(modules_GL)
			}),
		)
		desk.SetSystemTrayMenu(menu)
		desk.SetSystemTrayIcon(icon)
	}

	// Minimize to tray on close
	my_window_GL.SetCloseIntercept(func() {
		if !UtilsSWA.GetValueREGISTRY(ClientRegKeys.K_MINIMIZE_TO_TRAY).GetBool(true) {
			quitApp(modules_GL)

			return
		}

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
				quitApp(modules_GL)

				return
			}

			Utils.WriteSettingsFile(true)

			time.Sleep(5 * time.Second)
		}
	}()

	// Show and run the application
	my_window_GL.ShowAndRun()
}

func showScreen(uid string) {
	switch uid {
		case "":
			content_container_GL.Objects = []fyne.CanvasObject{}
		case Screens.ID_HOME:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.Home()}
		case Screens.ID_MOD_MOD_MANAGER:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModModulesManager(modules_GL)}
		case Screens.ID_MOD_SPEECH:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModSpeech()}
		case Screens.ID_MOD_RSS_FEED_NOTIFIER:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModRSSFeedNotifier()}
		case Screens.ID_MOD_GPT_COMM:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModGPTCommunicator()}
		case Screens.ID_MOD_TASKS_EXECUTOR:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModTasksExecutor()}
		case Screens.ID_MOD_USER_LOCATOR:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModUserLocator()}
		case Screens.ID_ONLINE_INFO_CHK:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModOnlineInfoChk()}
		case Screens.ID_MOD_SYS_CHECKER:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModSystemChecker()}
		case Screens.ID_SMART_CHECKER:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModSMARTChecker()}
		case Screens.ID_REGISTRY:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.Registry()}
		case Screens.ID_GOOGLE_MANAGER:
			content_container_GL.Objects = []fyne.CanvasObject{Screens.ModGoogleManager()}
	}
	content_container_GL.Refresh()
}

/*
processCommsChannel processes in a different thread the communications channel.
*/
func processCommsChannel() {
	for {
		var comms_map map[string]any = Utils.GetFromCommsChannel(true, Utils.NUM_MOD_VISOR, 0, -1)
		if comms_map == nil {
			return
		}

		if map_value, ok := comms_map["Notification"]; ok {
			var notif_info []string = map_value.([]string)
			notification := fyne.NewNotification(notif_info[0], notif_info[1])
			my_app_GL.SendNotification(notification)

			time.Sleep(5 * time.Second)
		} else if _, ok = comms_map["ShowApp"]; ok {
			showWindow()
		} else if _, ok = comms_map["Redraw"]; ok {
			showScreen(Screens.Current_screen_GL)
		}
	}
}

func isOpenGLSupported() bool {
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
	// Hide too because in case the window is shown but behind other apps, it won't show. So hiding and showing does it.
	// Maybe this happens because RequestFocus doesn't always work? Who knows. But this fixes whatever the problem is.
	my_window_GL.Hide()
	my_window_GL.Show()
	my_window_GL.RequestFocus()

	lockApp()
}

func quitApp(modules []Utils.Module) {
	if quitting_GL {
		return
	}
	quitting_GL = true

	Utils.CloseCommsChannels()
	Utils.SignalModulesStopMODULES(modules)

	// Ignore a "runtime error: invalid memory address or nil pointer dereference" that happens who knows why.
	Tcef.Tcef{
		Try: func() {
			my_app_GL.Quit()
		},
	}.Do()
}

func lockApp() {
	showScreen("")

	createPinDialog(func() {
		showScreen(Screens.Current_screen_GL)
	})
}

func createPinDialog(callback func()) {
	if Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Pin == "" {
		callback()

		return
	}

	var entry_pin *widget.Entry = widget.NewPasswordEntry()
	entry_pin.SetPlaceHolder("PIN")
	entry_pin.Validator = validation.NewRegexp(`^\d+$`, "PIN must be numberic")

	var form_items []*widget.FormItem = []*widget.FormItem{
		widget.NewFormItem("PIN", entry_pin),
	}
	dialog.ShowForm("Insert PIN", "Unlock", "Cancel", form_items, func(b bool) {
		if !b {
			lockApp()

			return
		}

		if entry_pin.Text == Utils.GetUserSettings(Utils.LOCK_UNLOCK).General.Pin {
			callback()
		} else {
			dialog.ShowInformation("Wrong PIN", "The PIN you entered is wrong.", my_window_GL)
			lockApp()
		}
	}, my_window_GL)
}
