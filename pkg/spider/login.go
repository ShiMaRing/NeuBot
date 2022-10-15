package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/neucn/neugo"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	Url       = "https://pass.neu.edu.cn/tpass/login?service=https://portal.neu.edu.cn/tp_up/"
	Timeout   = 5 * time.Second
)

/*
rsa:20206759xgs583719992LT-368544-jTeSIlu77aCfmqfRcGnj10qOe1LXAu-tpass
ul:8
pl:12
lt:LT-368544-jTeSIlu77aCfmqfRcGnj10qOe1LXAu-tpass
execution:e1s1
_eventId:submit


*/

type Form struct {
	rsa       string `json:"rsa,omitempty"`
	ul        string `json:"ul,omitempty"`
	pl        string `json:"pl,omitempty"`
	lt        string `json:"lt,omitempty"`
	execution string `json:"execution,omitempty"`
	_eventId  string `json:"_eventId,omitempty"`
}

// AuthWithAccount  使用账号密码进行认证，同时返回token
func AuthWithAccount(stdNum string, password string) (success bool, token string, err error) {
	c := colly.NewCollector(
		colly.UserAgent(UserAgent),
	)
	c.SetRequestTimeout(Timeout) //设置超时时间
	var lt string
	var execution string
	//获取lt
	c.OnHTML("#lt", func(element *colly.HTMLElement) {
		lt = element.Attr("value")
	})
	c.OnHTML("input[name='execution']", func(element *colly.HTMLElement) {
		execution = element.Attr("value")
	})
	var jsession_id string
	c.OnResponse(func(response *colly.Response) {
		cookies := response.Headers.Values("Set-Cookie")
		for i := range cookies {
			cookie := cookies[i]
			if strings.Contains(cookie, "jsessionid_tpass") {
				res := strings.Split(cookie, ";")
				jsession_id = res[0]
				break
			}
		}
	})
	err = c.Visit(Url)
	if err != nil {
		return false, "", err
	}
	form := buildForm(stdNum, password, lt, execution)
	//用反射构建表单
	formData := make(map[string]string)
	v := reflect.ValueOf(form).Elem()
	t := reflect.TypeOf(form).Elem()
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		field := v.Field(i).String()
		formData[name] = field
	}
	url := "https://pass.neu.edu.cn/tpass/login;" + jsession_id + "?service=https://portal.neu.edu.cn/tp_up/"
	c.OnRequest(func(req *colly.Request) {
		req.Headers.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Headers.Add("Referer", "https://pass.neu.edu.cn/tpass/login")
	})
	c.OnResponse(func(response *colly.Response) {
		cookie := response.Request.Headers.Get("Cookie")
		//解析token
		tmps := strings.Split(cookie, ";")
		for i := range tmps {
			if strings.Contains(tmps[i], "tp_up") {
				tmp := tmps[i]
				token = strings.Split(tmp, "=")[1]
				break
			}
		}
	})
	err = c.Post(url, formData)
	if err != nil {
		return false, "", err
	}
	if token == "" {
		return false, "", nil
	}
	return true, token, nil
}

func Auth(stdNum string, password string) (success bool, token string, err error) {
	client := neugo.NewSession()
	err = neugo.Use(client).WithAuth(stdNum, password).Login(neugo.CAS)
	if err != nil {
		return false, "", err
	}
	token = neugo.About(client).Token(neugo.CAS)
	return true, token, nil
}

func buildForm(stdNum string, password string, lt string, execution string) *Form {
	lt = strings.TrimSpace(lt)
	rsa := fmt.Sprintf("%s%s%s", stdNum, password, lt)
	ul := len(stdNum)
	pl := len(password)
	return &Form{
		rsa:       rsa,
		ul:        strconv.Itoa(ul),
		pl:        strconv.Itoa(pl),
		lt:        lt,
		execution: execution,
		_eventId:  "submit",
	}
}
