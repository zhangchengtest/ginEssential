package service

import (
	"ginEssential/dao"
	"ginEssential/model/constants"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"strings"

	"ginEssential/model"
)

var BookImgService = newBookImgService()

func newBookImgService() *bookImgService {
	return &bookImgService{}
}

type bookImgService struct {
}

//func (c *bookImgService) SelectPageList(queryVo model.BookImgDTO) (*model.PageResponse[model.BookImg], error) {
//	p := &dao.Page[model.BookImg]{
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
//	err := bookImgModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *bookImgService) Get(id string) *model.BookImg {
	return dao.BookImgDao.Get(sqls.DB(), id)
}

func (s *bookImgService) Take(where ...interface{}) *model.BookImg {
	return dao.BookImgDao.Take(sqls.DB(), where...)
}

func (s *bookImgService) Find(cnd *sqls.Cnd) []model.BookImg {
	return dao.BookImgDao.Find(sqls.DB(), cnd)
}

func (s *bookImgService) FindOne(cnd *sqls.Cnd) *model.BookImg {
	return dao.BookImgDao.FindOne(sqls.DB(), cnd)
}

func (s *bookImgService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.BookImg], error) {
	return dao.BookImgDao.FindPageByParams(sqls.DB(), params)
}

func (s *bookImgService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.BookImg], error) {
	return dao.BookImgDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *bookImgService) Create(t *model.BookImg) error {
	return dao.BookImgDao.Create(sqls.DB(), t)
}

func (s *bookImgService) Update(t *model.BookImg) error {
	if err := dao.BookImgDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *bookImgService) Updates(id string, columns map[string]interface{}) error {
	return dao.BookImgDao.Updates(sqls.DB(), id, columns)
}

func (s *bookImgService) UpdateAll(id string, columns *model.BookImg) error {
	return dao.BookImgDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *bookImgService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.BookImgDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *bookImgService) Delete(id string) {
	dao.BookImgDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *bookImgService) Autocomplete(input string) []model.BookImg {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.BookImgDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *bookImgService) GetByName(name string) *model.BookImg {
	return dao.BookImgDao.GetByName(name)
}

func (s *bookImgService) GetBookImgs(url string) []model.BookImg {
	list := dao.BookImgDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url))

	return list
}

func (s *bookImgService) GetBookImgInIds(bookImgIds []int64) []model.BookImg {
	return dao.BookImgDao.GetBookImgInIds(bookImgIds)
}
