package model

	type Article struct {
	Id int64 `gorm:"type:bigint(20);not null;unique"`
	Chapter      int32 `gorm:"type:int(11);not null"`
	Title      string `gorm:"type:varchar(200);not null"`
	Content      string `gorm:"type:varchar(1024);not null"`
	ReadCount      int32 `gorm:"type:int(11);not null"`
}

