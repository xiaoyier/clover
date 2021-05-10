package model

import (
	"clover/pkg/log"
	"clover/setting"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db           *gorm.DB
	AutoMigrates = make([]interface{}, 0)
)

func InitDB(conf *setting.MysqlConf, timeZone string) {

	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", conf.User, conf.Passwd, conf.Host, conf.DB, timeZone)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.WithCategory("model").WithError(err).Error("InitDB: dialect failed")
		panic(err)
	}

	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(conf.MaxConnLifeTime) * time.Second)
	db.LogMode(conf.DebugMode)

	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(AutoMigrates...)
	db.AutoMigrate()
}

func CloseMysql() {
	err := db.Close()
	if err != nil {
		log.WithCategory("mysql").WithError(err).Error("CloseMysql: failed")
	}
}

func GetDB() *gorm.DB {
	return db
}

func RegisterAutoMigrates(migrates interface{}) {
	AutoMigrates = append(AutoMigrates, migrates)
}
