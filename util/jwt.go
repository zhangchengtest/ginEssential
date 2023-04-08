package util

import (
	"ginEssential/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("dcw#zc#hht#qyj#wh#cb#yl#czz#zzy#dmj#rc#yzh#yc#zyy#nhb#yc#sz#lb@cunw.com.cn&$(YHK!UGHJGS&^%!^&!JBMN!(*^#(*Y!KJ#GH!JGS&^!$%R@TYFVGHVDYTR!&^#RHNVDKAL:!())(123IUH#IU!(*&!&^HK!^*^*()(&*F$$#^^FLKLHFIUEHK)_APQY^#!")

type Claims struct {
	Data string
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 有效期7天

	params, _ := ToJSONString(map[string]interface{}{
		"userId": user.UserId,
	})
	claims := &Claims{
		Data: params,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 发放时间
			Issuer:    "oceanlearn.tech",     // 发放者
			Subject:   "user token",          // 主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
