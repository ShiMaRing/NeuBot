package main

import (
	"NeuBot/api"
	"NeuBot/configs"
	"NeuBot/handler"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func init() {
	//初始化配置文件
	configs.ConfigInit()
	handler.Init()
}

//程序入口
func main() {
	//启动定时任务
	err := handler.StartSchedule()
	if err != nil {
		log.Fatalln("启动定时任务失败", err)
	}
	router := gin.Default()
	router.POST("/", api.GetMsg)
	if err := router.Run(fmt.Sprintf(":%d", configs.BotConf.Port)); err != nil {
		log.Fatalln(err)
	}
}
