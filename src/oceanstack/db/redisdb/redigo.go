package redisdb

import (
	"oceanstack/conf"
	"oceanstack/logging"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

var once sync.Once
var redis_pool redis.Pool

func InitRedisConnection() {
	var err error
	once.Do(func() {
		redis_pool = redis.Pool{
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

func TokenIssue(userinfo string) (string, error) {
	token := uuid.NewV4().String()
	conn := redis_pool.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", token, conf.REDIS_EXPIRE, userinfo)
	if nil != err {
		return "", err
	}
	return token, nil
}

func ValidToken(token string) bool {
	conn := redis_pool.Get()
	defer conn.Close()
	reply, err := conn.Do("GET", token)
	if nil != err {
		logging.LOG.Errorf("ERROR:%v\n", err)
		return false
	}
	if _, err = redis.String(reply, err); nil != err {
		logging.LOG.Errorf("ERROR:%v\n", err)
		return false
	}
	return true
}
