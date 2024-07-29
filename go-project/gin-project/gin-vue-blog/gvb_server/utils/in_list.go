package utils

// 判断元素是否存在列表中
func InList(key string, list []string) bool {

	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}
