package model

type MusicBookDTO struct {
	PageInfo
	BookTitle string
	Author    string
}

type BookDetailDTO struct {
	PageInfo
	BookId string
}
