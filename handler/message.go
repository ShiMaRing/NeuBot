package handler

import (
	"NeuBot/configs"
	"NeuBot/internal/service"
	"NeuBot/model"
	"NeuBot/pkg/spider"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	Menu = "1.显示菜单\n2.绑定学号\n3.订阅课程提醒\n4.取消课程提醒\n5.自动健康上报\n6.取消自动上报\n7.解绑账号\n输入对应序号执行指令\n反馈请携带前缀 feedback"
)

const (
	MenuMessage        string = "1"
	LoginMessage              = "2"
	SubCourseMessage          = "3"
	UnSubCourseMessage        = "4"
	SubHealthMessage          = "5"
	UnSubHealthMessage        = "6"
	LogoutMessage             = "7"

	FeedBackPrefix = "feedback"
)

const (
	MinStdNumber = 20170000
	MaxStdNumber = 99999999
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
		h.handleMenuMessage(msg.UserID)
	case LoginMessage:
		h.handleLoginMessage(msg.UserID)
	case SubCourseMessage:
	case UnSubCourseMessage:
	case SubHealthMessage:
	case UnSubHealthMessage:
	case LogoutMessage:
		h.handleLogoutMessage(msg.UserID)
	default:
		//说明用户在输入其他内容，此时需要查询用户状态进行判断
		h.handleUnknownMessage(msg)
	}
}

func (h *MessageHandler) handleMenuMessage(qqNumber int64) {
	ReplyMsg(qqNumber, Menu, false)
}

func (h *MessageHandler) handleLoginMessage(id int64) {
	user, err := h.srv.GetUser(id)
	//如果没找到或者查询过程中出现错误
	if err != nil {
		if errors.Is(err, model.UserNotFoundError) {
			usr := &model.User{
				QQ:    id,
				State: model.LOGOUT, //未登陆
			}
			//如果找不到用户，说明允许注册，此时将数据持久化至数据库
			err = h.srv.SetUser(usr)

			if err != nil {
				h.loginFail(id, err)
				return
			}

			ReplyMsg(id, "请输入东北大学学生账号密码\n格式：账号 密码\n以空格分割")
		} else {
			h.loginFail(id, err)
		}
		return
	}
	//如果已经注册过了
	if user.State&model.Logined != 0 {
		ReplyMsg(id, "请勿重复绑定")
		return
	}

	ReplyMsg(id, "请输入东北大学学生账号密码\n格式：账号 密码\n以空格分割")
}

//解除绑定
func (h *MessageHandler) handleLogoutMessage(id int64) {
	user, err := h.srv.GetUser(id)
	if err != nil {
		if errors.Is(err, model.UserNotFoundError) {
			ReplyMsg(id, "解绑失败，尚未绑定账号")
			logError(err)
		} else {
			h.loginFail(id, err)
		}
		return
	}
	//检查一下user状态
	if user == nil {
		h.logoutFail(id, fmt.Errorf("无法找到对应的绑定账号"))
	}
	if user.State&model.Logined == 0 {
		ReplyMsg(id, "解绑失败，尚未绑定账号")
	} else {
		err := h.srv.DeleteUser(id)
		if err != nil {
			h.logoutFail(id, err)
			return
		}
		ReplyMsg(id, "解绑成功")
	}
}

//处理未知消息,可能是反馈，可能是登陆消息
func (h *MessageHandler) handleUnknownMessage(msg *model.MsgReq) {
	//检查一下前缀，是否带有feedback
	message := strings.TrimSpace(msg.Message) //删除前后空格
	if strings.HasPrefix(message, FeedBackPrefix) {
		//转发给master
		ReplyMsg(configs.BotConf.MasterId, buildTransmitMsg(msg))
		ReplyMsg(msg.UserID, "反馈成功")
		return
	}

	id := msg.UserID
	//剩下只可能是绑定信息，检查一下消息格式
	split := strings.Split(message, " ")
	tmp := make([]string, 0)
	for i := range split {
		if split[i] != "" {
			tmp = append(tmp, split[i])
		}
	}
	//通过对学号合法性进行校验
	stdNumber, err := strconv.Atoi(tmp[0])
	if err != nil || len(tmp) != 2 {
		ReplyMsg(id, "请不要发送无关消息")
		return
	}

	account := tmp[0]
	password := tmp[1]

	if stdNumber < MinStdNumber || stdNumber > MaxStdNumber {
		ReplyMsg(id, "学号不合法")
		return
	}
	//用户在加入好友，以及输入2时，都会将自身加入缓存中，以及加入数据库中
	user, err := h.srv.GetUser(msg.UserID)
	if err != nil {
		ReplyMsg(msg.UserID, fmt.Sprintf("暂时无法绑定账号\n 错误原因：\n %v", err))
		return
	}

	ReplyMsg(msg.UserID, fmt.Sprintf("正在验证请稍后"))

	success, err := spider.Auth(account, password)

	if err != nil {
		ReplyMsg(msg.UserID, fmt.Sprintf("验证失败\n 错误原因：\n %v", err))
		return
	}
	if !success {
		ReplyMsg(msg.UserID, "账号密码错误")
		return
	}
	user.StdNumber = account
	user.Password = password
	user.State = model.Logined
	h.srv.UpdateUser(user) //更新用户信息
	ReplyMsg(msg.UserID, "绑定账号成功")
	return

}

func buildTransmitMsg(msg *model.MsgReq) string {
	nickName := msg.Sender.Nickname
	var sender string
	if nickName == "" {
		sender = strconv.Itoa(int(msg.UserID))
	} else {
		sender = fmt.Sprintf("%s(%d)", nickName, msg.UserID)
	}
	message := strings.TrimSpace(msg.Message)
	message = message[len(FeedBackPrefix):]
	return fmt.Sprintf("收到来自 %s 的反馈:\n %s", sender, message)
}

func (h *MessageHandler) loginFail(id int64, err error) {
	ReplyMsg(id, fmt.Sprintf("系统错误，暂时无法注册，错误信息\n %v", err))
	logError(err)
}

func (h *MessageHandler) logoutFail(id int64, err error) {
	ReplyMsg(id, fmt.Sprintf("系统错误，暂时无法解绑，错误信息\n %v", err))
	logError(err)
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}
