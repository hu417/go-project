package main


// 定义 Message 结构体
type Message struct {
	// 标题、内容、信息类型
	title, message, messageType string
 
	// 账号
	account     string
	accountList []string
 
	// token
	token     string
	tokenList []string
 }
 
 // 定义 MessageOptionm选项函数类型，用于接收 Message 参数的函数
 type MessageOption func(*Message)
 
 /*
	定义 NewMessage 函数，用于创建一个 Message 指针变量，在 NewMessage 函数中，固定参数包括 title、message 和 messageType，它们是必需的参数。
	然后，通过可选参数 opts ...MessageOption 来接收一系列的函数选项
 */
 func NewMessage(title, message, messageType string, opts ...MessageOption) *Message {
	msg := &Message{
	   title:       title,
	   message:     message,
	   messageType: messageType,
	}
 
	for _, opt := range opts {
	   opt(msg)
	}
 
	return msg
 }
 
 /*
 	定义四个选项函数: WithAccount、WithAccountList、WithToken 和 WithTokenList。
	这些选项函数分别用于设置被推送消息的账号、账号列表、令牌和令牌列表
 */
 func WithAccount(account string) MessageOption {
	return func(message *Message) {
	   message.account = account
	}
 }
 
 func WithAccountList(accountList []string) MessageOption {
	return func(message *Message) {
	   message.accountList = accountList
	}
 }
 
 func WithToken(token string) MessageOption {
	return func(message *Message) {
	   message.token = token
	}
 }
 
 func WithTokenList(tokenList []string) MessageOption {
	return func(message *Message) {
	   message.tokenList = tokenList
	}
 }
 

 func main() {
	// 单账号推送
	/*
		创建单账号推送的消息，通过调用 NewMessage 并传递相应的参数和选项函数（WithAccount）来配置消息
	*/
	_ = NewMessage(
	   "来自陈明勇的信息",
	   "你好，我是陈明勇",
	   "单账号推送",
	   WithAccount("123456"),
	)
 
	// 多账号推送
	/*
		创建多账号推送的消息，同样通过调用 NewMessage 并使用不同的选项函数（WithAccountList）来配置消息
	*/
	_ = NewMessage(
	   "来自陈明勇的信息",
	   "你好，我是陈明勇",
	   "多账号推送",
	   WithAccountList([]string{"123456", "654321"}),
	)
 }
 