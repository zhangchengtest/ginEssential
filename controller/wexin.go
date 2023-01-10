package controller

import (
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/sjsdfg/common-lang-in-go/StringUtils"
	"github.com/zhangchengtest/simple/sqls"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func LoginByWeixinCode(ctx *gin.Context) {

	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &miniConfig.Config{
		AppID:     "wx029106fe29ab6dde",
		AppSecret: "c54f17c5a7cb10246225a17ce3f43d7d",
		//EncodingAESKey: "xxxx",
		Cache: memory,
	}

	wc := wechat.NewWechat()
	mini := wc.GetMiniProgram(cfg)
	a := mini.GetAuth()

	var queryVo model.Code2SessionRequest
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	resp, err := a.Code2Session(queryVo.JsCode)
	if err != nil {
		log.Printf("send sms failed resp=%v err=%v", resp, err)
		panic(err)
	}
	fmt.Println(resp)

	DB := sqls.DB()

	var olduser model.User
	DB.Where("openid = ?", resp.OpenID).First(&olduser)
	if olduser.UserId == "" {
		// 加密密码
		hasedPassword := util.MD5("123456")
		// 创建用户
		newUser := model.User{
			UserId:        util.Myuuid(),
			CreateDt:      time.Now(),
			UpdateDt:      nil,
			UserName:      "test",
			NickName:      "test",
			Email:         "test@qq.com",
			Pwd:           hasedPassword,
			Openid:        resp.OpenID,
			LastWrongPwDt: nil,
			LastLoginDt:   time.Now(),
			AvatarUrl:     "https://thirdwx.qlogo.cn/mmopen/vi_32/POgEwh4mIHO4nibH0KlMECNjjGxQUq24ZEaGT4poC6icRiccVGKSyXwibcPq4BWmiaIGuG1icwxaQX6grC9VemZoJ8rg/132",
		}

		var s = util.Worker1{}

		var sysUserRole = model.SysUserRole{
			Id:     int(s.GetId()),
			UserId: newUser.UserId,
			RoleId: 2,
		}

		DB.Create(&sysUserRole)
		DB.Create(&newUser)

		// 创建图
		autvo := model.WechatAuthVO{
			Username:  newUser.UserName,
			NickName:  newUser.NickName,
			AvatarUrl: newUser.AvatarUrl,
		}

		token, _ := util.ReleaseToken(newUser)

		model.Success(ctx, gin.H{"userInfo": autvo, "token": token}, "查询成功")
	} else {
		// 创建图
		autvo := model.WechatAuthVO{
			Username:  olduser.UserName,
			NickName:  olduser.NickName,
			AvatarUrl: olduser.AvatarUrl,
		}
		token, _ := util.ReleaseToken(olduser)
		model.Success(ctx, gin.H{"userInfo": autvo, "token": token}, "查询成功")
	}

}

func UserDetail(ctx *gin.Context) {

	olduser, _ := ctx.MustGet("user").(model.User)
	// 创建图
	autvo := model.WechatAuthVO{
		Username:  olduser.UserName,
		NickName:  olduser.NickName,
		AvatarUrl: olduser.AvatarUrl,
	}

	model.Success(ctx, gin.H{"userInfo": autvo}, "查询成功")

}

func ModifyUser(ctx *gin.Context) {

	olduser, _ := ctx.MustGet("user").(model.User)

	var userDTO model.UserDTO
	if e := ctx.ShouldBindJSON(&userDTO); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	result := map[string]interface{}{}
	if !StringUtils.IsEmpty(userDTO.NickName) {
		result["nick_name"] = userDTO.NickName
	}
	if !StringUtils.IsEmpty(userDTO.AvatarUrl) {
		result["avatar_url"] = userDTO.AvatarUrl
	}

	DB := sqls.DB()
	DB.Model(&olduser).Where("user_id", olduser.UserId).Updates(result)
	// 创建图
	autvo := model.WechatAuthVO{
		Username:  olduser.UserName,
		NickName:  olduser.NickName,
		AvatarUrl: userDTO.AvatarUrl,
	}

	model.Success(ctx, gin.H{"userInfo": autvo}, "查询成功")
}

func UploadFile(ctx *gin.Context) {
	file1, header, err := ctx.Request.FormFile("upfile")
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}

	filename := header.Filename

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	sourceFile1, err := os.Create(config.Instance.Uploader.Local.LogoPath + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	// 将file的内容拷贝到out
	_, err = io.Copy(sourceFile1, file1)
	if err != nil {
		log.Fatal(err)
	}

	defer sourceFile1.Close()

	model.Success(ctx, gin.H{"url": config.Instance.Uploader.Local.Host + "logoPath" + "/" + filename}, "查询成功")
}
