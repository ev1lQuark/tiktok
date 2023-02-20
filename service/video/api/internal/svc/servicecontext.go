package svc

import (
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/commentclient"
	"github.com/ev1lQuark/tiktok/service/like/rpc/likeclient"
	"github.com/ev1lQuark/tiktok/service/user/rpc/userclient"
	"github.com/ev1lQuark/tiktok/service/video/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/video/query"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Query       *query.Query
	MinioClient *minio.Client
	UserRpc     userclient.User
	CommentRpc  commentclient.Comment
	LikeRpc     likeclient.Like
}

func NewServiceContext(c config.Config) *ServiceContext {

	mc, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(c.Minio.AccessKey, c.Minio.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:      c,
		Query:       query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
		MinioClient: mc,
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
