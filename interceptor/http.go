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

package interceptor

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/vicanso/forest/asset"
	"github.com/vicanso/go-axios"
	"go.uber.org/atomic"
)

type httpInterceptors struct {
	scripts    *atomic.Value
	baseScript string
}
type httpRequest struct {
	// 请求地址
	URI string `json:"uri"`
	// HTTP Cookies
	Cookies map[string]string `json:"cookies"`
	// IP地址
	IP string `json:"ip"`
	// 路由参数
	Params map[string]string `json:"params"`
	// 是否有修改路由params
	ModifiedParams bool `json:"modifiedParams"`
	// 查询的query
	Query map[string]string `json:"query"`
	// 是否有修改query
	ModifiedQuery   bool `json:"modifiedQuery"`
	isURLValuesBody bool
	// POST的参数
	Body map[string]any `json:"body"`
	// 是否有修改body
	ModifiedBody bool `json:"modifiedBody"`
}

// http服务响应数据
type httpResponse struct {
	// 响应状态码
	Status int `json:"status"`
	// 响应头
	Header map[string]string `json:"header"`
	// 响应数据
	Body map[string]any `json:"body"`
}

func newHTTPResponse() *httpResponse {
	return &httpResponse{
		Header: make(map[string]string),
		Body:   make(map[string]any),
	}
}

func (resp *httpResponse) ToAxiosResponse() *axios.Response {
	data := &axios.Response{
		Status: resp.Status,
	}
	header := make(url.Values)
	for k, v := range resp.Header {
		header.Set(k, v)
	}
	data.Data, _ = json.Marshal(resp.Body)
	return data
}

type httpInterceptor struct {
	Before func() (*httpResponse, error)
	After  func() (*httpResponse, error)
}

func newHTTPInterceptors() *httpInterceptors {
	script, _ := asset.GetFS().ReadFile("http_interceptor.js")
	return &httpInterceptors{
		scripts:    &atomic.Value{},
		baseScript: string(script),
	}
}
func getScripts[T *httpServerInterceptorScript | *httpRequestInterceptorScript](intercetpros *httpInterceptors) map[string]T {
	value := intercetpros.scripts.Load()
	if value == nil {
		return nil
	}
	scripts, ok := value.(map[string]T)
	if !ok {
		return nil
	}
	return scripts
}

func getScript[T *httpServerInterceptorScript | *httpRequestInterceptorScript](intercetpros *httpInterceptors, router string) T {
	scripts := getScripts[T](intercetpros)
	if scripts == nil {
		return nil
	}
	return scripts[router]
}

func (interceptors *httpInterceptors) script(script string) string {
	return fmt.Sprintf(`%s;(function() {
		%s
	})();
	`, interceptors.baseScript, script)
}
