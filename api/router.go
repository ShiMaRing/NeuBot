package api

import (
	"NeuBot/handler"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
)

const (
	Message   = "message"    //消息, 例如, 群聊消息
	Request   = "request"    //请求, 例如, 好友申请
	Notice    = "notice"     //通知, 例如, 群成员增加
	MetaEvent = "meta_event" //元事件, 例如, go-cqhttp 心跳包
	AddFriend = "friend_add"
)

// GetMsg Bot Api 接受来自cqhttp上报的信息,bot入口
// 获取得到消息之后对消息类型进行判断，启动goroutine进行处理，需要使用线程池
func GetMsg(c *gin.Context) {
	//对接收到的消息进行解析
	body := c.Request.Body
	if body == nil {
		return
	}
	defer body.Close()
	reqBody, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
		return
	}
	baseReq := &BaseReq{}
	err = json.Unmarshal(reqBody, baseReq)
	if err != nil {
		log.Println(err)
	}
	switch baseReq.PostType {
	case MetaEvent:
	//元数据直接抛弃
	case Request:
	//请求，暂时不实现
	case Notice:

		noticeReq := &NoticeReq{}
		err := json.Unmarshal(reqBody, noticeReq)
		if err != nil {
			log.Println(err)
			return
		}
		if noticeReq.NoticeType == AddFriend {
			go func() {
				noticeHandler := handler.NewNoticeHandler(noticeReq)
				err = noticeHandler.Greet()
				if err != nil {
					log.Println(err)
					return
				}
			}()
		}
	case Message:

	}

}
