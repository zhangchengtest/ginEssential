package model

type Food struct {
	Id       string `gorm:"type:varchar(50);not null;unique"`
	FoodName string `gorm:"type:varchar(11);not null"`
	Category string `gorm:"type:varchar(200);not null"`
	Material string `gorm:"type:varchar(1024);not null"`
	Url      string `gorm:"type:varchar(11);not null"`
}
