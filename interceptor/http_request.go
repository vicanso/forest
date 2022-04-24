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

package interceptor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dop251/goja"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/go-axios"
)

type httpRequestInterceptorScript struct {
	Service string `json:"service"`
	Method  string `json:"method"`
	Route   string `json:"route"`
	Before  string `json:"before"`
	After   string `json:"after"`
}

var currentHTTPRequestInterceptors = newHTTPInterceptors()

func UpdateHTTPRequest(arr []string) {
	scripts := make(map[string]*httpRequestInterceptorScript)
	for _, item := range arr {
		script := httpRequestInterceptorScript{}
		_ = json.Unmarshal([]byte(item), &script)
		if script.Method == "" ||
			script.Route == "" {
			continue
		}
		router := fmt.Sprintf("%s %s %s", script.Service, script.Method, script.Route)
		scripts[router] = &script
	}
	currentHTTPRequestInterceptors.scripts.Store(scripts)
}

func newHTTPRequest(config *axios.Config) *httpRequest {

	req := &httpRequest{
		Params: config.Params,
		// 简单的处理为query不允许相同的参数
		// 因此只取第一个
		Query: util.URLValuesToMap(config.Query),
		Body:  make(map[string]any),
	}
	if config.Body != nil {
		switch data := config.Body.(type) {
		case url.Values:
			req.isURLValuesBody = true
			for key, values := range data {
				req.Body[key] = values
			}
		default:
			// 出错均忽略
			buf, _ := json.Marshal(data)
			_ = json.Unmarshal(buf, &req.Body)
		}
	}
	return req
}

func (req *httpRequest) getBody() any {
	if req.isURLValuesBody {
		return util.MapAnyToURLValues(req.Body)
	}
	return req.Body
}

func (req *httpRequest) getQuery() url.Values {
	return util.MapToURLValues(req.Query)
}

func NewHTTPRequest(service string, conf *axios.Config) (*httpInterceptor, error) {
	if conf == nil {
		return nil, nil
	}
	method := conf.Method
	if method == "" {
		method = http.MethodGet
	}
	route := conf.Route
	if route == "" {
		route = conf.URL
	}
	router := fmt.Sprintf("%s %s %s", service, method, route)
	script := getScript[*httpRequestInterceptorScript](currentHTTPRequestInterceptors, router)
	if script == nil {
		return nil, nil
	}
	req := newHTTPRequest(conf)

	vm := goja.New()

	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

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
			_, err := vm.RunString(currentHTTPRequestInterceptors.script(script.Before))
			if err != nil {
				return nil, err
			}
			if req.ModifiedQuery {
				conf.Query = req.getQuery()
			}
			if req.ModifiedParams {
				conf.Params = req.Params
			}
			if req.ModifiedBody {
				conf.Body = req.getBody()
			}
			if resp.Status != 0 {
				conf.Response = resp.ToAxiosResponse()
			}
			return nil, nil
		},
		After: func() (*httpResponse, error) {
			if script.After == "" {
				return nil, nil
			}
			err := json.Unmarshal(conf.Response.Data, &resp.Body)
			if err != nil {
				return nil, err
			}
			err = vm.Set("resp", resp)
			if err != nil {
				return nil, err
			}
			_, err = vm.RunString(currentHTTPRequestInterceptors.script(script.After))
			if err != nil {
				return nil, err
			}
			if resp.Status != 0 {
				conf.Response.Data, _ = json.Marshal(resp.Body)
			}
			return nil, nil
		},
	}
	return inter, nil
}
