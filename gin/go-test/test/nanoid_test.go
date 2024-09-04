package test

import (
	"testing"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func Test_Nanoid(t *testing.T) {
	// 生成一个默认长度的 ID
	id, err := gonanoid.New()
	if err != nil {
		t.Logf("Error generating ID: %v", err)
		return
	}
	t.Logf("Generated ID: %v", id)

	// 生成一个自定义长度的 ID
	customLengthID, err := gonanoid.Generate("abcde1234567890", 10) // 生成一个10个字符长的ID
	if err != nil {
		t.Logf("Error generating custom length ID: %v", err)
		return
	}
	t.Logf("Generated custom length ID: %v", customLengthID)
}
