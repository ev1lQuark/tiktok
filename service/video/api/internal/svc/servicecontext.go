package svc

import (
	"github.com/ev1lQuark/tiktok/service/video/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/video/query"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	Query       *query.Query
	MinioClient *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	query := query.Use(db)

	mc, err := minio.New(c.Minio.Endpoint, &minio.Options{})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:      c,
		Query:       query,
		MinioClient: mc,
	}
}
