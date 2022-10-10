package controller

import (
	"encoding/base64"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
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
	// 创建用户
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

func QueryPuzzle(ctx *gin.Context) {

	list := service.PuzzlePieceService.GetPuzzlePiecesGroup()
	var size = len(list)
	rand.Seed(time.Now().UnixNano())
	var num = rand.Intn(size)
	puzzlePiece := list[num]
	list2 := service.PuzzlePieceService.GetPuzzlePieces(puzzlePiece.Url)

	model.Success(ctx, list2, "查询成功")
}
