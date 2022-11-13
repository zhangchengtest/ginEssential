package dao

import (
	"ginEssential/model"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"
)

var BookPieceDao = newBookPieceDao()

func newBookPieceDao() *bookPieceDao {
	return &bookPieceDao{}
}

type bookPieceDao struct {
}

func (r *bookPieceDao) Get(db *gorm.DB, id string) *model.BookPiece {
	ret := &model.BookPiece{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookPieceDao) Take(db *gorm.DB, where ...interface{}) *model.BookPiece {
	ret := &model.BookPiece{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *bookPieceDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.BookPiece) {
	cnd.Find(db, &list)
	return
}

func (r *bookPieceDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.BookPiece {
	ret := &model.BookPiece{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *bookPieceDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.BookPiece], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *bookPieceDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.BookPiece], error) {

	page := &model.Page[model.BookPiece]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.BookPiece{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.BookPiece{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *bookPieceDao) Create(db *gorm.DB, t *model.BookPiece) (err error) {
	err = db.Create(t).Error
	return
}

func (r *bookPieceDao) Update(db *gorm.DB, t *model.BookPiece) (err error) {
	err = db.Save(t).Error
	return
}

func (r *bookPieceDao) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.BookPiece{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *bookPieceDao) UpdateAll(db *gorm.DB, id string, columns *model.BookPiece) (err error) {
	err = db.Model(&model.BookPiece{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *bookPieceDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.BookPiece{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *bookPieceDao) Delete(db *gorm.DB, phase_id string) {
	db.Delete(&model.BookPiece{}, "phase_id = ?", phase_id)
}

func (r *bookPieceDao) GetBookPieceInIds(bookPieceIds []int64) []model.BookPiece {
	if len(bookPieceIds) == 0 {
		return nil
	}
	var bookPieces []model.BookPiece
	sqls.DB().Where("id in (?)", bookPieceIds).Find(&bookPieces)
	return bookPieces
}

func (r *bookPieceDao) GetByName(name string) *model.BookPiece {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
