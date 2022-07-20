package controller

import (
	"ginEssential/model"
	"github.com/gin-gonic/gin"
)

func Readorc(ctx *gin.Context) {
	//
	//// 获取所有图片
	//form, err := ctx.MultipartForm()
	//if err != nil {
	//	return
	//}
	//if len(form.File) <= 0 {
	//	return
	//}
	//
	//client := gosseract.NewClient()
	//defer client.Close()
	//var ret string;
	//for _, files := range form.File {
	//	for _, file := range files {
	//
	//		if err := ctx.SaveUploadedFile(file, file.Filename); err != nil {
	//			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
	//			return
	//		}
	//
	//		client.SetImage(file.Filename)
	//		text, _ := client.Text()
	//		arr := strings.Split(text, "\n")
	//		fmt.Println(text)
	//
	//		for _, s := range arr {
	//
	//			if strings.TrimSpace(s) == ""{
	//				continue
	//			}
	//			arr1 := strings.Split(s, " ")
	//			//fmt.Println(len(arr1))
	//			if(len(arr1) > 3){
	//				chapter, err := strconv.ParseInt(arr1[len(arr1)-2], 10, 32)
	//				if err == nil {
	//					fmt.Println(chapter)
	//					ret += arr1[len(arr1)-2] +"\r\n"
	//				}
	//			}
	//		}
	//		ret += "\r\n"
	//	}
	//}

	model.Success2(ctx, "ret", "新增成功")
}
