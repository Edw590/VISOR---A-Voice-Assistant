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

package UtilsSWA

import "Utils"

/*
EncryptBytesCRYPTOENDECRYPT encrypts the given data using the parameters defined on the file doc (encode with UTF-7
first).

Check if the device is running low on memory before calling this function! It needs some memory to calculate the keys!

ATTENTION: the passwords' order must NOT be changed once the passwords are used to encrypt some data! Use them always in
the same order they were entered!

-----------------------------------------------------------

– Params:
  - raw_password1 – the first character sequence to calculate the 2 keys from
  - raw_password2 – the second character sequence to calculate the 2 keys from
  - raw_data – the data to encrypt encoded in UTF-7
  - raw_aad_suffix – additional not encrypted metadata suffix to include in the encrypted message, right after
    _RAW_AAD_PREFIX; nil if not to be used

– Returns:
  - the encrypted message using the mentioned method, or nil if the chosen algorithm was unable to process the data
    provided.
*/
func EncryptBytesCRYPTOENDECRYPT(raw_password1 []byte, raw_password2 []byte, raw_data []byte, raw_aad_suffix []byte) []byte {
	return Utils.EncryptBytesCRYPTOENDECRYPT(raw_password1, raw_password2, raw_data, raw_aad_suffix)
}

/*
DecryptBytesCRYPTOENDECRYPT decrypts the given data using the parameters defined on the file doc (decode with UTF-7
afterwards).

Check if the device is running low on memory before calling this function! It needs some memory to calculate the keys!

ATTENTION: the passwords order must NOT be changed once the passwords are used to encrypt some data! Use them always in
the same order they were entered!

-----------------------------------------------------------

– Params:
  - raw_password1 – the first character sequence to calculate the 2 keys from
  - raw_password2 – the second character sequence to calculate the 2 keys from
  - raw_data – the data to encrypt encoded in UTF-7
  - raw_aad_suffix – the associated authenticated data suffix used with the encrypted message; or nil if not to be used

– Returns:
  - the original message text; nil in case either the message was not encrypted using the parameters defined in the
	file doc or in case it has been tampered with.
*/
func DecryptBytesCRYPTOENDECRYPT(raw_password1 []byte, raw_password2 []byte, raw_data []byte, raw_aad_suffix []byte) []byte {
	return Utils.DecryptBytesCRYPTOENDECRYPT(raw_password1, raw_password2, raw_data, raw_aad_suffix)
}
