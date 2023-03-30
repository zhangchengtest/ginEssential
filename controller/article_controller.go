package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"ginEssential/model"
	"ginEssential/util"
	"github.com/gin-gonic/gin"
	"github.com/zhangchengtest/simple/sqls"
	"gorm.io/gorm"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func AddArticle(ctx *gin.Context) {
	DB := sqls.DB()
	// 1. 使用map获取application/json请求的参数
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	// fmt.Printf("requestMap：%v", requestMap)

	// 2. 使用结构体获取application/json请求的参数
	// var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	// fmt.Printf("requestUser：%v", requestUser)

	// 3. gin自带的bind获取application/json请求的参数
	var articleDTO = model.Article{}
	ctx.Bind(&articleDTO)
	fmt.Printf("articleDTO：%v", articleDTO)

	// 获取参数
	chapter := articleDTO.Chapter
	title := articleDTO.Title
	content := articleDTO.Content
	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	var s = util.Worker1{}
	// 创建用户
	article := model.Article{
		Id:       s.GetId(),
		Chapter:  chapter,
		Category: articleDTO.Category,
		Title:    title,
		Content:  content,
		CreateDt: time.Now(),
	}

	DB.Create(&article)

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func AddDinary(ctx *gin.Context) {
	DB := sqls.DB()

	var articleDTO = model.Article{}
	ctx.Bind(&articleDTO)
	fmt.Printf("articleDTO：%v", articleDTO)

	// 获取参数
	chapter := articleDTO.Chapter
	title := articleDTO.Title
	content := articleDTO.Content
	category := articleDTO.Category

	old := model.Article{}

	DB.Where("title = ? and category = ?", title, category).First(&old)
	if old.Id != 0 {

		content := old.Content + "\n" + articleDTO.Content
		// User 的 ID 是 `111`
		DB.Model(&old).Update("content", content)

	} else {

		var s = util.Worker1{}
		// 创建用户
		article := model.Article{
			Id:       s.GetId(),
			Chapter:  chapter,
			Category: articleDTO.Category,
			Title:    title,
			Content:  content,
			CreateDt: time.Now(),
		}

		DB.Create(&article)
	}

	model.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func SeeDinary(ctx *gin.Context) {
	DB := sqls.DB()

	title := ctx.Query("title")
	category := ctx.Query("category")

	old := model.Article{}

	var arr []model.Article

	if title == "分类" {
		DB.Select("category").Group("category").Find(&arr)
		content := ""
		for i, data := range arr {
			i++
			content = content + util.IntToString(i) + " " + data.Category + "\n"
			old.Content = content
		}
		model.Success(ctx, old, "查询成功")
		return

	}
	DB.Where("title = ? and category = ?", title, category).First(&old)

	if category == "行程" {
		arr := strings.Split(old.Content, "\n")
		content := ""
		for i, data := range arr {
			i++
			content = content + util.IntToString(i) + " " + data + "\n"
		}
		old.Content = content
	} else if category == "日记" {
		data := strings.ReplaceAll(old.Content, "\n", "")
		data = strings.ReplaceAll(data, " ", "")
		length := len(data)
		old.Content = old.Content + " " + util.IntToString(length)
	}

	model.Success(ctx, old, "查询成功")
}

func RandomArticle(ctx *gin.Context) {
	DB := sqls.DB()

	// 创建用户
	newUser := model.Article{}

	articleVO := model.ArticleVO{}

	rand.Seed(time.Now().UnixNano())
	chapter := rand.Intn(80) + 1
	DB.Where("chapter = ?", chapter).First(&newUser)

	// User 的 ID 是 `111`
	DB.Model(&newUser).Update("read_count", gorm.Expr("read_count + ?", 1))

	content := newUser.Content
	arr := strings.Split(content, "，")
	random := rand.Intn(len(arr))
	ret := strings.Replace(content, arr[random], "_______", -1)
	newUser.Content = ret

	util.SimpleCopyProperties(&articleVO, &newUser)
	articleVO.Question = arr[random]

	model.Success(ctx, gin.H{"article": articleVO}, "查询成功")
}

type MyArticle struct {
	Id             interface{} `json:"id"`
	ArticleTitle   string      `json:"articleTitle"`
	ArticleContent string      `json:"articleContent"`
	ArticleCover   string      `json:"articleCover"`
	CategoryName   string      `json:"categoryName"`
	TagNames       []string    `json:"tagNames"`
	IsTop          int         `json:"isTop"`
	Type           int         `json:"type"`
	Status         int         `json:"status"`
	IsFeatured     int         `json:"isFeatured"`
}

func CopyArticle(ctx *gin.Context) {
	DB := sqls.DB()

	// 创建用户
	var articles []model.Article

	DB.Where("title = ?", "道德经").Find(&articles)

	//{
	//"id": null,
	//"articleTitle": "2023-03-22",
	//"articleContent": "ccccccc",
	//"articleCover": "https://cheng-resource.oss-cn-hangzhou.aliyuncs.com/articles/ec6b56cb53f08d27e487e1442b36f581.png",
	//"categoryName": "测试",
	//"tagNames": [
	//"前端"
	//],
	//"isTop": 0,
	//"type": 1,
	//"status": 1,
	//"isFeatured": 0
	//}
	go test(articles)

	model.Success2(ctx, "ok", "")
}

func test(articles []model.Article) {
	for _, article := range articles {

		my := MyArticle{
			ArticleTitle:   article.Title + util.Int32ToString(article.Chapter),
			ArticleContent: article.Content,
			ArticleCover:   "https://cheng-resource.oss-cn-hangzhou.aliyuncs.com/articles/87a19c5bfbca6d057f244014f85f9881.jpg",
			CategoryName:   "道德经",
			TagNames:       []string{"道德经"},
			IsTop:          0,
			Type:           1,
			Status:         1,
			IsFeatured:     0,
		}
		data, _ := json.Marshal(my)
		result := util.Post("https://apemgr.punengshuo.com/api/admin/articles", data, "application/json")
		fmt.Printf(result)
	}
}

func AddArticleFromFile(ctx *gin.Context) {

	//firstFile := javabean.FirstFile
	//secondFile :=javabean.SecondFile

	file1, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("get file error: %s", err)
		model.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
		return
	}

	filename := header.Filename

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	sourceFile1, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer sourceFile1.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(sourceFile1, file1)
	if err != nil {
		log.Fatal(err)
	}

	sourceFile1.Seek(0, 0)
	readFile(sourceFile1)

	model.Success2(ctx, "ok", "")
}

func readFile(f1 *os.File) {
	sc1 := bufio.NewScanner(f1)
	DB := sqls.DB()
	for {
		sc1Bool := sc1.Scan()
		if !sc1Bool {
			break
		}
		if strings.TrimSpace(sc1.Text()) == "" {
			continue
		}

		arr := strings.Split(sc1.Text(), ".")
		chapter, err := strconv.ParseInt(arr[0], 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		var s = util.Worker1{}
		// 创建用户
		newUser := model.Article{
			Id:      s.GetId(),
			Chapter: int32(chapter),
			Title:   "道德经",
			Content: arr[1],
		}

		DB.Create(&newUser)
	}
}
