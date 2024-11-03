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
	"Utils/UtilsSWA"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"strings"
)

func Settings() fyne.CanvasObject {
	Current_screen_GL = ID_SETTINGS

	return container.NewAppTabs(
		container.NewTabItem("Local settings", settingsCreateMainTab()),
	)
}

func settingsCreateMainTab() *container.Scroll {
	var objects []fyne.CanvasObject = nil
	var values []*UtilsSWA.Value = UtilsSWA.GetValuesREGISTRY()
	for i := len(values) - 1; i >= 0; i-- {
		var value *UtilsSWA.Value = values[i]
		if !value.Auto_set && strings.HasPrefix(value.Pretty_name, "General - ") {
			objects = append(objects, createValueChooserUTILS(value))
		}
	}

	return createMainContentScrollUTILS(objects...)
}
