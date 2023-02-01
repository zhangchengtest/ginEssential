package push

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/spider"
	"ginEssential/util"
	jsoniter "github.com/json-iterator/go"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	log "github.com/sirupsen/logrus"
	"github.com/zhangchengtest/simple/sqls"
	"net/url"
	"reflect"
	"strings"
)

var funcMap = make(map[string]func(msg string))

type Push struct {
	Label string `mapstructure:"label"`
	Value string `mapstructure:"value"`
}

type Msg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var ddUrl = "https://oapi.dingtalk.com/robot/send?access_token="

func (token Push) Dd(msg Msg) bool {
	content := `{"msgtype": "text",
		"text": {"content": "` + msg.Content + `"}
	}`
	if token.Value == "" {
		log.Error("dd token is empty!")
		return false
	}
	return spider.PostJson(ddUrl+token.Value, content)
}

func (token Push) Wechat(msg Msg) bool {
	sendWechat(msg.Content)
	return true
}

const wxToken = "cheng12345678"

func sendWechat(msg string) {
	log.Info("send to remote")
	DB := sqls.DB()
	var users []model.User
	DB.Where("unionid = ?", "123456").Find(&users)

	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "wx70711c9b88f9c12f",
		AppSecret: "20993710aa48342888d3a0b1755af9d6",
		Token:     wxToken,
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	for _, user := range users {

		data := message.MediaText{
			Content: msg,
		}
		customerMessage := message.CustomerMessage{
			Msgtype: message.MsgTypeText,
			Text:    &data,
			ToUser:  user.Openid,
		}
		officialAccount.GetCustomerMessageManager().Send(&customerMessage)
	}

}

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

func (token Push) Hook(msg Msg) bool {
	b, err := jsonIterator.MarshalToString(msg)
	if err != nil {
		log.Error(err)
		return false
	}
	return spider.PostJson(token.Value, b)
}

func (token Push) Console(msg Msg) bool {
	fmt.Println(msg.Content)
	return true
}

func (token Push) ServerChan(msg Msg) bool {
	scUrl := "https://sc.ftqq.com/" + token.Value + ".send"
	return strings.Contains(spider.GetResponseBody(scUrl+"?text="+url.QueryEscape(msg.Title)+"&desp="+url.QueryEscape(msg.Content)), "success")
}

func callReflect(any interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	if v := reflect.ValueOf(any).MethodByName(name); v.String() == "<invalid Value>" {
		return nil
	} else {
		return v.Call(inputs)
	}
}

func (token Push) Push(msg Msg) bool {
	value := callReflect(&token, util.Capitalize(token.Label), msg)
	if value != nil {
		return value[0].Bool()
	}
	return false
}
