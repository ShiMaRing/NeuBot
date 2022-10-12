package model

import (
	"gorm.io/gorm"
	"net/http/cookiejar"
	"time"
)

type State int

const (
	LOGOUT   State = iota //未登陆
	Logining              //正在登录
	Logined               //已登陆
	FeedBack              //正在填写回馈状态
)

const (
	HealthPerm = 1 << iota
	CoursePerm = 1 << iota
)

// User 用户类
type User struct {
	gorm.Model               //id字段
	QQ         int64         //qq号
	StdNumber  string        //学号
	Password   string        //密码
	State      int           //当前状态
	Perm       int           //用户权限
	LastSend   time.Time     //上一次发送消息的时间，时间不应该
	jar        cookiejar.Jar //记录登陆的jar，放便下次请求直接携带jar
}
