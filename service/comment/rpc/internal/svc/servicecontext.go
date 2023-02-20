package svc

import (
	"github.com/ev1lQuark/tiktok/service/comment/query"
	"github.com/ev1lQuark/tiktok/service/comment/rpc/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Query  *query.Query
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	query := query.Use(db)

	return &ServiceContext{
		Config: c,
		Query:  query,
	}
}
