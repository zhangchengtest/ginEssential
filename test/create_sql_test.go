package test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestCreateSql(t *testing.T) {

	b, err := ioutil.ReadFile("data.txt") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	arr := strings.Split(str, "\n")

	for i, w := range arr {
		sql := "INSERT INTO `topic` VALUES ('%s', '%s', 1, '', '2023-02-05 04:33:25', '', '2023-02-05 04:33:25');\n"
		data := strings.Replace(w, "\r", "", -1)
		dd := fmt.Sprintf(sql, strconv.Itoa(i+1), data)
		//fmt.Print(sql)  // print the content as a 'string'
		fmt.Println(dd) // print the content as a 'string'
	}

	// 以字符串的形式打印消息
}
