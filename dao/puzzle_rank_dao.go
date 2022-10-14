package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var PuzzleRankDao = newPuzzleRankDao()

func newPuzzleRankDao() *puzzleRankDao {
	return &puzzleRankDao{}
}

type puzzleRankDao struct {
}

func (r *puzzleRankDao) Get(db *gorm.DB, id string) *model.PuzzleRank {
	ret := &model.PuzzleRank{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *puzzleRankDao) Take(db *gorm.DB, where ...interface{}) *model.PuzzleRank {
	ret := &model.PuzzleRank{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *puzzleRankDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.PuzzleRank) {
	cnd.Find(db, &list)
	return
}

func (r *puzzleRankDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.PuzzleRank {
	ret := &model.PuzzleRank{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *puzzleRankDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.PuzzleRank], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *puzzleRankDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.PuzzleRank], error) {

	page := &model.Page[model.PuzzleRank]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.PuzzleRank{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.PuzzleRank{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *puzzleRankDao) Create(db *gorm.DB, t *model.PuzzleRank) (err error) {
	err = db.Create(t).Error
	return
}

func (r *puzzleRankDao) Update(db *gorm.DB, t *model.PuzzleRank) (err error) {
	err = db.Save(t).Error
	return
}

func (r *puzzleRankDao) Updates(db *gorm.DB, id string, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PuzzleRank{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *puzzleRankDao) UpdateAll(db *gorm.DB, id string, columns *model.PuzzleRank) (err error) {
	err = db.Model(&model.PuzzleRank{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *puzzleRankDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PuzzleRank{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *puzzleRankDao) Delete(db *gorm.DB, id string) {
	db.Delete(&model.PuzzleRank{}, "book_id = ?", id)
}

func (r *puzzleRankDao) GetPuzzleRankInIds(puzzleRankIds []int64) []model.PuzzleRank {
	if len(puzzleRankIds) == 0 {
		return nil
	}
	var puzzleRanks []model.PuzzleRank
	sqls.DB().Where("id in (?)", puzzleRankIds).Find(&puzzleRanks)
	return puzzleRanks
}

func (r *puzzleRankDao) GetByName(name string) *model.PuzzleRank {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
