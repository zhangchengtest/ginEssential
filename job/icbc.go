package job

import (
	"fmt"
	"ginEssential/model"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zhangchengtest/simple/sqls"
	"math"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

// URLGold 工商的黄金价格URL
var URLGold = "http://www.icbc.com.cn/ICBCDynamicSite/Charts/GoldTendencyPicture.aspx"

// Cache 缓存, 设置 告警的阈值 , +-0.5
type Cache struct {
	Alarm float64
}

var mycache = &Cache{0}

// IcbcGold 查询黄金价格
func IcbcGold() {
	var (
		res        *http.Response
		err        error
		doc        *goquery.Document
		httpClient http.Client
		jar        *cookiejar.Jar
		price      float64
	)

	// 处理cookies, 这里用不到保持session
	jar, _ = cookiejar.New(nil)
	httpClient.Jar = jar

	if res, err = httpClient.Get(URLGold); err != nil {
		log.Error().Msgf("请求失败, %v", err.Error())
		return
	}
	defer res.Body.Close()
	if doc, err = goquery.NewDocumentFromReader(res.Body); err != nil {
		log.Error().Msgf("goquery解析失败, %v", err.Error())
		return
	}
	// Attr 获取属性
	flag := false
	doc.Find(`#TABLE1 > tbody > tr:nth-child(2) > td:nth-child(3)`).Each(func(i int, s *goquery.Selection) {
		flag = true
		price, _ = strconv.ParseFloat(strings.TrimSpace(s.Text()), 64)
		log.Info().Msgf("当前黄金价格: %v, 告警阈值: %v", price, mycache.Alarm)

		t := time.Now()
		date := strftime.Format(t, "%Y-%m-%d %H:%M:%S")

		go sendGold(fmt.Sprintf(" %v黄金价格: %v", date, price))
		Alarm(price)
	})
	if !flag {
		log.Error().Msgf("没有获取到黄金价格")
	}
	// fmt.Println(doc.Find("#TABLE1"))
}

// Alarm 判断价格
func Alarm(price float64) {

	inc := int(math.Abs(price-mycache.Alarm) / 0.5)
	if inc >= 1 {
		if price-mycache.Alarm > 0 {
			if mycache.Alarm != 0 {
				go sendGold(fmt.Sprintf("当前价格: %v [上升]", price))
				log.Info().Msgf("当前价格: %v [上升]", price)
			}
			mycache.Alarm = mycache.Alarm + float64(inc)*0.5

		} else {
			mycache.Alarm = mycache.Alarm - float64(inc)*0.5
			log.Info().Msgf("当前价格: %v [下降]", price)
			go sendGold(fmt.Sprintf("当前价格: %v [下降]", price))
		}
	}

}

const wxToken = "cheng12345678" // 这里填微信开发平台里设置的 Token

func sendGold(msg string) {
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

//func wechat(msg string) {
//	var (
//		res *http.Response
//		err error
//	)
//	if res, err = http.PostForm("http://api.xx.com/weixin", url.Values{"msg": {msg}}); err != nil {
//		log.Error().Msg("发送失败")
//	} else {
//		log.Info().Msgf("发送成功")
//	}
//	res.Body.Close()
//}
