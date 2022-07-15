package model

import (
	"time"
)

type User struct {
	UserId        string    `json:"userId"`
	UserName      string    `json:"userName"`
	NickName      string    `json:"nickName"`
	Mobile        string    `json:"mobile"`
	Email         string    `json:"email"`
	Salt          string    `json:"salt"`
	Pwd           string    `json:"pwd"`
	Autologin     string    `json:"autologin"`
	UserStatus    string    `json:"userStatus"`
	UpdateBy      string    `json:"updateBy"`
	UpdateDt      time.Time `json:"updateDt"`
	CreateDt      time.Time `json:"createDt"`
	CreateBy      string    `json:"createBy"`
	IsLock        string    `json:"isLock"`
	Retry         string    `json:"retry"`
	LastWrongPwDt time.Time `json:"lastWrongPwDt"`
	LastLoginDt   time.Time `json:"lastLoginDt"`
}

type UserDto struct {
	UserId       string `json:"userId"`
	UserName     string `json:"userName"`
	Mobile       string `json:"mobile"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	RoleCode     string `json:"roleCode"`
}

// 自定义表名
func (User) TableName() string {
	return "sys_user"
}

type JavaBean struct {
	OriginText string
	TableName  string
}

type MyFile struct {
	FirstFile  string
	SecondFile string
}
