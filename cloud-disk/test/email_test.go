package test

import (
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)


func TestEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "mason <sh941118@163.com>"
	e.To = []string{"2224742726@qq.com"}

	e.Subject = "Awesome Subject"
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.163.com:25", smtp.PlainAuth("", "sh941118@163.com", "VSCYFUSEUZXKMYMQ", "smtp.163.com"))
	if err != nil {
		t.Fatal(err)
	}
}