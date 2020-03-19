package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"stock/config"
)

func GormOpenDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("mysql", config.Config.MysqlDB.Dsn())
	return
}
