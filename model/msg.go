package model

//user_id	int64	-	对方 QQ 号
//group_id	int64	-	主动发起临时会话群号(机器人本身必须是管理员/群主)
//message	message	-	要发送的内容
//auto_escape	boolean	false	消息内容是否作为纯文本发送 ( 即不解析 CQ 码 ) , 只在 message 字段是字符串时有效

// ReplyMessage 返回给cqhttp的消息
type ReplyMessage struct {
	UserId     int64  `json:"user_id,omitempty"`
	GroupId    int64  `json:"group_id,omitempty"`
	Message    string `json:"message,omitempty"`
	AutoEscape bool   `json:"auto_escape,omitempty"`
}
