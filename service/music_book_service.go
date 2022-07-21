package service

import (
	"ginEssential/dao"
	"ginEssential/model/constants"
	"ginEssential/util"
	"strings"

	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"

	"ginEssential/model"
)

var MusicBookService = newMusicBookService()

func newMusicBookService() *musicBookService {
	return &musicBookService{}
}

type musicBookService struct {
}

//func (c *musicBookService) SelectPageList(queryVo model.MusicBookDTO) (*model.PageResponse[model.MusicBook], error) {
//	p := &dao.Page[model.MusicBook]{
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
//	err := musicBookModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *musicBookService) Get(id string) *model.MusicBook {
	return dao.MusicBookDao.Get(sqls.DB(), id)
}

func (s *musicBookService) Take(where ...interface{}) *model.MusicBook {
	return dao.MusicBookDao.Take(sqls.DB(), where...)
}

func (s *musicBookService) Find(cnd *sqls.Cnd) []model.MusicBook {
	return dao.MusicBookDao.Find(sqls.DB(), cnd)
}

func (s *musicBookService) FindOne(cnd *sqls.Cnd) *model.MusicBook {
	return dao.MusicBookDao.FindOne(sqls.DB(), cnd)
}

func (s *musicBookService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.MusicBook], error) {
	return dao.MusicBookDao.FindPageByParams(sqls.DB(), params)
}

func (s *musicBookService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.MusicBook], error) {
	return dao.MusicBookDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *musicBookService) Create(t *model.MusicBook) error {
	return dao.MusicBookDao.Create(sqls.DB(), t)
}

func (s *musicBookService) Update(t *model.MusicBook) error {
	if err := dao.MusicBookDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *musicBookService) Updates(id string, columns map[string]interface{}) error {
	return dao.MusicBookDao.Updates(sqls.DB(), id, columns)
}

func (s *musicBookService) UpdateAll(id string, columns *model.MusicBook) error {
	return dao.MusicBookDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *musicBookService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.MusicBookDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *musicBookService) Delete(id int64) {
	dao.MusicBookDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *musicBookService) Autocomplete(input string) []model.MusicBook {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.MusicBookDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *musicBookService) GetByName(name string) *model.MusicBook {
	return dao.MusicBookDao.GetByName(name)
}

func (s *musicBookService) GetMusicBooks() []model.MusicBookVO {
	list := dao.MusicBookDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ?", constants.StatusOk))

	var musicBooks []model.MusicBookVO
	for _, musicBook := range list {

		bookvo := model.MusicBookVO{}
		util.SimpleCopyProperties(&bookvo, &musicBook)
		musicBooks = append(musicBooks, bookvo)
	}
	return musicBooks
}

func (s *musicBookService) GetMusicBookInIds(musicBookIds []int64) []model.MusicBook {
	return dao.MusicBookDao.GetMusicBookInIds(musicBookIds)
}
