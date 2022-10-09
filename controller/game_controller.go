package controller

import (
	"encoding/base64"
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"image"
	"image/png"
	"log"
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

	var ret string
	for _, files := range form.File {
		for _, file := range files {

			if err := ctx.SaveUploadedFile(file, config.Instance.Uploader.Local.Path+"/"+file.Filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			ret = file.Filename
		}
	}
	var i = 0
	i++
	savePuzzle(ctx.Request.FormValue("piece1"), ret, i)
	i++
	fmt.Printf("insert piece2 ")
	savePuzzle(ctx.Request.FormValue("piece2"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece3"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece4"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece5"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece6"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece7"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece8"), ret, i)
	i++
	savePuzzle(ctx.Request.FormValue("piece9"), ret, i)
	i++

	//title := ctx.Request.FormValue("title")

	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func savePuzzle(piece string, ret string, sort int) {

	//ddd, _ := base64.StdEncoding.DecodeString(piece) //成图片文件并把文件写入到buffer
	//
	//err := ioutil.WriteFile(, ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	//if err != nil {
	//	fmt.Println("Error:", err)
	//}

	//unbased, err := base64.StdEncoding.DecodeString(piece)
	//if err != nil {
	//	panic("Cannot decode b64")
	//}
	//
	//r := bytes.NewReader(unbased)
	//im, err := png.Decode(r)
	//if err != nil {
	//	panic("Bad png")
	//}
	//
	//f, err := os.OpenFile(config.Instance.Uploader.Local.Path+"/"+"output"+strconv.Itoa(sort)+".png", os.O_WRONLY|os.O_CREATE, 0777)
	//if err != nil {
	//	panic("Cannot open file")
	//}
	//
	//png.Encode(f, im)

	piece = strings.Replace(piece, "data:image/png;base64,", "", 1)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(piece))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	//bounds := m.Bounds()
	//fmt.Println(bounds, formatString)

	//Encode from image format to writer
	pngFilename := config.Instance.Uploader.Local.Path + "/" + "output" + strconv.Itoa(sort) + ".png"
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
		Content:  config.Instance.Uploader.Local.Host + "images/" + "output" + strconv.Itoa(sort) + ".png",
		Title:    ret,
		Url:      config.Instance.Uploader.Local.Host + "images/" + ret,
		Sort:     sort,
		CreateDt: time.Now(),
		CreateBy: "",
	}

	DB.Create(&newUser)
}
