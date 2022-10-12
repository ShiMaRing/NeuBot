package handler

import (
	"NeuBot/configs"
	"NeuBot/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	Timeout         = 5 * time.Second
	JSONContentType = "application/json"
)

var (
	once       sync.Once
	client     *http.Client
	sendMsgUrl string
)

func Init() {
	initClient()
	initUrl()
}

//后期可以根据性能进行优化
func initClient() {
	client = &http.Client{
		Timeout: Timeout,
	}
}

func initUrl() {
	address := configs.CqhttpConf.Address
	port := configs.CqhttpConf.Port
	sendMsgUrl = fmt.Sprintf("http://%s:%d/send_private_msg", address, port)
}

// ReplyMsg 发送消息，需要指定是否需要进行转义，true则表示作为纯文本发送,默认选项
func ReplyMsg(receiver int64, msg string, autoEscape ...bool) ([]byte, error) {

	var isAutoEscape bool
	//不传参默认为false，即不转解析直接发送
	if autoEscape == nil {
		isAutoEscape = true
	} else {
		isAutoEscape = autoEscape[0]
	}
	//如果要解析的话，需要做逃逸处理

	message := model.ReplyMessage{
		UserId:     receiver,
		Message:    msg,
		AutoEscape: isAutoEscape,
	}
	jsonBody, err := json.Marshal(message)

	if err != nil {
		log.Println(err) //此处应当进行日志记录，后续需要集成Zap
		return nil, err
	}
	data := bytes.NewReader(jsonBody)
	res, err := client.Post(sendMsgUrl, JSONContentType, data)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		data, _ := ioutil.ReadAll(res.Body)
		return data, fmt.Errorf("response code is %d, body:%s", res.StatusCode, string(data))
	}

	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	d, _ := ioutil.ReadAll(res.Body)

	return d, nil
}
