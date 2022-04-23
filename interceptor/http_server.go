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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dop251/goja"
	"github.com/samber/lo"
	"github.com/vicanso/elton"
)

type httpServerInterceptorScript struct {
	Router string `json:"router"`
	Before string `json:"before"`
	After  string `json:"after"`
	IP     string `json:"ip"`
	Cookie string `json:"cookie"`
}

func UpdateHTTPServer(arr []string) {
	scripts := make(map[string]*httpServerInterceptorScript)
	for _, item := range arr {
		script := httpServerInterceptorScript{}
		_ = json.Unmarshal([]byte(item), &script)
		router := script.Router
		// 只根据是否有router来判断是否正确
		if router == "" {
			continue
		}
		scripts[router] = &script
	}
	currentHTTPServerInterceptors.scripts.Store(scripts)
}

var currentHTTPServerInterceptors = newHTTPInterceptors()

func newHTTPServerRequest(c *elton.Context) *httpRequest {
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		query[key] = values[0]
	}
	cookies := make(map[string]string)
	for _, cookie := range c.Request.Cookies() {
		if cookie.Value == "" {
			continue
		}
		cookies[cookie.Name] = cookie.Value
	}

	req := &httpRequest{
		URI:     c.Request.URL.RequestURI(),
		Params:  c.Params.ToMap(),
		Query:   query,
		Body:    make(map[string]any),
		IP:      c.ClientIP(),
		Cookies: cookies,
	}
	if len(c.RequestBody) != 0 {
		_ = json.Unmarshal(c.RequestBody, &req.Body)
	}
	return req
}

// 设置响应数据
func (resp *httpResponse) SetResponse(c *elton.Context) {
	c.StatusCode = resp.Status
	for k, v := range resp.Header {
		c.SetHeader(k, v)
	}
	buf, _ := json.Marshal(resp.Body)
	c.BodyBuffer = bytes.NewBuffer(buf)
}

func NewHTTPServer(c *elton.Context) (*httpInterceptor, error) {
	script := getScript[*httpServerInterceptorScript](currentHTTPServerInterceptors, c.Request.Method+" "+c.Route)
	if script == nil {
		return nil, nil
	}
	// 如果指定了IP但客户非此IP
	if script.IP != "" && c.ClientIP() != script.IP {
		return nil, nil
	}
	if script.Cookie != "" {
		arr := strings.Split(script.Cookie, ";")
		cookies := lo.Map[*http.Cookie, string](c.Request.Cookies(), func(cookie *http.Cookie, _ int) string {
			return fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
		})
		for _, item := range arr {
			// 如果不包括，则直接返回
			if !lo.Contains[string](cookies, item) {
				return nil, nil
			}
		}
	}

	vm := goja.New()

	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	req := newHTTPServerRequest(c)
	err := vm.Set("req", req)
	if err != nil {
		return nil, err
	}
	resp := newHTTPResponse()
	err = vm.Set("resp", resp)
	if err != nil {
		return nil, err
	}

	inter := &httpInterceptor{
		Before: func() (*httpResponse, error) {
			if script.Before == "" {
				return nil, nil
			}
			_, err := vm.RunString(currentHTTPServerInterceptors.script(script.Before))
			if err != nil {
				return nil, err
			}
			// 如果有修改query，则重新生成
			if req.ModifiedQuery {
				query := c.Request.URL.Query()
				for k, v := range req.Query {
					query.Set(k, v)
				}
				c.Request.URL.RawQuery = query.Encode()
				c.Request.RequestURI = c.Request.URL.RequestURI()
			}
			// 如果有修改body
			if req.ModifiedBody {
				c.RequestBody, _ = json.Marshal(req.Body)
			}
			return resp, nil
		},
		After: func() (*httpResponse, error) {
			if script.After == "" {
				return nil, nil
			}
			_ = json.Unmarshal(c.BodyBuffer.Bytes(), &resp.Body)
			_, err := vm.RunString(currentHTTPServerInterceptors.script(script.After))
			if err != nil {
				return nil, err
			}
			return resp, nil
		},
	}
	return inter, nil
}
