package svc

import (
	"github.com/ev1lQuark/tiktok/service/comment/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/comment/query"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	Query    *query.Query
	VideoRpc videoclient.Video
	UserRpc  userclient.User
	LikeRpc  likeclient.Like
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	query := query.Use(db)

	return &ServiceContext{
		Config:   c,
		Query:    query,
		VideoRpc: videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
		UserRpc:  userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		LikeRpc:  likeclient.NewLike(zrpc.MustNewClient(c.LikeRpc)),
	}
}
