package service

import (
	"NeuBot/configs"
	"NeuBot/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	configs.ConfigInit()
	service, err := NewHealthUpdateService()
	assert.NoError(t, err)

	user := &model.User{
		StdNumber: "test",
		Password:  "test",
	}
	err = service.Insert(user)
	assert.NoError(t, err)
}

func TestCancel(t *testing.T) {
	configs.ConfigInit()
	service, err := NewHealthUpdateService()
	assert.NoError(t, err)

	user := &model.User{
		StdNumber: "test",
		Password:  "test",
	}
	err = service.Cancel(user)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	configs.ConfigInit()
	service, err := NewHealthUpdateService()
	assert.NoError(t, err)

	res, err := service.GetUser("20206759")
	assert.NoError(t, err)
	assert.Equal(t, true, res)

	res, err = service.GetUser("20546")
	assert.NoError(t, err)
	assert.Equal(t, false, res)

}
