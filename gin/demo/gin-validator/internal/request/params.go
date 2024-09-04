package request

// UserParams 注册参数
type Signup struct {
	// 必填，并且其长度至少为 4 个字符，最多为 20 个字符。
	Username string `json:"username" binding:"required,min=4,max=20,username"`
	// 必需的，并且应符合一个名为 passwd 的自定义验证规则
	Password string `form:"password" json:"password" binding:"required,password"`
	// 必需的，并且其值必须与 Password 字段相等
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	// 如果字段为空或未设置，则在序列化为 JSON 时不包括它
	// 如果字段非空，则应满足一个名为 age 的验证规则，这可能是指定年龄范围
	Age int `json:"age,omitempty" binding:"omitempty,gte=0,lte=130"`
	// 如果字段非空，则其值必须是 1 或 0，通常 1 可能代表男性，0 代表女性。
	Gender int `json:"gender,omitempty" binding:"omitempty,oneof=1 0"`
	// 必填项，有效的电子邮件地址
	Email string `json:"email,omitempty" binding:"omitempty,required,email"`
	// 必填项，有效的手机号
	Mobile string `json:"mobile,omitempty" binding:"omitempty,required,mobile"`
}



// GetMessages 自定义错误信息
func (register Signup) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":      "用户名称不能为空,长度为4-20位字符",
		"password.required":  "用户密码不能为空",
		"password.password":  "密码格式不正确",
		"rePassword.eqfield": "两次密码不一致",
		"email.required":     "邮箱不能为空",
		"email.email":        "邮箱格式不正确",
		"mobile.required":    "手机号码不能为空",
		"mobile.mobile":      "手机号码格式不正确",
	}
}
