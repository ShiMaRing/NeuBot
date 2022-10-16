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

// DeleteUser 逻辑为删除用户账号密码并切换状态为未登陆
func (s *UserService) DeleteUser(qqNumber int64) error {
	return s.dao.DeleteUser(qqNumber)
}

func (s *UserService) UpdateUser(user *model.User) error {
	return s.dao.UpdateUser(user)
}

func (s *UserService) GetAllUser() ([]*model.User, error) {
	return s.dao.GetAllUser()

}
