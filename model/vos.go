package model

import "time"

type MusicBookVO struct {
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

type BookDetailVO struct {
	Id          string    `json:"id"`
	BookId      string    `json:"bookId"`
	BookContent string    `json:"bookContent"`
	Lyric       string    `json:"lyric"`
	ShowClass   string    `json:"showClass"`
	Order       int       `json:"order"`
	UpdateDt    time.Time `json:"updateDt"`
	CreateDt    time.Time `json:"createDt"`
}

type PuzzlePieceVO struct {
	Id       int64     `json:"id"`
	Content  string    `json:"content"`
	Title    string    `json:"title"`
	Url      string    `json:"url"`
	Sort     int       `json:"sort"`
	CreateDt time.Time `json:"createDt"`
	CreateBy string    `json:"createBy"`
}

type PuzzleRankVO struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	Url       string    `json:"url"`
	Sort      int       `json:"sort"`
	SpendTime int       `json:"spendTime"`
	Step      int       `json:"step"`
	CreateDt  time.Time `json:"createDt"`
	CreateBy  string    `json:"createBy"`
}
type PlaneRankVO struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Sort     int    `json:"sort"`
	Coin     int    `json:"coin"`
	CreateDt string `json:"createDt"`
	CreateBy string `json:"createBy"`
}

type PuzzlePieceVO2 struct {
	Url     string   `json:"url"`
	Piecces []string `json:"piecces"`
	Orders  []int    `json:"orders"`
}

type ArticleVO struct {
	Id        int64  `json:"id"`
	Chapter   int32  `json:"chapter"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ReadCount int32  `json:"readCount"`
	Question  string `json:"question"`
}

type TempOssVO struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
	UploadUrl       string `json:"uploadUrl"`
	Token           string `json:"token"`
}
