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
	return r
}
