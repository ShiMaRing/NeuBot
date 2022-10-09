package configs

import (
	"fmt"
	"testing"
)

//测试初始化方法
func TestConfigInit(t *testing.T) {
	ConfigInit()
	fmt.Println(*RedisConf)
	fmt.Println(*MysqlConf)
	fmt.Println(*CqhttpConf)
}
