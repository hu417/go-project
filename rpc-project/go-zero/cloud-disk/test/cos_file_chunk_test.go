package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"

	"fmt"

	"math"
	"os"
	"strconv"

	"net/http"
	"net/url"

	// "strconv"

	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// 文件分片
func TestCosChunkNumFile(t *testing.T) {
	// 读取文件
	fileInfo, err := os.Stat("../img/vue.png")
	if err != nil {
		t.Fatal(err)
	}

	// 分片个数 = 文件大小 / 分片大小
	// 390 / 100 ==> 4, 向上取整
	chunkNum := math.Ceil(float64(fileInfo.Size()) / chunkSize)
	// 只读方式打开文件
	myFile, err := os.OpenFile("../img/vue.png", os.O_RDONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// 存放每一次的分片数据
	b := make([]byte, chunkSize)
	// 遍历所有分片
	for i := 0; i < int(chunkNum); i++ {
		// 指定读取文件的起始位置
		myFile.Seek(int64(i*chunkSize), 0)
		// 最后一次的分片数据不一定是整除下来的数据
		// 例如: 文件 120M, 第一次读了 100M, 剩下只有 20M
		if chunkSize > fileInfo.Size()-int64(i*chunkSize) {
			b = make([]byte, fileInfo.Size()-int64(i*chunkSize))
		}
		myFile.Read(b)

		f, err := os.OpenFile("./"+strconv.Itoa(i)+".chunk", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		f.Write(b)
		f.Close()
	}
	defer myFile.Close()

}

// 分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID
	fmt.Println(UploadID) // 16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107
}

// 分片上传
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	UploadID := "16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107"
	f, err := os.ReadFile("0.chunk")
	if err != nil {
		t.Fatal(err)
	}
	// opt可选
	resp, err := client.Object.UploadPart(
		context.Background(), key, UploadID, 1, bytes.NewReader(f), nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	PartETag := resp.Header.Get("ETag") // chunk的Md5值: 82b9c7a5a3f405032b1db71a25f67021
	fmt.Println(PartETag)
}

// 分片上传完成
func TestPartUploadComplete(t *testing.T) {
	u, _ := url.Parse(define.BucketUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.COS_SECRETID,
			SecretKey: define.COS_SECRETKEY,
		},
	})

	key := "cloud-disk/111.png"
	UploadID := "16798358340b5b37b4bc47a02a4dca772975088fd07c8e27493e89c5cd7ad6c06eec2ee107"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{
		PartNumber: 1, ETag: "82b9c7a5a3f405032b1db71a25f67021"},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}
