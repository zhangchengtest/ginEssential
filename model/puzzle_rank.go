package model

import "time"

type PuzzleRank struct {
	Id        int64     `gorm:"type:bigint(20);not null;unique"`
	Username  string    `gorm:"type:varchar(200);not null"`
	Title     string    `gorm:"type:varchar(200);not null"`
	Url       string    `gorm:"type:varchar(200);not null"`
	SpendTime int       `gorm:"type:int(11);not null"`
	Step      int       `gorm:"type:int(11);not null"`
	CreateDt  time.Time `json:"createDt"`
	CreateBy  string    `json:"createBy"`
}

func (PuzzleRank) TableName() string {
	return "puzzle_rank"
}
