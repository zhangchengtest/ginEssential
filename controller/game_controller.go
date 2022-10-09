package controller

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"net/http"
	"time"
)

func UploadSplitImages(ctx *gin.Context) {
	DB := sqls.DB()
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

			if err := ctx.SaveUploadedFile(file, "D:/test/images/"+file.Filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			ret += "\r\n"
		}
	}

	piece := ctx.Request.FormValue("piece1")
	title := ctx.Request.FormValue("title")

	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	var s = util.Worker1{}
	// 创建用户
	newUser := model.PuzzlePiece{
		Id:       s.GetId(),
		Content:  piece,
		Title:    title,
		Sort:     1,
		CreateDt: time.Now(),
		CreateBy: "",
	}

	DB.Create(&newUser)

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}
