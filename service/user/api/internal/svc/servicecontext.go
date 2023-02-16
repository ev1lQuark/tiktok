package svc

import (
	"github.com/ev1lQuark/tiktok/service/user/api/internal/config"
	"github.com/ev1lQuark/tiktok/service/user/query"
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
