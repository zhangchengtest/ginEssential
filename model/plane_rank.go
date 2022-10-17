package model

import "time"

type PlaneRank struct {
	Id       int64     `gorm:"type:bigint(20);not null;unique"`
	Username string    `gorm:"type:varchar(200);not null"`
	Coin     int       `gorm:"type:int(11);not null"`
	CreateDt time.Time `json:"createDt"`
	CreateBy string    `json:"createBy"`
}

func (PlaneRank) TableName() string {
	return "plane_rank"
}
