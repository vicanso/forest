// Copyright 2021 tree xie
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

package log

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/forest/cs"
)

func init() {
	SetLogRegexp(cs.MaskRegExp, nil)
}

func TestCutOrMask(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("***", cutOrMaskString("password", "abcd"))
	assert.Equal(logFieldValueMaxSize+3, len(cutOrMaskString("abc", "0123456789012345678901234567890123456789")))

	assert.Equal("***", cutOrMaskInterface("password", map[string]string{}))
	assert.Equal(true, cutOrMaskInterface("bool", true))
}

func TestLogStruct(t *testing.T) {
	assert := assert.New(t)

	type T struct {
		Name string
	}

	e := Struct(&T{
		Name: "abc",
	})
	assert.NotNil(e)

	e = Struct([]*T{
		{
			Name: "abc",
		}, {
			Name: "def",
		},
	})
	assert.NotNil(e)

	e = Struct(url.Values{
		"key": []string{
			"a",
			"b",
		},
		"name": []string{
			"a",
		},
	})
	assert.NotNil(e)
}
