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
	r.GET("/api/auth/loadUserByEmail", controller.LoadUserByEmail)

	//r.GET("/api/auth/randomImage", controller.RandomImage)

	r.GET("/api/auth/redirectU", controller.RedirectTOUnsplash)
	r.GET("/api/auth/backFromU", controller.BackFromUnsplash)

	r.GET("/api/auth/redirectW", controller.RedirectTOWechat)
	r.GET("/api/auth/backFromW", controller.BackFromWechat)

	r.POST("/api/javatosql", controller.Javatosql)
	r.POST("/api/compareFile", controller.CompareFile)
	//r.POST("/api/testThread", controller.TestThread)
	r.POST("/api/addArticle", controller.AddArticle)
	r.POST("/api/addDinary", controller.AddDinary)
	r.GET("/api/seeDinary", controller.SeeDinary)

	r.GET("/api/copyArticle", controller.CopyArticle)

	r.POST("/api/addArticleFromFile", controller.AddArticleFromFile)
	r.GET("/api/randomArticle", controller.RandomArticle)
	r.GET("/api/randomNovel", controller.RandomNovel)

	r.POST("/api/readorc", controller.Readorc)

	r.POST("/wx/auth/login_by_weixin", controller.LoginByWeixinCode)
	r.POST("/wx/auth/modify", middleware.AuthMiddleware(), controller.ModifyUser)
	r.GET("/wx/auth/detail", middleware.AuthMiddleware(), controller.UserDetail)
	r.POST("/wx/dfs/upload/file", controller.UploadFile)

	r.GET("/wx/testTemplate", middleware.AuthMiddleware(), controller.TestTemplate)

	r.GET("/api/wx/share", controller.WeixinShare)

	r.GET("/api/randomFood", controller.RandomFood)
	r.GET("/wx/searchFood", controller.SearchFood)

	r.POST("/api/musicBook/save", middleware.AuthMiddleware(), controller.SaveMusicBook)
	r.POST("/api/musicBook/search", controller.SearchMusicBook)
	r.GET("/api/musicBook/detail/:id", controller.DetailMusicBook)
	r.POST("/api/musicBook/delete/:id", controller.DeleteMusicBook)

	r.POST("/api/musicBook/uploadBookImg", controller.UploadBookImg)

	r.POST("/api/musicBookDetail/searchByPage", controller.SearchMusicBookDetail)
	r.POST("/api/musicBookDetail/searchOne", controller.SearchOneMusicBookDetail)
	r.POST("/api/musicBookDetail/updateContent", controller.UpdateMusicBookConent)
	r.POST("/api/musicBookDetail/updateLyric", controller.UpdateMusicBookLyric)

	r.POST("/api/musicBookPiece/deletePiece", controller.DeletePiece)
	r.POST("/api/musicBookPiece/updateContent", controller.UpdateBookPiece)
	r.POST("/api/musicBookPiece/searchPieces", controller.SearchPieces)
	r.POST("/api/musicBookPiece/searchPiecesByPhaseId", controller.SearchPiecesByPhaseId)
	r.POST("/api/musicBookPiece/copy", controller.CopyPiece)
	r.POST("/api/musicBookPiece/stickUp", controller.StickUp)

	r.POST("/api/musicBookPiece/test", controller.TestBookPiece)

	r.POST("/api/game/uploadSplitImages", controller.UploadSplitImages)

	r.GET("/api/game/queryPuzzle", controller.QueryPuzzle)
	r.GET("/api/game/queryPuzzleByUrl", controller.QueryPuzzleByUrl)

	r.GET("/api/game/queryPuzzleRank", controller.QueryPuzzleRank)
	r.POST("/api/game/savePuzzleRank", controller.SavePuzzleRank)

	r.GET("/api/game/queryPlaneRank", controller.QueryPlaneRank)
	r.GET("/api/game/visit", controller.Visit)
	r.GET("/api/game/nickname", controller.Nickname)
	r.POST("/api/game/savePlaneRank", controller.SavePlaneRank)
	r.POST("/api/game/modifyUsername", controller.ModifyUsername)

	r.GET("/api/game/push", controller.PushTest)

	r.GET("/api/oss/authTemp", controller.AuthTemp)
	r.POST("/api/oss/temp2formal", controller.Copy)

	return r
}
