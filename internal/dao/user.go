package dao

import (
	"NeuBot/internal/cache"
	"NeuBot/model"
	"gorm.io/gorm"
)

// UserDao dao,暂时实现假的业务逻辑
type UserDao struct {
	db    *gorm.DB
	cache *cache.UserCache
}

// NewUserDao  构造方法，需要读取配置连接数据库，暂时不需要实现，暂时使用缓存来代替
func NewUserDao() (*UserDao, error) {
	return &UserDao{
		db:    &gorm.DB{},
		cache: cache.NewUserCache(),
	}, nil
}

func (u *UserDao) GetUser(qqNumber int64) (*model.User, error) {
	return u.cache.GetUser(qqNumber)
}

func (u *UserDao) SetUser(user *model.User) error {
	return u.cache.SetUser(user)
}

func (u *UserDao) DeleteUser(qqNumber int64) error {
	return u.cache.DeleteUser(qqNumber)
}
