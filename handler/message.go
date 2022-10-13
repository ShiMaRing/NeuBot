package handler

import (
	"NeuBot/internal/service"
	"NeuBot/model"
)

const (
	MenuMessage     string = "1"
	LoginMessage           = "2"
	CourseMessage          = "3"
	HealthMessage          = "4"
	LogoutMessage          = "5"
	FeedbackMessage        = "6"
)

type MessageHandler struct {
	srv *service.UserService
}

// NewMessageHandler 构造函数
func NewMessageHandler() (*MessageHandler, error) {
	userService, err := service.NewUserService()
	if err != nil {
		return nil, err
	}
	return &MessageHandler{srv: userService}, nil
}

// HandleMessage 对请求进行处理，提取消息，并进行回复
// 需要根据用户当前的状态判断逻辑
func (h *MessageHandler) HandleMessage(msg *model.MsgReq) {
	//先判断一下消息类型
	message := msg.Message
	//此时可能为数字，乱发的消息
	switch message {
	case MenuMessage:
	case LoginMessage:
	case CourseMessage:
	case HealthMessage:
	case LogoutMessage:
	case FeedbackMessage:
	}

}

func (h *MessageHandler) HandleMenuMessage() {

}

func (h *MessageHandler) HandleLoginMessage() {

}

// ReplyErrorMsg 发送报错信息
func (h *MessageHandler) ReplyErrorMsg(qqNumber int64, msg string) error {
	_, err := ReplyMsg(qqNumber, msg, false)
	return err
}
