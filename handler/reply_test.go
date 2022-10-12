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
	data, err := ReplyMsg(1150840779, "你好", false)
	assert.NoError(t, err)
	t.Log(string(data))
}
