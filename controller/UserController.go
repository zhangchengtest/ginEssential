package controller

import (
	"fmt"
	"ginEssential/common"
	"ginEssential/dto"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()
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

	// 获取参数
	name := ginBindUser.Name
	telephone := ginBindUser.Telephone
	password := ginBindUser.Password
	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		// ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		// 这个gin.H实际是type H map[string]interface{}，所以也可以写成下面这样
		// ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"code": 422, "msg": "手机号必须为11位"})

		// 自己封装过后
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 如果name为空，就随机一个10位的字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	// 加密密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	// 创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}

	DB.Create(&newUser)

	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	response.Success(ctx, gin.H{"token": token}, "注册成功")
}

func Login(ctx *gin.Context) {
	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	// 判断手机号是否存在
	DB := common.GetDB()
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回结果
	// ctx.JSON(200, gin.H{
	// 	"code":    200,
	// 	"data":    gin.H{"token": token},
	// 	"message": "注册成功",
	// })
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"code": 200,
	// 	"data": gin.H{
	// 		"user": dto.ToUserDto(user.(model.User)), // user转dto
	// 	},
	// })
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
}

func Javatosql(ctx *gin.Context) {
	//1. 使用map获取application/json请求的参数


	var javabean = model.JavaBean{}
	ctx.Bind(&javabean)
	fmt.Printf("javabean：%v", javabean)
	originText := javabean.OriginText
	tableName :=javabean.TableName

	originText = strings.Trim(originText," ")
	var arr []string
	if strings.Contains(originText, "\n") {
		fmt.Printf("change 1")
		arr = strings.Split(originText, "\n")
	}else if strings.Contains(originText, "\n\r") {
		fmt.Printf("change 2")
		arr = strings.Split(originText, "\n\r")
	}else {
		arr = strings.Split(originText, "\n\r")
	}


	var ret string
	for _, s := range arr {
		ret += split(tableName, s)+"\r\n"
	}

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"code": 200,
	// 	"data": gin.H{
	// 		"user": dto.ToUserDto(user.(model.User)), // user转dto
	// 	},
	// })
	response.Success2(ctx, ret, "")
}

func split(tableName string, originText string) string {
	originText = strings.Trim(originText," ")
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
	originText = strings.Trim(originText," ")
	arr := strings.Split(originText, " ")
	fmt.Printf("%q\n", arr)
	var ret string
	if arr[1] == "String" {
		/* 如果条件为 true 则执行以下语句 */
		ret =   strings.Replace(arr[2], ";", "", -1 )
		ret = fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` varchar(10) DEFAULT NULL COMMENT '%s';", tableName, ret, comment);
	}
	return ret
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	return user.ID != 0
}
