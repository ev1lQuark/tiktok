package svc

import (
	"github.com/ev1lQuark/tiktok/service/like/query"
	"github.com/ev1lQuark/tiktok/service/like/rpc/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Query  *query.Query
	Redis  *redis.Client
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
		Redis:  redis.NewClient(&redis.Options{Addr: c.CustomRedis.Addr, DB: c.CustomRedis.DB}),
	}
}
