package utils

import (
	"github.com/go-redis/redis/v8"
)

// RedisConf holds redis config
type RedisConf struct {
	Addresses map[string]string
	// endpoints array for universal client
	Endpoints []string
	Username  string
	Password  string
	DB        int
}

// NewRedisRing returns a new redis ring
func NewRedisRing(conf *RedisConf) *redis.Ring {
	options := &redis.RingOptions{
		Addrs: conf.Addresses,
		NewClient: func(name string, opt *redis.Options) *redis.Client {
			opt.Username = conf.Username
			opt.Password = conf.Password
			opt.DB = conf.DB
			return redis.NewClient(opt)
		},
	}
	return redis.NewRing(options)
}
