package controller

import (
	"fmt"
	"ginEssential/dao"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var musicBookService service.MusicBookService

func AddMusicBook(ctx *gin.Context) {
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
	var book = model.MusicBook{}

	ctx.Bind(&book)
	book.BookId = util.Myuuid()
	if book.CreateBy == "" {
		book.CreateBy = "1"
	}
	book.CreateDt = time.Now()
	book.UpdateDt = time.Now()

	fmt.Printf("book：%v", book)

	DB.Create(&book)

	response.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func SearchMusicBook(ctx *gin.Context) {

	var queryVo model.MusicBookDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	pageResponse, e := musicBookService.SelectPageList(queryVo)
	if e != nil {
		ctx.JSON(http.StatusOK, model.Response{Code: 400, Msg: "操作失败"})
		return
	}
	ctx.JSON(
		http.StatusOK,
		model.Response{Code: 200, Msg: "操做成功", Data: pageResponse},
	)
}

func DetailMusicBook(ctx *gin.Context) {
	DB := dao.GetDB()
	var book model.MusicBook

	if err := DB.Where("book_id = ?", ctx.Param("id")).First(&book).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	ctx.JSON(
		http.StatusOK,
		model.Response{Code: 200, Msg: "操做成功", Data: book},
	)
}
