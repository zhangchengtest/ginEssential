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

var PieceContentService = newPieceContentService()

func newPieceContentService() *pieceContentService {
	return &pieceContentService{}
}

type pieceContentService struct {
}

//func (c *pieceContentService) SelectPageList(queryVo model.PieceContentDTO) (*model.PageResponse[model.PieceContent], error) {
//	p := &dao.Page[model.PieceContent]{
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
//	err := pieceContentModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *pieceContentService) Get(id string) *model.PieceContent {
	return dao.PieceContentDao.Get(sqls.DB(), id)
}

func (s *pieceContentService) GetByPhaseId(id string) *model.PieceContent {
	return dao.PieceContentDao.GetByPhaseId(sqls.DB(), id)
}

func (s *pieceContentService) Take(where ...interface{}) *model.PieceContent {
	return dao.PieceContentDao.Take(sqls.DB(), where...)
}

func (s *pieceContentService) Find(cnd *sqls.Cnd) []model.PieceContent {
	return dao.PieceContentDao.Find(sqls.DB(), cnd)
}

func (s *pieceContentService) FindOne(cnd *sqls.Cnd) *model.PieceContent {
	return dao.PieceContentDao.FindOne(sqls.DB(), cnd)
}

func (s *pieceContentService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.PieceContent], error) {
	return dao.PieceContentDao.FindPageByParams(sqls.DB(), params)
}

func (s *pieceContentService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.PieceContent], error) {
	return dao.PieceContentDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *pieceContentService) Create(t *model.PieceContent) error {
	return dao.PieceContentDao.Create(sqls.DB(), t)
}

func (s *pieceContentService) Update(t *model.PieceContent) error {
	if err := dao.PieceContentDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *pieceContentService) Updates(id int64, columns map[string]interface{}) error {
	return dao.PieceContentDao.Updates(sqls.DB(), id, columns)
}

func (s *pieceContentService) UpdateAll(id string, columns *model.PieceContent) error {
	return dao.PieceContentDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *pieceContentService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.PieceContentDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *pieceContentService) Delete(id int64) {
	dao.PieceContentDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *pieceContentService) Autocomplete(input string) []model.PieceContent {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.PieceContentDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *pieceContentService) GetByName(name string) *model.PieceContent {
	return dao.PieceContentDao.GetByName(name)
}

func (s *pieceContentService) GetPieceContents(url string) []model.PieceContent {
	list := dao.PieceContentDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url))

	return list
}

func (s *pieceContentService) GetPieceContentsGroup() []model.PieceContentVO {

	t := time.Now()
	date := strftime.Format(t, "%Y-%m-%d")
	fmt.Println(date)
	t2, _ := strftime.Parse(date+" 00:00:00", "%Y-%m-%d %H:%M:%S")

	var list []model.PieceContent
	sqls.DB().Model(&model.PieceContent{}).Select("url").Where("create_dt > ?", t2).Group("url").Find(&list)
	var pieceContents []model.PieceContentVO
	for _, pieceContent := range list {

		bookvo := model.PieceContentVO{}
		util.SimpleCopyProperties(&bookvo, &pieceContent)
		pieceContents = append(pieceContents, bookvo)
	}
	return pieceContents
}

func (s *pieceContentService) GetPieceContentInIds(pieceContentIds []int64) []model.PieceContent {
	return dao.PieceContentDao.GetPieceContentInIds(pieceContentIds)
}

func (s *pieceContentService) FindByContentType(book_id string, content_type int) []model.BookPiece {

	return dao.PieceContentDao.SelectByContentType(sqls.DB(), book_id, content_type)
}

func (s *pieceContentService) FindByBreakFlag(book_id string, start_order int, end_order int) []model.BookPiece {

	return dao.PieceContentDao.SelectByBreakFlag(sqls.DB(), book_id, start_order, end_order)
}

func (r *pieceContentService) SelectMax(book_id string) int {
	return dao.PieceContentDao.SelectMax(sqls.DB(), book_id)
}
