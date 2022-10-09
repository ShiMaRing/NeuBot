package main

import (
	"NeuBot/api"
	"NeuBot/configs"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	//初始化配置文件
	configs.ConfigInit()
}

//程序入口
func main() {
	router := gin.Default()
	router.POST("/", api.GetMsg)
	if err := router.Run(fmt.Sprintf(":%d", configs.BotConf.Port)); err != nil {
		log.Fatalln(err)
	}
}
