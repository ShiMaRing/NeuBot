package model

import "fmt"

// DupStdError qq 重复绑定学生账号
var DupStdError = fmt.Errorf("qq duplicate bind")

// UserNotFoundError 无法在缓存中查找到User
var UserNotFoundError = fmt.Errorf("user not found")

// MaxLoginError  已到达最大注册数
var MaxLoginError = fmt.Errorf("get max login num")

var PasswordIncorrectError = fmt.Errorf("password incorrect ")
