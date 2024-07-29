package test

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/jordan-wright/email"
)

// 1、普通发送
func testSimpleEmail() {
	e := email.NewEmail()
	e.From = "lym <1154041111@qq.com>"
	e.To = []string{"15074941111@163.com"}
	e.Subject = "【测试】QQ邮箱给163邮箱发送邮件"
	e.Text = []byte("测试邮件，收到可以忽略")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1154041111@qq.com", "yyy", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}

// 2、添加抄送
func testCarbonCopy() {
	e := email.NewEmail()
	e.From = "lym <1154041111@qq.com>"
	e.To = []string{"15074941111@163.com"}
	// 主要就是在这里添加了抄送人以及秘密抄送人
	e.Cc = []string{"1315381111@qq.com"}
	e.Bcc = []string{"1315381111@qq.com"}
	e.Subject = "【测试】QQ邮箱发送邮件，并抄送和秘密抄送发送"
	e.Text = []byte("测试抄送和密码抄送邮件，收到可以忽略")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1154041111@qq.com", "yyy", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}

// 3、发送html
func testHTMLEmail() {
	e := email.NewEmail()
	e.From = "lym <1154041111@qq.com>"
	e.To = []string{"15074941111@163.com"}
	e.Cc = []string{"1154041111@qq.com"}
	e.Subject = "【测试】发送HTML邮件"

	// 主要就是在这里添加HTML字节数组内容
	e.HTML = []byte(`
  <ul>
<li><a "https://pkg.go.dev/github.com/jordan-wright/email/">Go邮件库github.com/jordan-wright/email文档地址</a></li>
<li><a "https://darjun.github.io/2020/01/10/godailylib/go-flags/">Go邮件库github.com/jordan-wright/email Github地址</a></li>
</ul>
  `)
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1154041111@qq.com", "yyy", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}

// 4、发送附件
func testAttachFile() {
	e := email.NewEmail()
	e.From = "lym <1154041111@qq.com>"
	e.To = []string{"15074941111@163.com"}
	e.Subject = "【测试】邮件携带附件"
	e.Text = []byte("测试邮件，请查看是否有附件")

	// 主要就是在这个位置添加了附件，注意：这里AttachFile是Email结构体的方法而不是字段，所以不是用等于而是括号
	e.AttachFile("30-go-email/test.txt")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "1154041111@qq.com", "yyy", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}

// 5、多次发送
func testEmailPool() {
	ch := make(chan *email.Email, 10)

	// 创建连接池，这里的参数和Send方法的参数几乎一致，不过多了一个连接池数量的参数
	p, err := email.NewPool(
		"smtp.qq.com:587",
		4,
		smtp.PlainAuth("", "1154041111@qq.com", "yyy", "smtp.qq.com"),
	)

	if err != nil {
		log.Fatal("failed to create pool:", err)
	}

	// 开启四个协程充当消费者，发送邮件
	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				// 这里的Send是连接池对象的方法，第一个参数是Email对象，第二个参数表示超时时间，这个时间内没有发送成功会报错
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}

	for i := 0; i < 10; i++ {
		e := email.NewEmail()
		e.From = "lym <1154041111@qq.com>"
		e.To = []string{"15074941111@163.com"}
		e.Subject = fmt.Sprintf("测试连接池：%d", i+1)
		e.Text = []byte(fmt.Sprintf("测试连接池：%d", i+1))
		ch <- e
	}

	close(ch) // 发送完毕后，关闭通道
	wg.Wait() // 通道关闭后,ch中消费完后，4个消费者会退出协程，从而放行
}

func main() {
	testEmailPool()
}
