package utils

import "strconv"

func StrToInt(str string) (int64, error) {

	return strconv.ParseInt(str, 10, 64)
}
