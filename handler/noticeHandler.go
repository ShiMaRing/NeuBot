package handler

import "NeuBot/api"

// NoticeHandler 对于新添加的好友的处理逻辑
type NoticeHandler struct {
	noticeReq *api.NoticeReq
}

func NewNoticeHandler(noticeReq *api.NoticeReq) *NoticeHandler {
	return &NoticeHandler{noticeReq: noticeReq}
}

// Greet 给新好友打招呼
func (handler *NoticeHandler) Greet() error {

	return nil
}
