package test

import (
	"crypto/tls"
	"fmt"

	"math/rand"
	"net/smtp"
	"strconv"
	"testing"

	"cloud-disk/core/define"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Get <hu729919300@163.com>"   // 来自
	e.To = []string{"hu729919300@163.com"} // 发送
	// e.Cc = []string{"xxx@126.com"} // 抄送
	e.Subject = "验证码测试发送" // 主题

	num := rand.Intn(100000) + 100000
	fmt.Println("num = ", num)

	e.HTML = []byte("你的验证码是: <h3>" + strconv.Itoa(num) + "</h3>") // 内容
	//err := e.Send("smtp.163.com:465", smtp.PlainAuth("", "hu729919300@163", define.Pwd, "smtp.163.com"))
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "hu729919300@163.com", define.Pwd, "smtp.163.com"), &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "smtp.163.com",
	})

	if err != nil {
		t.Fatal("failed to send email: ", err)
	}
}
