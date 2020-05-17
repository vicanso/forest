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

func TestConfigurationValidate(t *testing.T) {
	assert := assert.New(t)
	t.Run("xConfigName", func(t *testing.T) {
		type xConfigName struct {
			Value string `json:"value,omitempty" validate:"xConfigName"`
		}
		x := xConfigName{}

		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xConfigName.Value' Error:Field validation for 'Value' failed on the 'xConfigName' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": "abcd"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "测试"}`))
		assert.Equal(`Key: 'xConfigName.Value' Error:Field validation for 'Value' failed on the 'xConfigName' tag`, err.Error())

	})

	t.Run("xConfigCategory", func(t *testing.T) {
		type xConfigCategory struct {
			Value string `json:"value,omitempty" validate:"xConfigCategory"`
		}

		x := xConfigCategory{}

		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xConfigCategory.Value' Error:Field validation for 'Value' failed on the 'xConfigCategory' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": "abcd"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "测试"}`))
		assert.Equal(`Key: 'xConfigCategory.Value' Error:Field validation for 'Value' failed on the 'xConfigCategory' tag`, err.Error())
	})

	t.Run("xConfigData", func(t *testing.T) {
		type xConfigData struct {
			Value string `json:"value,omitempty" validate:"xConfigData"`
		}

		x := xConfigData{}
		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "测试"}`))
		assert.Nil(err)
	})

	t.Run("xConfigStatus", func(t *testing.T) {
		type xConfigStatus struct {
			Value int `json:"value,omitempty" validate:"xConfigStatus"`
		}
		x := xConfigStatus{}
		err := doValidate(&x, []byte(`{"value": 0}`))
		assert.Equal(`Key: 'xConfigStatus.Value' Error:Field validation for 'Value' failed on the 'xConfigStatus' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": 1}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": 3}`))
		assert.Equal(`Key: 'xConfigStatus.Value' Error:Field validation for 'Value' failed on the 'xConfigStatus' tag`, err.Error())
	})
}
