package model

import "time"

type ClockVO struct {
	Days        int
	EventTime   int
	EventType   int
	RealDate    string
	Description string
}

type MusicBookVO struct {
	BookId      string      `json:"bookId"`
	BookContent string      `json:"bookContent"`
	Lyric       string      `json:"lyric"`
	BookType    int         `json:"bookType"`
	BookTitle   string      `json:"bookTitle"`
	MusicKey    string      `json:"musicKey"`
	MusicTime   string      `json:"musicTime"`
	Author      string      `json:"author"`
	Composer    string      `json:"composer"`
	Singer      string      `json:"singer"`
	UpdateDt    time.Time   `json:"updateDt"`
	CreateDt    time.Time   `json:"createDt"`
	CreateBy    string      `json:"createBy"`
	PieceAll    BookPieceVO `json:"pieceAll"`
	MyUrl       string      `json:"myUrl"`
	OtherUrl    string      `json:"otherUrl"`
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

type BookPieceVO struct {
	Id          string          `json:"id"`
	BookId      string          `json:"bookId"`
	BookContent string          `json:"bookContent"`
	Lyric       string          `json:"lyric"`
	List        []PieceDetailVO `json:"list"`
	PhaseId     string          `json:"phaseId"`
	BookOrder   int             `json:"bookOrder"`
	ShowClass   string          `json:"showClass"`
}

type PieceDetailVO struct {
	BookContent      string `json:"bookContent"`
	Connection       string `json:"connection"`
	Connectionxstart []int  `json:"connectionxstart"`
	Connectionxstop  []int  `json:"connectionxstop"`
	UpPoints         string `json:"upPoints"`
	Line1            string `json:"line1"`
	Line2            string `json:"line2"`
	Line1xstart      []int  `json:"line1xstart"`
	Line1xstop       []int  `json:"line1xstop"`
	Line2xstart      []int  `json:"line2xstart"`
	Line2xstop       []int  `json:"line2xstop"`
	DownPoints       string `json:"downPoints"`
	Indent           string `json:"indent"`
	Lyric            string `json:"lyric"`
}

type BookPieceVO2 struct {
	Id1    string `json:"id1"`
	Id2    string `json:"id2"`
	Id3    string `json:"id3"`
	Id4    string `json:"id4"`
	Val1   string `json:"val1"`
	Val2   string `json:"val2"`
	Val3   string `json:"val3"`
	Val4   string `json:"val4"`
	Lyric1 string `json:"lyric1"`
	Lyric2 string `json:"lyric2"`
	Lyric3 string `json:"lyric3"`
	Lyric4 string `json:"lyric4"`
}

type BookPieceVO3 struct {
	Val0 string `json:"val0"`
	Val1 string `json:"val1"`
	Val2 string `json:"val2"`
	Val3 string `json:"val3"`
	Val4 string `json:"val4"`
}

type PieceContentVO struct {
	Id          int64      `json:"id"`
	PieceId     string     `json:"pieceId"`
	Content     string     `json:"content"`
	ContentType int        `json:"contentType"`
	UpdateDt    *time.Time `json:"updateDt"`
	CreateDt    time.Time  `json:"createDt"`
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

type FoodVO struct {
	Id       string `json:"id"`
	FoodName string `json:"foodName"`
	Category string `json:"category"`
	Material string `json:"material"`
	Url      string `json:"url"`
}

type TempOssVO struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	BucketName      string `json:"bucketName"`
	UploadUrl       string `json:"uploadUrl"`
	Token           string `json:"token"`
}

type WechatAuthVO struct {
	Username  string `json:"username"`
	NickName  string `json:"nickName"`
	AvatarUrl string `json:"avatarUrl"`
}
