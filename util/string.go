// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const digitBytes = "0123456789"

// randomString create a random string
func randomString(baseLetters string, n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(baseLetters) {
			b[i] = baseLetters[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// RandomString create a random string
func RandomString(n int) string {
	return randomString(letterBytes, n)
}

// RandomDigit create a random digit string
func RandomDigit(n int) string {
	return randomString(digitBytes, n)
}

var entropy = rand.New(rand.NewSource(time.Unix(0, 0).UnixNano()))

// GenUlid generate ulid
func GenUlid() string {
	t := time.Now()
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}

// Sha256 gen sha256 string
func Sha256(str string) string {
	hash := sha256.New()
	_, _ = hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashBytes)
}

// ContainsString check the string slice contain the string
func ContainsString(arr []string, str string) (found bool) {
	for _, v := range arr {
		if found {
			break
		}
		if v == str {
			found = true
		}
	}
	return
}

// UserRoleIsValid check user rols is valid
func UserRoleIsValid(validRoles []string, userRoles []string) bool {
	valid := false
	for _, role := range validRoles {
		if ContainsString(userRoles, role) {
			valid = true
			break
		}
	}
	return valid
}

// Encrypt encrypt
// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64
func Encrypt(key, text []byte) ([]byte, error) {
	// 需要注意 key的长度必须为32字节
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt decrypt
// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
