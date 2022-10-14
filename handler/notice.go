package handler

import (
	"NeuBot/model"
	"log"
)

const (
	Menu = "1.显示菜单\n2.绑定学号\n3.订阅课程提醒\n4.取消课程提醒\n5.自动健康上报\n6.取消自动上报\n7.解绑账号\n输入对应序号使用命令\n反馈请携带前缀 feedback"
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
