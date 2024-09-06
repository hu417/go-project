package req

// UserRegister 用户注册
type UserRegister struct {
	Name     string `json:"name" binding:"required,min=4,max=20"`                 // 用户名称,长度为4-20位字符
	Password string `form:"password" json:"password" binding:"required,password"` // 必需的，并且应符合一个名为 passwd 的自定义验证规则，密码格式为 6-20 位数字和字母组合
	Mobile   string `json:"mobile,omitempty" binding:"omitempty,required,mobile"` // 必需的，并且是有效的手机号码
}

// GetMessages 自定义错误信息
func (register UserRegister) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":     "用户名称不能为空,长度为4-20位字符",
		"password.required": "用户密码不能为空",
		"password.password": "密码格式不正确,6-20位数字和大小写字母组合",
		//"rePassword.eqfield": "两次密码不一致",
		// "email.required":     "邮箱不能为空",
		// "email.email":        "邮箱格式不正确",
		"mobile.required": "手机号码不能为空",
		"mobile.mobile":   "手机号码格式不正确",
	}
}

// UserLogin 用户登录
type UserLogin struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`       // 必需的，并且是有效的手机号码
	Password string `form:"password" json:"password" binding:"required,password"` // 必需的，并且应符合一个名为 passwd 的自定义验证规则，密码格式为 6-20 位数字和字母组合
}

// GetMessages 自定义错误信息
func (login UserLogin) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"mobile.required":   "用户名称不能为空",
		"mobile.mobile":     "手机号码格式不正确,11位手机号码",
		"password.required": "密码不能为空",
		"password.password": "密码格式不正确,6-20位数字和大小写字母组合",
	}
}
