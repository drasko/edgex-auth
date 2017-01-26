package services

import "github.com/garyburd/redigo/redis"

const (
	maxIdle   = 10
	maxActive = 50
)

var cache *redis.Pool

// StartCaching creates new redis connection pool with fixed number of allowed
// inactive connections.
func StartCaching(url string) {
	cache = &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

// StopCaching terminates the redis pool.
func StopCaching() error {
	return cache.Close()
}
