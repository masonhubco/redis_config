package redisconfig

import (
	"strconv"

	"github.com/gobuffalo/envy"
	"github.com/gomodule/redigo/redis"
)

//RedisConfig stores the information necessary to connect to Redis
type RedisConfig struct {
	maxActive      int
	maxIdle        int
	maxConcurrency int
	name           string
	environment    string
	pool           *redis.Pool
}

//Initiate stores local copies of the envy variables for retrieval
func Initiate() *RedisConfig {
	r := &RedisConfig{}
	if maxActive, err := strconv.Atoi(envy.Get("REDIS_MAX_ACTIVE", "150")); err != nil {
		r.maxActive = 150
	} else {
		r.maxActive = maxActive
	}

	if maxIdle, err := strconv.Atoi(envy.Get("REDIS_MAX_IDLE", "75")); err != nil {
		r.maxIdle = 75
	} else {
		r.maxIdle = maxIdle
	}

	if maxConcurrency, err := strconv.Atoi(envy.Get("REDIS_MAX_CONCURRENCY", "75")); err == nil {
		r.maxConcurrency = 75
	} else {
		r.maxConcurrency = maxConcurrency
	}

	r.environment = envy.Get("GO_ENV", "development")
	if r.environment == "test" { //covers the lack of .env file in test environment
		r.name = "test"
	}
	r.name = envy.Get("REDIS_QUEUE_NAME", "staging") //defaulting to staging so any forgotten config won't pollute or destroy data in sandbox or production

	return r
}

//RedisPool returns the pool as defined in by the envy variable
func (r *RedisConfig) RedisPool() *redis.Pool {
	return &redis.Pool{
		MaxActive: r.maxActive,
		MaxIdle:   r.maxIdle,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			if r.environment == "development" || r.environment == "test" {
				return redis.Dial("tcp", ":6379")
			}
			return redis.DialURL(envy.Get("REDIS_URL", ""))
		},
	}
}

//MaxConcurrency returns the max concurrency set in initiation
func (r *RedisConfig) MaxConcurrency() int {
	return r.maxConcurrency
}

//RedisName returns the name of the Redis DB set in initiation
func (r *RedisConfig) RedisName() string {
	return r.name
}
