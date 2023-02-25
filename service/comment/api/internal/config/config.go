package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr   string
		DB     int
		ExpireTime int
	}
	RocketMQ struct {
		NameServer string
		Topic      string
		Group      string
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc  zrpc.RpcClientConf
	VideoRpc zrpc.RpcClientConf
	LikeRpc  zrpc.RpcClientConf
	DelayTime int
}
