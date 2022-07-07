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
	"strconv"
	"strings"
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
	arr := strings.Split(text, "\n")
	fmt.Println(text)
	for _, s := range arr {

		if strings.TrimSpace(s) == ""{
			continue
		}
		arr1 := strings.Split(s, " ")
		//fmt.Println(len(arr1))
		if(len(arr1) > 3){
			fmt.Println(arr1[len(arr1)-2])
			chapter, err := strconv.ParseInt(arr1[len(arr1)-2], 10, 32)
			if err == nil {
				fmt.Println(chapter)
			}

		}
	}



	response.Success(ctx, gin.H{"status": "ok"}, "新增成功")
}



