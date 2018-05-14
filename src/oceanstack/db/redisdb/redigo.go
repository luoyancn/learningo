package redisdb

import (
	"oceanstack/conf"
	"oceanstack/logging"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var once sync.Once
var REDIS_POOL redis.Pool

func InitRedisConnection() {
	var err error
	once.Do(func() {
		REDIS_POOL = redis.Pool{
			MaxIdle:         conf.REDIS_MAX_IDLE,
			MaxActive:       conf.REDIS_MAX_ACTIVE,
			MaxConnLifetime: conf.REDIS_MAX_CONN_LIFETIME,
			IdleTimeout:     conf.REDIS_IDLE_TIMEOUT,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", conf.REDIS_CONNECTION)
				if nil != err {
					logging.LOG.Fatalf(
						"Cannot connect to redis server:%s\n", err)
				}
				if _, err = conn.Do(
					"SELECT", conf.REDIS_DATABASE); nil != err {
					conn.Close()
					logging.LOG.Fatalf(
						"Cannot select the redis database:%s\n", err)
				}
				return conn, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < time.Minute {
					return nil
				}
				_, err = c.Do("PING")
				return err
			},
		}
	})
}
