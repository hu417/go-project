package utils

import (
	"fmt"
	"os"
)

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) error {
	if path == "" {
		return fmt.Errorf("path is empty")
	}
	_, err := os.Stat(path)

	if err == nil {
		return nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在

		// return false,nil

		/*
			// 创建单层目录
			err := os.Mkdir("testdir", 0755)
			if err != nil {
				fmt.Println("create directory error:", err)
				return
			}
		*/

		// 递归创建多层目录
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("create directory [%s] error => %w", path,err)
		}

	}
	// return false,err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
	return err
}

func GetFilename() {
	// 遍历目录文件信息
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("read directory error:", err)
		return
	}

	// 打印文件信息
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
	}

	// 遍历目录文件名
	fileNames, err := os.ReadDir(".")
	if err != nil {
		fmt.Println("read directory names error:", err)
		return
	}

	// 打印文件名
	for _, fileName := range fileNames {
		fmt.Println(fileName)
	}

}

// 删除
func DeletePath() {
	// 删除单个目录
	err := os.Remove("testdir")
	if err != nil {
		fmt.Println("remove directory error:", err)
		return
	}

	// 递归删除多层目录
	err = os.RemoveAll("testdir")
	if err != nil {
		fmt.Println("remove directory error:", err)
		return
	}

}
