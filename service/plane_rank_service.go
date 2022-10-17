package service

import (
	"ginEssential/dao"
	"ginEssential/model/constants"
	"ginEssential/util"
	strftime "github.com/itchyny/timefmt-go"
	"github.com/zhangchengtest/simple/sqls"
	"github.com/zhangchengtest/simple/web/params"
	"strings"

	"ginEssential/model"
)

var PlaneRankService = newPlaneRankService()

func newPlaneRankService() *planeRankService {
	return &planeRankService{}
}

type planeRankService struct {
}

//func (c *planeRankService) SelectPageList(queryVo model.PlaneRankDTO) (*model.PageResponse[model.PlaneRank], error) {
//	p := &dao.Page[model.PlaneRank]{
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
//	err := planeRankModel.SelectPageList(p, querystr, args, "create_dt desc")
//	if err != nil {
//		return nil, err
//	}
//	pageResponse := model.NewPageResponse(p)
//	return pageResponse, err
//}

func (s *planeRankService) Get(id string) *model.PlaneRank {
	return dao.PlaneRankDao.Get(sqls.DB(), id)
}

func (s *planeRankService) Take(where ...interface{}) *model.PlaneRank {
	return dao.PlaneRankDao.Take(sqls.DB(), where...)
}

func (s *planeRankService) Find(cnd *sqls.Cnd) []model.PlaneRank {
	return dao.PlaneRankDao.Find(sqls.DB(), cnd)
}

func (s *planeRankService) FindOne(cnd *sqls.Cnd) *model.PlaneRank {
	return dao.PlaneRankDao.FindOne(sqls.DB(), cnd)
}

func (s *planeRankService) FindPageByParams(params *params.QueryParams) (*model.PageResponse[model.PlaneRank], error) {
	return dao.PlaneRankDao.FindPageByParams(sqls.DB(), params)
}

func (s *planeRankService) FindPageByCnd(cnd *sqls.Cnd) (*model.PageResponse[model.PlaneRank], error) {
	return dao.PlaneRankDao.FindPageByCnd(sqls.DB(), cnd)
}

func (s *planeRankService) Create(t *model.PlaneRank) error {
	return dao.PlaneRankDao.Create(sqls.DB(), t)
}

func (s *planeRankService) Update(t *model.PlaneRank) error {
	if err := dao.PlaneRankDao.Update(sqls.DB(), t); err != nil {
		return err
	}
	return nil
}

func (s *planeRankService) Updates(id string, columns map[string]interface{}) error {
	return dao.PlaneRankDao.Updates(sqls.DB(), id, columns)
}

func (s *planeRankService) UpdateAll(id string, columns *model.PlaneRank) error {
	return dao.PlaneRankDao.UpdateAll(sqls.DB(), id, columns)
}

//
// func (s *planeRankService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return dao.PlaneRankDao.UpdateColumn(sqls.DB(), id, name, value)
// }
//
func (s *planeRankService) Delete(id string) {
	dao.PlaneRankDao.Delete(sqls.DB(), id)
}

// 自动完成
func (s *planeRankService) Autocomplete(input string) []model.PlaneRank {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return dao.PlaneRankDao.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *planeRankService) GetByName(name string) *model.PlaneRank {
	return dao.PlaneRankDao.GetByName(name)
}

func (s *planeRankService) GetRanks() []model.PlaneRankVO {
	list := dao.PlaneRankDao.Find(sqls.DB(), sqls.NewCnd().Desc("coin"))
	var planeRanks []model.PlaneRankVO
	sort := 1
	for _, planeRank := range list {

		bookvo := model.PlaneRankVO{}
		util.SimpleCopyProperties(&bookvo, &planeRank)
		date := strftime.Format(planeRank.CreateDt, "%m-%d(%H:%M)")
		bookvo.CreateDt = date
		bookvo.Sort = sort
		sort++
		planeRanks = append(planeRanks, bookvo)
	}
	return planeRanks
}

func (s *planeRankService) GetPlaneRankInIds(planeRankIds []int64) []model.PlaneRank {
	return dao.PlaneRankDao.GetPlaneRankInIds(planeRankIds)
}
