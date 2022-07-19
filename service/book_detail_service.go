package service

import (
	"ginEssential/dao"
	"ginEssential/model/constants"
	"strings"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

	"ginEssential/cache"
	"ginEssential/model"
)

var BookDetailService = newBookDetailService()

func newBookDetailService() *bookDetailService {
	return &bookDetailService{}
}

type bookDetailService struct {
}

func (s *bookDetailService) Get(id int64) *model.BookDetail {
	return dao.BookDetailDao.Get(sqls.DB(), id)
}

func (s *bookDetailService) Take(where ...interface{}) *model.BookDetail {
	return dao.BookDetailDao.Take(sqls.DB(), where...)
}

func (s *bookDetailService) Find(cnd *sqls.Cnd) []model.BookDetail {
	return dao.BookDetailDao.Find(sqls.DB(), cnd)
}

func (s *bookDetailService) FindOne(cnd *sqls.Cnd) *model.BookDetail {
	return dao.BookDetailDao.FindOne(sqls.DB(), cnd)
}

func (s *bookDetailService) FindPageByParams(params *params.QueryParams) (list []model.BookDetail, paging *sqls.Paging) {
	return dao.BookDetailDao.FindPageByParams(sqls.DB(), params)
}

func (s *bookDetailService) FindPageByCnd(cnd *sqls.Cnd) (list []model.BookDetail, paging *sqls.Paging) {
	return dao.BookDetailDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *bookDetailService) Create(t *model.BookDetail) error {
	return dao.BookDetailDao.Create(sqls.DB(), t)
}

func (s *bookDetailService) Update(t *model.BookDetail) error {
	if err := dao.BookDetailDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.BookDetailCache.Invalidate(t.Id)
	return nil
}

// func (s *bookDetailService) Updates(id int64, columns map[string]interface{}) error {
// 	return dao.BookDetailDao.Updates(sqls.DB(), id, columns)
// }
//
// func (s *bookDetailService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.BookDetailDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
// func (s *bookDetailService) Delete(id int64) {
// 	dao.BookDetailDao.Delete(sqls.DB(), id)
// }

// 自动完成
func (s *bookDetailService) Autocomplete(input string) []model.BookDetail {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.BookDetailDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *bookDetailService) GetByName(name string) *model.BookDetail {
	return dao.BookDetailDao.GetByName(name)
}

func (s *bookDetailService) GetBookDetails() []model.BookDetailResponse {
	list := dao.BookDetailDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ?", constants.StatusOk))

	var bookDetails []model.BookDetailResponse
	for _, bookDetail := range list {

		bookDetails = append(bookDetails,
			model.BookDetailResponse{
				Id:    bookDetail.Id,
				Lyric: bookDetail.Lyric,
			},
		)
	}
	return bookDetails
}

func (s *bookDetailService) GetBookDetailInIds(bookDetailIds []int64) []model.BookDetail {
	return dao.BookDetailDao.GetBookDetailInIds(bookDetailIds)
}

// 扫描
func (s *bookDetailService) Scan(callback func(bookDetails []model.BookDetail)) {
	var cursor int64
	for {
		list := dao.BookDetailDao.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
