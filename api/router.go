package api

import (
	"NeuBot/handler"
	"NeuBot/model"
	"NeuBot/tools"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"time"
)

const (
	Message   = "message"       //消息, 例如, 群聊消息
	Request   = "request"       //请求, 例如, 好友申请
	Notice    = "notice"        //通知, 例如, 群成员增加
	MetaEvent = "meta_event"    //元事件, 例如, go-cqhttp 心跳包
	Friend    = "friend"        //请求，请求添加好友
	AddFriend = "friend_add"    //通知，添加了新好友
	Duration  = 2 * time.Second //限流默认间隔
)

var noticeHandler *handler.NoticeHandler
var msgHandler *handler.MessageHandler
var limiter *tools.Limiter

func init() {
	var err error
	noticeHandler = handler.NewNoticeHandler()
	msgHandler, err = handler.NewMessageHandler()
	if err != nil {
		log.Fatalln(err)
	}
	limiter = tools.NewLimiter(Duration)
}

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
	baseReq := &model.BaseReq{}
	err = json.Unmarshal(reqBody, baseReq)
	if err != nil {
		log.Println(err)
	}
	switch baseReq.PostType {
	case MetaEvent:
	//元数据直接抛弃,不做处理
	case Request:
		request := &model.RequestReq{}
		err := json.Unmarshal(reqBody, request)
		if err != nil {
			log.Fatalln(err)
		}
		//加好友请求，直接同意即可
		if request.RequestType == Friend {
			data, _ := json.Marshal(struct {
				Approve bool   `json:"approve,omitempty"`
				Remark  string `json:"remark,omitempty"`
			}{
				Approve: true,
			})
			c.Data(200, handler.JSONContentType, data)
			return
		}
	case Notice:
		noticeReq := &model.NoticeReq{}
		err := json.Unmarshal(reqBody, noticeReq)
		if err != nil {
			log.Fatalln(err)
		}
		if noticeReq.PostType == Notice && noticeReq.NoticeType == AddFriend {
			go func() {
				noticeHandler.Greet(noticeReq)
			}()
			return
		}
	case Message:
		message := &model.MsgReq{}
		err := json.Unmarshal(reqBody, message)
		if err != nil {
			log.Fatalln(err)
		}
		if !limiter.Filter(message.UserID) {
			handler.ReplyMsg(message.UserID, "请不要过于频繁发送消息", false)
			return
		}
		//得到消息请求后，需要交给handler处理
		go func() {
			msgHandler.HandleMessage(message)
		}()
	default:
		return
	}
}
