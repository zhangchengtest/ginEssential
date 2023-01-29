package job

import (
	"ginEssential/model"
	"github.com/rs/zerolog/log"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zhangchengtest/simple/sqls"
	"math/rand"
	"time"
)

func RandomArticle() {
	DB := sqls.DB()

	// 创建用户
	var mywords []model.Words

	var theme model.Theme

	var ss string

	rand.Seed(time.Now().UnixNano())
	aid := rand.Intn(19) + 1
	DB.Where("id = ?", aid).Find(&theme)
	ss = "主题：" + theme.Name + "\n"
	ramids := make([]int, 20)
	for i := 0; i < 20; i++ {
		ramId := rand.Intn(70000) + 1
		ramids[i] = ramId
	}
	DB.Where("id in ?", ramids).Find(&mywords)

	for _, w := range mywords {
		ss = ss + " " + w.Name
	}
	sendArticle(ss)
}

func sendArticle(msg string) {
	log.Info().Msgf("send to remote")
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
