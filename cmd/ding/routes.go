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
	r.GET("/api/auth/mockinfo", controller.MockInfo)
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
	r.POST("/api/musicBook/delete/:id", controller.DeleteMusicBook)

	r.POST("/api/musicBookDetail/searchByPage", controller.SearchMusicBookDetail)
	r.POST("/api/musicBookDetail/searchOne", controller.SearchOneMusicBookDetail)
	r.POST("/api/musicBookDetail/updateContent", controller.UpdateMusicBookConent)
	r.POST("/api/musicBookDetail/updateLyric", controller.UpdateMusicBookLyric)

	r.POST("/api/game/uploadSplitImages", controller.UploadSplitImages)

	r.GET("/api/game/queryPuzzle", controller.QueryPuzzle)
	r.GET("/api/game/queryPuzzleByUrl", controller.QueryPuzzleByUrl)
	r.GET("/api/game/queryPuzzleRank", controller.QueryPuzzleRank)
	r.POST("/api/game/savePuzzleRank", controller.SavePuzzleRank)

	r.GET("/api/oss/authTemp", controller.AuthTemp)
	r.POST("/api/oss/temp2formal", controller.Copy)

	return r
}
