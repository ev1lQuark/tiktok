package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlConn(dataSource string, cfg *gorm.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dataSource), cfg)
	if err != nil {
		panic(err)
	}
	return db
}
