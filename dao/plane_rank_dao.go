package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var PlaneRankDao = newPlaneRankDao()

func newPlaneRankDao() *planeRankDao {
	return &planeRankDao{}
}

type planeRankDao struct {
}

func (r *planeRankDao) Get(db *gorm.DB, id string) *model.PlaneRank {
	ret := &model.PlaneRank{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *planeRankDao) Take(db *gorm.DB, where ...interface{}) *model.PlaneRank {
	ret := &model.PlaneRank{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *planeRankDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.PlaneRank) {
	cnd.Find(db, &list)
	return
}

func (r *planeRankDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.PlaneRank {
	ret := &model.PlaneRank{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *planeRankDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.PlaneRank], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *planeRankDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.PlaneRank], error) {

	page := &model.Page[model.PlaneRank]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.PlaneRank{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.PlaneRank{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *planeRankDao) Create(db *gorm.DB, t *model.PlaneRank) (err error) {
	err = db.Create(t).Error
	return
}

func (r *planeRankDao) Update(db *gorm.DB, t *model.PlaneRank) (err error) {
	err = db.Save(t).Error
	return
}

func (r *planeRankDao) Updates(db *gorm.DB, id string, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PlaneRank{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *planeRankDao) UpdateAll(db *gorm.DB, id string, columns *model.PlaneRank) (err error) {
	err = db.Model(&model.PlaneRank{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *planeRankDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PlaneRank{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *planeRankDao) Delete(db *gorm.DB, id string) {
	db.Delete(&model.PlaneRank{}, "book_id = ?", id)
}

func (r *planeRankDao) GetPlaneRankInIds(planeRankIds []int64) []model.PlaneRank {
	if len(planeRankIds) == 0 {
		return nil
	}
	var planeRanks []model.PlaneRank
	sqls.DB().Where("id in (?)", planeRankIds).Find(&planeRanks)
	return planeRanks
}

func (r *planeRankDao) GetByName(name string) *model.PlaneRank {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
