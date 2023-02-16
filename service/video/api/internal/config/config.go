package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Video struct {
		NumberLimit int
	}
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	Minio struct {
		Endpoint    string
		VideoBucket string
		ImageBucket string
	}
}
