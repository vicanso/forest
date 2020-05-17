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

func TestUserValidate(t *testing.T) {
	assert := assert.New(t)
	t.Run("xUserAccount", func(t *testing.T) {
		type xUserAccount struct {
			Value string `json:"value" validate:"xUserAccount"`
		}
		x := xUserAccount{}
		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xUserAccount.Value' Error:Field validation for 'Value' failed on the 'xUserAccount' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": "abcd"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "测试"}`))
		assert.Equal(`Key: 'xUserAccount.Value' Error:Field validation for 'Value' failed on the 'xUserAccount' tag`, err.Error())
	})

	t.Run("xUserPassword", func(t *testing.T) {
		type xUserPassword struct {
			Value string `json:"value" validate:"xUserPassword"`
		}
		x := xUserPassword{}
		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xUserPassword.Value' Error:Field validation for 'Value' failed on the 'xUserPassword' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": "fEqNCco3Yq9h5ZUglD3CZJT4lBsfEqNCco31Yq9h5ZUB"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "123"}`))
		assert.Equal(`Key: 'xUserPassword.Value' Error:Field validation for 'Value' failed on the 'xUserPassword' tag`, err.Error())
	})

	t.Run("xUserRole", func(t *testing.T) {
		type xUserRole struct {
			Value string `json:"value" validate:"xUserRole"`
		}
		x := xUserRole{}

		err := doValidate(&x, []byte(`{"value": ""}`))
		assert.Equal(`Key: 'xUserRole.Value' Error:Field validation for 'Value' failed on the 'xUserRole' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": "su"}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": "test"}`))
		assert.Equal(`Key: 'xUserRole.Value' Error:Field validation for 'Value' failed on the 'xUserRole' tag`, err.Error())
	})

	t.Run("xUserRoles", func(t *testing.T) {
		type xUserRoles struct {
			Value []string `json:"value" validate:"xUserRoles"`
		}
		x := xUserRoles{}
		err := doValidate(&x, []byte(`{"value": []}`))
		assert.Equal(`Key: 'xUserRoles.Value' Error:Field validation for 'Value' failed on the 'xUserRoles' tag`, err.Error())

		err = doValidate(&x, []byte(`{"value": ["su"]}`))
		assert.Nil(err)

		err = doValidate(&x, []byte(`{"value": ["test"]}`))
		assert.Equal(`Key: 'xUserRoles.Value' Error:Field validation for 'Value' failed on the 'xUserRoles' tag`, err.Error())
	})
}
