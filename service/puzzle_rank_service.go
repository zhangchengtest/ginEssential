package service

import (
	"ginEssential/dao"
	"ginEssential/model/constants"
	"ginEssential/util"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"strings"

	"ginEssential/model"
)

var PuzzleRankService = newPuzzleRankService()

func newPuzzleRankService() *puzzleRankService {
	return &puzzleRankService{}
}

type puzzleRankService struct {
}

//func (c *puzzleRankService) SelectPageList(queryVo model.PuzzleRankDTO) (*model.PageResponse[model.PuzzleRank], error) {
//	p := &dao.Page[model.PuzzleRank]{
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
//	err := puzzleRankModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *puzzleRankService) Get(id string) *model.PuzzleRank {
	return dao.PuzzleRankDao.Get(sqls.DB(), id)
}

func (s *puzzleRankService) Take(where ...interface{}) *model.PuzzleRank {
	return dao.PuzzleRankDao.Take(sqls.DB(), where...)
}

func (s *puzzleRankService) Find(cnd *sqls.Cnd) []model.PuzzleRank {
	return dao.PuzzleRankDao.Find(sqls.DB(), cnd)
}

func (s *puzzleRankService) FindOne(cnd *sqls.Cnd) *model.PuzzleRank {
	return dao.PuzzleRankDao.FindOne(sqls.DB(), cnd)
}

func (s *puzzleRankService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.PuzzleRank], error) {
	return dao.PuzzleRankDao.FindPageByParams(sqls.DB(), params)
}

func (s *puzzleRankService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.PuzzleRank], error) {
	return dao.PuzzleRankDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *puzzleRankService) Create(t *model.PuzzleRank) error {
	return dao.PuzzleRankDao.Create(sqls.DB(), t)
}

func (s *puzzleRankService) Update(t *model.PuzzleRank) error {
	if err := dao.PuzzleRankDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *puzzleRankService) Updates(id string, columns map[string]interface{}) error {
	return dao.PuzzleRankDao.Updates(sqls.DB(), id, columns)
}

func (s *puzzleRankService) UpdateAll(id string, columns *model.PuzzleRank) error {
	return dao.PuzzleRankDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *puzzleRankService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.PuzzleRankDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *puzzleRankService) Delete(id string) {
	dao.PuzzleRankDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *puzzleRankService) Autocomplete(input string) []model.PuzzleRank {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.PuzzleRankDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *puzzleRankService) GetByName(name string) *model.PuzzleRank {
	return dao.PuzzleRankDao.GetByName(name)
}

func (s *puzzleRankService) GetPuzzleRanks(url string) []model.PuzzleRankVO {
	list := dao.PuzzleRankDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url).Asc("spend_time").Asc("step"))
	var puzzleRanks []model.PuzzleRankVO
	sort := 1
	for _, puzzleRank := range list {

		bookvo := model.PuzzleRankVO{}
		util.SimpleCopyProperties(&bookvo, &puzzleRank)
		bookvo.Sort = sort
		sort++
		puzzleRanks = append(puzzleRanks, bookvo)
	}
	return puzzleRanks
}

func (s *puzzleRankService) GetPuzzleRankInIds(puzzleRankIds []int64) []model.PuzzleRank {
	return dao.PuzzleRankDao.GetPuzzleRankInIds(puzzleRankIds)
}
