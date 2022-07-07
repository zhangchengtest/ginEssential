package controller

import (
	"fmt"
	"ginEssential/response"
	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
	"io"
	"log"
	"net/http"
	"os"
)

func Readorc(ctx *gin.Context) {

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

	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(filename)
	text, _ := client.Text()
	fmt.Println(text)


	response.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}



