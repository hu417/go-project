
package main

import "fmt"

type A interface {
	FunctionA()
}

type B interface {
	FunctionB()
}

// 定义一个新的接口C，继承接口A和B。
type C interface {
	A
	B
	FunctionC()
}

//注意：使用关键字“interface”后可在接口定义中使用“{ }”进行多继承。

// 实现接口C中的方法。
type MyStruct struct{}

func NewMyStruct() *MyStruct {
	
	return &MyStruct{}
}

func (s *MyStruct) FunctionA() {
	// 实现A接口中的方法
	fmt.Println("这是a方法")
}

func (s *MyStruct) FunctionB() {
	// 实现B接口中的方法
	fmt.Println("这是b方法")
}

func (s *MyStruct) FunctionC() {
	// 实现C接口中的方法
	fmt.Println("这是c方法")
}

func main() {
	// 使用实现了接口C的结构体对象
	var obj C = NewMyStruct()
	obj.FunctionA()
	obj.FunctionB()
	obj.FunctionC()
}
