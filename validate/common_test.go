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

package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonValidate(t *testing.T) {
	assert := assert.New(t)

	t.Run("xLimit", func(t *testing.T) {
		type xLimit struct {
			Value int `json:"value" validate:"xLimit"`
		}

		x := xLimit{}

		err := doValidate(&x, []byte(`{"value": 0}`))
		assert.Equal(`Key: 'xLimit.Value' Error:Field validation for 'Value' failed on the 'xLimit' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": 10}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": 11}`))
		assert.Equal(`Key: 'xLimit.Value' Error:Field validation for 'Value' failed on the 'xLimit' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": -1}`))
		assert.Equal(`Key: 'xLimit.Value' Error:Field validation for 'Value' failed on the 'xLimit' tag`, err.Error())

	})

	t.Run("xDuration", func(t *testing.T) {
		type xDuration struct {
			Value string `json:"value" validate:"xDuration"`
		}

		x := xDuration{}

		err := doValidate(&x, []byte(`{"value": "1s"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "1m"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "1h"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xDuration.Value' Error:Field validation for 'Value' failed on the 'xDuration' tag`, err.Error())
	})
}
