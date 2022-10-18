package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"sync"
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
	gorm.Model            //id字段
	QQ         int64      `gorm:"unique"` //qq号
	StdNumber  string     //学号
	Password   string     //密码
	State      int        //当前状态
	Perm       int        //用户权限
	TimeTable  []*Course  //用户持有当前星期的课表
	Token      string     `gorm:"-:all"` //用户的查询token,不进行持久化
	Mu         sync.Mutex `gorm:"-:all"` //互斥锁，主要保护课程信息不会被异步操作，不进行持久化
}

func (u User) String() string {
	bytes, _ := json.Marshal(u)
	return string(bytes)
}
