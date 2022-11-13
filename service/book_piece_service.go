package service

import (
	"fmt"
	"ginEssential/dao"
	"ginEssential/model/constants"
	"ginEssential/util"
	strftime "github.com/itchyny/timefmt-go"
	"strings"
	"time"

	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"

	"ginEssential/model"
)

var BookPieceService = newBookPieceService()

func newBookPieceService() *bookPieceService {
	return &bookPieceService{}
}

type bookPieceService struct {
}

//func (c *bookPieceService) SelectPageList(queryVo model.BookPieceDTO) (*model.PageResponse[model.BookPiece], error) {
//	p := &dao.Page[model.BookPiece]{
//		CurrentPage: queryVo.PageNum,
//		PageSize:    queryVo.PageSize,
//	}
//	var query []string
//	var args []interface{}
//	if queryVo.BookTitle != "" {
//		query = append(query, "book_title LIKE ?")
//		args = append(args, "%"+queryVo.BookTitle+"%")
//	}
//	if queryVo.Author != "" {
//		query = append(query, "author LIKE ?")
//		args = append(args, "%"+queryVo.Author+"%")
//	}
//	var querystr string
//
//	for index, value := range query {
//		if index == len(query)-1 {
//			querystr += value
//		} else {
//			querystr += value + " and "
//		}
//		fmt.Printf("index: %d value: %s\n", index, value)
//	}
//
//	err := bookPieceModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *bookPieceService) Get(id string) *model.BookPiece {
	return dao.BookPieceDao.Get(sqls.DB(), id)
}

func (s *bookPieceService) Take(where ...interface{}) *model.BookPiece {
	return dao.BookPieceDao.Take(sqls.DB(), where...)
}

func (s *bookPieceService) Find(cnd *sqls.Cnd) []model.BookPiece {
	return dao.BookPieceDao.Find(sqls.DB(), cnd)
}

func (s *bookPieceService) FindOne(cnd *sqls.Cnd) *model.BookPiece {
	return dao.BookPieceDao.FindOne(sqls.DB(), cnd)
}

func (s *bookPieceService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.BookPiece], error) {
	return dao.BookPieceDao.FindPageByParams(sqls.DB(), params)
}

func (s *bookPieceService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.BookPiece], error) {
	return dao.BookPieceDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *bookPieceService) Create(t *model.BookPiece) error {
	return dao.BookPieceDao.Create(sqls.DB(), t)
}

func (s *bookPieceService) Update(t *model.BookPiece) error {
	if err := dao.BookPieceDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *bookPieceService) Updates(id int64, columns map[string]interface{}) error {
	return dao.BookPieceDao.Updates(sqls.DB(), id, columns)
}

func (s *bookPieceService) UpdateAll(id string, columns *model.BookPiece) error {
	return dao.BookPieceDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *bookPieceService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.BookPieceDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *bookPieceService) Delete(phase_id string) {
	dao.BookPieceDao.Delete(sqls.DB(), phase_id)
}

// 自动完成
func (s *bookPieceService) Autocomplete(input string) []model.BookPiece {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.BookPieceDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *bookPieceService) GetByName(name string) *model.BookPiece {
	return dao.BookPieceDao.GetByName(name)
}

func (s *bookPieceService) GetBookPieces(url string) []model.BookPiece {
	list := dao.BookPieceDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url))

	return list
}

func (s *bookPieceService) GetBookPiecesGroup() []model.BookPieceVO {

	t := time.Now()
	date := strftime.Format(t, "%Y-%m-%d")
	fmt.Println(date)
	t2, _ := strftime.Parse(date+" 00:00:00", "%Y-%m-%d %H:%M:%S")

	var list []model.BookPiece
	sqls.DB().Model(&model.BookPiece{}).Select("url").Where("create_dt > ?", t2).Group("url").Find(&list)
	var bookPieces []model.BookPieceVO
	for _, bookPiece := range list {

		bookvo := model.BookPieceVO{}
		util.SimpleCopyProperties(&bookvo, &bookPiece)
		bookPieces = append(bookPieces, bookvo)
	}
	return bookPieces
}

func (s *bookPieceService) GetBookPieceInIds(bookPieceIds []int64) []model.BookPiece {
	return dao.BookPieceDao.GetBookPieceInIds(bookPieceIds)
}
