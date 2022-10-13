package handler

import (
	"NeuBot/model"
	"log"
)

const (
	Menu = "1.显示菜单\n2.绑定学号\n3.订阅或取消每日课程提醒\n4.订阅或取消每日健康上报\n5.解除绑定\n6.反馈\n输入对应序号使用命令"
)

// NoticeHandler 对于新添加的好友的处理逻辑
type NoticeHandler struct {
}

func NewNoticeHandler() *NoticeHandler {
	return &NoticeHandler{}
}

// Greet 给新好友打招呼
func (h *NoticeHandler) Greet(req *model.NoticeReq) {
	_, err := ReplyMsg(req.UserId, "欢迎使用NEU-BOT，发送序号指令使用，目前不支持南湖同学", false)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = ReplyMsg(req.UserId, Menu, false)
	if err != nil {
		log.Fatalln(err)
	}
}
