package utils

import "regexp"

// IsValidUsername 检查字符串是否为有效的用户名
func IsValidUsername(username string) bool {
	// 正则表达式
	// ^[\p{L}\d]{1,15}$
	// 解释：
	// - ^: 匹配字符串的开始
	//   - \pL: 匹配字母
	// - {1,15}: 匹配1到15个字符
	// - $: 匹配字符串的结束
	regex := regexp.MustCompile(`^[\p{L}\d]{1,15}$`)
	return regex.MatchString(username)
}

// IsValidEmail 检查字符串是否为有效的电子邮件地址
func IsValidEmail(email string) bool {
	// 使用正则表达式验证电子邮件地址的格式
	// 示例格式：[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}
	// 匹配规则：
	// - [a-zA-Z0-9._%+-]+: 匹配电子邮件地址的前缀部分，可以是字母、数字、下划线、百分号、加号、减号和点
	// - @: 匹配"@"符号
	// - [a-zA-Z0-9.-]+: 匹配电子邮件地址的主体部分，可以是字母、数字、点或减号
	// - \.: 匹配"."符号
	// - [a-zA-Z]{2,}: 匹配电子邮件地址的后缀部分，必须是至少两个字母
	// 如果匹配成功，则返回true，否则返回false
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// IsValidPhoneNumber 检查字符串是否为有效的中国大陆电话号码
func IsValidPhoneNumber(phone string) bool {
	// 这里可以根据你的需求调整正则表达式
	// 示例格式：1[3-9]\d{9}
	// 匹配规则：
	// - 1: 匹配"1"
	// - [3-9]: 匹配3到9之间的数字
	// - \d{9}: 匹配9个数字
	// 如果匹配成功，则返回true，否则返回false
	regex := regexp.MustCompile(`^(1[3-9]\d{9})$`)
	return regex.MatchString(phone)
}

// IsValidPassword 验证密码格式
func IsValidPassword(password string) bool {
	// 这里可以根据你的需求调整正则表达式
	// 示例格式：[a-zA-Z0-9]{8,20}
	// 匹配规则：
	// - [a-zA-Z0-9]{8,20}: 匹配8到20个字母或数字
	// 如果匹配成功，则返回true，否则返回false
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{8,20}$`)
	return regex.MatchString(password)
}
