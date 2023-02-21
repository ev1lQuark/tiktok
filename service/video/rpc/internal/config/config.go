package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Minio struct {
		Endpoint    string
		VideoBucket string
		ImageBucket string
		AccessKey   string
		SecretKey   string
	}
}
