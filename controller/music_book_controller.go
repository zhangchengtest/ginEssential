package controller

import (
	"fmt"
	"ginEssential/model"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/web/params"
	"net/http"
	"strings"
	"time"
)

func AddMusicBook(ctx *gin.Context) {
	var book = model.MusicBook{}

	user := ctx.MustGet("user").(model.User)

	ctx.Bind(&book)
	if book.BookId != "" {
		old := service.MusicBookService.Get(book.BookId)
		if old.BookId == "" {
			ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "not found"})
			return
		}

		book.UpdateDt = time.Now()
		//DB.Where("book_id = ?", book.BookId).Updates(&book)
		service.MusicBookService.UpdateAll(book.BookId, &book)

		model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
	} else {
		book.BookId = util.Myuuid()
		book.CreateDt = time.Now()
		book.UpdateDt = time.Now()
		book.CreateBy = user.UserId

		fmt.Printf("book：%v", book)

		service.MusicBookService.Create(&book)

		arr := strings.Split(book.Lyric, "\n")
		worker := util.NewSnow(55)
		num := 1
		for _, tt := range arr {

			bookdetail := model.BookDetail{
				Id:        worker.GetId(),
				BookId:    book.BookId,
				Lyric:     tt,
				CreateDt:  time.Now(),
				UpdateDt:  nil,
				BookOrder: num,
			}
			num++
			service.BookDetailService.Create(&bookdetail)
		}

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

func SearchMusicBookDetail(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	params := params.NewQueryParams(ctx)
	if queryVo.BookId != "" {
		params.Like("book_id", queryVo.BookId)
	}
	params.Page(queryVo.PageNum, queryVo.PageSize).Asc("book_order")

	pageResponse, e := service.BookDetailService.FindPageByParams(params)

	if e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	model.Success(ctx, pageResponse, "查询成功")
}
