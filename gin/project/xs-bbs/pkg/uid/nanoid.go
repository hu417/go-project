package uid

import (
	"fmt"
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func Nanoid(t *testing.T) string {
	// 生成一个默认长度的 ID
	id, err := gonanoid.New()
	if err != nil {
		fmt.Printf("Error generating ID: %v", err)
	}
	return id

	// // 生成一个自定义长度的 ID
	// customLengthID, err := gonanoid.Generate("abcde1234567890", 10) // 生成一个10个字符长的ID
	// if err != nil {
	// 	fmt.Printf("Error generating custom length ID: %v", err)
	// }
	// return customLengthID
}
