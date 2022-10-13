package handler

import (
	"NeuBot/model"
	"log"
)

const (
	Menu = `1.显示菜单
            2.注册账号
            3.订阅或取消每日课程提醒
  			4.订阅或取消每日健康上报
            5.注销账号
            6.反馈`
)

// NoticeHandler 对于新添加的好友的处理逻辑
type NoticeHandler struct {
}

func NewNoticeHandler() *NoticeHandler {
	return &NoticeHandler{}
}

// Greet 给新好友打招呼
func (h *NoticeHandler) Greet(req *model.RequestReq) {

	_, err := ReplyMsg(req.UserID, "欢迎使用NEU-BOT，发送序号指令使用，目前不支持南湖同学", false)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = ReplyMsg(req.UserID, Menu, false)
	if err != nil {
		log.Fatalln(err)
	}
}
