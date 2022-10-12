package handler

import (
	"NeuBot/model"
	"fmt"
)

const (
	Menu = `1.显示菜单
            2.注册账号
            3.订阅或取消每日课程提醒
  			4.订阅或取消每日健康上报
			6.订阅或取消每周东B要闻
            7.注销账号
            8.反馈`
)

// NoticeHandler 对于新添加的好友的处理逻辑
type NoticeHandler struct {
	noticeReq *model.RequestReq
}

func NewNoticeHandler(noticeReq *model.RequestReq) *NoticeHandler {
	return &NoticeHandler{noticeReq: noticeReq}
}

// Greet 给新好友打招呼
func (handler *NoticeHandler) Greet() error {
	fmt.Println(*handler.noticeReq)
	_, err := ReplyMsg(int64(handler.noticeReq.UserID), "欢迎使用NEU-BOT，发送序号指令使用，目前不支持南湖同学", false)
	if err != nil {
		return err
	}
	_, err = ReplyMsg(int64(handler.noticeReq.UserID), Menu, false)

	if err != nil {
		return err
	}
	return nil
}
