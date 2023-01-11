package controller

import (
	"encoding/json"
	"flag"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/render"
	"ginEssential/util"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	PHONE_NUM = 15216771668
)

func RandomFood(ctx *gin.Context) {
	DB := sqls.DB()

	var foods []model.Food

	rand.Seed(time.Now().UnixNano())

	DB.Where("category = ?", "rising").Find(&foods)
	chapter := rand.Intn(len(foods))

	articleVO := model.FoodVO{}

	util.SimpleCopyProperties(&articleVO, &foods[chapter])

	model.Success(ctx, gin.H{"food": articleVO}, "查询成功")

}

func SearchFood(ctx *gin.Context) {
	DB := sqls.DB()

	var foods []model.Food

	rand.Seed(time.Now().UnixNano())

	DB.Where("category = ?", "rising").Order("create_dt desc").Limit(30).Find(&foods)

	vos := render.BuildFoods(foods)

	model.Success(ctx, gin.H{"foods": vos}, "查询成功")

}

func SendSms(ctx *gin.Context) {
	DB := sqls.DB()

	var foods []model.Food

	rand.Seed(time.Now().UnixNano())

	DB.Where("category = ?", "rising").Find(&foods)
	chapter := rand.Intn(len(foods))

	articleVO := model.FoodVO{}

	util.SimpleCopyProperties(&articleVO, &foods[chapter])

	var regionId = flag.String("regionId", "cn-hangzhou", "区域标识")
	var accessKeyId = flag.String("id", config.Instance.DAYU.APP_KEY, "accessKeyId")
	var accessKeySecret = flag.String("secret", config.Instance.DAYU.APP_SECRET, "accessKeySecret")
	var verifyCode = flag.String("code", "1234", "验证码")
	var phoneNumbers = flag.Int("phonenumbers", PHONE_NUM, "手机号")
	flag.Parse()

	if *phoneNumbers <= 0 {
		panic(fmt.Errorf("invalid phonenumbers"))
	}

	client, err := dysmsapi.NewClientWithAccessKey(*regionId, *accessKeyId, *accessKeySecret)
	if err != nil {
		panic(err)
	}

	params, _ := json.Marshal(map[string]interface{}{
		"code": verifyCode,
	})

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.TemplateCode = config.Instance.DAYU.SMS_TEMPLATE_CODE
	request.SignName = "新云网"
	request.TemplateParam = string(params)
	request.PhoneNumbers = strconv.Itoa(*phoneNumbers)

	resp, err := client.SendSms(request)
	if err != nil {
		log.Printf("send sms failed resp=%v err=%v", resp, err)
		panic(err)
	}

	if !resp.IsSuccess() {
		log.Printf("send sms failed resp=%v err=%v", resp, err)
		panic(fmt.Errorf("failed: unknown reason"))
	}

	model.Success(ctx, gin.H{"food": articleVO}, "查询成功")

}
