package svc

import (
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/video/query"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Query       *query.Query
	MinioClient *minio.Client
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
		MinioClient: mc,
		Query:       query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
	}
}
