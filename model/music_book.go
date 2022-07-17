package model

import (
	"ginEssential/dao"
	"time"
)

type MusicBookDTO struct {
	PageInfo
	BookTitle string
	Author    string
}

type MusicBook struct {
	BookId      string    `json:"bookId"`
	BookContent string    `json:"bookContent"`
	Lyric       string    `json:"lyric"`
	BookType    int       `json:"bookType"`
	BookTitle   string    `json:"bookTitle"`
	MusicKey    string    `json:"musicKey"`
	MusicTime   string    `json:"musicTime"`
	Author      string    `json:"author"`
	Composer    string    `json:"composer"`
	Singer      string    `json:"singer"`
	UpdateDt    time.Time `json:"updateDt"`
	CreateDt    time.Time `json:"createDt"`
	CreateBy    string    `json:"createBy"`
}

func (c *MusicBook) SelectPageList(p *dao.Page[MusicBook], query interface{}, args []interface{}, order string) error {
	return p.SelectPage(query, args, order)
}
