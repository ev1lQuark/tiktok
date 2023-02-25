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
		Addr string
		DB   int
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc    zrpc.RpcClientConf
	CommentRpc zrpc.RpcClientConf
	VideoRpc   zrpc.RpcClientConf
}
