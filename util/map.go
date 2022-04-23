// Copyright 2022 tree xie
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
	"net/url"

	"github.com/spf13/cast"
)

// URLValuesToMap convert url values to map[string]string
func URLValuesToMap(values url.Values) map[string]string {
	m := make(map[string]string)
	for k, v := range values {
		m[k] = v[0]
	}
	return m
}

// MapToURLValues convert map[string]string to url values
func MapToURLValues(m map[string]string) url.Values {
	values := make(url.Values)
	for k, v := range m {
		values.Set(k, v)
	}
	return values
}

// MapAnyToURLValues convert map[string]any to url values
func MapAnyToURLValues(m map[string]any) url.Values {
	values := make(url.Values)
	for k, v := range m {
		values.Set(k, cast.ToString(v))
	}
	return values
}
