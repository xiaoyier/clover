package setting

import (
	"clover/pkg/log"
	"flag"

	"github.com/spf13/viper"
)

type AppConfig struct {
	EndPoint  string
	TimeZone  string
	MachineID uint16
	Mysql     MysqlConf
	Redis     RedisConf
}

type MysqlConf struct {
	Host            string
	User            string
	Passwd          string
	DB              string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxConnLifeTime int
	DebugMode       bool
}

type RedisConf struct {
	Host     string
	DB       int
	PoolSize int
}

var confFile = flag.String("conf_file", "./config/develop.yaml", "the gloabl config file for application")

var appConf AppConfig

func InitAppSettings() {

	var err error
	viper.SetConfigFile(*confFile)
	err = viper.ReadInConfig()
	if err != nil {
		log.WithCategory("setting").WithError(err).Error("InitAppSettings: read config failed")
		panic(err)
	}

	err = viper.Unmarshal(&appConf)
	if err != nil {
		log.WithCategory("setting").WithError(err).Error("InitAppSettings: unmarshal config failed")
		panic(err)
	}
}

func GetAppConfig() *AppConfig {
	return &appConf
}

func GetTimeZone() string {
	return appConf.TimeZone
}

func GetEndpoint() string {
	return appConf.EndPoint
}

func GetMachineID() uint16 {
	return appConf.MachineID
}

func GetMysqlConfig() *MysqlConf {
	return &appConf.Mysql
}

func GetRedisConfig() *RedisConf {
	return &appConf.Redis
}
