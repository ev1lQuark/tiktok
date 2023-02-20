package svc

import (
	"github.com/ev1lQuark/tiktok/common/db"
	"github.com/ev1lQuark/tiktok/service/video/query"
	"github.com/ev1lQuark/tiktok/service/video/rpc/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Query  *query.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Query:  query.Use(db.NewMysqlConn(c.Mysql.DataSource, &gorm.Config{})),
	}
}
