package dao

import (
	"NeuBot/configs"
	"NeuBot/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var user *model.User
var userDao *UserDao

func testInit() error {
	configs.ConfigInit()
	var err error
	userDao, err = NewUserDao()
	if err != nil {
		return err
	}
	//添加数据
	user = &model.User{
		QQ:        1150840779,
		StdNumber: "20206759",
		Password:  "xgs583719992",
		State:     model.Logined,
		Perm:      model.CoursePerm,
		TimeTable: []*model.Course{
			{
				WeekDay:      1,
				ClassName:    "hello",
				Start:        1,
				ClassLength:  3,
				Place:        "103",
				Teacher:      "wuwu",
				IsSubmission: false,
			},
			{
				WeekDay:      2,
				ClassName:    "毛概",
				Start:        2,
				ClassLength:  4,
				Place:        "A389",
				Teacher:      "无",
				IsSubmission: false,
			},
		},
	}
	return nil
}

func TestTableInit(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	err = userDao.InitTable()
	assert.NoError(t, err)
}

func TestCreateUser(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	userDao.SetUser(user)
	fmt.Println(user)
}

func TestGetUser(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	user, err := userDao.GetUser(1150840779)
	assert.NoError(t, err)
	fmt.Println(user)
}

func TestUpdateUser(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	assert.NoError(t, err)
	user.StdNumber = "hdakdajd"
	user.Password = "dadada"
	user.TimeTable[0].Teacher = "xgs"
	userDao.UpdateUser(user)
	assert.NoError(t, err)
	fmt.Println(user)
}

func TestGetAllUser(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	allUser, err := userDao.GetAllUser()
	assert.NoError(t, err)
	for i := range allUser {
		fmt.Println(*allUser[i])
	}
}

func TestCleanUp(t *testing.T) {
	err := testInit()
	assert.NoError(t, err)
	err = userDao.CleanUp()
	assert.NoError(t, err)

}
