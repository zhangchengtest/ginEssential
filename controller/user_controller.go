package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"ginEssential/common"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/redis"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/sjsdfg/common-lang-in-go/StringUtils"
	"github.com/zhangchengtest/simple/sqls"
	"gorm.io/gorm"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	aliutil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	qrcode "github.com/skip2/go-qrcode"
)

func Register(ctx *gin.Context) {
	DB := sqls.DB()
	// 1. 使用map获取application/json请求的参数
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	// fmt.Printf("requestMap：%v", requestMap)

	// 2. 使用结构体获取application/json请求的参数
	// var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	// fmt.Printf("requestUser：%v", requestUser)

	// 3. gin自带的bind获取application/json请求的参数
	var ginBindUser = model.User{}
	ctx.Bind(&ginBindUser)
	fmt.Printf("ginBindUser：%v", ginBindUser)

	email := ginBindUser.Email
	password := ginBindUser.Pwd
	userName := ginBindUser.UserName
	// name := ctx.PostForm("name")
	// mobile := ctx.PostForm("mobile")
	// password := ctx.PostForm("password")

	// 数据验证
	if !util.VerifyEmailFormat(*email) {
		// 自己封装过后
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "邮箱格式不对")
		return
	}

	if len(password) < 6 {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	log.Println(email, password)
	// 判断手机号是否存在
	if isEmailExist(DB, *email) {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	// 加密密码
	hasedPassword := util.MD5(password)
	// 创建用户
	newUser := model.User{
		UserId:      util.Myuuid(),
		CreateDt:    time.Now(),
		UpdateDt:    nil,
		UserName:    userName,
		LastLoginDt: time.Now(),
		NickName:    userName,
		Email:       email,
		Pwd:         hasedPassword,
	}

	var s = util.Worker1{}

	var sysUserRole = model.SysUserRole{
		Id:     int(s.GetId()),
		UserId: newUser.UserId,
		RoleId: 2,
	}
	DB.Create(&sysUserRole)
	DB.Create(&newUser)

	common.SendRegister(userName, *email)

	model.Success(ctx, gin.H{"userName": userName}, "注册成功")
}

func Login(ctx *gin.Context) {
	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	account := userLoginDTO.Account
	password := userLoginDTO.Pwd

	if len(password) < 6 {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	DB := sqls.DB()
	var user model.User
	// 数据验证
	if util.VerifyEmailFormat(account) {
		DB.Where("email = ?", account).First(&user)
		if user.UserId == "" {
			model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	} else {
		DB.Where("user_name = ?", account).First(&user)
		if user.UserId == "" {
			model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	}

	newSig := util.MD5(password) //转成加密编码
	// 将编码转换为字符串
	log.Printf("newSig : %v", newSig)
	// 判断密码是否正确
	if user.Pwd != newSig {
		model.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		model.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	uservo := model.UserVO{}

	util.SimpleCopyProperties(&uservo, &user)
	uservo.AccessToken = token
	uservo.Avatar = user.AvatarUrl

	var res []model.SysRole
	DB.Table("sys_role").Select("sys_role.code").
		Joins("left join sys_user_role on sys_role.id = sys_user_role.role_id").Where("user_id = ?", user.UserId).Scan(&res)
	fmt.Println(res)
	uservo.RoleCode = res[0].Code

	// 返回结果
	// ctx.JSON(200, gin.H{
	// 	"code":    200,
	// 	"data":    gin.H{"token": token},
	// 	"message": "注册成功",
	// })
	uservo.Pwd = ""
	model.Success(ctx, uservo, "登录成功")
}

func LoginThird(ctx *gin.Context) {
	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	account := userLoginDTO.Account
	password := userLoginDTO.Pwd

	if len(password) < 6 {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	DB := sqls.DB()
	var user model.User
	// 数据验证
	if util.VerifyEmailFormat(account) {
		DB.Where("email = ?", account).First(&user)
		if user.UserId == "" {
			model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	} else {
		DB.Where("user_name = ?", account).First(&user)
		if user.UserId == "" {
			model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	}

	newSig := util.MD5(password) //转成加密编码
	// 将编码转换为字符串
	log.Printf("newSig : %v", newSig)
	// 判断密码是否正确
	if user.Pwd != newSig {
		model.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		model.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	uservo := model.UserVO{}

	util.SimpleCopyProperties(&uservo, &user)
	uservo.AccessToken = token
	code := util.RandomString(10)
	redis.Set("code"+code, uservo.UserId, 300)
	uservo.RedirectUrl = userLoginDTO.RedirectUrl + "?code=" + code
	uservo.Pwd = ""
	model.Success(ctx, uservo, "登录成功")
}

func LoginPhone(ctx *gin.Context) {
	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	account := userLoginDTO.Account
	smsCaptcha := userLoginDTO.SmsCaptcha

	if len(smsCaptcha) != 6 {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	smsCaptcha_exist := redis.Get("code" + account)

	// 将编码转换为字符串
	log.Printf("smsCaptcha_exist : %v", smsCaptcha_exist)
	// 判断密码是否正确
	if smsCaptcha != smsCaptcha_exist {
		model.Response(ctx, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}

	DB := sqls.DB()
	var user model.User
	// 数据验证
	DB.Where("mobile = ?", account).First(&user)
	if user.UserId == "" {
		//create user

		randomStr := util.GenerateCode()
		password := randomStr
		userName := "游客" + randomStr

		// 加密密码
		hasedPassword := util.MD5(password)
		// 创建用户
		newUser := model.User{
			UserId:      util.Myuuid(),
			CreateDt:    time.Now(),
			UpdateDt:    nil,
			UserName:    userName,
			LastLoginDt: time.Now(),
			NickName:    userName,
			Mobile:      account,
			//Openid:      userLoginDTO.Openid,
			Email: nil,
			Pwd:   hasedPassword,
		}

		var s = util.Worker1{}

		var sysUserRole = model.SysUserRole{
			Id:     int(s.GetId()),
			UserId: newUser.UserId,
			RoleId: 2,
		}
		DB.Create(&sysUserRole)
		DB.Create(&newUser)

		user = newUser

	} else {
		//update user
		//result := map[string]interface{}{}

		//if !StringUtils.IsEmpty(userLoginDTO.Openid) {
		//	result["openid"] = userLoginDTO.Openid
		//	DB := sqls.DB()
		//	DB.Model(&user).Where("user_id", user.UserId).Updates(result)
		//}
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		model.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	uservo := model.UserVO{}

	util.SimpleCopyProperties(&uservo, &user)
	uservo.AccessToken = token
	code := util.RandomString(10)
	redis.Set("code"+code, uservo.UserId, 300)
	if strings.Contains(userLoginDTO.RedirectUrl, "?") {
		uservo.RedirectUrl = userLoginDTO.RedirectUrl + "&code=" + code
	} else {
		uservo.RedirectUrl = userLoginDTO.RedirectUrl + "?code=" + code
	}

	uservo.Pwd = ""
	model.Success(ctx, uservo, "登录成功")
}

func LoginPhoneAndWechat(ctx *gin.Context) {
	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	account := userLoginDTO.Account
	smsCaptcha := userLoginDTO.SmsCaptcha
	uuid := userLoginDTO.Uuid

	openid := redis.Get("qrcode" + uuid + "openid")

	if len(smsCaptcha) != 6 {
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	smsCaptcha_exist := redis.Get("code" + account)

	// 将编码转换为字符串
	log.Printf("smsCaptcha_exist : %v", smsCaptcha_exist)
	// 判断密码是否正确
	if smsCaptcha != smsCaptcha_exist {
		model.Response(ctx, http.StatusBadRequest, 400, nil, "验证码错误")
		return
	}

	DB := sqls.DB()
	var user model.User
	// 数据验证
	DB.Where("mobile = ?", account).First(&user)
	if user.UserId == "" {
		//create user

		randomStr := util.GenerateCode()
		password := randomStr
		userName := "游客" + randomStr

		// 加密密码
		hasedPassword := util.MD5(password)
		// 创建用户
		newUser := model.User{
			UserId:      util.Myuuid(),
			CreateDt:    time.Now(),
			UpdateDt:    nil,
			UserName:    userName,
			LastLoginDt: time.Now(),
			NickName:    userName,
			Mobile:      account,
			Openid:      openid,
			Email:       nil,
			Pwd:         hasedPassword,
		}

		var s = util.Worker1{}

		var sysUserRole = model.SysUserRole{
			Id:     int(s.GetId()),
			UserId: newUser.UserId,
			RoleId: 2,
		}
		DB.Create(&sysUserRole)
		DB.Create(&newUser)

		user = newUser

	} else {
		//update user
		result := map[string]interface{}{}

		if !StringUtils.IsEmpty(openid) {
			result["openid"] = openid
			DB := sqls.DB()
			DB.Model(&user).Where("user_id", user.UserId).Updates(result)
		}
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		model.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	uservo := model.UserVO{}

	util.SimpleCopyProperties(&uservo, &user)
	uservo.AccessToken = token

	userBytes, _ := json.Marshal(uservo)
	redis.Set("qrcode"+uuid, "done", 5000)
	redis.SetBytes("qrcode"+uuid+"user", userBytes, 5000)

	code := util.RandomString(10)
	redis.Set("code"+code, uservo.UserId, 300)
	redirectUrl := redis.Get("qrcode"+uuid+"redirectUrl") + "?code=" + code
	uservo.RedirectUrl = redirectUrl
	uservo.Pwd = ""
	model.Success(ctx, uservo, "登录成功")
}

//func SendSms(ctx *gin.Context) {
//
//	// 获取参数
//	var userLoginDTO = model.UserLoginDTO{}
//	ctx.Bind(&userLoginDTO)
//	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
//	// 输出换行符
//	fmt.Printf("\n")
//
//	account := userLoginDTO.Account
//
//	code := rand.Intn(900000) + 100000
//
//	redis.Set("code"+account, util.IntToString(code), 300)
//
//	var regionId = flag.String("regionId", "cn-hangzhou", "区域标识")
//	var accessKeyId = flag.String("id", config.Instance.DAYU.APP_KEY, "accessKeyId")
//	var accessKeySecret = flag.String("secret", config.Instance.DAYU.APP_SECRET, "accessKeySecret")
//	var verifyCode = flag.String("code", util.IntToString(code), "验证码")
//	var phoneNumbers = flag.Int("phonenumbers", util.StringToInt(account), "手机号")
//	flag.Parse()
//
//	if *phoneNumbers <= 0 {
//		panic(fmt.Errorf("invalid phonenumbers"))
//	}
//
//	client, err := dysmsapi.NewClientWithAccessKey(*regionId, *accessKeyId, *accessKeySecret)
//	if err != nil {
//		panic(err)
//	}
//
//	params, _ := json.Marshal(map[string]interface{}{
//		"code": verifyCode,
//	})
//
//	request := dysmsapi.CreateSendSmsRequest()
//	request.Scheme = "https"
//	request.TemplateCode = config.Instance.DAYU.SMS_TEMPLATE_CODE
//	request.SignName = "新云网"
//	request.TemplateParam = string(params)
//	request.PhoneNumbers = strconv.Itoa(*phoneNumbers)
//
//	resp, err := client.SendSms(request)
//	if err != nil {
//		log.Printf("send sms failed resp=%v err=%v", resp, err)
//		panic(err)
//	}
//
//	if !resp.IsSuccess() {
//		log.Printf("send sms failed resp=%v err=%v", resp, err)
//		panic(fmt.Errorf("failed: unknown reason"))
//	}
//
//	model.Success(ctx, gin.H{}, "查询成功")
//
//}

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: accessKeyId,
		// 必填，您的 AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
func SendSmsAsync(account string, code string) {
	client, _err := CreateClient(tea.String(config.Instance.DAYU.APP_KEY), tea.String(config.Instance.DAYU.APP_SECRET))
	if _err != nil {
		log.Printf("send sms failed err=%v", _err)
		return
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(account),
		TemplateCode:  tea.String(config.Instance.DAYU.SMS_TEMPLATE_CODE),
		TemplateParam: tea.String("{\"code\":\"" + code + "\"}"),
		SignName:      tea.String("猿盟"),
	}
	runtime := &aliutil.RuntimeOptions{}
	resp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)

	if err != nil {
		log.Printf("send sms failed err=%v", err)
		return
	}

	log.Printf("send sms success resp=%v", resp)
}

func SendSms(ctx *gin.Context) {
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	account := userLoginDTO.Account
	code := util.GenerateCode()
	redis.Set("code"+account, code, 300)

	go SendSmsAsync(account, code)

	model.Success(ctx, gin.H{}, "查询成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	data := user.(model.User)
	fmt.Printf("user [%+v]", user)
	uservo := model.UserVO{UserName: data.UserName}
	fmt.Printf("uservo [%+v]", uservo)
	model.Success(ctx, uservo, "")
}

func LoadUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	DB := sqls.DB()
	var user model.User
	DB.Where("email = ?", email).First(&user)

	uservo := model.UserVO{}
	util.SimpleCopyProperties(&uservo, &user)
	uservo.Avatar = user.AvatarUrl

	var res []model.SysRole
	DB.Table("sys_role").Select("sys_role.code").
		Joins("left join sys_user_role on sys_role.id = sys_user_role.role_id").Where("user_id = ?", user.UserId).Scan(&res)
	fmt.Println(res)
	uservo.RoleCode = res[0].Code

	model.Success(ctx, uservo, "")
}

func LoadUserById(ctx *gin.Context) {
	userId := ctx.Query("userId")
	DB := sqls.DB()
	var user model.User
	DB.Where("user_id = ?", userId).First(&user)

	uservo := model.UserVO{}
	util.SimpleCopyProperties(&uservo, &user)
	uservo.Avatar = user.AvatarUrl

	var res []model.SysRole
	DB.Table("sys_role").Select("sys_role.code").
		Joins("left join sys_user_role on sys_role.id = sys_user_role.role_id").Where("user_id = ?", user.UserId).Scan(&res)
	fmt.Println(res)
	uservo.RoleCode = res[0].Code

	model.Success(ctx, uservo, "")
}

func LoadUserByOpenId(ctx *gin.Context) {
	openId := ctx.Query("openId")
	DB := sqls.DB()
	var user model.User
	DB.Where("openid = ?", openId).First(&user)

	uservo := model.UserVO{}
	util.SimpleCopyProperties(&uservo, &user)
	uservo.Avatar = user.AvatarUrl

	var res []model.SysRole
	DB.Table("sys_role").Select("sys_role.code").
		Joins("left join sys_user_role on sys_role.id = sys_user_role.role_id").Where("user_id = ?", user.UserId).Scan(&res)
	fmt.Println(res)
	uservo.RoleCode = res[0].Code

	model.Success(ctx, uservo, "")
}

func LoadUserByIds(ctx *gin.Context) {
	result := ctx.Query("userIds")
	DB := sqls.DB()
	var users []model.User
	userIds := userIdsToArray(result)
	DB.Where("user_id in ?", userIds).Find(&users)

	uservos := copyUserProps(users)

	model.Success(ctx, uservos, "")
}

func userIdsToArray(userIds string) []int {
	userIdsArr := strings.Split(userIds, ",")
	result := make([]int, len(userIdsArr))
	for i, id := range userIdsArr {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			// 处理转换失败的情况
			panic(err)
		}
		result[i] = idInt
	}
	return result
}

func copyUserProps(users []model.User) []model.UserVO {
	uservos := make([]model.UserVO, len(users))
	for i, user := range users {
		uservo := model.UserVO{}
		util.SimpleCopyProperties(&uservo, &user)
		uservo.Avatar = user.AvatarUrl
		uservos[i] = uservo
	}
	return uservos
}

func MockInfo(ctx *gin.Context) {

	uservo := model.UserVO{}
	model.Success(ctx, uservo, "")
}

func RedirectTOUnsplash(ctx *gin.Context) {

	domain := "https://unsplash.com/oauth/authorize?"
	url2 := "client_id=uwKjSmclPhET8snMTq38-TwQqKNHDd8SWACTk-Vr9mg"
	redirect_uri := "https://pgw.punengshuo.com/api/auth/backFromU"
	url2 += "&redirect_uri=" + redirect_uri
	url2 += "&response_type=code"
	url2 += "&scope=public+read_photos"

	fmt.Printf(domain + url2)
	fmt.Println()
	fmt.Println()
	ss := url.QueryEscape(url2)
	fmt.Printf("data: s%", domain+ss)

	ctx.Redirect(http.StatusFound, domain+url2)
}

func RedirectTOWechat(ctx *gin.Context) {

	uuid := ctx.Query("uuid")

	domain := "https://open.weixin.qq.com/connect/oauth2/authorize?"
	url2 := "appid=wx70711c9b88f9c12f"
	redirect_uri := "https://clock.cuiyi.club/api/auth/backFromW/" + uuid
	url2 += "&redirect_uri=" + redirect_uri
	url2 += "&response_type=code"
	url2 += "&scope=snsapi_base#wechat_redirect"

	fmt.Printf(domain + url2)
	fmt.Println()
	fmt.Println()
	ss := url.QueryEscape(url2)
	fmt.Printf("data: s%", domain+ss)

	ctx.Redirect(http.StatusFound, domain+url2)
}

func BackFromWechat(ctx *gin.Context) {

	uuid := ctx.Param("uuid")

	inputs, err := RequestInputs(ctx)
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}
	code := inputs["code"].(string)
	fmt.Printf("data--------: %v", inputs)

	domain := "https://api.weixin.qq.com/sns/oauth2/access_token?"
	url2 := "appid=wx70711c9b88f9c12f"
	url2 += "&secret=20993710aa48342888d3a0b1755af9d6"
	url2 += "&code=" + code
	url2 += "&grant_type=authorization_code"

	content := util.Get(domain + url2)
	fmt.Println()
	fmt.Printf("token--------: %s", content)
	var wechatToken model.WechatToken

	err2 := json.Unmarshal([]byte(content), &wechatToken)
	if err2 != nil {
		fmt.Println("error:", err2)
	}
	fmt.Printf("%+v", wechatToken)

	DB := sqls.DB()
	var user model.User
	DB.Where("openid = ?", wechatToken.Openid).First(&user)

	if user.UserId != "" {

		userBytes, _ := json.Marshal(user)
		redis.Set("qrcode"+uuid, "done", 5000)
		redis.SetBytes("qrcode"+uuid+"user", userBytes, 5000)

		mycode := util.RandomString(10)
		redis.Set("code"+mycode, user.UserId, 300)

		fmt.Println("qrcode" + uuid + "redirectUrl")
		redirectUrl := redis.Get("qrcode"+uuid+"redirectUrl") + "?code=" + mycode

		fmt.Println(redirectUrl)

		ctx.Redirect(http.StatusFound, redirectUrl)
		return
	}
	redis.Set("qrcode"+uuid+"openid", wechatToken.Openid, 6000)

	var loginUrl = fmt.Sprintf("https://sso.punengshuo.com/wechat?uuid=%s", uuid)

	ctx.Redirect(http.StatusFound, loginUrl)
}

func BackFromUnsplash(ctx *gin.Context) {

	inputs, err := RequestInputs(ctx)
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}
	code := inputs["code"].(string)
	fmt.Printf("data: %v", inputs)

	posturl := "https://unsplash.com/oauth/token"
	redirect_uri := "https://pgw.punengshuo.com/api/auth/backFromU"
	jsonStr := []byte(`{ "client_id": "uwKjSmclPhET8snMTq38-TwQqKNHDd8SWACTk-Vr9mg", "client_secret": "B4_p5ZZqDLKKFF4V6XyRHsqzzLoCJ7f9tlFfFECJ_H4", 
		"redirect_uri": "` + redirect_uri + `", "code": "` + code + `", "grant_type": "authorization_code" }`)
	content := util.Post(posturl, jsonStr, "application/json")
	fmt.Printf("data: s%", content)
	uservo := model.UserVO{}
	model.Success(ctx, uservo, "")
}

func GetByCodeForPuzzle(ctx *gin.Context) {

	inputs, err := RequestInputs(ctx)
	if err != nil {
		log.Printf(" error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 500, nil, "东西不对啊1")
		return
	}
	code := inputs["code"].(string)
	fmt.Printf("data: %v", inputs)

	posturl := "https://api.punengshuo.com/api/auth/token"
	jsonStr := []byte(`{ "client_id": "puzzle_xxx", "client_secret": "bbbbb", 
             "code": "` + code + `", "grant_type": "authorization_code" }`)
	content := util.Post(posturl, jsonStr, "application/json")
	fmt.Printf("data: %s", content)
	var result model.ResultVO

	err2 := json.Unmarshal([]byte(content), &result)
	if err2 != nil {
		fmt.Println("error:", err2)
	}
	if result.Code != 200 {
		model.Fail(ctx, gin.H{}, "")
		return
	}
	model.Success(ctx, result.Data, "")
}

func GetToken(ctx *gin.Context) {

	inputs, err := RequestInputs(ctx)
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "东西不对啊2")
		return
	}
	code := inputs["code"].(string)
	client_id := inputs["client_id"].(string)
	client_secret := inputs["client_secret"].(string)

	if util.IsEmptyString(code) {
		model.Response(ctx, http.StatusBadRequest, 500, nil, "东西不对啊3")
		return
	}
	if !util.CompareStrings(client_id, "puzzle_xxx") {
		model.Response(ctx, http.StatusBadRequest, 500, nil, "东西不对啊4")
		return
	}

	if !util.CompareStrings(client_secret, "bbbbb") {
		model.Response(ctx, http.StatusBadRequest, 500, nil, "东西不对啊5")
		return
	}
	userId := redis.Get("code" + code)
	if util.IsEmptyString(userId) {
		model.Response(ctx, http.StatusBadRequest, 500, nil, "东西不对啊6")
		return
	}
	user := model.User{
		UserId: userId,
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		model.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	model.Success(ctx, token, "")
}

func GetQrcode(c *gin.Context) {

	width, _ := strconv.Atoi(c.Query("width"))
	height, _ := strconv.Atoi(c.Query("height"))
	uuid := c.Query("uuid")
	var url = fmt.Sprintf("https://clock.cuiyi.club/api/auth/redirectW?uuid=%s", uuid)

	qrCode, err := GenerateQRCode(url, width, height)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	redis.Set("qrcode"+uuid, "init", 3000)

	c.Header("Content-Type", "image/png")
	c.Writer.Write(qrCode)
}

func GetUUID(ctx *gin.Context) {

	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	ctx.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	redirectUrl := userLoginDTO.RedirectUrl

	uuid := util.Myuuid()
	ctx.SetCookie("uuid", uuid, 0, "/", "", false, false)
	var url = fmt.Sprintf("https://clock.cuiyi.club/api/auth/redirectW?uuid=%s", uuid)
	redis.Set("qrcode"+uuid, "init", 3000)
	redis.Set("qrcode"+uuid+"redirectUrl", redirectUrl, 3000)
	model.Success(ctx, url, "success")
}

func CheckQrcode(c *gin.Context) {

	// 获取参数
	var userLoginDTO = model.UserLoginDTO{}
	c.Bind(&userLoginDTO)
	fmt.Printf("userLoginDTO：%+v", userLoginDTO)
	// 输出换行符
	fmt.Printf("\n")

	uuid := userLoginDTO.Uuid

	status := redis.Get("qrcode" + uuid)
	if status == "done" {

		userBytes := redis.GetBytes("qrcode" + uuid + "user")

		user := model.User{}
		err := json.Unmarshal(userBytes, &user)
		if err != nil {
			model.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
			log.Printf("json generate error : %v", err)
			return
		}

		// 发放token
		token, err := util.ReleaseToken(user)
		if err != nil {
			model.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
			log.Printf("token generate error : %v", err)
			return
		}

		uservo := model.UserVO{}

		util.SimpleCopyProperties(&uservo, &user)
		uservo.AccessToken = token
		code := util.RandomString(10)
		redis.Set("code"+code, uservo.UserId, 300)
		redirectUrl := redis.Get("qrcode"+uuid+"redirectUrl") + "?code=" + code
		fmt.Println("CheckQrcode " + redirectUrl)
		uservo.RedirectUrl = redirectUrl
		uservo.Pwd = ""
		model.Success(c, uservo, "登录成功")
		return
	}

	model.Response(c, http.StatusOK, 201, nil, "还未扫码")
}

// GenerateQRCode generates QR code and returns the image buffer
func GenerateQRCode(data string, width, height int) ([]byte, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	img := qr.Image(width)

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// RequestInputs 获取所有参数
func RequestInputs(c *gin.Context) (map[string]interface{}, error) {

	const defaultMemory = 32 << 20
	contentType := c.ContentType()

	var (
		dataMap  = make(map[string]interface{})
		queryMap = make(map[string]interface{})
		postMap  = make(map[string]interface{})
	)

	// @see gin@v1.7.7/binding/query.go ==> func (queryBinding) Bind(req *http.Request, obj interface{})
	for k := range c.Request.URL.Query() {
		queryMap[k] = c.Query(k)
	}

	if "application/json" == contentType {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// @see gin@v1.7.7/binding/json.go ==> func (jsonBinding) Bind(req *http.Request, obj interface{})
		if c.Request != nil && c.Request.Body != nil {
			if err := json.NewDecoder(c.Request.Body).Decode(&postMap); err != nil {
				return nil, err
			}
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	} else if "multipart/form-data" == contentType {
		// @see gin@v1.7.7/binding/form.go ==> func (formMultipartBinding) Bind(req *http.Request, obj interface{})
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	} else {
		// ParseForm 解析 URL 中的查询字符串，并将解析结果更新到 r.Form 字段
		// 对于 POST 或 PUT 请求，ParseForm 还会将 body 当作表单解析，
		// 并将结果既更新到 r.PostForm 也更新到 r.Form。解析结果中，
		// POST 或 PUT 请求主体要优先于 URL 查询字符串（同名变量，主体的值在查询字符串的值前面）
		// @see gin@v1.7.7/binding/form.go ==> func (formBinding) Bind(req *http.Request, obj interface{})
		if err := c.Request.ParseForm(); err != nil {
			return nil, err
		}
		if err := c.Request.ParseMultipartForm(defaultMemory); err != nil {
			if err != http.ErrNotMultipart {
				return nil, err
			}
		}
		for k, v := range c.Request.PostForm {
			if len(v) > 1 {
				postMap[k] = v
			} else if len(v) == 1 {
				postMap[k] = v[0]
			}
		}
	}

	var mu sync.RWMutex
	for k, v := range queryMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}
	for k, v := range postMap {
		mu.Lock()
		dataMap[k] = v
		mu.Unlock()
	}

	return dataMap, nil
}

func Javatosql(ctx *gin.Context) {
	var javabean = model.JavaBean{}
	ctx.Bind(&javabean)
	fmt.Printf("javabean：%v", javabean)
	originText := javabean.OriginText
	tableName := javabean.TableName

	originText = strings.Trim(originText, " ")
	var arr []string
	if strings.Contains(originText, "\n") {
		arr = strings.Split(originText, "\n")
	} else if strings.Contains(originText, "\n\r") {
		arr = strings.Split(originText, "\n\r")
	} else {
		arr = strings.Split(originText, "\n\r")
	}

	var ret string
	for _, s := range arr {
		ret += split(tableName, s) + "\r\n"
	}

	model.Success2(ctx, ret, "")
}

func CompareFile(ctx *gin.Context) {

	//firstFile := javabean.FirstFile
	//secondFile :=javabean.SecondFile

	file1, header, err := ctx.Request.FormFile("file1")
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}

	filename := header.Filename

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	sourceFile1, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer sourceFile1.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(sourceFile1, file1)
	if err != nil {
		log.Fatal(err)
	}

	file2, header, err := ctx.Request.FormFile("file2")
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}

	filename2 := header.Filename

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	sourceFile2, err := os.Create(filename2)
	if err != nil {
		log.Fatal(err)
	}

	defer sourceFile2.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(sourceFile2, file2)
	if err != nil {
		log.Fatal(err)
	}

	sourceFile1.Seek(0, 0)
	sourceFile2.Seek(0, 0)
	list := compareFileByLine(sourceFile1, sourceFile2)

	model.Success2(ctx, list, "")
}

//func TestThread(ctx *gin.Context) {
//
//	//firstFile := javabean.FirstFile
//	//secondFile :=javabean.SecondFile
//	sws := util.GetInstance()
//	//change := util.Change{Add: "ssss"}
//	sws.AddChange("ssss")
//	model.Success2(ctx, "ok", "")
//}

func compareFileByLine(f1, f2 *os.File) string {
	sc1 := bufio.NewScanner(f1)
	sc2 := bufio.NewScanner(f2)

	var s1 string
	var s2 string
	for {
		sc1Bool := sc1.Scan()
		sc2Bool := sc2.Scan()
		if !sc1Bool && !sc2Bool {
			break
		}
		s1 += sc1.Text() + "\n\r"
		s2 += sc2.Text() + "\n\r"
	}
	s := util.Diff(s1, s2)
	return s
}

func split(tableName string, originText string) string {
	originText = strings.Trim(originText, " ")
	arr := strings.Split(originText, "//")

	var ret string
	if len(arr) == 2 {
		ret = change(tableName, arr[0], arr[1])
	} else {
		ret = change(tableName, arr[0], "")
	}

	return ret
}

func change(tableName string, originText string, comment string) string {
	fmt.Printf("originText：%v", originText)
	originText = strings.Trim(originText, " ")
	arr := strings.Split(originText, " ")
	fmt.Printf("%q\n", arr)
	var ret string
	if arr[1] == "String" {
		/* 如果条件为 true 则执行以下语句 */
		ret = strings.Replace(arr[2], ";", "", -1)
		ret = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` varchar(10) DEFAULT NULL COMMENT '%s';", tableName, ret, comment)
	}
	return ret
}

func isTelephoneExist(db *gorm.DB, mobile string) bool {
	var user model.User
	db.Where("mobile = ?", mobile).First(&user)
	return user.UserId != ""
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.UserId != ""
}
