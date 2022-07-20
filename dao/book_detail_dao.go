package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var BookDetailDao = newBookDetailDao()

func newBookDetailDao() *bookDetailDao {
	return &bookDetailDao{}
}

type bookDetailDao struct {
}

func (r *bookDetailDao) Get(db *gorm.DB, id int64) *model.BookDetail {
	ret := &model.BookDetail{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookDetailDao) Take(db *gorm.DB, where ...interface{}) *model.BookDetail {
	ret := &model.BookDetail{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookDetailDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.BookDetail) {
	cnd.Find(db, &list)
	return
}

func (r *bookDetailDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.BookDetail {
	ret := &model.BookDetail{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *bookDetailDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []model.BookDetail, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *bookDetailDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []model.BookDetail, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.BookDetail{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *bookDetailDao) Create(db *gorm.DB, t *model.BookDetail) (err error) {
	err = db.Create(t).Error
	return
}

func (r *bookDetailDao) Update(db *gorm.DB, t *model.BookDetail) (err error) {
	err = db.Save(t).Error
	return
}

func (r *bookDetailDao) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.BookDetail{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *bookDetailDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.BookDetail{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *bookDetailDao) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.BookDetail{}, "id = ?", id)
}

func (r *bookDetailDao) GetBookDetailInIds(bookDetailIds []int64) []model.BookDetail {
	if len(bookDetailIds) == 0 {
		return nil
	}
	var bookDetails []model.BookDetail
	sqls.DB().Where("id in (?)", bookDetailIds).Find(&bookDetails)
	return bookDetails
}

func (r *bookDetailDao) GetByName(name string) *model.BookDetail {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
