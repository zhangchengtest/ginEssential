package service

import (
	"fmt"
	"ginEssential/dao"
	"ginEssential/model"
)

type MusicBookService struct{}

var musicBookModel model.MusicBook

func (c *MusicBookService) SelectPageList(queryVo model.MusicBookDTO) (*model.PageResponse[model.MusicBook], error) {
	p := &dao.Page[model.MusicBook]{
		CurrentPage: queryVo.PageNum,
		PageSize:    queryVo.PageSize,
	}
	var query []string
	var args []interface{}
	if queryVo.BookTitle != "" {
		query = append(query, "book_title LIKE ?")
		args = append(args, "%"+queryVo.BookTitle+"%")
	}
	if queryVo.Author != "" {
		query = append(query, "author LIKE ?")
		args = append(args, "%"+queryVo.Author+"%")
	}
	var querystr string

	for index, value := range query {
		if index == len(query)-1 {
			querystr += value
		} else {
			querystr += value + " and "
		}
		fmt.Printf("index: %d value: %s\n", index, value)
	}

	err := musicBookModel.SelectPageList(p, querystr, args, "create_dt desc")
	if err != nil {
		return nil, err
	}
	pageResponse := model.NewPageResponse(p)
	return pageResponse, err
}
