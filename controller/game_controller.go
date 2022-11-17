package controller

import (
	"encoding/base64"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/redis"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/zhangchengtest/simple/sqls"
	"image"
	"image/png"
	"log"
	"math/rand"
	"net/http"
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
	// 创建用户
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

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
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

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
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
	fmt.Printf(val)

	//var queryVo model.PuzzlePieceDTO
	//if e := ctx.ShouldBindJSON(&queryVo); e != nil {
	//	ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
	//	return
	//}
	list := service.PuzzlePieceService.GetPuzzlePiecesRandom(val)

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
