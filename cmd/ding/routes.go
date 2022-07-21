package main

import (
	"ginEssential/controller"
	"ginEssential/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/javatosql", controller.Javatosql)
	r.POST("/api/compareFile", controller.CompareFile)
	r.POST("/api/testThread", controller.TestThread)
	r.POST("/api/addArticle", controller.AddArticle)
	r.POST("/api/addArticleFromFile", controller.AddArticleFromFile)
	r.GET("/api/randomArticle", controller.RandomArticle)
	r.POST("/api/readorc", controller.Readorc)

	r.POST("/api/musicBook/save", middleware.AuthMiddleware(), controller.SaveMusicBook)
	r.POST("/api/musicBook/search", controller.SearchMusicBook)
	r.GET("/api/musicBook/detail/:id", controller.DetailMusicBook)

	r.POST("/api/musicBookDetail/searchByPage", controller.SearchMusicBookDetail)
	r.POST("/api/musicBookDetail/searchOne", controller.SearchOneMusicBookDetail)
	r.POST("/api/musicBookDetail/update", controller.UpdateMusicBookDetail)

	return r
}
