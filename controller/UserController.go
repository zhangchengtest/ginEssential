package controller

import (
	"bufio"
	"fmt"
	"ginEssential/common"
	"ginEssential/dao"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Register(ctx *gin.Context) {
	DB := dao.GetDB()
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
	if !util.VerifyEmailFormat(email) {
		// 自己封装过后
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "邮箱格式不对")
		return
	}

	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	log.Println(email, password)
	// 判断手机号是否存在
	if isEmailExist(DB, email) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	// 加密密码
	hasedPassword := util.MD5(password)
	// 创建用户
	newUser := model.User{
		UserId:   util.Myuuid(),
		CreateDt: time.Now(),
		UpdateDt: nil,
		UserName: userName,
		NickName: userName,
		Email:    email,
		Pwd:      hasedPassword,
	}

	var s = util.Worker1{}

	var sysUserRole = model.SysUserRole{
		Id:     int(s.GetId()),
		UserId: newUser.UserId,
		RoleId: 2,
	}
	DB.Create(&sysUserRole)
	DB.Create(&newUser)

	common.SendRegister(userName, email)

	response.Success(ctx, gin.H{"userName": userName}, "注册成功")
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
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	DB := dao.GetDB()
	var user model.User
	// 数据验证
	if util.VerifyEmailFormat(account) {
		DB.Where("email = ?", account).First(&user)
		if user.UserId == "" {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	} else {
		DB.Where("user_name = ?", account).First(&user)
		if user.UserId == "" {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
			return
		}
	}

	newSig := util.MD5(password) //转成加密编码
	// 将编码转换为字符串
	log.Printf("newSig : %v", newSig)
	// 判断密码是否正确
	if user.Pwd != newSig {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}

	uservo := model.UserVO{}

	util.SimpleCopyProperties(&uservo, &user)
	uservo.AccessToken = token

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
	response.Success(ctx, uservo, "登录成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	uservo := model.UserVO{}
	util.SimpleCopyProperties(&uservo, &user)
	response.Success(ctx, uservo, "")
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

	response.Success2(ctx, ret, "")
}

func CompareFile(ctx *gin.Context) {

	//firstFile := javabean.FirstFile
	//secondFile :=javabean.SecondFile

	file1, header, err := ctx.Request.FormFile("file1")
	if err != nil {
		log.Printf("get file error: %s", err)
		response.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
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
		response.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
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

	response.Success2(ctx, list, "")
}

func TestThread(ctx *gin.Context) {

	//firstFile := javabean.FirstFile
	//secondFile :=javabean.SecondFile
	sws := util.GetInstance()
	//change := util.Change{Add: "ssss"}
	sws.AddChange("ssss")
	response.Success2(ctx, "ok", "")
}

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
