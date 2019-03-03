package middleware

import (
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/global"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/hes"
	"go.uber.org/zap"

	concurrentLimiter "github.com/vicanso/cod-concurrent-limiter"
)

var (
	errTooManyRequest = &hes.Error{
		StatusCode: http.StatusTooManyRequests,
		Message:    "too many request",
		Category:   errLimitCategory,
	}
	errTooFrequently = &hes.Error{
		StatusCode: http.StatusBadRequest,
		Message:    "request to frequently",
		Category:   errLimitCategory,
	}
)

const (
	defaultRequestLimit      = 2048
	concurrentLimitKeyPrefix = "mid-concurrent-limit"
	ipLimitKeyPrefix         = "mid-ip-limit"
	errLimitCategory         = "request-limit"
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

// createConcurrentLimitLock 创建并发限制的lock函数
func createConcurrentLimitLock(prefix string, ttl time.Duration, withDone bool) concurrentLimiter.Lock {
	return func(key string, c *cod.Context) (success bool, done func(), err error) {
		k := concurrentLimitKeyPrefix + "-" + prefix + "-" + key
		done = nil
		if withDone {
			success, redisDone, err := service.LockWithDone(k, ttl)
			done = func() {
				err := redisDone()
				if err != nil {
					log.Default().Error("redis done fail",
						zap.String("key", k),
						zap.Error(err),
					)
				}
			}
			return success, done, err
		}
		success, err = service.Lock(k, ttl)
		return
	}
}

// NewConcurrentLimit create a concurrent limit
func NewConcurrentLimit(keys []string, ttl time.Duration, prefix string) cod.Handler {
	return concurrentLimiter.New(concurrentLimiter.Config{
		Lock: createConcurrentLimitLock(prefix, ttl, false),
		Keys: keys,
	})
}

// NewIPLimit create a limit middleware by ip address
func NewIPLimit(maxCount int64, ttl time.Duration, prefix string) cod.Handler {
	return func(c *cod.Context) (err error) {
		key := ipLimitKeyPrefix + "-" + prefix + "-" + c.RealIP()
		count, err := service.IncWithTTL(key, ttl)
		if err != nil {
			return
		}
		if count > maxCount {
			err = errTooFrequently
			return
		}
		return c.Next()
	}
}
