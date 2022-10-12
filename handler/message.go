package handler

import "NeuBot/model"

const (
	MenuMessage     string = "1"
	LoginMessage           = "2"
	CourseMessage          = "3"
	HealthMessage          = "4"
	LogoutMessage          = "5"
	FeedbackMessage        = "6"
)

type MessageHandler struct {
	message *model.MsgReq
}

// NewMessageHandler 构造函数
func NewMessageHandler(message *model.MsgReq) *MessageHandler {
	return &MessageHandler{message: message}
}

// HandleMessage 对请求进行处理，提取消息，并进行回复
// 需要根据用户当前的状态判断逻辑
func (h *MessageHandler) HandleMessage() {

}
