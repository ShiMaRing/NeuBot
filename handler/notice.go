package handler

import (
	"NeuBot/internal/service"
	"NeuBot/model"
	"log"
)

// NoticeHandler 对于新添加的好友的处理逻辑
type NoticeHandler struct {
	srv *service.UserService
}

func NewNoticeHandler() (*NoticeHandler, error) {
	userService, err := service.NewUserService()
	if err != nil {
		return nil, err
	}
	return &NoticeHandler{
		srv: userService,
	}, nil
}

// Greet 给新好友打招呼
func (h *NoticeHandler) Greet(req *model.NoticeReq) {
	ReplyMsg(req.UserId, "欢迎使用NEU-BOT，发送序号指令使用，目前不支持南湖同学", false)
	ReplyMsg(req.UserId, Menu, false)
	//将用户进行保存
	err := h.srv.SetUser(&model.User{QQ: req.UserId, State: model.LOGOUT})
	if err != nil {
		log.Println(err)
	}
}
