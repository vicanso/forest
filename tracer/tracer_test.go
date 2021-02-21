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

package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetID(t *testing.T) {
	assert := assert.New(t)
	id := getID()
	ch := make(chan struct{})
	var newID uintptr
	go func() {
		newID = getID()
		ch <- struct{}{}
	}()
	<-ch
	assert.NotEmpty(id)
	assert.NotEmpty(newID)
	assert.NotEqual(id, newID)
}

func TestTracerInfo(t *testing.T) {
	assert := assert.New(t)
	account := "test"
	traceID := "123"

	SetTracerInfo(TracerInfo{
		Account: account,
		TraceID: traceID,
	})

	info := GetTracerInfo()
	assert.Equal(account, info.Account)
	assert.Equal(traceID, info.TraceID)
}
