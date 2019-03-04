package helper

import (
	"net/http"
	"net/url"
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/dusk"
	"github.com/vicanso/hes"

	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/util"

	jsoniter "github.com/json-iterator/go"

	"go.uber.org/zap"
)

var (
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
	logger = log.Default()
	// DefaultHTTPClient default http client
	DefaultHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:           100,
			IdleConnTimeout:        90 * time.Second,
			TLSHandshakeTimeout:    5 * time.Second,
			ExpectContinueTimeout:  1 * time.Second,
			MaxResponseHeaderBytes: 5 * 1024,
		},
	}
)

const (
	errCategoryHTTPRequest = "http-request"
	contextID              = "cid"

	xForwardedForHeader = "X-Forwarded-For"
)

// httpConvertResponse convert http response
func httpConvertResponse(d *dusk.Dusk) {
	statusCode := d.Response.StatusCode
	if statusCode < 400 {
		return
	}
	// 对于状态码大于400的，转化为 hes.Error
	he := &hes.Error{
		StatusCode: statusCode,
		Category:   json.Get(d.Body, "category").ToString(),
		Message:    json.Get(d.Body, "message").ToString(),
	}
	if he.Category != "" {
		he.Category = errCategoryHTTPRequest + "-" + he.Category
	} else {
		he.Category = errCategoryHTTPRequest
	}
	if he.Message == "" {
		he.Message = "unknown error"
	}
	if d.ConvertError != nil {
		d.Error = d.ConvertError(he, d)
		return
	}
	d.Error = he
}

// httpDoneEvent http请求完成的触发，用于统计、日志等输出
func httpDoneEvent(d *dusk.Dusk) {
	req := d.Request
	resp := d.Response
	err := d.Error
	uri := req.URL.RequestURI()
	stats := d.GetTimelineStats()
	use := ""
	if stats != nil {
		use = stats.Total.String()
	}
	statusCode := 0
	if err != nil {
		he, ok := err.(*hes.Error)
		if ok {
			statusCode = he.StatusCode
		}
	}
	if resp != nil {
		statusCode = resp.StatusCode
	}
	cid := ""
	cidValue := d.GetValue(contextID)
	if cidValue != nil {
		cid = cidValue.(string)
	}

	// TODO 是否将POST参数也记录（有可能会有敏感信息）
	// TODO 是否将响应数据输出（有可能敏感信息以及数据量较大），或者写入缓存数据库，保存较短时间方便排查
	if resp == nil || err != nil {
		logger.Error("http request fail",
			zap.String("method", req.Method),
			zap.String("host", req.Host),
			zap.String("uri", uri),
			zap.String("cid", cid),
			zap.Int("status", statusCode),
			zap.String("use", use),
			zap.Error(err),
		)
		return
	}
	logger.Info("http request done",
		zap.String("method", req.Method),
		zap.String("host", req.Host),
		zap.String("uri", uri),
		zap.String("cid", cid),
		zap.Int("status", statusCode),
		zap.String("use", use),
	)
}

// httpErrorConvert convert http error
func httpErrorConvert(err error, d *dusk.Dusk) error {
	he, ok := err.(*hes.Error)
	resp := d.Response
	req := d.Request
	if !ok {
		he = hes.NewWithError(err)
		statusCode := http.StatusInternalServerError
		if resp != nil {
			statusCode = resp.StatusCode
		}
		if ue, ok := err.(*url.Error); ok {
			// 请求超时中断
			if ue.Timeout() {
				statusCode = http.StatusRequestTimeout
			}
		}
		he.StatusCode = statusCode
		he.Category = errCategoryHTTPRequest
	}
	// 仅在测试中输出请求 url至 hes中（避免将重要信息输出）
	if !util.IsProduction() {
		extra := he.Extra
		if extra == nil {
			extra = make(map[string]interface{})
		}
		url := req.URL
		extra["uri"] = url.RequestURI()
		extra["host"] = url.Host
		extra["method"] = req.Method
		he.Extra = extra
	}
	return he
}

// NewRequest new request
func NewRequest() *dusk.Dusk {
	d := dusk.New()
	d.Client = DefaultHTTPClient
	d.EnableTimelineTrace = true
	d.ConvertError = httpErrorConvert
	d.On(dusk.EventResponse, httpConvertResponse)
	d.On(dusk.EventDone, httpDoneEvent)
	return d
}

// NewRequestWithContext 使用context 初始化 request
func NewRequestWithContext(c *cod.Context) *dusk.Dusk {
	d := NewRequest()
	d.SetValue(contextID, c.ID)
	d.On(dusk.EventRequest, func(d *dusk.Dusk) {
		// 设置x-forwarded-for
		v := c.GetRequestHeader(xForwardedForHeader)
		if v == "" {
			v = c.RealIP()
		}
		d.Request.Header.Set(xForwardedForHeader, v)
	})
	return d
}
