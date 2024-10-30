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
	"fyne.io/fyne/v2/widget"
	"time"
)

func Registry() fyne.CanvasObject {
	Current_screen_GL = ID_REGISTRY

	return container.NewAppTabs(
		container.NewTabItem("All values", registryCreateAllValuesTab()),
	)
}

func registryCreateAllValuesTab() *container.Scroll {
	var registry_text *widget.Label = widget.NewLabel("")
	registry_text.Wrapping = fyne.TextWrapWord

	go func() {
		time.Sleep(500 * time.Millisecond)
		for {
			if Current_screen_GL == ID_REGISTRY {
				registry_text.SetText(UtilsSWA.GetRegistryTextREGISTRY())
			} else {
				break
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return createMainContentScrollUTILS(registry_text)
}
