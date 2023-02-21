package svc

import (
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/commentclient"
	"github.com/ev1lQuark/tiktok/service/like/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/like/query"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/rpc/videoclient"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	Query      *query.Query
	UserRpc    userclient.User
	CommentRpc commentclient.Comment
	VideoRpc   videoclient.Video
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:     c,
		Query:      query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		UserRpc:    userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		CommentRpc: commentclient.NewComment(zrpc.MustNewClient(c.CommentRpc)),
		VideoRpc:   videoclient.NewVideo(zrpc.MustNewClient(c.VideoRpc)),
	}
}
