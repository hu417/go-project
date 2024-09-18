package utils

import "strconv"

// str转int
func StrToInt(str string) (int) {
	num,err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return num
}

func StrToInt64(str string) (int64) {
	num,err := strconv.ParseInt(str,10,64)
	if err != nil {
		panic(err)
	}
	return num
}