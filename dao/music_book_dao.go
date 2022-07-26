package dao

import (
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"

	"ginEssential/model"
)

var MusicBookDao = newMusicBookDao()

func newMusicBookDao() *musicBookDao {
	return &musicBookDao{}
}

type musicBookDao struct {
}

func (r *musicBookDao) Get(db *gorm.DB, id string) *model.MusicBook {
	ret := &model.MusicBook{}
	if err := db.First(ret, "book_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *musicBookDao) Take(db *gorm.DB, where ...interface{}) *model.MusicBook {
	ret := &model.MusicBook{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *musicBookDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.MusicBook) {
	cnd.Find(db, &list)
	return
}

func (r *musicBookDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.MusicBook {
	ret := &model.MusicBook{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *musicBookDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.MusicBook], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *musicBookDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.MusicBook], error) {

	page := &model.Page[model.MusicBook]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.MusicBook{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.MusicBook{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *musicBookDao) Create(db *gorm.DB, t *model.MusicBook) (err error) {
	err = db.Create(t).Error
	return
}

func (r *musicBookDao) Update(db *gorm.DB, t *model.MusicBook) (err error) {
	err = db.Save(t).Error
	return
}

func (r *musicBookDao) Updates(db *gorm.DB, id string, columns map[string]interface{}) (err error) {
	err = db.Model(&model.MusicBook{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *musicBookDao) UpdateAll(db *gorm.DB, id string, columns *model.MusicBook) (err error) {
	err = db.Model(&model.MusicBook{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *musicBookDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.MusicBook{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *musicBookDao) Delete(db *gorm.DB, id string) {
	db.Delete(&model.MusicBook{}, "book_id = ?", id)
}

func (r *musicBookDao) GetMusicBookInIds(musicBookIds []int64) []model.MusicBook {
	if len(musicBookIds) == 0 {
		return nil
	}
	var musicBooks []model.MusicBook
	sqls.DB().Where("id in (?)", musicBookIds).Find(&musicBooks)
	return musicBooks
}

func (r *musicBookDao) GetByName(name string) *model.MusicBook {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}
