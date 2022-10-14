package model

import (
	"gorm.io/gorm"
	"net/http/cookiejar"
)

const (
	LOGOUT  int = 1 << iota //未登陆
	Logined                 //已登陆
)

const (
	HealthPerm = 1 << iota //开启健康上报
	CoursePerm             //开启课程提醒
)

// User 用户类
type User struct {
	gorm.Model               //id字段
	QQ         int64         //qq号
	StdNumber  string        //学号
	Password   string        //密码
	State      int           //当前状态
	Perm       int           //用户权限
	TimeTable  *TimeTable    //用户持有当前星期的课表
	jar        cookiejar.Jar //记录登陆的jar，放便下次请求直接携带jar
}
