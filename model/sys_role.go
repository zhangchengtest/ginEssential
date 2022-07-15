package model

import (
	"time"
)

type SysRole struct {
	Id       int       `json:"id"`
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Status   int       `json:"status"`
	UpdateDt time.Time `json:"updateDt"`
	CreateDt time.Time `json:"createDt"`
}

// 自定义表名
func (SysRole) TableName() string {
	return "sys_role"
}
