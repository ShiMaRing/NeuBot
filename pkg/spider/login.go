package spider

//爬虫逻辑

// Auth 登陆校园网关进行验证
func Auth(stdNum string, password string) (bool, error) {
	if stdNum == "20206759" && password == "1234" {
		return true, nil
	}
	return false, nil
}
