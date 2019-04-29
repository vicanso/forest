package helper

import (
	"net/http"
	"net/url"
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/dusk"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/util"
	"github.com/vicanso/hes"
	"go.uber.org/zap"

	jsoniter "github.com/json-iterator/go"
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

func init() {
	dusk.AddRequestListener(func(_ *http.Request, d *dusk.Dusk) (newReq *http.Request, newErr error) {
		if d.GetClient() == nil {
			d.SetClient(DefaultHTTPClient)
		}
		d.EnableTrace()
		return
	}, dusk.EventTypeBefore)
	dusk.AddResponseListener(httpConvertResponse, dusk.EventTypeAfter)
	dusk.AddDoneListener(httpDoneEvent)
	dusk.AddErrorListener(httpErrorConvert)
}

// httpConvertResponse convert http response
func httpConvertResponse(resp *http.Response, d *dusk.Dusk) (newResp *http.Response, newErr error) {
	statusCode := resp.StatusCode
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

	return nil, he
}

// httpDoneEvent http请求完成的触发，用于统计、日志等输出
func httpDoneEvent(d *dusk.Dusk) error {
	req := d.Request
	resp := d.Response
	err := d.Err
	uri := req.URL.RequestURI()
	unescapeURI, _ := url.QueryUnescape(uri)
	if unescapeURI != "" {
		uri = unescapeURI
	}
	ht := d.GetHTTPTrace()
	use := ""
	if ht != nil {
		use = ht.Stats().Total.String()
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
	fields := make([]zap.Field, 0, 6)
	fields = append(fields, zap.String("host", req.Host))
	fields = append(fields, zap.String("method", req.Method))
	fields = append(fields, zap.String("path", d.GetPath()))
	fields = append(fields, zap.String("uri", uri))
	fields = append(fields, zap.String("cid", cid))
	fields = append(fields, zap.Int("status", statusCode))
	fields = append(fields, zap.String("use", use))
	if resp == nil || err != nil {
		fields = append(fields, zap.Error(err))
		logger.Error("http request fail", fields...)
		return nil
	}
	logger.Info("http request done", fields...)
	return nil
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

// AttachContext attach dusk with context
func AttachContext(d *dusk.Dusk, c *cod.Context) {
	if c != nil {
		if c.ID != "" {
			d.SetValue(contextID, c.ID)
		}
		// 设置x-forwarded-for
		v := c.GetRequestHeader(xForwardedForHeader)
		if v == "" {
			v = c.RealIP()
		}
		d.Set(xForwardedForHeader, v)
	}
}
