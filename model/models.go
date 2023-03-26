package model

import (
	"time"
)

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

type BookDetail struct {
	Id          int64      `json:"id"`
	BookId      string     `json:"bookId"`
	BookContent string     `json:"bookContent"`
	Lyric       string     `json:"lyric"`
	BookOrder   int        `json:"book_order"`
	Status      int        `json:"status"`
	UpdateDt    *time.Time `json:"updateDt"`
	CreateDt    time.Time  `json:"createDt"`
}

type BookPiece struct {
	Id          int64      `json:"id"`
	BookId      string     `json:"bookId"`
	PhaseId     string     `json:"phaseId"`
	BookContent string     `json:"bookContent"`
	ContentType int        `json:"contentType"`
	Phase       int        `json:"phase"`
	UpdateDt    *time.Time `json:"updateDt"`
	CreateDt    time.Time  `json:"createDt"`
}

type PieceContent struct {
	Id        int64      `json:"id"`
	BookId    string     `json:"bookId"`
	PhaseId   string     `json:"phaseId"`
	BreakFlag int        `json:"breakFlag"`
	BookOrder int        `json:"bookOrder"`
	UpdateDt  *time.Time `json:"updateDt"`
	CreateDt  time.Time  `json:"createDt"`
}

type BookImg struct {
	Id         int64     `gorm:"type:bigint(20);not null;unique"`
	BookId     string    `gorm:"type:varchar(64);not null"`
	Title      string    `gorm:"type:varchar(200);not null"`
	Url        string    `gorm:"type:varchar(200);not null"`
	AuthorType int       `gorm:"type:int(11);not null"`
	CreateDt   time.Time `json:"createDt"`
	CreateBy   string    `json:"createBy"`
}

type Clock struct {
	Id               string `gorm:"type:varchar(36);not null;unique"`
	EventType        int    `gorm:"type:int(11);not null"`
	EventDescription string `gorm:"type:varchar(200);not null"`
	NotifyDate       string `gorm:"type:varchar(200);not null"`
	LunarFlag        int
	Status           int       `gorm:"type:int(11);not null"`
	CreateDate       time.Time `json:"createDate"`
	CreateBy         string    `json:"createBy"`
	UpdateDate       time.Time `json:"updateDate"`
	UpdateBy         string    `json:"updateBy"`
}

type Calendar struct {
	Id               string `gorm:"type:varchar(36);not null;unique"`
	EventType        int    `gorm:"type:int(11);not null"`
	EventDescription string `gorm:"type:varchar(200);not null"`
	NotifyDate       time.Time
	LunarFlag        int
	Status           int       `gorm:"type:int(11);not null"`
	CreateDate       time.Time `json:"createDate"`
	CreateBy         string    `json:"createBy"`
	UpdateDate       time.Time `json:"updateDate"`
	UpdateBy         string    `json:"updateBy"`
}

type Topic struct {
	Id         int       `gorm:"type:int(11);not null;unique"`
	Name       string    `gorm:"type:int(11);not null"`
	Status     int       `gorm:"type:int(11);not null"`
	CreateDate time.Time `json:"createDate"`
	CreateBy   string    `json:"createBy"`
	UpdateDate time.Time `json:"updateDate"`
	UpdateBy   string    `json:"updateBy"`
}

type Sport struct {
	Id         int       `gorm:"type:int(11);not null;unique"`
	Name       string    `gorm:"type:int(11);not null"`
	Status     int       `gorm:"type:int(11);not null"`
	CreateDate time.Time `json:"createDate"`
	CreateBy   string    `json:"createBy"`
	UpdateDate time.Time `json:"updateDate"`
	UpdateBy   string    `json:"updateBy"`
}

func (BookImg) TableName() string {
	return "book_img"
}

type Words struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Theme struct {
	Id       int64     `gorm:"type:int(11);not null;unique"`
	Name     string    `gorm:"type:varchar(64);not null"`
	CreateDt time.Time `json:"createDt"`
	CreateBy string    `json:"createBy"`
}

type Poetry struct {
	Id           int64 `gorm:"type:int(11);not null;unique"`
	Content      string
	Translate    string
	Notes        string
	Appreciation string
	Pinyin       string
	Name         string
	Dynasty      string
	Poet         string
	linkId       string
}

type Tag struct {
	Id       int64 `gorm:"type:int(11);not null;unique"`
	PoetryId string
	Tag      string
}
