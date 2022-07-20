package controller

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"net/http"
	"time"
)

func AddMusicBook(ctx *gin.Context) {
	DB := sqls.DB()
	var book = model.MusicBook{}

	user := ctx.MustGet("user").(model.User)

	ctx.Bind(&book)
	if book.BookId != "" {
		old := model.MusicBook{}
		DB.Where("book_id = ?", book.BookId).First(&old)
		if old.BookId == "" {
			ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "not found"})
			return
		}

		book.UpdateDt = time.Now()
		DB.Where("book_id = ?", book.BookId).Updates(&book)
		model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
	} else {
		book.BookId = util.Myuuid()
		book.CreateDt = time.Now()
		book.UpdateDt = time.Now()
		book.CreateBy = user.UserId

		fmt.Printf("book：%v", book)

		DB.Create(&book)
		model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
	}

}

func SearchMusicBook(ctx *gin.Context) {

	var queryVo model.MusicBookDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	params := params.NewQueryParams(ctx)
	if queryVo.BookTitle != "" {
		params.Like("book_title", queryVo.BookTitle)
	}
	params.Page(queryVo.PageNum, queryVo.PageSize).Desc("create_dt")

	pageResponse, e := service.MusicBookService.FindPageByParams(params)

	if e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	model.Success(ctx, pageResponse, "查询成功")
}

func DetailMusicBook(ctx *gin.Context) {

	book := service.MusicBookService.Get(ctx.Param("id"))

	model.Success(ctx, book, "查询成功")
}
