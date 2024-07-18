package email

// 邮件验证码
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"reflect"
	"text/template"

	"github.com/k3a/html2text"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
)

var (
	Email *MailOptions
	ctx   = context.Background()

	// 邮箱配置
	EmailHost = "smtp.163.com"
	EmailPort = 465
	EmailUser = "hu729919300@163.com"
	EmailPass = "BYDLCRTRLUFACRXY"
	EmailTo   = "hu729919300@163.com"
	EmailCc   = ""
)

type MailOptions struct {
	MailHost    string
	MailPort    int
	MailUser    string          // 发件人
	MailPass    string          // 发件人密码: BYDLCRTRLUFACRXY
	MailTo      string          // 收件人 多个用,分割
	Subject     string          // 邮件主题
	Body        string          // 邮件内容
	MailDialer  *gomail.Dialer  // 邮件对象
	MailMessage *gomail.Message // 消息
	MailBody    bytes.Buffer    // 内容

}

// Emailer 构造函数
func NewEmailer(host string, port int, user, pass string) *MailOptions {
	Email = &MailOptions{
		MailHost:    host,
		MailPort:    port,
		MailUser:    user,
		MailPass:    pass,
		MailDialer:  gomail.NewDialer(host, port, user, pass),
		MailMessage: gomail.NewMessage(),
	}
	return Email
}

// Setup 初始化邮件函数
func (e *MailOptions) Setup() *gomail.Dialer {
	if e.MailDialer == nil {
		// 实例化gomail邮件连接对象
		emailhost := EmailHost
		emailport := EmailPort
		emailuser := EmailUser
		emailpass := EmailPass

		// 若有错误就关闭连接
		var s gomail.SendCloser
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("gomail.NewDialer connect failed!")
				s.Close()
				panic(err)
			}
		}()

		// 创建 smtp 实例
		// QQ 邮箱：SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/587, 非SSL协议端口：25）
		// 163 邮箱：SMTP 服务器地址：smtp.163.com（SSL协议端口：465/994，非SSL协议端口：25）
		e.MailDialer = gomail.NewDialer(emailhost, emailport, emailuser, emailpass)
		// 允许跳过不安全的认证
		// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		return e.MailDialer
	}
	return e.MailDialer
}

// MailDialAndSend 发送邮件后立即关闭连接。
func (e *MailOptions) MailDialAndSend(msg *gomail.Message) (string, error) {
	if err := e.MailDialer.DialAndSend(msg); err != nil {
		return "邮件发送失败!", err
	} else {
		msg.Reset()
		return "邮件发送成功!", nil
	}
}

// SendMsg text文本以及HTML格式邮件信息发送
func (e *MailOptions) SendMsg(to []string, cc string, subject, body, sendtype string) (string, error) {
	// 初始化连接验证
	e.Setup()

	// 构建一个 Message 对象也就是邮件对象
	e.MailMessage = gomail.NewMessage()

	// 设置发信人，收信人、抄送
	e.MailMessage.SetHeader("From", EmailUser)
	e.MailMessage.SetHeader("To", to...)
	// e.MailMessage.SetHeader("Cc",cc...)
	if cc != "" {
		// chaosong := strings.Split(cc, ",")
		e.MailMessage.SetHeader("Cc", cc)
		// e.MailMessage.SetAddressHeader("Cc", cc)
	}

	// 设置邮件标题与正文
	if subject != "" && body != "" {
		// 邮件标题
		e.MailMessage.SetHeader("Subject", subject)
		// 判断发送邮件的类型设置对应正文
		if sendtype == "text" {
			e.MailMessage.SetBody("text/plain", body)
		} else if sendtype == "html" {
			e.MailMessage.SetBody("text/html", body)
		}
	} else {
		return "邮件 subject 或 body 字段不能为空", errors.New("Email Message The subject or body field cannot be empty")
	}

	// 发送邮件
	if err := e.MailDialer.DialAndSend(e.MailMessage); err != nil {
		return "邮件发送失败!", err
	} else {
		e.MailMessage.Reset()
		return "邮件发送成功!", err
	}
}

// 发送邮件，包含附件
func (e *MailOptions) SendAttachMsg(to []string, cc string, subject, body, file, filename string, sendtype string) (string, error) {
	// 初始化连接验证
	e.Setup()

	// 构建一个 Message 对象也就是邮件对象
	e.MailMessage = gomail.NewMessage()

	// 设置发信人，收信人、抄送
	e.MailMessage.SetHeader("From", EmailUser)
	e.MailMessage.SetHeader("To", to...)
	// e.MailMessage.SetAddressHeader("Cc", cc...)
	if cc != "" {
		// chaosong := strings.Split(cc, ",")
		e.MailMessage.SetHeader("Cc", cc)
		// e.MailMessage.SetAddressHeader("Cc", cc)
	}

	// 设置邮件标题与正文
	if subject != "" && body != "" {
		// 邮件标题
		e.MailMessage.SetHeader("Subject", subject)
		// 判断发送邮件的类型设置对应正文
		if sendtype == "text" {
			e.MailMessage.SetBody("text/plain", body)
		} else if sendtype == "html" {
			e.MailMessage.SetBody("text/html", body)
		}
	} else {
		return "邮件 subject 或 body 字段不能为空", errors.New("Email Message The subject or body field cannot be empty")
	}

	// 添加附件，注意附件是需要传入完整路径
	if file != "" && len(file) > 0 {
		e.MailMessage.Attach(file,
			gomail.Rename(filename),
			gomail.SetHeader(map[string][]string{
				"Content-Disposition": {
					fmt.Sprintf(`attachment; filename="%s"`, mime.QEncoding.Encode("UTF-8", filename)),
				},
			}),
		)
	}

	// 发送邮件
	if err := e.MailDialer.DialAndSend(e.MailMessage); err != nil {
		e.MailMessage.Reset()
		return "邮件发送失败!", err
	} else {
		e.MailMessage.Reset()
		return "邮件发送成功!", nil
	}
}

// 使用template定制化html
type EmailData struct {
	SiteName     string
	UserName     string
	UserCode     string
	UserCodeTime string
	SiteAddr     string
}

// 1、判断目录是否存在
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}


// 此处值得学习利用反射实现传入string字符串解析为对应函数以及call参数进行执行，
func (e *MailOptions) reflectFunc(tpl *template.Template, tplname string, body ...string) string {
	//var t MailOptions
	ref := reflect.ValueOf(&e)
	refFunc := ref.MethodByName(tplname)
	fmt.Printf("Kind : %s, Type : %s\n", refFunc.Kind(), refFunc.Type())
	refVal := make([]reflect.Value, 0)
	refVal = append(refVal, reflect.ValueOf(tpl))
	for _, v := range body {
		fmt.Println(v)
		refVal = append(refVal, reflect.ValueOf(v))
	}
	fmt.Println(refVal)
	res := refFunc.Call(refVal)
	return res[0].String()

}

// 2、模板信息发送
// func (e *MailOptions) TemplateMsg(to []string, tplname, subject string, bodys ...string) (string, error) {
func (e *MailOptions) TemplateMsg(to []string, temphtml,subject string, data *EmailData) (string, error) {
	// 初始化连接验证
	e.Setup()
	fmt.Println("初始化完成")
	// 构建一个 Message 对象也就是邮件对象
	e.MailMessage = gomail.NewMessage()

	// 设置发信人，收信人、抄送
	e.MailMessage.SetHeader("From", EmailUser)
	e.MailMessage.SetHeader("To", to...)
	e.MailMessage.SetHeader("Subject", subject)

	// 模板读取并返回template对象
	var body bytes.Buffer
	template, err := ParseTemplateDir("../../template/email")
	if err != nil {
		fmt.Printf("err => %v", err)
		return "传入的模板文件名称有误", err
		//return errors.New("could not parse template")
	}
	template.ExecuteTemplate(&body, fmt.Sprintf("%s.html",temphtml), &data)
	htmlString := body.String()
	prem, _ := premailer.NewPremailerFromString(htmlString, nil)
	htmlInline, err := prem.Transform()

	// 邮件正文
	e.MailMessage.SetBody("text/html", htmlInline)
	e.MailMessage.AddAlternative("text/plain", html2text.HTML2Text(htmlString))

	// 邮件发送
	res, err := e.MailDialAndSend(e.MailMessage)
	fmt.Println("邮件发送中")
	return res, err
}
