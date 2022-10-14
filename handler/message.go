package handler

import (
	"NeuBot/internal/service"
	"NeuBot/model"
	"errors"
	"fmt"
	"log"
)

const (
	MenuMessage        string = "1"
	LoginMessage              = "2"
	SubCourseMessage          = "3"
	UnSubCourseMessage        = "4"
	SubHealthMessage          = "5"
	UnSubHealthMessage        = "6"
	LogoutMessage             = "7"
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

	}

}

func (h *MessageHandler) handleMenuMessage(qqNumber int64) error {
	_, err := ReplyMsg(qqNumber, Menu, false)
	return err
}

func (h *MessageHandler) handleLoginMessage(id int64) {
	//收到了登录请求，先将其保存至缓存中，只有在验证账号密码正确的前提下才会写入数据库
	user, err := h.srv.GetUser(id)
	//如果没找到或者查询过程中出现错误
	if err != nil {
		if errors.Is(err, model.UserNotFoundError) {
			usr := &model.User{
				QQ:    id,
				State: model.Logining, //状态切换
			}
			//如果找不到用户，说明允许注册，此时将数据持久化至数据库
			err = h.srv.SetUser(usr)

			if err != nil {
				h.loginFail(id, err)
				return
			}
			_, err = ReplyMsg(id, "请输入东北大学学生账号密码\n格式：账号 密码\n以空格分割")
			logError(err)
		} else {
			h.loginFail(id, err)
		}
		return
	}
	//如果已经注册过了
	if user.State&model.Logined != 0 {
		_, err = ReplyMsg(id, "请勿重复注册")
		logError(err)
		return
	}
	_, err = ReplyMsg(id, "请输入东北大学学生账号密码\n格式：账号 密码\n以空格分割")
	logError(err)
}

//解除绑定
func (h *MessageHandler) handleLogoutMessage(id int64) {
	user, err := h.srv.GetUser(id)
	if err != nil {
		if errors.Is(err, model.UserNotFoundError) {
			ReplyMsg(id, "解绑失败，尚未绑定账号")
		} else {
			h.loginFail(id, err)
		}
		return
	}
	//检查一下user状态
	if user == nil {
		h.logoutFail(id, fmt.Errorf("无法找到对应的绑定账号"))
	}
	if user.State&model.Logined != 1 {
		ReplyMsg(id, "解绑失败，尚未绑定账号")
	} else {
		err := h.srv.DeleteUser(id)
		if err != nil {
			h.loginFail(id, err)
			return
		}
		ReplyMsg(id, "解绑成功")
	}
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
