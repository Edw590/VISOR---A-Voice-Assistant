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

package Utils

import (
	"github.com/zalando/go-keyring"
)

/*
SavePasswordCREDENTIALS saves the password in the keyring.

-----------------------------------------------------------

– Params:
  - password – the password to save

– Returns:
  - an error if the password could not be saved, nil otherwise
 */
func SavePasswordCREDENTIALS(password string) error {
	return keyring.Set("VISOR", "user", password)
}

/*
GetPasswordCREDENTIALS returns the password from the keyring.

-----------------------------------------------------------

– Returns:
  - the password from the keyring

– Returns:
  - the password from the keyring or an empty string if the password could not be retrieved
 */
func GetPasswordCREDENTIALS() string {
	password, err := keyring.Get("VISOR", "user")
	if err != nil {
		return ""
	}

	return password
}

/*
DeletePasswordCREDENTIALS deletes the password from the keyring.

-----------------------------------------------------------

– Returns:
  - an error if the password could not be deleted, nil otherwise
 */
func DeletePasswordCREDENTIALS() error {
	return keyring.Delete("VISOR", "user")
}
