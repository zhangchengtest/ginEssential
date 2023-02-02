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
	"strconv"
	"strings"
	"time"
)

func RandomPoetry() {
	DB := sqls.DB()

	// 创建用户
	var mywords []model.Tag

	var theme model.Poetry

	var article model.Article

	var ss string

	rand.Seed(time.Now().UnixNano())
	aid := rand.Intn(30000) + 1
	DB.Where("id = ?", aid).Find(&theme)
	ss = theme.Name + "\n"
	ss = ss + theme.Dynasty + "\n"
	ss = ss + theme.Poet + "\n"

	DB.Where("poetry_id = ?", aid).Find(&mywords)

	for _, w := range mywords {
		ss = ss + " " + w.Tag
	}
	dd := strings.Replace(theme.Content, "<br/>", "\n", -1)
	dd = strings.Replace(dd, "<br>", "\n", -1)
	ss = ss + "\n" + dd + "\n"

	aid = rand.Intn(80) + 1
	DB.Where("chapter = ?", aid).Find(&article)

	ss = ss + strconv.Itoa(aid) + "\n"
	ss = ss + article.Content + "\n"

	sendArticle(ss)
}

func sendPoetry(msg string) {
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
