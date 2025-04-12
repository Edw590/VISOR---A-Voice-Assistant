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

package UtilsSWA

import (
	"Utils"
)

/*
CompressBytes compresses bytes.

-----------------------------------------------------------

– Params:
	- to_compress – the bytes to compress

– Returns:
	- the compressed bytes or nil if an error occurred
*/
func CompressBytes(to_compress []byte) []byte {
	return Utils.CompressBytes(to_compress)
}

/*
DecompressBytes decompresses bytes.

-----------------------------------------------------------

– Params:
	- to_decompress – the bytes to decompress

– Returns:
	- the decompressed bytes or nil if an error occurred
*/
func DecompressBytes(to_decompress []byte) []byte {
	return Utils.DecompressBytes(to_decompress)
}
