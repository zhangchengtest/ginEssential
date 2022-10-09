package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var PuzzlePieceDao = newPuzzlePieceDao()

func newPuzzlePieceDao() *puzzlePieceDao {
	return &puzzlePieceDao{}
}

type puzzlePieceDao struct {
}

func (r *puzzlePieceDao) Get(db *gorm.DB, id string) *model.PuzzlePiece {
	ret := &model.PuzzlePiece{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *puzzlePieceDao) Take(db *gorm.DB, where ...interface{}) *model.PuzzlePiece {
	ret := &model.PuzzlePiece{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *puzzlePieceDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.PuzzlePiece) {
	cnd.Find(db, &list)
	return
}

func (r *puzzlePieceDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.PuzzlePiece {
	ret := &model.PuzzlePiece{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *puzzlePieceDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.PuzzlePiece], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *puzzlePieceDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.PuzzlePiece], error) {

	page := &model.Page[model.PuzzlePiece]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.PuzzlePiece{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.PuzzlePiece{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *puzzlePieceDao) Create(db *gorm.DB, t *model.PuzzlePiece) (err error) {
	err = db.Create(t).Error
	return
}

func (r *puzzlePieceDao) Update(db *gorm.DB, t *model.PuzzlePiece) (err error) {
	err = db.Save(t).Error
	return
}

func (r *puzzlePieceDao) Updates(db *gorm.DB, id string, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PuzzlePiece{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *puzzlePieceDao) UpdateAll(db *gorm.DB, id string, columns *model.PuzzlePiece) (err error) {
	err = db.Model(&model.PuzzlePiece{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *puzzlePieceDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PuzzlePiece{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *puzzlePieceDao) Delete(db *gorm.DB, id string) {
	db.Delete(&model.PuzzlePiece{}, "book_id = ?", id)
}

func (r *puzzlePieceDao) GetPuzzlePieceInIds(puzzlePieceIds []int64) []model.PuzzlePiece {
	if len(puzzlePieceIds) == 0 {
		return nil
	}
	var puzzlePieces []model.PuzzlePiece
	sqls.DB().Where("id in (?)", puzzlePieceIds).Find(&puzzlePieces)
	return puzzlePieces
}

func (r *puzzlePieceDao) GetByName(name string) *model.PuzzlePiece {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
