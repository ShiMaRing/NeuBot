package handler

import (
	"NeuBot/configs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testInit() {
	configs.ConfigInit()
	Init()
}

func TestReplyMsg(t *testing.T) {
	testInit()
	data, err := replyMsg(1150840779, "你好")
	assert.NoError(t, err)
	t.Log(string(data))
}

func TestReplyImage(t *testing.T) {
	testInit()
	ReplyMsg(1150840779, BuildImageMessage("onland.jpg"), false)
}

func TestRefreshCourses(t *testing.T) {
	configs.ConfigInit()
	handler, err := newSchedulerHandler()
	handler.srv.RebuildUsers()
	err = handler.refreshCourse()
	assert.NoError(t, err)
}
