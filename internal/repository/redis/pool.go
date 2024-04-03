package repository

import "github.com/gomodule/redigo/redis"

func NewPool(protocol, address, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   100,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(protocol, address, redis.DialPassword(password))
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}
