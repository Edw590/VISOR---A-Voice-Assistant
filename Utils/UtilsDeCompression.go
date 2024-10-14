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
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
)

/*
CompressString compresses a string.

-----------------------------------------------------------

– Params:
	- to_compress – the string to compress

– Returns:
	- the compressed string or nil if an error occurred
 */
func CompressString(to_compress string) []byte {
	var buffer bytes.Buffer

	writer := brotli.NewWriterLevel(nil, 99999)

	// Reset the compressor and encode from some input stream.
	writer.Reset(&buffer)
	if _, err := io.WriteString(writer, to_compress); err != nil {
		return nil
	}
	if err := writer.Close(); err != nil {
		return nil
	}

	return buffer.Bytes()
}

/*
DecompressString decompresses a string.

-----------------------------------------------------------

– Params:
	- to_decompress – the string to decompress

– Returns:
	- the decompressed string or an empty string if an error occurred
 */
func DecompressString(to_decompress []byte) string {
	reader := brotli.NewReader(bytes.NewReader(to_decompress))

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		return ""
	}

	return buffer.String()
}
