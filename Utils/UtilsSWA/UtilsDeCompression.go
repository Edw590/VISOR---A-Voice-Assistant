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

package UtilsSWA

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io"
	"log"
)

/*
CompressString compresses a string.

-----------------------------------------------------------

– Params:
	- to_compress – the string to compress

– Returns:
	- the compressed string
 */
func CompressString(to_compress string) []byte {
	var buffer bytes.Buffer

	writer := brotli.NewWriterLevel(nil, 99999)

	// Reset the compressor and encode from some input stream.
	writer.Reset(&buffer)
	if _, err := io.WriteString(writer, to_compress); err != nil {
		log.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}

	return buffer.Bytes()
}

/*
DecompressString decompresses a string.

-----------------------------------------------------------

– Params:
	- to_decompress – the string to decompress

– Returns:
	- the decompressed string
 */
func DecompressString(to_decompress []byte) string {
	reader := brotli.NewReader(bytes.NewReader(to_decompress))

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		log.Fatal(err)
	}

	return buffer.String()
}
