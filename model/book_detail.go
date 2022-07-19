package model

import (
	"time"
)

type BookDetail struct {
	Id          int64     `json:"id"`
	BookId      string    `json:"bookId"`
	BookContent string    `json:"bookContent"`
	Lyric       string    `json:"lyric"`
	order       int64     `json:"order"`
	UpdateDt    time.Time `json:"updateDt"`
	CreateDt    time.Time `json:"createDt"`
}

type BookDetailResponse struct {
	Id          int64     `json:"id"`
	BookId      string    `json:"bookId"`
	BookContent string    `json:"bookContent"`
	Lyric       string    `json:"lyric"`
	order       int64     `json:"order"`
	UpdateDt    time.Time `json:"updateDt"`
	CreateDt    time.Time `json:"createDt"`
}
