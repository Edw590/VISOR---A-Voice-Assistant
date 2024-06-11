/*******************************************************************************
 * Copyright 2023-2023 Edw590
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
	"math"
)

/*
AdjustGainBufferAUDIO adjusts the volume of raw audio data.

Note that this function does not care about noise or anything at all - it just applies the same gain to ALL bytes. If
the gain would generate an overflow, the maximum value for the number of provided bits is used instead. If gain is 1,
this function is a no-op.

-----------------------------------------------------------

– Params:
  - audio_bytes – the audio data
  - gain – the gain to apply to the data
  - n_bits – the number of bits per sample (e.g. 16 for 16-bit PCM)

– Returns:
  - the adjusted audio data
 */
func AdjustGainBufferAUDIO(audio_bytes[] byte, gain float64, n_bits int32) {
	if gain == 1 {
		return
	}

	var audio_length int = len(audio_bytes)
	for i := 0; i < audio_length; i++ {
		audio_bytes[i] = byte(math.Min(float64(audio_bytes[i])*gain, float64(int(1)<<n_bits - 1)))
	}
}
