package cache

import (
	"NeuBot/model"
	"sync"
	"time"
)

const MaxSize = 100

//回收器，定时回收更新cache中的TimeTable
type janitor struct {
	Interval time.Duration
	stop     chan bool
}

// UserCache 使用qq号与一个User类对应
type UserCache struct {
	mu      sync.RWMutex
	users   map[int64]*model.User
	janitor *janitor //清理器，暂时没什么用，之后需要定时运行来清理垃圾
}

func NewUserCache() *UserCache {
	return &UserCache{
		users: make(map[int64]*model.User, MaxSize),
	}
}

// SetUser 设置User
func (c *UserCache) SetUser(user *model.User) error {
	c.mu.Lock()
	if _, ok := c.users[user.QQ]; ok {
		c.mu.Unlock()
		return model.DupStdError
	} else {
		//检查是否到达最大注册人数
		if len(c.users) >= MaxSize {
			return model.MaxLoginError
		}
		c.users[user.QQ] = user
		c.mu.Unlock()
		return nil
	}
}

func (c *UserCache) UpdateUser(user *model.User) error {
	c.mu.Lock()
	if _, ok := c.users[user.QQ]; ok {
		c.users[user.QQ] = user
		c.mu.Unlock()
		return nil
	} else {
		c.mu.Unlock()
		return model.UserNotFoundError
	}
}

// GetUser 获取User实例
func (c *UserCache) GetUser(qqNumber int64) (*model.User, error) {
	c.mu.Lock()
	user, ok := c.users[qqNumber]
	c.mu.Unlock()
	if ok {
		return user, nil
	}
	return nil, model.UserNotFoundError
}

// DeleteUser 注销方法，注销用户
func (c *UserCache) DeleteUser(qqNumber int64) error {
	c.mu.Lock()
	if _, ok := c.users[qqNumber]; ok {
		delete(c.users, qqNumber)
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()
	return model.UserNotFoundError
}

// CleanTimeTable 清理所有的课表
func (c *UserCache) CleanTimeTable() {
	c.mu.Lock()
	for _, user := range c.users {
		user.TimeTable = nil
	}
	c.mu.Unlock()
}
