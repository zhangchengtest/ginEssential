package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var BookImgDao = newBookImgDao()

func newBookImgDao() *bookImgDao {
	return &bookImgDao{}
}

type bookImgDao struct {
}

func (r *bookImgDao) Get(db *gorm.DB, id string) *model.BookImg {
	ret := &model.BookImg{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookImgDao) Take(db *gorm.DB, where ...interface{}) *model.BookImg {
	ret := &model.BookImg{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookImgDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.BookImg) {
	cnd.Find(db, &list)
	return
}

func (r *bookImgDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.BookImg {
	ret := &model.BookImg{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *bookImgDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.BookImg], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *bookImgDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.BookImg], error) {

	page := &model.Page[model.BookImg]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.BookImg{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.BookImg{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *bookImgDao) Create(db *gorm.DB, t *model.BookImg) (err error) {
	err = db.Create(t).Error
	return
}

func (r *bookImgDao) Update(db *gorm.DB, t *model.BookImg) (err error) {
	err = db.Save(t).Error
	return
}

func (r *bookImgDao) Updates(db *gorm.DB, id string, columns map[string]interface{}) (err error) {
	err = db.Model(&model.BookImg{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *bookImgDao) UpdateAll(db *gorm.DB, id string, columns *model.BookImg) (err error) {
	err = db.Model(&model.BookImg{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *bookImgDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.BookImg{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *bookImgDao) Delete(db *gorm.DB, id string) {
	db.Delete(&model.BookImg{}, "book_id = ?", id)
}

func (r *bookImgDao) GetBookImgInIds(bookImgIds []int64) []model.BookImg {
	if len(bookImgIds) == 0 {
		return nil
	}
	var bookImgs []model.BookImg
	sqls.DB().Where("id in (?)", bookImgIds).Find(&bookImgs)
	return bookImgs
}

func (r *bookImgDao) GetByName(name string) *model.BookImg {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
