package controller

import (
	"fmt"
	"ginEssential/config"
	"ginEssential/model"
	"ginEssential/render"
	"ginEssential/service"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/sjsdfg/common-lang-in-go/Cast"
	"github.com/sjsdfg/common-lang-in-go/StringUtils"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func UploadBookImg(ctx *gin.Context) {

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

	_, err = os.Stat(config.Instance.Uploader.Local.BookPath + "/" + dir)
	if os.IsNotExist(err) {
		os.Mkdir(config.Instance.Uploader.Local.BookPath+"/"+dir, os.ModePerm)
	}
	var ret string
	for _, files := range form.File {
		for _, file := range files {

			if err := ctx.SaveUploadedFile(file, config.Instance.Uploader.Local.BookPath+"/"+dir+"/"+file.Filename); err != nil {
				ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			ret = file.Filename
		}
	}

	bookId := ctx.PostForm("bookId")
	authorType := Cast.ToInt(ctx.PostForm("authorType"))

	DB := sqls.DB()
	var s = util.Worker1{}
	// 创建图
	newUser := model.BookImg{
		Id:         s.GetId(),
		BookId:     bookId,
		Title:      ret,
		Url:        config.Instance.Uploader.Local.Host + "musicBook/" + dir + "/" + ret,
		AuthorType: authorType,
		CreateDt:   time.Now(),
		CreateBy:   "",
	}

	DB.Create(&newUser)

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func buildContent(ctx *gin.Context, arr []model.BookPiece) model.PieceDetailVO {

	var contents []model.BookPiece
	var lyrics []model.BookPiece
	var uppoints_arr []model.BookPiece
	var line1s []model.BookPiece
	var line2s []model.BookPiece
	var connections []model.BookPiece

	var downpoints_arr []model.BookPiece

	var indents []model.BookPiece
	for _, tag := range arr {
		if tag.ContentType == 1 {
			contents = append(contents, tag)
		} else if tag.ContentType == 2 {
			lyrics = append(lyrics, tag)
		} else if tag.ContentType == 3 {
			uppoints_arr = append(uppoints_arr, tag)
		} else if tag.ContentType == 4 {
			line1s = append(line1s, tag)
		} else if tag.ContentType == 5 {
			line2s = append(line2s, tag)
		} else if tag.ContentType == 7 {
			connections = append(connections, tag)
		} else if tag.ContentType == 6 {
			downpoints_arr = append(downpoints_arr, tag)
		} else if tag.ContentType == 8 {
			indents = append(indents, tag)
		}

	}

	detailvo := model.PieceDetailVO{}

	detailvo.BookContent = buildpiece(ctx, contents)
	detailvo.Lyric = buildpiece(ctx, lyrics)
	detailvo.UpPoints = buildpiece(ctx, uppoints_arr)

	xstart1, xstop1 := buildline(ctx, line1s)
	detailvo.Line1xstart = xstart1
	detailvo.Line1xstop = xstop1

	xstart2, xstop2 := buildline(ctx, line2s)
	detailvo.Line2xstart = xstart2
	detailvo.Line2xstop = xstop2

	xstartConnection, xstopConnection := buildconnection(ctx, connections)
	detailvo.Connectionxstart = xstartConnection
	detailvo.Connectionxstop = xstopConnection

	detailvo.DownPoints = buildpiece(ctx, downpoints_arr)
	detailvo.Indent = buildpiece(ctx, indents)

	return detailvo

}
func DetailMusicBook(ctx *gin.Context) {

	book := service.MusicBookService.Get(ctx.Param("id"))
	book2 := render.BuildBook(book)

	bookvo := model.BookPieceVO{}

	params1 := params.NewQueryParams(ctx)
	params1.Eq("book_id", ctx.Param("id"))
	params1.Eq("break_flag", 1)
	params1.Asc("book_order")

	list := service.PieceContentService.Find(&params1.Cnd)

	start_order := 0
	end_order := 100000
	var details []model.PieceDetailVO
	for _, tag := range list {

		end_order = tag.BookOrder
		arr := service.PieceContentService.FindByBreakFlag(ctx.Param("id"), start_order, end_order)

		detailvo := buildContent(ctx, arr)

		details = append(details, detailvo)

		start_order = tag.BookOrder
	}

	max_order := getMaxOrder(ctx.Param("id"))
	if end_order < max_order-1 {
		arr := service.PieceContentService.FindByBreakFlag(ctx.Param("id"), end_order, max_order)

		detailvo := buildContent(ctx, arr)

		details = append(details, detailvo)
	}

	bookvo.List = details
	book2.PieceAll = bookvo
	getMyUrl(ctx, book2)
	log.Println("aaaa")
	log.Println(book2.MyUrl)
	getOtherUrl(ctx, book2)
	model.Success(ctx, book2, "查询成功")
}

//func main() {
//	book := test2()
//	test1(*book)
//
//}

func getMyUrl(ctx *gin.Context, book2 *model.MusicBookVO) {
	params2 := params.NewQueryParams(ctx)
	params2.Eq("book_id", ctx.Param("id"))
	params2.Eq("author_type", 0)
	params2.Desc("create_dt")
	imgs := service.BookImgService.Find(&params2.Cnd)

	if imgs != nil && len(imgs) > 0 {
		book2.MyUrl = imgs[0].Url
		log.Println("bbbbb")
		log.Println(book2.MyUrl)
	}
}

func getOtherUrl(ctx *gin.Context, book2 *model.MusicBookVO) {
	params2 := params.NewQueryParams(ctx)
	params2.Eq("book_id", ctx.Param("id"))
	params2.Eq("author_type", 1)
	params2.Desc("create_dt")
	imgs := service.BookImgService.Find(&params2.Cnd)

	if imgs != nil && len(imgs) > 0 {
		book2.OtherUrl = imgs[0].Url
	}

}

func chunk(array []model.BookPiece, size int) [][]model.BookPiece {
	length := len(array)
	if length < 1 {
		return [][]model.BookPiece{}
	}
	index := 0
	resIndex := 0
	aa := length/size + 1
	result := make([][]model.BookPiece, aa)

	for {
		if index < (length - 1) {
			if index+size < (length - 1) {
				result[resIndex] = array[index:(index + size)]
				resIndex++
				index = index + size
			} else {
				result[resIndex] = array[index:(length - 1)]
				resIndex++
				index = length - 1
			}

		} else {
			break
		}

	}

	return result
}

func buildParent(ctx *gin.Context, content_type int) [][]model.BookPiece {
	list := service.PieceContentService.FindByContentType(ctx.Param("id"), content_type)
	arr := chunk(list, 16)
	return arr
}

func buildpiece(ctx *gin.Context, list []model.BookPiece) string {

	str := ""
	for index, tag := range list {
		str = str + tag.BookContent

		if (index+1)%4 == 0 {
			str = str + "~"
		} else {
			str = str + "@"
		}
	}

	return str
}

func buildconnection(ctx *gin.Context, list []model.BookPiece) ([]int, []int) {

	prelen := 0
	var xstart []int
	var xstop []int
	for _, tag := range list {
		bb := strings.Split(tag.BookContent, "")

		for index2, tag2 := range bb {
			if tag2 != "*" {
				xstart = append(xstart, prelen+index2)
				num, _ := strconv.Atoi(tag2)
				xstop = append(xstop, prelen+index2+num)
			}
		}
		prelen = prelen + len(tag.BookContent) + 1
	}

	if xstart == nil {
		xstart = []int{}
	}

	if xstop == nil {
		xstop = []int{}
	}

	return xstart, xstop
}

func buildline(ctx *gin.Context, list []model.BookPiece) ([]int, []int) {

	prelen := 0
	var xstart []int
	var xstop []int
	for _, tag := range list {
		bb := strings.Split(tag.BookContent, "")
		isstart := false

		for index2, tag2 := range bb {
			if tag2 == "_" && !isstart {
				isstart = true
				xstart = append(xstart, prelen+index2)
			}
			if tag2 == "_" && isstart {
				if index2+1 == len(bb) {
					xstop = append(xstop, prelen+index2)
				}
				continue
			}
			if tag2 == "*" && isstart {
				isstart = false
				xstop = append(xstop, prelen+index2-1)
			}
		}
		prelen = prelen + len(tag.BookContent) + 1
	}

	return xstart, xstop
}

func UpdateBookPiece(ctx *gin.Context) {

	var piece = model.BookPieceDTO{}
	ctx.ShouldBindJSON(&piece)
	fmt.Printf("piece：%+v", piece)

	if StringUtils.IsNotEmpty(piece.PhaseId) {

		theone := service.PieceContentService.GetByPhaseId(piece.PhaseId)

		service.PieceContentService.Updates(theone.Id, map[string]interface{}{
			"break_flag": piece.BreakFlag,
			"update_dt":  time.Now(),
		})

		params := params.NewQueryParams(ctx)
		params.Eq("phase_id", piece.PhaseId)
		params.Asc("content_type")

		list := service.BookPieceService.Find(&params.Cnd)

		updateBook(list[0].Id, piece.BookContent)
		updateBook(list[1].Id, piece.Lyric)
		updateBook(list[2].Id, piece.UpPoints)
		updateBook(list[3].Id, piece.Line1)
		updateBook(list[4].Id, piece.Line2)
		updateBook(list[5].Id, piece.DownPoints)
		updateBook(list[6].Id, piece.Connection)
		updateBook(list[7].Id, piece.Indent)

		model.Success(ctx, gin.H{"status": "ok"}, "更新成功")
		return
	}

	if StringUtils.IsEmpty(piece.BookContent) {
		log.Println("data not correct1")
		model.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "data not correct")
		return
	}
	phaseId := util.Myuuid()
	order := getMaxOrder(piece.BookId)
	createContent(piece.BookId, phaseId, order, piece.BreakFlag)
	createBook(piece.BookContent, 1, piece.BookId, phaseId)
	createBook(piece.Lyric, 2, piece.BookId, phaseId)
	createBook(piece.UpPoints, 3, piece.BookId, phaseId)
	createBook(piece.Line1, 4, piece.BookId, phaseId)
	createBook(piece.Line2, 5, piece.BookId, phaseId)
	createBook(piece.DownPoints, 6, piece.BookId, phaseId)
	createBook(piece.Connection, 7, piece.BookId, phaseId)
	createBook(piece.Indent, 8, piece.BookId, phaseId)

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func TestBookPiece(ctx *gin.Context) {

	var queryVo model.BookDetailDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	list := service.PieceContentService.FindByContentType(queryVo.BookId, 6)
	for _, tag := range list {
		str := tag.BookContent
		data := strings.Replace(str, "•", "1", -1)
		log.Println(data)
		service.BookPieceService.Updates(tag.Id, map[string]interface{}{
			"book_content": data,
			"update_dt":    time.Now(),
		})
	}

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func StickUp(ctx *gin.Context) {

	var piece = model.PieceContentDTO{}
	ctx.ShouldBindJSON(&piece)
	fmt.Printf("piece：%+v", piece)

	params := params.NewQueryParams(ctx)
	params.Eq("book_id", piece.BookId)
	params.Asc("book_order")

	list := service.PieceContentService.Find(&params.Cnd)

	var arr []model.PieceContent

	var theone model.PieceContent
	for _, tag := range list {

		if tag.PhaseId != piece.PhaseId {
			arr = append(arr, tag)
		} else {
			theone = tag
		}
	}

	service.PieceContentService.Updates(theone.Id, map[string]interface{}{
		"book_order": 1,
		"update_dt":  time.Now(),
	})

	for index, tag := range arr {
		service.PieceContentService.Updates(tag.Id, map[string]interface{}{
			"book_order": index + 2,
			"update_dt":  time.Now(),
		})
	}

	model.Success(ctx, gin.H{"status": "ok"}, "置顶成功")
}

func createContent(bookId string, phaseId string, order int, breakFlag int) {
	worker := util.NewSnow(55)
	bookdetail := model.PieceContent{
		Id:        worker.GetId(),
		BookId:    bookId,
		PhaseId:   phaseId,
		BreakFlag: breakFlag,
		BookOrder: order,
		CreateDt:  time.Now(),
		UpdateDt:  nil,
	}
	service.PieceContentService.Create(&bookdetail)
}

func createBook(data string, dataType int, bookId string, phaseId string) {
	worker := util.NewSnow(55)
	bookdetail := model.BookPiece{
		Id:          worker.GetId(),
		BookId:      bookId,
		BookContent: data,
		ContentType: dataType,
		PhaseId:     phaseId,
		CreateDt:    time.Now(),
		UpdateDt:    nil,
	}
	service.BookPieceService.Create(&bookdetail)
}

func updateBook(id int64, data string) {
	service.BookPieceService.Updates(id, map[string]interface{}{
		"book_content": data,
		"update_dt":    time.Now(),
	})
}

func getMaxOrder(book_id string) int {

	num := service.PieceContentService.SelectMax(book_id)
	return num
}

func SearchPieces(ctx *gin.Context) {

	var queryVo model.BookPieceDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}
	fmt.Printf("BookOrder：%+v", queryVo.BookOrder)

	list := service.PieceContentService.FindByContentType(queryVo.BookId, 1)

	list2 := service.PieceContentService.FindByContentType(queryVo.BookId, 2)

	log.Println("len----")
	log.Println(len(list))
	ll := len(list)/4 + 1
	arr := make([][]model.BookPieceVO, ll)

	num := 0
	var arr2 []model.BookPieceVO
	for index, tag := range list {
		log.Println(index)
		aa := *BuildBookPiece(&tag)
		aa.Lyric = list2[index].BookContent
		arr2 = append(arr2, aa)
		if (index+1)%4 == 0 {
			log.Println("huan hang----")
			arr[num] = make([]model.BookPieceVO, 4)
			copy(arr[num], arr2)
			num++
			arr2 = []model.BookPieceVO{}
		}
		log.Println(arr2)
	}
	log.Println(arr)
	if num == ll-1 {
		log.Println("huan hang----")
		arr[num] = make([]model.BookPieceVO, 4)
		copy(arr[num], arr2)
		num++
		arr2 = []model.BookPieceVO{}
	}
	log.Println("arr----")
	log.Println(arr)

	var books []model.BookPieceVO2

	for _, arr_child := range arr {
		bookvo := model.BookPieceVO2{}
		for index, dd := range arr_child {
			if index == 0 {
				bookvo.Val1 = dd.BookContent
				bookvo.Id1 = dd.PhaseId
				bookvo.Lyric1 = dd.Lyric
			} else if index == 1 {
				bookvo.Val2 = dd.BookContent
				bookvo.Id2 = dd.PhaseId
				bookvo.Lyric2 = dd.Lyric
			} else if index == 2 {
				bookvo.Val3 = dd.BookContent
				bookvo.Id3 = dd.PhaseId
				bookvo.Lyric3 = dd.Lyric
			} else if index == 3 {
				bookvo.Val4 = dd.BookContent
				bookvo.Id4 = dd.PhaseId
				bookvo.Lyric4 = dd.Lyric
			}
		}
		books = append(books, bookvo)
	}

	model.Success(ctx, gin.H{
		"contents": books,
	}, "查询成功")
}

func DeletePiece(ctx *gin.Context) {

	var queryVo model.BookPieceDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	service.BookPieceService.Delete(queryVo.PhaseId)
	model.Success(ctx, gin.H{"status": "ok"}, "删除成功")
}

func CopyPiece(ctx *gin.Context) {

	var queryVo model.BookPieceDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	params := params.NewQueryParams(ctx)
	params.Eq("phase_id", queryVo.PhaseId)
	params.Asc("content_type")

	list := service.BookPieceService.Find(&params.Cnd)

	phaseId := util.Myuuid()
	order := getMaxOrder(list[0].BookId)

	//theone := service.PieceContentService.GetByPhaseId(piece.PhaseId)

	createContent(list[0].BookId, phaseId, order, 0)
	createBook(list[0].BookContent, 1, list[0].BookId, phaseId)
	createBook(list[1].BookContent, 2, list[1].BookId, phaseId)
	createBook(list[2].BookContent, 3, list[2].BookId, phaseId)
	createBook(list[3].BookContent, 4, list[3].BookId, phaseId)
	createBook(list[4].BookContent, 5, list[4].BookId, phaseId)
	createBook(list[5].BookContent, 6, list[5].BookId, phaseId)
	createBook(list[6].BookContent, 7, list[6].BookId, phaseId)
	createBook(list[7].BookContent, 8, list[7].BookId, phaseId)

	model.Success(ctx, gin.H{"status": "ok"}, "复制成功")

}

func SearchPiecesByPhaseId(ctx *gin.Context) {

	var queryVo model.BookPieceDTO
	if e := ctx.ShouldBindJSON(&queryVo); e != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 300, "msg": "参数错误"})
		return
	}

	theone := service.PieceContentService.GetByPhaseId(queryVo.PhaseId)

	params := params.NewQueryParams(ctx)
	params.Eq("phase_id", queryVo.PhaseId)
	params.Asc("content_type")

	list := service.BookPieceService.Find(&params.Cnd)
	//books := BuildBookPieces(list)

	var books []model.BookPieceVO3
	chooseVal := list[0].BookContent

	books = append(books, createPiece(list[6], "连音"))
	books = append(books, createPiece(list[2], "重音"))
	books = append(books, createPiece(list[3], "一线"))
	books = append(books, createPiece(list[4], "二线"))
	books = append(books, createPiece(list[5], "低音"))
	books = append(books, createPiece(list[1], "歌词"))
	books = append(books, createPiece(list[7], "缩进"))

	model.Success(ctx, gin.H{
		"chooseVal": chooseVal,
		"breakFlag": theone.BreakFlag,
		"contents":  books,
	}, "查询成功")
}

func createPiece(tag model.BookPiece, name string) model.BookPieceVO3 {
	bb := strings.Split(tag.BookContent, "")
	bookvo := model.BookPieceVO3{
		Val0: name,
	}
	if len(bb) > 0 {
		if bb[0] != "0" {
			bookvo.Val1 = bb[0]
		}
	}
	if len(bb) > 1 {
		if bb[1] != "0" {
			bookvo.Val2 = bb[1]
		}
	}
	if len(bb) > 2 {
		if bb[2] != "0" {
			bookvo.Val3 = bb[2]
		}
	}
	if len(bb) > 3 {
		if bb[3] != "0" {
			bookvo.Val4 = bb[3]
		}
	}

	return bookvo

}

func BuildBookPiece(tag *model.BookPiece) *model.BookPieceVO {
	if tag == nil {
		return nil
	}
	bookvo := model.BookPieceVO{
		Id:          strconv.FormatInt(tag.Id, 10),
		BookContent: tag.BookContent,
		PhaseId:     tag.PhaseId,
	}
	bookvo.ShowClass = "normal_span"
	return &bookvo
}
