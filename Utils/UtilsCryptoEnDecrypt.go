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
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"math"

	"golang.org/x/crypto/scrypt"
)

/*
 * This utility file encrypts and decrypts the given data using the method Scrypt-65536-8-1 + AES-256/CBC/PKCS#5 +
 * HMAC-SHA512.
 *
 * Additional information:
 * - Initialization Vector of 128 bits (16 bytes) for both AES-CBC-PKCS#5 and HMAC.
 * - Key of 256 bits (32 bytes) to encrypt the data using AES-CBC-PKCS#5.
 * - Key of 512 bits (64 bytes) to calculate the HMAC tag, which has a length of 512 bits (64 bytes).
 * - The keys are wiped from memory as soon as they're no longer needed.
 * - Constant-time preventions not taken in consideration.
 *
 * Note: to calculate the keys for AES (256 bits) and HMAC (512 bits), the SHA-512 hashes of 2 passwords are used.
 * Those hashes are then used to calculate the keys using Scrypt (N = 65536; r = 8; p = 1).
 *
 * Various things learned from: Security Best Practices: Symmetric Encryption with AES in Java and Android
 * (https://proandroiddev.com/security-best-practices-symmetric-encryption-with-aes-in-java-7616beaaade9) and
 * Security Best Practices: Symmetric Encryption with AES in Java and Android: Part 2: AES-CBC + HMAC
 * (https://proandroiddev.com/security-best-practices-symmetric-encryption-with-aes-in-java-and-android-part-2-b3b80e99ad36).
 */

const _IV_LENGTH_BYTES int = aes.BlockSize // 16 bytes, 128 bits
const _HMAC_TAG_LENGTH_BYTES int = 64      // 512 bits
const _AES_KEY_SIZE int = 32               // 256 bits (for use with AES-256)
const _HMAC_KEY_SIZE int = 64              // 512 bits --> ALWAYS AT LEAST THE OUTPUT LENGTH FOR SECURITY (rules)
const _INIT_LEN_ENCRYPTED_MSG int = 1 + _IV_LENGTH_BYTES + 1 + _HMAC_TAG_LENGTH_BYTES

const _RAW_AAD_SEPARATOR_STR string = " \\\\\\/// "
var _RAW_AAD_PREFIX [] byte = []byte("Scrypt-65536-8-1 + AES-256/CBC/PKCS#5 + HMAC-SHA512" + _RAW_AAD_SEPARATOR_STR)

// Scrypt parameters - tested a bit slow on MiTab Advance and slower on BV9500 (wtf). Good enough, I think.
// Memory required to calculate the keys: 128 * N * r * p --> 128 * 8_192 * 8 * 1 = 8_388_608 B = 8.4 MB
// EDIT: changed from 16_384 to 8_192 to be faster on BV9500 (takes a few seconds still, maybe 30).

const _SCRYPT_N int = 65_536
const _SCRYPT_R int = 8
const _SCRYPT_P int = 1













// TODO: IMPLEMENT TIME-CONSTANT PREVENTIONS!!!!














/*
EncryptBytesCRYPTOENDECRYPT encrypts the given data using the parameters defined on the file doc.

Check if the device is running low on memory before calling this function! It needs some memory to calculate the keys!

ATTENTION: the passwords' order must NOT be changed once the passwords are used to encrypt some data! Use them always in
the same order they were entered!

-----------------------------------------------------------

– Params:
  - raw_password1 – the first character sequence to calculate the 2 keys from
  - raw_password2 – the second character sequence to calculate the 2 keys from
  - raw_data – the data to encrypt
  - raw_aad_suffix – additional not encrypted metadata suffix to include in the encrypted message, right after
    _RAW_AAD_PREFIX; nil if not to be used

– Returns:
  - the encrypted message using the mentioned method, or nil if the chosen algorithm was unable to process the data
    provided.
 */
func EncryptBytesCRYPTOENDECRYPT(raw_password1 []byte, raw_password2 []byte, raw_data []byte, raw_aad_suffix []byte) []byte {
	var raw_aad_ready []byte = nil
	if raw_aad_suffix != nil {
		raw_aad_ready = getAADReady(raw_aad_suffix)
	}

	iv, err := getIv()
	keys, err := getKeys(raw_password1, raw_password2)
	if err != nil {
		return nil
	}

	cipher_block, err := getCipher(keys[0], iv, true)
	if err != nil {
		return nil
	}

	randomized_raw_data_padded := pkcs5Padding(raw_data, cipher_block.BlockSize())

	cipher_text := make([]byte, len(randomized_raw_data_padded))
	cipher_block.CryptBlocks(cipher_text, randomized_raw_data_padded)

	hmac_tag := getHmacTag(keys[1], iv, cipher_text, raw_aad_ready)

	encrypted_msg := make([]byte, 1 + _IV_LENGTH_BYTES + 1 + _HMAC_TAG_LENGTH_BYTES+ len(cipher_text))
	encrypted_msg[0] = byte(_IV_LENGTH_BYTES)
	copy(encrypted_msg[1:], iv)
	encrypted_msg[1 + _IV_LENGTH_BYTES] = byte(_HMAC_TAG_LENGTH_BYTES)
	copy(encrypted_msg[1 + _IV_LENGTH_BYTES + 1:], hmac_tag)
	copy(encrypted_msg[1 + _IV_LENGTH_BYTES + 1 +_HMAC_TAG_LENGTH_BYTES:], cipher_text)

	return encrypted_msg
}

/*
DecryptBytesCRYPTOENDECRYPT decrypts the given data using the parameters defined on the file doc.

Check if the device is running low on memory before calling this function! It needs some memory to calculate the keys!

ATTENTION: the passwords order must NOT be changed once the passwords are used to encrypt some data! Use them always in
the same order they were entered!

-----------------------------------------------------------

– Params:
  - raw_password1 – the first character sequence to calculate the 2 keys from
  - raw_password2 – the second character sequence to calculate the 2 keys from
  - raw_data – the data to encrypt
  - raw_aad_suffix – the associated authenticated data suffix used with the encrypted message; or nil if not to be used

– Returns:
  - the original message text; nil in case either the message was not encrypted using the parameters defined in the
	file doc or in case it has been tampered with.
 */
func DecryptBytesCRYPTOENDECRYPT(raw_password1 []byte, raw_password2 []byte, raw_encrypted_message []byte, raw_aad_suffix []byte) []byte {
	if len(raw_encrypted_message) <= _INIT_LEN_ENCRYPTED_MSG {
		return nil
	}

	iv_length := int(raw_encrypted_message[0])
	if iv_length != _IV_LENGTH_BYTES {
		return nil
	}

	hmac_tag_length := int(raw_encrypted_message[1 + _IV_LENGTH_BYTES])
	if hmac_tag_length != _HMAC_TAG_LENGTH_BYTES {
		return nil
	}

	iv := raw_encrypted_message[1 : 1 + _IV_LENGTH_BYTES]
	message_hmac := raw_encrypted_message[(1 + _IV_LENGTH_BYTES + 1):(1 + _IV_LENGTH_BYTES + 1 + _HMAC_TAG_LENGTH_BYTES)]
	cipher_text := raw_encrypted_message[1 + _IV_LENGTH_BYTES + 1 + _HMAC_TAG_LENGTH_BYTES:]

	var raw_aad_ready []byte = nil
	if raw_aad_suffix != nil {
		raw_aad_ready = getAADReady(raw_aad_suffix)
	}
	keys, err := getKeys(raw_password1, raw_password2)
	if err != nil {
		return nil
	}

	supposed_hmac := getHmacTag(keys[1], iv, cipher_text, raw_aad_ready)
	if !bytes.Equal(supposed_hmac, message_hmac) {
		return nil
	}

	cipher_block, err := getCipher(keys[0], iv, false)
	if err != nil {
		return nil
	}

	raw_data_padded := make([]byte, len(cipher_text))
	cipher_block.CryptBlocks(raw_data_padded, cipher_text)

	return pkcs5Unpadding(raw_data_padded)
}

/*
randomizeData randomizes UTF-7 encoded data for even safer encryption using this file' method.

This is done by adding a random byte value from 128 to 255 in a random index each 16 elements. So, for example, on index
5 a 173 byte is added, and on index 18 a 209 byte is added.

This means the new byte array will have
	length = data.length + Math.ceil(data.length/16)

Call derandomizeData() to undo this.

-----------------------------------------------------------

– Params:
  - data – UFT-7 encoded data to be randomized

– Returns:
  - the randomized data
 */
func randomizeData(data []byte) []byte {
	var b []byte = make([]byte, 1)
	var random_index int = -1

	var randomized_data_length int = len(data) + int(math.Ceil(float64(len(data))/16.0))
	var randomized_data[] byte = make([]byte, randomized_data_length)

	var num_added_bytes int = 0
	for i := 0; i < randomized_data_length; i++ {
		if i % 16 == 0 {
			_, _ = rand.Read(b)
			random_index = i + int(b[0]) % int(math.Min(float64(randomized_data_length-i), 16))
		}
		if i == random_index {
			_, _ = rand.Read(b)
			randomized_data[i] = b[0]%128 + 128
			num_added_bytes++
		} else {
			randomized_data[i] = data[i-num_added_bytes]
		}
	}

	return randomized_data
}

/*
derandomizeData undoes the randomization done by randomizeData().

-----------------------------------------------------------

– Params:
  - randomized_data – the randomized data

– Returns:
  - the original data
 */
func derandomizeData(randomizedData []byte) []byte {
	var dataLength int = int(float64(len(randomizedData)) / (1.0 + (1.0 / 16.0)))
	var data []byte = make([]byte, dataLength)

	var i int = 0
	for _, b := range randomizedData {
		if b <= 127 {
			data[i] = b
			i++
		}
	}

	return data[:i]
}

/*
getAADReady gets the Authenticated Associated Data (AAD) ready for use with correct prefixes.

-----------------------------------------------------------

– Params:
  - raw_aad_suffix – the additional AAD to add to the final one

– Returns:
  - the ready AAD byte vector
 */
func getAADReady(raw_aad_suffix []byte) []byte {
	return append(_RAW_AAD_PREFIX, raw_aad_suffix...)
}

/*
getIv creates an initialization vector as randomly as possible, with length of {@link #IV_LENGTH_BYTES} bytes to use
with AES.

-----------------------------------------------------------

– Returns:
  - the initialization vector
 */
func getIv() ([]byte, error) {
	iv := make([]byte, _IV_LENGTH_BYTES)
	_, _ = rand.Read(iv)

	return iv, nil
}

/*
getKeys gets the keys to be used with AES and MAC according to the file doc.

ATTENTION: the passwords' order must NOT be changed once the passwords are used to encrypt some data! Use them always in
the same order they were entered!

-----------------------------------------------------------

– Params:
  - password1 – one of the passwords to use to create the keys
  - password2 – the other password to use to create the keys

– Returns:
  - 1st index AES key; 2nd index, MAC key; nil if there's not enough memory for the keys calculation
 */
func getKeys(password1 []byte, password2 []byte) ([][]byte, error) {
	password1_sha512 := GetHashBytesOfBytesCRYPTOHASHING(password1)
	password2_sha512 := GetHashBytesOfBytesCRYPTOHASHING(password2)

	key_aes, err := scrypt.Key(password1_sha512, password2_sha512, _SCRYPT_N, _SCRYPT_R, _SCRYPT_P, _AES_KEY_SIZE)
	if err != nil {
		return nil, err
	}

	key_hmac, err := scrypt.Key(password2_sha512, password1_sha512, _SCRYPT_N, _SCRYPT_R, _SCRYPT_P, _HMAC_KEY_SIZE)
	if err != nil {
		return nil, err
	}

	return [][]byte{key_aes, key_hmac}, nil
}

/*
getCipher gets the cipher to be used with AES according to the file doc.

-----------------------------------------------------------

– Params:
  - key_aes – the AES key to use with the cipher
  - iv – the initialization vector to use with the cipher
  - encrypt – true to prepare the Cipher to encrypt, false to prepare to decrypt

– Returns:
  - the cipher to be used with AES
 */
func getCipher(key_aes []byte, iv []byte, encrypt bool) (cipher.BlockMode, error) {
	block, err := aes.NewCipher(key_aes)
	if err != nil {
		return nil, err
	}

	var mode cipher.BlockMode
	if encrypt {
		mode = cipher.NewCBCEncrypter(block, iv)
	} else {
		mode = cipher.NewCBCDecrypter(block, iv)
	}

	return mode, nil
}

/*
getHmacTag creates an HMAC tag from the given data bytes using the algorithm defined in the file doc.

-----------------------------------------------------------

– Params:
  - key_mac – the key to create the tag
  - iv – the initialization vector to get the tag from
  - cipher_text – the encrypted data to get the tag from
  - associated_authed_data – additional not encrypted metadata to add to the tag, or nil if not to use

– Returns:
  - the tag bytes
 */
func getHmacTag(key_mac, iv, cipher_text, associated_authed_data []byte) []byte {
	h := hmac.New(sha512.New, key_mac)

	h.Write(iv)
	h.Write(cipher_text)
	if associated_authed_data != nil {
		h.Write(associated_authed_data)
	}

	return h.Sum(nil)
}

/*
pkcs5Padding adds padding to the given data to be used with AES.

-----------------------------------------------------------

– Params:
  - src – the data to add padding to
  - block_size – the block size to use with AES

– Returns:
  - the padded data
 */
func pkcs5Padding(src []byte, block_size int) []byte {
	padding := block_size - len(src)%block_size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padtext...)
}

/*
pkcs5Unpadding removes padding from the given data to be used with AES.

-----------------------------------------------------------

– Params:
  - src – the data to remove padding from

– Returns:
  - the unpadded data
 */
func pkcs5Unpadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
