package model

import (
	"time"
)

type User struct {
	UserId        string     `json:"userId"`
	UserName      string     `json:"userName"`
	NickName      string     `json:"nickName"`
	Mobile        string     `json:"mobile"`
	Email         *string    `json:"email"`
	Salt          string     `json:"salt"`
	Pwd           string     `json:"pwd"`
	Autologin     int        `json:"autologin"`
	UserStatus    int        `json:"userStatus"`
	IsLock        int        `json:"isLock"`
	Retry         int        `json:"retry"`
	LastWrongPwDt *time.Time `json:"lastWrongPwDt"`
	LastLoginDt   time.Time  `json:"lastLoginDt"`
	Openid        string     `json:"openid"`
	Unionid       string     `json:"unionid"`
	AvatarUrl     string     `json:"avatarUrl"`
	UpdateBy      string     `json:"updateBy"`
	UpdateDt      *time.Time `json:"updateDt"`
	CreateDt      time.Time  `json:"createDt"`
	CreateBy      string     `json:"createBy"`
}

type UserVO struct {
	UserId       string   `json:"userId"`
	UserName     string   `json:"userName"`
	Mobile       string   `json:"mobile"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
	RoleCode     string   `json:"roleCode"`
	Pwd          string   `json:"pwd"`
	IsLock       int      `json:"isLock"`
	Avatar       string   `json:"avatar"`
	RedirectUrl  string   `json:"redirectUrl"`
	Punches      []string `json punches`
}

type WechatToken struct {
	Scope        string `json:"scope"`
	Openid       string `json:"openid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginDTO struct {
	Account     string `json:"account"`
	SmsCaptcha  string `json:"smsCaptcha"`
	Pwd         string `json:"pwd"`
	RedirectUrl string `json:"redirectUrl"`
	Uuid        string `json:"uuid"`
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
