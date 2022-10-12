package handler

import (
	"NeuBot/model"
	"fmt"
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
	err := ReplyMsg(int64(handler.noticeReq.UserID), "你好，欢迎使用NEUBot")
	if err != nil {
		return err
	}
	return nil
}
