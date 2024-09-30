package utils

import (
	"fmt"
	"go-admin/global"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEmail(to string, subject string, text string) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = global.EmailFrom
	// 设置接收方的邮箱
	e.To = []string{to}
	//设置主题
	e.Subject = subject
	//设置文件发送的内容
	e.Text = []byte(text)
	//设置服务器相关的配置
	emailAddr := global.EmailHost + ":" + global.EmailPort
	err := e.Send(emailAddr, smtp.PlainAuth("", global.EmailFrom, global.EmailPassWord, global.EmailHost))
	if err != nil {
		fmt.Println(err)
	}
}
