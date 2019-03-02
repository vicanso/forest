package controller

import (
	"time"

	"github.com/vicanso/cod"
	concurrentLimiter "github.com/vicanso/cod-concurrent-limiter"
	"github.com/vicanso/forest/log"
	m "github.com/vicanso/forest/middleware"
	"github.com/vicanso/forest/service"
	"github.com/vicanso/forest/util"
	"go.uber.org/zap"
)

var (
	logger  = log.Default()
	noQuery = m.NoQuery
	waitFor = m.WaitFor
	now     = util.NowString
	// getUserSession = service.NewUserSession
	getTrackID = util.GetTrackID

	// usesSession = createSessionMiddleware()
)

const (
	concurrentLimitKeyPrefix = "mid-concurrent-limit"
)

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
					logger.Error("redis done fail",
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

func createConcurrentLimit(keys []string, ttl time.Duration, prefix string) cod.Handler {
	return concurrentLimiter.New(concurrentLimiter.Config{
		Lock: createConcurrentLimitLock(prefix, ttl, false),
		Keys: keys,
	})
}

// // createSessionMiddleware 创建session中间件
// func createSessionMiddleware() cod.Handler {
// 	scf := config.GetSessionConfig()
// 	options := &ss.Options{
// 		// session的缓存时间，按需要设置更长的值
// 		TTL: scf.TTL,
// 		Key: scf.Key,
// 		// 用于将id与密钥生成校验串，建议配置此参数，并注意保密
// 		SignKeys: scf.SignKeys,
// 		CookieOptions: &cookies.Options{
// 			HttpOnly: true,
// 			Path:     scf.CookiePath,
// 		},
// 		IDGenerator: util.GenUlid,
// 	}
// 	// TODO create store 应该可以出错
// 	redisClient := service.GetRedisClient()
// 	createStore := func(c *cod.Context) middleware.Store {
// 		rs := ss.NewRedisStore(redisClient)
// 		rs.SetOptions(options)
// 		// 设置context
// 		rs.SetContext(c)
// 		return rs
// 	}
// 	return middleware.NewSession(middleware.SessionConfig{
// 		CreateStore: createStore,
// 	})
// }
