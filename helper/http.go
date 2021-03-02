// Copyright 2020 tree xie
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

package helper

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/tidwall/gjson"
	"github.com/vicanso/elton"
	"github.com/vicanso/forest/cs"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

const (
	httpErrCategoryDNS     = "dns"
	httpErrCategoryTimeout = "timeout"
	httpErrCategoryAddr    = "addr"
	httpErrCategoryAborted = "aborted"
	httpErrCategoryRefused = "refused"
	httpErrCategoryReset   = "reset"
)

func newHTTPOnDone(serviceName string) axios.OnDone {
	return func(conf *axios.Config, resp *axios.Response, err error) {
		ht := conf.HTTPTrace

		reused := false
		addr := ""
		use := ""
		status := -1
		if resp != nil {
			status = conf.Response.Status
		}

		result := cs.ResultSuccess
		if err != nil {
			result = cs.ResultFail
		}
		tags := map[string]string{
			cs.TagService: serviceName,
			cs.TagRoute:   conf.Route,
			cs.TagMethod:  conf.Method,
			cs.TagResult:  strconv.Itoa(result),
		}
		fields := map[string]interface{}{
			cs.FieldURI:    conf.GetURL(),
			cs.FieldStatus: status,
		}
		if ht != nil {
			reused = ht.Reused
			addr = ht.Addr
			timelineStats := ht.Stats()
			use = timelineStats.String()
			fields[cs.FieldReused] = reused
			fields[cs.FieldUse] = int(timelineStats.Total.Milliseconds())
			dns := timelineStats.DNSLookup.Milliseconds()
			if dns != 0 {
				fields[cs.FieldDNSUse] = int(dns)
			}
			tcp := timelineStats.TCPConnection.Milliseconds()
			if tcp != 0 {
				fields[cs.FieldTCPUse] = int(tcp)
			}
			tls := timelineStats.TLSHandshake.Milliseconds()
			if tls != 0 {
				fields[cs.FieldTLSUse] = int(tls)
			}
			serverProcessing := timelineStats.ServerProcessing.Milliseconds()
			if serverProcessing != 0 {
				fields[cs.FieldProcessingUse] = int(serverProcessing)
			}
			contentTransfer := timelineStats.ContentTransfer.Milliseconds()
			if contentTransfer != 0 {
				fields[cs.FieldTransferUse] = int(contentTransfer)
			}
		}
		fields[cs.FieldAddr] = addr
		message := ""
		if err != nil {
			he := hes.Wrap(err)
			message = he.Error()
			fields[cs.FieldError] = message
			errCategory := he.Category
			if errCategory != "" {
				fields[cs.FieldErrCategory] = errCategory
			}
			if he.Exception {
				fields[cs.FieldException] = true
			}
		}
		// 输出响应数据，如果响应数据为隐私数据可不输出
		var data interface{}
		size := -1
		if resp != nil {
			data = resp.UnmarshalData
			size = len(resp.Data)
		}
		// 由于http请求是较频繁的操作，因此判断是否启用debug再输出
		if log.DebugEnabled() {
			log.Default().Debug(conf.CURL())
		}
		log.Default().Info("http request stats",
			zap.String("service", serviceName),
			zap.String("method", conf.Method),
			zap.String("route", conf.Route),
			zap.String("url", conf.GetURL()),
			zap.Any("params", conf.Params),
			zap.Any("query", conf.Query),
			zap.Any("data", data),
			zap.Int("size", size),
			zap.Int("status", status),
			zap.String("addr", addr),
			zap.Bool("reused", reused),
			zap.String("use", use),
			zap.String("error", message),
		)
		GetInfluxSrv().Write(cs.MeasurementHTTPRequest, tags, fields)
	}
}

// newHTTPConvertResponseToError 将http响应码为>=400的转换为出错
func newHTTPConvertResponseToError(serviceName string) axios.ResponseInterceptor {
	return func(resp *axios.Response) (err error) {
		if resp.Status >= 400 {
			message := gjson.GetBytes(resp.Data, "message").String()
			exception := false
			if message == "" {
				message = string(resp.Data)
				// 如果出错响应不符合，则认为是异常响应
				exception = true
			}
			he := hes.NewWithStatusCode(message, resp.Status)
			he.Exception = exception
			return he
		}
		return
	}
}

// getHTTPErrorCategory 获取出错的类型，主要分类DNS错误，addr错误以及一些系统调用的异常
func getHTTPErrorCategory(err error) string {

	netErr, ok := err.(net.Error)
	if ok && netErr.Timeout() {
		return httpErrCategoryTimeout
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return httpErrCategoryDNS
	}
	var addrErr *net.AddrError
	if errors.As(err, &addrErr) {
		return httpErrCategoryAddr
	}

	opErr, ok := netErr.(*net.OpError)
	if !ok {
		return ""
	}
	switch e := opErr.Err.(type) {
	// 针对以下几种系统调用返回对应类型
	case *os.SyscallError:
		if no, ok := e.Err.(syscall.Errno); ok {
			switch no {
			case syscall.ECONNREFUSED:
				return httpErrCategoryRefused
			case syscall.ECONNABORTED:
				return httpErrCategoryAborted
			case syscall.ECONNRESET:
				return httpErrCategoryReset
			case syscall.ETIMEDOUT:
				return httpErrCategoryTimeout
			}
		}
	}

	return ""
}

// newHTTPOnError 新建error的处理函数
func newHTTPOnError(serviceName string) axios.OnError {
	return func(err error, conf *axios.Config) (newErr error) {
		code := -1
		if conf.Response != nil {
			code = conf.Response.Status
		}
		he := hes.Wrap(err)
		if code >= http.StatusBadRequest {
			he.StatusCode = code
		}
		// 如果状态码>500，认为是异常（一般nginx之类才会使用502 503之类的状态码）
		if code > http.StatusInternalServerError {
			he.Exception = true
		}
		// 如果未设置http响应码，则设置为500
		if he.StatusCode < http.StatusBadRequest {
			he.StatusCode = http.StatusInternalServerError
		}

		if he.Extra == nil {
			he.Extra = make(map[string]interface{})
		}

		// 如果为空，则通过error获取
		if he.Category == "" {
			he.Category = getHTTPErrorCategory(err)
			// 如果返回错误类型，则认为异常
			// 因为返回出错均是网络连接上的异常
			if he.Category != "" {
				he.Exception = true
			}
		}

		if !util.IsProduction() {
			he.Extra["requestRoute"] = conf.Route
			he.Extra["requestService"] = serviceName
			he.Extra["requestCURL"] = conf.CURL()
		}
		return he
	}
}

// NewHTTPInstance 新建实例
func NewHTTPInstance(serviceName, baseURL string, timeout time.Duration) *axios.Instance {
	return axios.NewInstance(&axios.InstanceConfig{
		EnableTrace: true,
		Timeout:     timeout,
		OnError:     newHTTPOnError(serviceName),
		OnDone:      newHTTPOnDone(serviceName),
		BaseURL:     baseURL,
		ResponseInterceptors: []axios.ResponseInterceptor{
			newHTTPConvertResponseToError(serviceName),
		},
	})
}

// AttachWithContext 添加context中的cid至请求的config中
func AttachWithContext(conf *axios.Config, c *elton.Context) {
	if c == nil || conf == nil {
		return
	}
	if conf.Context == nil {
		conf.Context = c.Context()
	}
}
