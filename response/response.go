package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, message string) {
	ctx.JSON(httpStatus, gin.H{
		"code":      code,
		"data":      data,
		"message":   message,
		"isSuccess": false,
	})
}

func Success(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":      200,
		"data":      data,
		"message":   message,
		"isSuccess": true,
	})
}

func Success2(ctx *gin.Context, data string, message string) {
	ctx.JSON(200, gin.H{
		"code":    http.StatusOK,
		"data":    data,
		"message": message,
	})
}

func Fail(ctx *gin.Context, data gin.H, message string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":      400,
		"data":      data,
		"message":   message,
		"isSuccess": false,
	})
}
