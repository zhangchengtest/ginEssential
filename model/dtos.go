package model

type MusicBookDTO struct {
	PageInfo
	BookTitle string
	Author    string
}

type BookDetailDTO struct {
	PageInfo
	Id          string
	BookId      string
	BookContent string
	Lyric       string
}
