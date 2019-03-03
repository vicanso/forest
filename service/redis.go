package service

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/vicanso/forest/config"
	"github.com/vicanso/forest/log"
	"github.com/vicanso/hes"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var (
	redisClient   *redis.Client
	redisOkResult = "OK"
	redisNoop     = func() error {
		return nil
	}
	errRedisNotInit = &hes.Error{
		Message:    "redis client is not init",
		Category:   "service-redis",
		StatusCode: http.StatusBadRequest,
	}
)

type (
	// Done redis done function
	Done func() error
)

func init() {
	uri := config.GetString("redis")
	if uri != "" {
		c, err := newRedisClient(uri)
		if err != nil {
			panic(err)
		}
		redisClient = c
		_, err = redisClient.Ping().Result()

		logger := log.Default()
		mask := regexp.MustCompile(`redis://:(\S+)\@`)
		str := mask.ReplaceAllString(uri, "redis://:***@")
		if err != nil {
			logger.Error("redis ping fail",
				zap.String("uri", str),
				zap.Error(err),
			)
		} else {
			logger.Info("redis ping success",
				zap.String("uri", str),
			)
		}
	}
}

// newRedisClient new client
func newRedisClient(uri string) (client *redis.Client, err error) {
	info, err := url.Parse(uri)
	if err != nil {
		return
	}
	opts := &redis.Options{
		Addr: info.Host,
	}
	query := info.Query()
	db := query.Get("db")
	if db != "" {
		opts.DB, _ = strconv.Atoi(db)
	}
	poolSize := query.Get("poolSize")
	if poolSize != "" {
		opts.PoolSize, _ = strconv.Atoi(poolSize)
	}
	opts.Password, _ = info.User.Password()
	client = redis.NewClient(opts)
	return
}

// GetRedisClient get redis client
func GetRedisClient() *redis.Client {
	return redisClient
}

// Lock lock the key for ttl seconds
func Lock(key string, ttl time.Duration) (bool, error) {
	if redisClient == nil {
		return false, errRedisNotInit
	}
	return redisClient.SetNX(key, true, ttl).Result()
}

// LockWithDone lock the key for ttl, and with done function
func LockWithDone(key string, ttl time.Duration) (bool, Done, error) {
	success, err := Lock(key, ttl)
	// 如果lock失败，则返回no op 的done function
	if err != nil || !success {
		return false, redisNoop, err
	}
	done := func() error {
		_, err := redisClient.Del(key).Result()
		return err
	}
	return true, done, nil
}

// IncWithTTL inc value with ttl
func IncWithTTL(key string, ttl time.Duration) (count int64, err error) {
	if redisClient == nil {
		return 0, errRedisNotInit
	}
	pipe := redisClient.TxPipeline()
	// 保证只有首次会设置ttl
	pipe.SetNX(key, 0, ttl)
	incr := pipe.Incr(key)
	_, err = pipe.Exec()
	if err != nil {
		return
	}
	count = incr.Val()
	return
}
