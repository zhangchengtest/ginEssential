package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code":       code,
		"data":       data,
		"msg":        msg,
		"isSucccess": false,
	})
}

func Success(ctx *gin.Context, data interface{}, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":       200,
		"data":       data,
		"msg":        msg,
		"isSucccess": true,
	})
}

func Success2(ctx *gin.Context, data string, msg string) {
	ctx.JSON(200, gin.H{
		"code": http.StatusOK,
		"data": data,
		"msg":  msg,
	})
}

func Fail(ctx *gin.Context, data gin.H, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":       400,
		"data":       data,
		"msg":        msg,
		"isSucccess": false,
	})
}
