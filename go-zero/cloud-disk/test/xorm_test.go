package test

import (
	"bytes"
	"cloud-disk/models"
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

func TestXorm(t *testing.T) {
	// 1、建立连接
	engine, err := xorm.NewEngine("mysql", "root:123456@tcp(10.0.0.91:3306)/cloud-disk?charset=utf8")
	engine.ShowSQL(true)

	if err != nil {
		t.Fatal(err)
	}
	defer engine.Close()

	// 2、执行curd
	data := make([]*models.User_basic, 0)
	err = engine.Find(&data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("data => ", data) // [0xc0001942d0]
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	dst := new(bytes.Buffer)
	err = json.Indent(dst, b, "", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("dst => ", dst.String())
	fmt.Printf("data 类型: %T", dst)

}
