package service

import (
	"NeuBot/internal/dao"
	"NeuBot/model"
)

// UserService 提供service方法
type UserService struct {
	dao *dao.UserDao
}

func NewUserService() (*UserService, error) {
	userDao, err := dao.NewUserDao()
	if err != nil {
		return nil, err
	}
	return &UserService{
		dao: userDao,
	}, nil
}

func (s *UserService) GetUser(qqNumber int64) (*model.User, error) {
	return s.dao.GetUser(qqNumber)
}

func (s *UserService) SetUser(user *model.User) error {
	return s.dao.SetUser(user)
}

func (s *UserService) DeleteUser(qqNumber int64) error {
	return s.dao.DeleteUser(qqNumber)
}
