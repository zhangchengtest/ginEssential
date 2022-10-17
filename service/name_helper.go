package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

var firstName []string
var lastName []string
var lastNameLen int
var firstNameLen int

func init() {
	jsonConfigList := getName("adjective.json")
	firstName = deserializeJson(jsonConfigList)
	jsonConfigList = getName("noun.json")
	lastName = deserializeJson(jsonConfigList)
	lastNameLen = len(lastName)
	firstNameLen = len(firstName)
}

func getName(name string) string {
	b, err := ioutil.ReadFile(name) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	fmt.Println(str) // print the content as a 'string'
	return str
}
func deserializeJson(configJson string) []string {

	jsonAsBytes := []byte(configJson)
	configs := make([]string, 0)
	err := json.Unmarshal(jsonAsBytes, &configs)
	fmt.Printf("%#v", configs)
	if err != nil {
		panic(err)
	}
	return configs
}

func GetFullName() string {
	rand.Seed(time.Now().UnixNano())     //设置随机数种子
	var first string                     //名
	for i := 0; i <= rand.Intn(1); i++ { //随机产生2位或者3位的名
		first = fmt.Sprint(firstName[rand.Intn(firstNameLen-1)])
	}
	//返回姓名
	return fmt.Sprintf("%s的%s", first, fmt.Sprint(lastName[rand.Intn(lastNameLen-1)]))
}
