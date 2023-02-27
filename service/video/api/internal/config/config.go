package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Video struct {
		NumberLimit int
	}
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr   string
		DB     int
	}
	Auth       struct {
		AccessSecret string
		AccessExpire int64
	}
	Minio struct {
		Endpoint    string
		VideoBucket string
		ImageBucket string
		AccessKey   string
		SecretKey   string
	}
	UserRpc    zrpc.RpcClientConf
	CommentRpc zrpc.RpcClientConf
	LikeRpc    zrpc.RpcClientConf
	ContinuedTime int64
}
