package api

// BaseReq 所有请求都会携带的基本消息
type BaseReq struct {
	Time     int    `json:"time"`
	PostType string `json:"post_type"`
}

// MsgReq 消息结构体，接受来自好友的消息
type MsgReq struct {
	BaseReq
	Sender     `json:"sender"`
	MessageID  int    `json:"message_id"`
	UserID     int    `json:"user_id"`
	Message    string `json:"message"`
	RawMessage string `json:"raw_message"`
}

// NoticeReq 消息结构体，当添加好友时触发
type NoticeReq struct {
	BaseReq
	NoticeType string `json:"notice_type"`
	UserID     int    `json:"user_id"`
}

// Sender  消息发送者
type Sender struct {
	Nickname string `json:"nickname"`
	UserID   int    `json:"user_id"`
}
