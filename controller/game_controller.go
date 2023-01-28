package controller

import (
	"encoding/base64"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/redis"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/Scorpio69t/jpush-api-golang-client"
	"github.com/gin-gonic/gin"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/zhangchengtest/simple/sqls"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func UploadSplitImages(ctx *gin.Context) {

	// 获取所有图片
	form, err := ctx.MultipartForm()
	if err != nil {
		return
	}
	if len(form.File) <= 0 {
		return
	}

	t := time.Now()
	dir := strftime.Format(t, "%Y%m%d%H%M%S")

	_, err = os.Stat(config.Instance.Uploader.Local.Path + "/" + dir)
	if os.IsNotExist(err) {
		os.Mkdir(config.Instance.Uploader.Local.Path+"/"+dir, os.ModePerm)
	}

	var ret string
	for _, files := range form.File {
		for _, file := range files {

			if err := ctx.SaveUploadedFile(file, config.Instance.Uploader.Local.Path+"/"+dir+"/"+file.Filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			ret = file.Filename
		}
	}
	var i = 0
	i++
	savePuzzle(ctx.Request.FormValue("piece1"), ret, i, dir)
	i++
	fmt.Printf("insert piece2 ")
	savePuzzle(ctx.Request.FormValue("piece2"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece3"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece4"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece5"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece6"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece7"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece8"), ret, i, dir)
	i++
	savePuzzle(ctx.Request.FormValue("piece9"), ret, i, dir)
	i++

	//title := ctx.Request.FormValue("title")

	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func savePuzzle(piece string, ret string, sort int, dir string) {

	piece = strings.Replace(piece, "data:image/png;base64,", "", 1)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(piece))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	//bounds := m.Bounds()
	//fmt.Println(bounds, formatString)

	//Encode from image format to writer
	pngFilename := config.Instance.Uploader.Local.Path + "/" + dir + "/" + "output" + strconv.Itoa(sort) + ".png"
	f, err := os.OpenFile(pngFilename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	err = png.Encode(f, m)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Png file", pngFilename, "created")

	DB := sqls.DB()
	var s = util.Worker1{}
	// 创建图
	newUser := model.PuzzlePiece{
		Id:       s.GetId(),
		Content:  config.Instance.Uploader.Local.Host + "images/" + dir + "/" + "output" + strconv.Itoa(sort) + ".png",
		Title:    ret,
		Url:      config.Instance.Uploader.Local.Host + "images/" + dir + "/" + ret,
		Sort:     sort,
		CreateDt: time.Now(),
		CreateBy: "",
	}

	DB.Create(&newUser)
}

func SavePuzzleRank(ctx *gin.Context) {
	DB := sqls.DB()
	var rank = model.PuzzleRank{}
	ctx.Bind(&rank)
	fmt.Printf("rank：%v", rank)

	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	var s = util.Worker1{}
	newUser := model.PuzzleRank{
		Id:        s.GetId(),
		Username:  rank.Username,
		Title:     rank.Title,
		Url:       rank.Url,
		SpendTime: rank.SpendTime,
		Step:      rank.Step,
		CreateBy:  "system",
		CreateDt:  time.Now(),
	}

	DB.Create(&newUser)

	sendRank(rank.Username, rank.Url)

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func PushTest(ctx *gin.Context) {
	var pf jpush.Platform
	pf.Add(jpush.ANDROID)
	// pf.All()

	// Audience: tag
	var at jpush.Audience
	id := []string{"1"}
	at.SetID(id)
	// at.All()

	// Notification
	var n jpush.Notification
	n.SetAlert("alert")
	n.SetAndroid(&jpush.AndroidNotification{Alert: "alert", Title: "title"})

	// Message
	var m jpush.Message
	m.MsgContent = "This is a message"
	m.Title = "Hello"

	// PayLoad
	payload := jpush.NewPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&at)
	payload.SetNotification(&n)
	payload.SetMessage(&m)

	// Send
	c := jpush.NewJPushClient("f4d450cd2dec6b2568fa74b9", "48aea09d1a15f68a34f375eb ") // appKey and masterSecret can be gotten from https://www.jiguang.cn/
	data, err := payload.Bytes()
	fmt.Printf("%s\n", string(data))
	if err != nil {
		panic(err)
	}
	res, err := c.Push(data)
	if err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("ok: %v\n", res)
	}
}

func SavePlaneRank(ctx *gin.Context) {
	DB := sqls.DB()
	var rank = model.PlaneRank{}
	ctx.Bind(&rank)
	fmt.Printf("rank：%v", rank)

	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	var s = util.Worker1{}
	// 创建用户
	newUser := model.PlaneRank{
		Id:       s.GetId(),
		Username: rank.Username,
		Coin:     rank.Coin,
		CreateBy: "system",
		CreateDt: time.Now(),
	}

	DB.Create(&newUser)

	model.Success(ctx, gin.H{"status": "ok"}, "新增排名")
}

func sendRank(username, myurl string) {

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

		msg := username + "在这个图上拼的很好\n" + config.Instance.PuzzleUrl + "/#/puzzle/index?randomUrl=" + url.QueryEscape(myurl) + "&ginToken=" + url.QueryEscape(user.UserName)
		data := message.MediaText{
			Content: msg,
		}
		if user.UserName != username {
			customerMessage := message.CustomerMessage{
				Msgtype: message.MsgTypeText,
				Text:    &data,
				ToUser:  user.Openid,
			}
			officialAccount.GetCustomerMessageManager().Send(&customerMessage)
		}
	}

}

func QueryPuzzleRank(ctx *gin.Context) {
	val := ctx.Request.FormValue("url")
	fmt.Println()
	fmt.Printf(val)
	fmt.Println()
	list := service.PuzzleRankService.GetPuzzleRanks(val)

	model.Success(ctx, list, "查询成功")
}

func Visit(ctx *gin.Context) {
	val := ctx.Request.FormValue("name")
	redis.Visit(val)
	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func ModifyUsername(ctx *gin.Context) {
	model.Success(ctx, gin.H{"status": "ok"}, "SUCCESS")
}

func Nickname(ctx *gin.Context) {
	model.Success(ctx, service.GetFullName(), "新增成功")
}

func QueryPlaneRank(ctx *gin.Context) {
	list := service.PlaneRankService.GetRanks()

	model.Success(ctx, list, "查询成功")
}

func QueryPuzzle(ctx *gin.Context) {

	arr := service.PuzzlePieceService.GetPuzzlePiecesGroup()
	var size = len(arr)
	rand.Seed(time.Now().UnixNano())
	var num = rand.Intn(size)
	puzzlePiece := arr[num]
	list := service.PuzzlePieceService.GetPuzzlePieces(puzzlePiece.Url)

	var puzzlePieces []string
	var orders []int
	for _, puzzlePiece := range list {
		//if puzzlePiece.Sort != 9 {
		puzzlePieces = append(puzzlePieces, puzzlePiece.Content)
		orders = append(orders, puzzlePiece.Sort)
		//}
	}

	vo := model.PuzzlePieceVO2{
		Url:     puzzlePiece.Url,
		Piecces: puzzlePieces,
		Orders:  orders,
	}

	model.Success(ctx, vo, "查询成功")
}

func QueryPuzzleByUrl(ctx *gin.Context) {

	val := ctx.Request.FormValue("url")
	order_type := ctx.Request.FormValue("order_type")
	fmt.Printf(val)

	//var queryVo model.PuzzlePieceDTO
	//if e := ctx.ShouldBindJSON(&queryVo); e != nil {
	//	ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
	//	return
	//}
	var list []model.PuzzlePiece
	if order_type == "random" {
		list = service.PuzzlePieceService.GetPuzzlePiecesRandom(val)
	} else {
		list = service.PuzzlePieceService.GetPuzzlePieces(val)
	}

	var puzzlePieces []string
	var orders []int
	for _, puzzlePiece := range list {
		//if puzzlePiece.Sort != 9 {
		puzzlePieces = append(puzzlePieces, puzzlePiece.Content)
		orders = append(orders, puzzlePiece.Sort)
		//}
	}

	vo := model.PuzzlePieceVO2{
		Url:     val,
		Piecces: puzzlePieces,
		Orders:  orders,
	}

	model.Success(ctx, vo, "查询成功")
}
