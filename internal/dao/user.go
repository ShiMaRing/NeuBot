package dao

import (
	"NeuBot/configs"
	"NeuBot/internal/cache"
	"NeuBot/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserDao dao,暂时实现假的业务逻辑
type UserDao struct {
	db    *gorm.DB
	cache *cache.UserCache
}

// NewUserDao  构造方法，需要读取配置连接数据库，暂时不需要实现，暂时使用缓存来代替
func NewUserDao() (*UserDao, error) {
	dsn := dsnBuild()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &UserDao{
		db:    db,
		cache: cache.NewUserCache(),
	}, nil
}

func dsnBuild() string {
	c := configs.MysqlConf
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Address, c.Port, c.Database,
	)
	return dsn
}

// GetUser 获取user,先获取缓存，缓存击穿再由db获取
func (u *UserDao) GetUser(qqNumber int64) (*model.User, error) {
	cacheUser, err := u.cache.GetUser(qqNumber)
	if err == nil {
		return cacheUser, nil
	}
	var user model.User
	result := u.db.Where("qq= ?", qqNumber).Preload("TimeTable").Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	//找不到该用户
	if user.ID == 0 {
		return nil, model.UserNotFoundError
	}
	return &user, nil
}

// InitTable 设置表
func (u *UserDao) InitTable() error {
	if u.db.Migrator().HasTable(&model.User{}) {
		if err := u.db.Migrator().DropTable(&model.User{}, &model.Course{}); err != nil {
			return err
		}
	}
	return u.db.Migrator().AutoMigrate(&model.User{}, &model.Course{})
}

// SetUser 先保存至数据库，再添加至缓存中，保证在缓存中的数据一定在数据库中
func (u *UserDao) SetUser(user *model.User) error {
	result := u.db.Clauses(clause.OnConflict{DoNothing: true}).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return u.cache.SetUser(user)
}

// CacheUser  缓存用户
func (u *UserDao) CacheUser(user *model.User) error {
	return u.cache.SetUser(user)
}

// UnbindUser  注销，解除与账号密码的绑定
func (u *UserDao) UnbindUser(qqNumber int64) error {
	user, err := u.GetUser(qqNumber)
	if err != nil {
		return err
	}
	user.State = model.LOGOUT
	user.Perm = 0
	user.StdNumber = ""
	user.Password = ""
	table := user.TimeTable
	if table != nil {
		for i := range table {
			table[i].IsSubmission = true
		}
	}
	err = u.UpdateUser(user)
	if err != nil {
		return err
	}
	return u.cache.UnbindUser(qqNumber)
}

// UpdateUser 由于采用指针传递，缓存中的用户信息不需要改变,只需要改变数据库中的信息
func (u *UserDao) UpdateUser(user *model.User) error {
	result := u.db.Session(&gorm.Session{FullSaveAssociations: true}).Where("qq= ?", user.QQ).Save(user)
	return result.Error
}

// GetAllUser 获取所有用户实例，直接查询缓存
func (u *UserDao) GetAllUser() ([]*model.User, error) {
	return u.cache.GetAllUser()
}

// RebuildUsers 重新组件所有用户信息
func (u *UserDao) RebuildUsers() error {
	users := make([]*model.User, 0)
	res := u.db.Find(&users)
	if res.Error != nil {
		return res.Error
	}
	for i := range users {
		u.cache.SetUser(users[i])
	}
	return nil
}

// CleanAllCourse 清理所有课程信息
func (u *UserDao) CleanAllCourse() error {
	res := u.db.Where("1=1").Unscoped().Delete(&model.Course{})
	u.cache.CleanTimeTable()
	return res.Error
}

// CleanUp 清理已发送的数据
func (u *UserDao) CleanUp() error {
	res := u.db.Where("is_submission = ?", 1).Unscoped().Delete(&model.Course{})
	return res.Error
}

// UpdateCourse 更新课程信息
func (u *UserDao) UpdateCourse(course *model.Course) error {
	res := u.db.Where("id=?", course.ID).Save(course)
	return res.Error
}
