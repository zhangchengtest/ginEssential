package common

import (
	"bytes"
	"crypto/tls"
	"fmt"
	logging "github.com/ipfs/go-log"
	"github.com/jordan-wright/email"
	"html/template"
	"net/smtp"
)

var log = logging.Logger("common/email")

var host = "smtp.163.com"           // smtp服务host
var port = 465                      // ssl端口
var user = "zhangchengtest@163.com" // 发送方邮箱
var pass = "PDBDMQCSCTWZVPRD"       // 授权码

func SendRegister(userName string, email string) {

	templateData := struct {
		Name string
	}{
		Name: userName,
	}
	r := NewRequest([]string{email}, "新用户注册", "Hello, World!")

	if err := r.ParseTemplate("common/registrationEmail.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	} else {
		log.Errorf("read config failed ,err : %v", err)
		panic("read config failed, err: " + err.Error())
	}

}

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	e := email.NewEmail()
	e.From = "Pu Neng Shuo <zhangchengtest@163.com>"
	e.To = r.to
	e.Subject = r.subject
	e.HTML = []byte(r.body)

	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", user, pass, host),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		panic(err)
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
