package service

import (
	"fmt"
	"ginEssential/dao"
	"ginEssential/model/constants"
	"ginEssential/util"
	strftime "github.com/itchyny/timefmt-go"
	"math/rand"
	"strings"
	"time"

	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"

	"ginEssential/model"
)

var PuzzlePieceService = newPuzzlePieceService()

func newPuzzlePieceService() *puzzlePieceService {
	return &puzzlePieceService{}
}

type puzzlePieceService struct {
}

//func (c *puzzlePieceService) SelectPageList(queryVo model.PuzzlePieceDTO) (*model.PageResponse[model.PuzzlePiece], error) {
//	p := &dao.Page[model.PuzzlePiece]{
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
//	err := puzzlePieceModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *puzzlePieceService) Get(id string) *model.PuzzlePiece {
	return dao.PuzzlePieceDao.Get(sqls.DB(), id)
}

func (s *puzzlePieceService) Take(where ...interface{}) *model.PuzzlePiece {
	return dao.PuzzlePieceDao.Take(sqls.DB(), where...)
}

func (s *puzzlePieceService) Find(cnd *sqls.Cnd) []model.PuzzlePiece {
	return dao.PuzzlePieceDao.Find(sqls.DB(), cnd)
}

func (s *puzzlePieceService) FindOne(cnd *sqls.Cnd) *model.PuzzlePiece {
	return dao.PuzzlePieceDao.FindOne(sqls.DB(), cnd)
}

func (s *puzzlePieceService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.PuzzlePiece], error) {
	return dao.PuzzlePieceDao.FindPageByParams(sqls.DB(), params)
}

func (s *puzzlePieceService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.PuzzlePiece], error) {
	return dao.PuzzlePieceDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *puzzlePieceService) Create(t *model.PuzzlePiece) error {
	return dao.PuzzlePieceDao.Create(sqls.DB(), t)
}

func (s *puzzlePieceService) Update(t *model.PuzzlePiece) error {
	if err := dao.PuzzlePieceDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *puzzlePieceService) Updates(id string, columns map[string]interface{}) error {
	return dao.PuzzlePieceDao.Updates(sqls.DB(), id, columns)
}

func (s *puzzlePieceService) UpdateAll(id string, columns *model.PuzzlePiece) error {
	return dao.PuzzlePieceDao.UpdateAll(sqls.DB(), id, columns)
}

//	func (s *puzzlePieceService) UpdateColumn(id int64, name string, value interface{}) error {
//		return dao.PuzzlePieceDao.UpdateColumn(sqls.DB(), id, name, value)
//	}
func (s *puzzlePieceService) Delete(id string) {
	dao.PuzzlePieceDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *puzzlePieceService) Autocomplete(input string) []model.PuzzlePiece {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.PuzzlePieceDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *puzzlePieceService) GetByName(name string) *model.PuzzlePiece {
	return dao.PuzzlePieceDao.GetByName(name)
}

func (s *puzzlePieceService) GetPuzzlePieces(url string) []model.PuzzlePiece {
	list := dao.PuzzlePieceDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url))

	return list
}

func (s *puzzlePieceService) GetPuzzlePiecesById(id string) []model.PuzzlePiece {
	piece := dao.PuzzlePieceDao.Get(sqls.DB(), id)
	list := dao.PuzzlePieceDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", piece.Url))

	return list
}

func (s *puzzlePieceService) GetPuzzlePiecesRandom(url string) []model.PuzzlePiece {
	list := dao.PuzzlePieceDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", url))
	Shuffle(list)

	return list
}

func (s *puzzlePieceService) GetPuzzlePiecesRandomById(id string) []model.PuzzlePiece {
	piece := dao.PuzzlePieceDao.Get(sqls.DB(), id)
	list := dao.PuzzlePieceDao.Find(sqls.DB(), sqls.NewCnd().Where("url = ?", piece.Url))
	Shuffle(list)

	return list
}

func Shuffle(slice []model.PuzzlePiece) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func (s *puzzlePieceService) GetPuzzlePiecesGroup() []model.PuzzlePieceVO {

	t := time.Now()
	date := strftime.Format(t, "%Y-%m-%d")
	fmt.Println(date)
	//t2, _ := strftime.Parse(date+" 00:00:00", "%Y-%m-%d %H:%M:%S")

	var list []model.PuzzlePiece
	sqls.DB().Model(&model.PuzzlePiece{}).Select("url, create_dt").Order("create_dt desc").Group("url, create_dt").Limit(30).Find(&list)
	var puzzlePieces []model.PuzzlePieceVO
	for _, puzzlePiece := range list {

		bookvo := model.PuzzlePieceVO{}
		util.SimpleCopyProperties(&bookvo, &puzzlePiece)
		puzzlePieces = append(puzzlePieces, bookvo)
	}
	return puzzlePieces
}

func (s *puzzlePieceService) GetPuzzlePieceInIds(puzzlePieceIds []int64) []model.PuzzlePiece {
	return dao.PuzzlePieceDao.GetPuzzlePieceInIds(puzzlePieceIds)
}
