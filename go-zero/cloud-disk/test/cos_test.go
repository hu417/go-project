package test

import (
	"context"
	"fmt"
	"os"

	"net/url"

	"net/http"

	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

func TestCosInit(t *testing.T) {
	fmt.Println(os.Getenv("COS_SECRETID"))
	fmt.Println(os.Getenv("COS_SECRETKEY"))
	fmt.Println("--------------------------------")

	// 存储桶名称，由 bucketname-appid + 地域组成
	u, _ := url.Parse("https://cloud-1304907914.cos.ap-guangzhou.myqcloud.com")
	b := &cos.BaseURL{
		BucketURL: u,
	}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId
			SecretID: os.Getenv("COS_SECRETID"),
			// 环境变量 SECRETKEY 表示用户的 SecretKey
			SecretKey: os.Getenv("COS_SECRETKEY"),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	// 查询
	// opt := &cos.BucketGetOptions{
	// 	Prefix:    "cloud-disk", // prefix 表示要查询的文件夹
	// 	Delimiter: "/",          // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
	// 	MaxKeys:   3,            // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	// }
	// v, _, err := c.Bucket.Get(context.Background(), opt)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, c := range v.Contents {
	// 	fmt.Printf("%s, %d\n", c.Key, c.Size)
	// }

	// 上传,需要指定路径及文件名
	key := "cloud-disk/vue1.png"

	if _, _, err := c.Object.Upload(
		context.Background(), key, "../img/vue.png", nil,
	); err != nil {
		panic(err)
	}

	// 下载
	// key = "cloud-disk/vue.png"
	// file := "../img/vue.png"

	// opt := &cos.MultiDownloadOptions{
	// 	ThreadPoolSize: 5,
	// }
	// _, err := c.Object.Download(
	// 	context.Background(), key, file, opt,
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// 列出
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    "cloud-disk/", // prefix 表示要查询的文件夹
		Delimiter: "/",           // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000,          // 设置最大遍历出多少个对象, 一次 listobject 最大支持1000
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.Bucket.Get(context.Background(), opt)
		if err != nil {
			fmt.Println(err)
			break
		}
		// common prefix 表示表示所有以 Prefix 开头，被 delimiter的值 截断的路径, 如 delimter 设置为/, common prefix 则表示所有子目录的路径
		for _, commonPrefix := range v.CommonPrefixes {
			fmt.Printf("当前目录: %v\n", commonPrefix)
		}

		for _, content := range v.Contents {
			fmt.Printf("cos对象文件夹/文件: %v\n", content.Key)
		}

		isTruncated = v.IsTruncated // 是否还有数据
		marker = v.NextMarker       // 设置下次请求的起始 key
	}

	// 删除
	// 1、删除文件
	key = "cloud-disk/vue1.png"
	_, err := c.Object.Delete(context.Background(), key)
	if err != nil {
		panic(err)
	}
	// 2、删除文件夹
	keys := "test/"
	if _, err := c.Object.Delete(context.Background(), keys); err != nil {
		panic(err)
	}
}
