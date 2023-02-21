package svc

import (
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/user/query"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Query    *query.Query
	LikeRpc  likeclient.Like
	VideoRpc videoclient.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		Query:    query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		LikeRpc:  likeclient.NewLike(zrpc.MustNewClient(c.LikeRpc)),
		VideoRpc: videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
	}
}
