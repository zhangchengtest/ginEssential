package middleware

import (
	"ginEssential/dao"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := util.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 验证通过，获取claims中的userId
		userId := claims.UserId
		DB := dao.GetDB()
		var user model.User
		DB.Where("user_id = ?", userId).First(&user)

		// 用户
		if user.UserId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		// 用户存在，将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
