package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

func Int32ToString(i int32) string {
	return fmt.Sprint(i)
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToInt(i string) int {
	data, _ := strconv.Atoi(i)
	return data
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data []byte, contentType string) (content string) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Add("content-type", contentType)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiJjNjY0NjBkNjIyNjc0YzRlODlkNDg1YTliNzBjODQ5YSIsInN1YiI6IjEiLCJpc3MiOiJodWF3ZWltaWFuIn0.JpG5fSNet5jaIHCitzDli7_plbUV2Z-UlVKQVUCsWkY")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}
