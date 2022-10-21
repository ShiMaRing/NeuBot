package service

import (
	"NeuBot/configs"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
)

// ChatService 提供聊天功能
type ChatService struct {
	C *nlp.Client
}

func NewChatService() *ChatService {
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	if !configs.ChatConf.Enable {
		return nil
	}
	credential := common.NewCredential(
		configs.ChatConf.AppId,
		configs.ChatConf.AppSecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := nlp.NewClient(credential, "ap-guangzhou", cpf)
	return &ChatService{
		C: client,
	}
}

// Chat 提供聊天方法
func (c *ChatService) Chat(query string) (string, error) {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := nlp.NewChatBotRequest()
	request.Query = common.StringPtr(query) //设置用户参数

	// 返回的resp是一个ChatBotResponse的实例，与请求对象对应
	response, err := c.C.ChatBot(request)
	if err != nil {
		return "", err
	}
	//解析请求
	return *response.Response.Reply, nil
}
