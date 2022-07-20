package model

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

type PageResponse[T any] struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	Total       int64 `json:"total"`
	Pages       int   `json:"pages"` // 总页数
	Data        []T   `json:"data"`
}

func NewPageResponse[T any](page *Page[T]) *PageResponse[T] {
	return &PageResponse[T]{
		CurrentPage: page.CurrentPage,
		PageSize:    page.PageSize,
		Total:       page.Total,
		Pages:       page.Pages,
		Data:        page.Data,
	}
}

type Page[T any] struct {
	CurrentPage int
	PageSize    int
	Total       int64
	Pages       int
	Data        []T
}
