package controller

import (
	"encoding/json"
	"fmt"
	"ginEssential/model"
	"ginEssential/render"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sjsdfg/common-lang-in-go/StringUtils"
	"github.com/zhangchengtest/simple/web/params"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SaveMusicBook(ctx *gin.Context) {
	var book = model.MusicBook{}

	user := ctx.MustGet("user").(model.User)

	ctx.Bind(&book)
	ly := book.Lyric
	arrly := strings.Split(ly, "\n")
	var resly string
	for _, tt := range arrly {
		if StringUtils.IsEmpty(tt) {
			continue
		}
		resly += strings.TrimSpace(tt) + "\n"
	}
	book.Lyric = resly

	if book.BookId != "" {

		old := service.MusicBookService.Get(book.BookId)
		if old.BookId == "" {
			ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "not found"})
			return
		}
		params := params.NewQueryParams(ctx)
		params.Eq("book_id", book.BookId)
		params.Asc("book_order")
		bookdetails := service.BookDetailService.Find(&params.Cnd)
		if len(bookdetails) == 0 {
			create(book.Lyric, book.BookId)
			service.MusicBookService.UpdateAll(book.BookId, &book)
			model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
			return
		}

		num := 1
		arr2 := util.DiffToArr(old.Lyric, book.Lyric)
		var arr3 []string
		for _, tt := range arr2 {

			if strings.Contains(tt, "+") {
				if len(arr3) == 0 {
					arr3 = append(arr3, tt)
				} else {
					sss := arr3[len(arr3)-1]
					sss += "," + tt
					arr3[len(arr3)-1] = sss
				}
			} else {
				arr3 = append(arr3, tt)
			}
		}
		var arrFinal []model.BookDetail
		var arrDel []model.BookDetail
		for _, tt := range arr3 {

			//fmt.Printf(strconv.Itoa(num) + " " + tt + " " + bookdetails[num-1].BookId + " " + bookdetails[num-1].Lyric)
			if strings.Contains(tt, "+") || strings.Contains(tt, "-") {
				arr4 := strings.Split(tt, ",")
				for _, tt2 := range arr4 {
					fmt.Printf(tt)
					fmt.Printf("\n")
					fmt.Printf(tt2)
					fmt.Printf("\n")
					if tt2 == "" {
						continue
					}
					if strings.Contains(tt2, "-") {
						arrDel = append(arrDel, bookdetails[num-1])
					}
					if strings.Contains(tt2, "o") {
						if tt != bookdetails[num-1].Lyric {
							bookdetails[num-1].Lyric = strings.TrimSpace(strings.Replace(tt2, "o", "", -1))
						}
						arrFinal = append(arrFinal, bookdetails[num-1])
					}
					if strings.Contains(tt2, "+") {
						arrFinal = append(arrFinal, model.BookDetail{
							Lyric: strings.TrimSpace(strings.Replace(tt2, "+", "", -1)),
						})
					}
				}
			} else {
				fmt.Printf("hi")
				res := strings.TrimSpace(strings.Replace(tt, "o", "", -1))
				fmt.Printf(res)
				fmt.Printf("\n")
				fmt.Printf(strconv.Itoa(num - 1))
				fmt.Printf("\n")
				fmt.Printf("hi end")
				if len(bookdetails) <= num-1 {
					arrFinal = append(arrFinal, model.BookDetail{
						Lyric: res,
					})
				} else {
					bookdetails[num-1].Lyric = res
					arrFinal = append(arrFinal, bookdetails[num-1])
				}

			}
			fmt.Printf("\n")
			num++
		}

		//arr := strings.Split(old.Lyric, "\n")
		////worker := util.NewSnow(55)
		//for _, tt := range arr {
		//	fmt.Printf(tt)
		//	fmt.Printf("\n")
		//	num++
		//}

		for _, tt := range arrDel {
			service.BookDetailService.Delete(tt.Id)
			num++
		}

		book.UpdateDt = time.Now()
		//DB.Where("book_id = ?", book.BookId).Updates(&book)
		service.MusicBookService.UpdateAll(book.BookId, &book)
		worker := util.NewSnow(55)
		num = 1

		for _, tt := range arrFinal {
			if tt.Id != 0 {
				tt.BookOrder = num
				service.BookDetailService.Updates(tt.Id, map[string]interface{}{
					"book_order": num,
					"lyric":      tt.Lyric,
				})
			} else {
				tt.Id = worker.GetId()
				tt.BookId = book.BookId
				tt.CreateDt = time.Now()
				tt.UpdateDt = nil
				tt.BookOrder = num
				service.BookDetailService.Create(&tt)
			}

			num++
		}

		model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
	} else {
		book.BookId = util.Myuuid()
		book.CreateDt = time.Now()
		book.UpdateDt = time.Now()
		book.CreateBy = user.UserId

		fmt.Printf("book：%v", book)

		service.MusicBookService.Create(&book)
		create(book.Lyric, book.BookId)
		model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
	}

}

func create(lyric, booId string) {

	arr := strings.Split(lyric, "\n")
	worker := util.NewSnow(55)
	num := 1
	for _, tt := range arr {

		bookdetail := model.BookDetail{
			Id:        worker.GetId(),
			BookId:    booId,
			Lyric:     tt,
			CreateDt:  time.Now(),
			UpdateDt:  nil,
			BookOrder: num,
		}
		num++
		service.BookDetailService.Create(&bookdetail)
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

func DeleteMusicBook(ctx *gin.Context) {

	service.MusicBookService.Delete(ctx.Param("id"))

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", ctx.Param("id"))
	params.Asc("book_order")

	list := service.BookDetailService.Find(&params.Cnd)

	for _, book := range list {
		service.BookDetailService.Delete(book.Id)
	}

	model.Success(ctx, nil, "删除成功")
}

func SearchMusicBookDetail(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", queryVo.BookId)
	params.Page(queryVo.PageNum, queryVo.PageSize).Asc("book_order")

	pageResponse, e := service.BookDetailService.FindPageByParams(params)

	if e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	model.Success(ctx, pageResponse, "查询成功")
}

func SearchOneMusicBookDetail(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", queryVo.BookId)
	params.Asc("book_order")

	list := service.BookDetailService.Find(&params.Cnd)
	books := render.BuildBookDetails(list)
	var result model.BookDetailVO
	var prev model.BookDetailVO
	//var lyrics string
	if len(books) == 0 {

	} else if queryVo.Id == "" {
		for index, book := range books {
			if len(book.BookContent) <= 0 {
				result = books[index]
				if index-1 >= 0 {
					prev = books[index-1]
				}
				break
			}
		}
	} else {
		for index, book := range books {
			if book.Id == queryVo.Id {
				result = books[index]
				if index-1 >= 0 {
					prev = books[index-1]
				}
			}
		}
	}
	var newBooks []model.BookDetailVO
	for _, book := range books {
		if book.BookContent != "" {
			if result.Id == book.Id {
				book.ShowClass = "red_span"
			} else {
				book.ShowClass = "green_span"
			}

		} else {
			if result.Id == book.Id {
				book.ShowClass = "red_span"
			} else {
				book.ShowClass = "normal_span"
			}
		}
		newBooks = append(newBooks, book)
	}
	//fmt.Printf("v+%", books)
	//fmt.Printf("v+%", newBooks)
	//book := service.MusicBookService.Get(queryVo.BookId)

	model.Success(ctx, gin.H{
		"prev":       prev,
		"bookDetail": result,
		"lyrics":     newBooks,
	}, "查询成功")
}

func UpdateMusicBookConent(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	num, _ := strconv.ParseInt(queryVo.Id, 10, 64)
	bookDetail := service.BookDetailService.Get(num)

	service.BookDetailService.Updates(num, map[string]interface{}{
		"book_content": queryVo.BookContent,
	})

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", bookDetail.BookId)
	params.Asc("book_order")

	list := service.BookDetailService.Find(&params.Cnd)
	var content string
	for _, book := range list {
		content += book.BookContent + " "
	}

	service.MusicBookService.Updates(bookDetail.BookId, map[string]interface{}{
		"book_content": content,
	})

	model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
}

func UpdateMusicBookLyric(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	b, _ := json.Marshal(queryVo)
	logrus.Info(string(b))
	num, _ := strconv.ParseInt(queryVo.Id, 10, 64)
	bookDetail := service.BookDetailService.Get(num)

	service.BookDetailService.Updates(num, map[string]interface{}{
		"lyric": queryVo.Lyric,
	})

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", bookDetail.BookId)
	params.Asc("book_order")

	list := service.BookDetailService.Find(&params.Cnd)
	var lyric string
	for _, book := range list {
		if len(book.Lyric) > 0 {
			lyric += book.Lyric + "\n"
		}
	}

	service.MusicBookService.Updates(bookDetail.BookId, map[string]interface{}{
		"lyric": lyric,
	})

	model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
}
