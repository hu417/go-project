package email

import (
	"fmt"
	"testing"
)

func Test_email(t *testing.T) {
	// 执行 Emailer 的构造函数
	obj := NewEmailer("smtp.163.com", 465, "hu729919300@163.com", "BYDLCRTRLUFACRXY")

	// 定义要发送的邮箱地址数组
	mailTo := []string{
		"hu729919300@163.com",
	}
	// cc := "humeng.417@bytedance.com"
	//cc := ""

	// 调用示例1: 文本信息发送
	//obj.SendMsg(mailTo,cc,"文本邮件示例","来自GG的消息","text")

	// 调用示例2: HTML信息发送带
	// obj.EmailMsg(mailTo,"chaosong@weiyigeek.top","网页邮件示例","来自<b>【全栈工程师修炼指南】</b>公众号的消息","html")

	// 调用示例3: HTML信息发送带附件
	//obj.SendAttachMsg(mailTo,cc,"网页邮件示例","来自<b>【GG】</b>的消息","/Users/bytedance/1.txt","测试.txt","html")

	// 调用示例4: 模版化html
	// 结构体参数为html中的变量
	data := &EmailData{
		SiteName:     "测试",
		UserName:     "hhh",
		UserCode:     "1234",
		UserCodeTime: "3",
		SiteAddr:     "http://www.baidu.com",
	}
	str, err := obj.TemplateMsg(mailTo, "test" ,"模板邮件示例", data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str)
}
