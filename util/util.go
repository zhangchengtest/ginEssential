package util

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"image"
	"math/rand"
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

func RandomString(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
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
