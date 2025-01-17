package util

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"html/template"
	"image"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ConsumeIterator drains the iterator of images and returns them in a slice
// Note that consuming an entire iterator may cause heavy memory usage
// and usually is a bad idea
func ConsumeIterator(it ImageIterator) []image.Image {
	ms := []image.Image{}
	for it.Next() {
		ms = append(ms, it.Get())
	}
	return ms
}
func Add(l int, msg string) string {
	for len(msg) < l {
		msg = "0" + msg
	}
	return msg
}

func GetAllFiles(dirPath string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			files = append(files, filepath.Join(dirPath, file.Name()))
		}
	}

	return files, nil
}

func GetAllFiles2(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func WriteHTMLFile(content string, path string) error {
	// 定义HTML5模板，使用 %s 占位符存放内容
	template := `<!DOCTYPE html>
                <html lang="en">
                <head>
                    <meta charset="UTF-8">
                    <title>HTML5模板</title>
                </head>
                <body>
                    %s
                </body>
                </html>`

	// 将内容写入模板
	html := fmt.Sprintf(template, content)

	// 创建文件
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	// 将HTML写入到文件
	_, err = f.WriteString(html)
	if err != nil {
		return err
	}

	// 关闭文件
	defer f.Close()

	return nil
}

func TxtToHTML(txtContent string) (string, error) {
	// 定义模板
	tmpl := `<html>
	           <head>
	             <meta charset="UTF-8">
	             <title>{{ .Title }}</title>
	           </head>
	           <body>
	             {{ range .Lines }}<p>{{ . }}</p>{{ end }}
	           </body>
	         </html>`

	// 定义结构体
	type Page struct {
		Title string
		Lines []string
	}

	// 将纯文本分行读入
	lines := bytes.Split([]byte(txtContent), []byte("\n"))
	var strLines []string
	for _, line := range lines {
		strLines = append(strLines, string(line))
	}

	// 赋值给结构体
	pageData := Page{
		Title: "My Text File",
		Lines: strLines,
	}

	t, err := template.New("txtToHTML").Parse(tmpl)
	if err != nil {
		return "", err
	}

	html := new(bytes.Buffer)
	if err := t.Execute(html, pageData); err != nil {
		return "", err
	}

	return html.String(), nil
}

func EncodeURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	u.Path = encodePath(u.Path)
	return u.Scheme + "://" + u.Host + u.Path, nil
}

func CompareStrings(s1, s2 string) bool {
	if s1 == s2 {
		return true
	} else {
		return false
	}
}
func IsEmptyString(s string) bool {
	if len(s) == 0 {
		return true
	} else {
		return false
	}
}

func encodePath(path string) string {
	// 将路径按斜杠拆分
	splits := strings.Split(path, "/")
	// 对每个部分进行编码
	for i, s := range splits {
		if s != "" {
			splits[i] = url.PathEscape(s)
		}

	}
	// 拼接并返回编码后的路径
	return strings.Join(splits, "/")
}

func GetFileContent(filepath string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	reader := transform.NewReader(bytes.NewReader(data), simplifiedchinese.GBK.NewDecoder())
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	content = string(data)
	return content, nil
}

func RandomReadFile(path string, length int64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 获取文件大小
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := stat.Size()

	// 计算随机读取的起始位置
	rand.Seed(time.Now().Unix())
	startPos := rand.Int63n(fileSize - length)

	// 设置读取起始位置
	_, err = file.Seek(startPos, 0)
	if err != nil {
		return "", err
	}

	// 读取固定字数的内容
	buf := make([]byte, length)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func GetRandomString(strs []string) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(strs))
	return strs[index]
}

func GetFileName(filePath string) string {
	_, fileName := filepath.Split(filePath)
	return fileName
}

func ToJSONString(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 从JSON字符串解析出对象
func FromJSONString(jsonStr string, obj interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), obj)
	if err != nil {
		return err
	}
	return nil
}

func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02" // 日期格式
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func SubDay(t time.Time, day int) time.Time {
	return t.AddDate(0, 0, day)
}

func TimeToString(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}

func GetFileNameWithoutExt(fileName string) string {
	return filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])
}

// getYearMonthToDay 查询指定年份指定月份有多少天
// @params year int 指定年份
// @params month int 指定月份
func GetYearMonthToDay(year int, month int) int {
	// 有31天的月份
	day31 := map[int]struct{}{
		1:  struct{}{},
		3:  struct{}{},
		5:  struct{}{},
		7:  struct{}{},
		8:  struct{}{},
		10: struct{}{},
		12: struct{}{},
	}
	if _, ok := day31[month]; ok {
		return 31
	}
	// 有30天的月份
	day30 := map[int]struct{}{
		4:  struct{}{},
		6:  struct{}{},
		9:  struct{}{},
		11: struct{}{},
	}
	if _, ok := day30[month]; ok {
		return 30
	}
	// 计算是平年还是闰年
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		// 得出2月的天数
		return 29
	}
	// 得出2月的天数
	return 28
}

func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func ReadJson(name string) string {
	b, err := ioutil.ReadFile(name) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'
	return str
}

func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 生成一个六位数的随机验证码
func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06v", rand.Intn(1000000))
}

func Myuuid() string {
	uuid := uuid.New().String()
	uuidWithoutHyphens := strings.Replace(uuid, "-", "", -1)
	return uuidWithoutHyphens
}

func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
