package controller

import (
	"bufio"
	"fmt"
	"ginEssential/dao"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"ginEssential/vo"
	"github.com/gin-gonic/gin"
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
	DB := dao.GetDB()
	// 1. 使用map获取application/json请求的参数
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	// fmt.Printf("requestMap：%v", requestMap)

	// 2. 使用结构体获取application/json请求的参数
	// var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)
	// fmt.Printf("requestUser：%v", requestUser)

	// 3. gin自带的bind获取application/json请求的参数
	var ginBindUser = model.Article{}
	ctx.Bind(&ginBindUser)
	fmt.Printf("ginBindUser：%v", ginBindUser)

	// 获取参数
	chapter := ginBindUser.Chapter
	title := ginBindUser.Title
	content := ginBindUser.Content
	// name := ctx.PostForm("name")
	// telephone := ctx.PostForm("telephone")
	// password := ctx.PostForm("password")

	var s = util.Worker1{}
	// 创建用户
	newUser := model.Article{
		Id:      s.GetId(),
		Chapter: chapter,
		Title:   title,
		Content: content,
	}

	DB.Create(&newUser)

	response.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}

func RandomArticle(ctx *gin.Context) {
	DB := dao.GetDB()

	// 创建用户
	newUser := model.Article{}

	articleVO := vo.ArticleVO{}

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

	response.Success(ctx, gin.H{"article": articleVO}, "查询成功")
}

func AddArticleFromFile(ctx *gin.Context) {

	//firstFile := javabean.FirstFile
	//secondFile :=javabean.SecondFile

	file1, header, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("get file error: %s", err)
		response.Response(ctx, http.StatusBadRequest, 422, nil, "文件上传失败")
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

	response.Success2(ctx, "ok", "")
}

func readFile(f1 *os.File) {
	sc1 := bufio.NewScanner(f1)
	DB := dao.GetDB()
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
