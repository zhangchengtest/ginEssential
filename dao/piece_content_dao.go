package dao

import (
	"ginEssential/model"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"gorm.io/gorm"
	"log"
)

var PieceContentDao = newPieceContentDao()

func newPieceContentDao() *pieceContentDao {
	return &pieceContentDao{}
}

type pieceContentDao struct {
}

func (r *pieceContentDao) Get(db *gorm.DB, id string) *model.PieceContent {
	ret := &model.PieceContent{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pieceContentDao) GetByPhaseId(db *gorm.DB, id string) *model.PieceContent {
	ret := &model.PieceContent{}
	if err := db.First(ret, "phase_id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pieceContentDao) Take(db *gorm.DB, where ...interface{}) *model.PieceContent {
	ret := &model.PieceContent{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *pieceContentDao) Find(db *gorm.DB, cnd *sqls.Cnd) (list []model.PieceContent) {
	cnd.Find(db, &list)
	return
}

func (r *pieceContentDao) FindOne(db *gorm.DB, cnd *sqls.Cnd) *model.PieceContent {
	ret := &model.PieceContent{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *pieceContentDao) FindPageByParams(db *gorm.DB, params *params.QueryParams) (*model.PageResponse[model.PieceContent], error) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *pieceContentDao) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (*model.PageResponse[model.PieceContent], error) {

	page := &model.Page[model.PieceContent]{
		CurrentPage: cnd.Paging.Page,
		PageSize:    cnd.Paging.Limit,
	}
	cnd.Paging.Total = cnd.Count(db, &model.PieceContent{})
	page.Total = cnd.Paging.Total
	page.Pages = cnd.Paging.TotalPage()

	if page.Total == 0 {
		page.Data = []model.PieceContent{}
		pageResponse := model.NewPageResponse(page)
		return pageResponse, nil
	}
	cnd.Find(db, &page.Data)

	pageResponse := model.NewPageResponse(page)
	return pageResponse, nil
}

func (r *pieceContentDao) Create(db *gorm.DB, t *model.PieceContent) (err error) {
	err = db.Create(t).Error
	return
}

func (r *pieceContentDao) Update(db *gorm.DB, t *model.PieceContent) (err error) {
	err = db.Save(t).Error
	return
}

func (r *pieceContentDao) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.PieceContent{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *pieceContentDao) UpdateAll(db *gorm.DB, id string, columns *model.PieceContent) (err error) {
	err = db.Model(&model.PieceContent{}).Where("book_id = ?", id).Updates(columns).Error
	return
}

func (r *pieceContentDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.PieceContent{}).Where("book_id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *pieceContentDao) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.PieceContent{}, "id = ?", id)
}

func (r *pieceContentDao) GetPieceContentInIds(pieceContentIds []int64) []model.PieceContent {
	if len(pieceContentIds) == 0 {
		return nil
	}
	var pieceContents []model.PieceContent
	sqls.DB().Where("id in (?)", pieceContentIds).Find(&pieceContents)
	return pieceContents
}

func (r *pieceContentDao) GetByName(name string) *model.PieceContent {
	if len(name) == 0 {
		return nil
	}
	return r.Take(sqls.DB(), "name = ?", name)
}

func (r *pieceContentDao) SelectByContentType(db *gorm.DB, book_id string, content_type int) []model.BookPiece {

	var bookPieces []model.BookPiece
	db.Table("piece_content").Joins("join book_piece ON piece_content.phase_id = book_piece.phase_id").Select("piece_content.book_order, book_piece.*").Where("piece_content.book_id =? and book_piece.content_type = ?", book_id, content_type).Order("piece_content.book_order asc").Find(&bookPieces)

	return bookPieces
}

func (r *pieceContentDao) SelectByBreakFlag(db *gorm.DB, book_id string, start_order int, end_order int) []model.BookPiece {

	var bookPieces []model.BookPiece
	db.Table("piece_content").Joins("join book_piece ON piece_content.phase_id = book_piece.phase_id").Select("piece_content.book_order, book_piece.*").Where("piece_content.book_id =? and piece_content.book_order > ? and piece_content.book_order <= ?", book_id, start_order, end_order).Order("piece_content.book_order asc").Find(&bookPieces)

	return bookPieces
}

func (r *pieceContentDao) SelectMax(db *gorm.DB, book_id string) int {

	var count int
	row := db.Table("piece_content").Select("IFNULL(max(book_order) + 1,1) ").Where("book_id =?", book_id).Row()
	err := row.Scan(&count)
	if err != nil {
		log.Println(err)
	}
	log.Println("count---------")
	log.Println(count)
	return count
}
