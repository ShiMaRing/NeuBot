package configs

import (
	"github.com/spf13/viper"
	"log"
)

//读取配置文件

// RedisConf redis 配置
type redisConf struct {
	Address  string
	Port     int
	Password string
}

// MysqlConf  mysql配置
type mysqlConf struct {
	Address  string
	Port     int
	Database string
	User     string
	Password string
}

// CQHttpConf cqhttp配置
type cQHttpConf struct {
	Address string
	Port    int
}

type botConfig struct {
	Port     int
	MasterId int64
}

// ConfigInit 读取Config
func ConfigInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	redisMap := viper.GetStringMap("redis")
	enable := redisMap["enable"].(bool)

	if enable {
		RedisConf = &redisConf{
			Address:  redisMap["address"].(string),
			Port:     redisMap["port"].(int),
			Password: redisMap["password"].(string),
		}
	} else {
		RedisConf = &redisConf{}
	}

	mysqlMap := viper.GetStringMap("mysql")
	MysqlConf = &mysqlConf{
		Address:  mysqlMap["address"].(string),
		Port:     mysqlMap["port"].(int),
		Database: mysqlMap["database"].(string),
		User:     mysqlMap["user"].(string),
		Password: mysqlMap["password"].(string),
	}

	cqhttpMap := viper.GetStringMap("cqhttp")
	CqhttpConf = &cQHttpConf{
		Address: cqhttpMap["address"].(string),
		Port:    cqhttpMap["port"].(int),
	}

	botMap := viper.GetStringMap("bot")
	BotConf = &botConfig{
		Port:     botMap["port"].(int),
		MasterId: int64(botMap["master"].(int)),
	}

}
