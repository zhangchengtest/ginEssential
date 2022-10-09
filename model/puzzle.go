package model

import "time"

type PuzzlePiece struct {
	Id       int64     `gorm:"type:bigint(20);not null;unique"`
	Content  string    `gorm:"type:text;not null"`
	Title    string    `gorm:"type:varchar(200);not null"`
	Url      string    `gorm:"type:varchar(200);not null"`
	Sort     int       `gorm:"type:int(11);not null"`
	CreateDt time.Time `json:"createDt"`
	CreateBy string    `json:"createBy"`
}

func (PuzzlePiece) TableName() string {
	return "puzzle_piece"
}
