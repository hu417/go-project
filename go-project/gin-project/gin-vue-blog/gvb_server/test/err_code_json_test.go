package test

import (
	"encoding/json"
	"fmt"
	"gvb_server/models/res"
	"os"
	"strconv"
	"testing"
)

//############ 错误状态码的封装

const file = "../models/res/err_code.json"

type ErrMap map[string]string

func TestErrCode(t *testing.T) {
	byteDate, err := os.ReadFile(file)

	if err != nil {
		panic(err)
	}

	var errMap = ErrMap{}
	err = json.Unmarshal(byteDate, &errMap)
	if err != nil {
		panic(err)
	}

	fmt.Println(errMap)
	fmt.Println(errMap[strconv.Itoa(int(res.SettingsError))])
}
