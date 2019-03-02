package middleware

import (
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/global"
	"github.com/vicanso/hes"
)

var (
	errTooManyRequest = &hes.Error{
		StatusCode: http.StatusTooManyRequests,
		Message:    "too many request",
		Category:   "request-limit",
	}
)

const (
	defaultRequestLimit = 2048
)

// NewLimiter create a limit middleware
func NewLimiter() cod.Handler {
	maxRequestLimit := uint32(config.GetIntDefault("requestLimit", defaultRequestLimit))
	var connectingCount uint32
	errTooManyRequest.Message += ("(" + strconv.Itoa(int(maxRequestLimit)) + ")")
	return func(c *cod.Context) (err error) {
		// 处理请求数+1/-1
		defer atomic.AddUint32(&connectingCount, ^uint32(0))
		v := atomic.AddUint32(&connectingCount, 1)
		global.SaveConnectingCount(v)
		if v >= maxRequestLimit {
			err = errTooManyRequest
			return
		}
		return c.Next()
	}
}
