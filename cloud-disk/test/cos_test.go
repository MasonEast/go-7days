package test

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/tencentyun/cos-go-sdk-v5"
)

func TestCosFielUpload(t *testing.T) {
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
    // 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
    u, _ := url.Parse("https://1-1301230924.cos.ap-chengdu.myqcloud.com")
    b := &cos.BaseURL{BucketURL: u}
    fmt.Println(define.CosSecretID, define.CosSecretKey)
    client := cos.NewClient(b, &http.Client{
        Transport: &cos.AuthorizationTransport{
            // 通过环境变量获取密钥
            // 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
            SecretID: define.CosSecretID,
            // 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
            SecretKey: define.CosSecretKey,
        },
    })

    // 存储桶路径
    key := "cloud-disk/test.jpg"

    _, _, err := client.Object.Upload(
        context.Background(), key, "./img/1.jpg", nil,
    )
    if err != nil {
        t.Fatal(err)
    }
}

// 分片上传初始化
func TestInitPartUpload(t *testing.T) {
	u, _ := url.Parse("https://getcharzp-1256268070.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.CosSecretID,
			SecretKey: define.CosSecretKey,
		},
	})
	key := "cloud-disk/exampleobject.jpeg"
	v, _, err := client.Object.InitiateMultipartUpload(context.Background(), key, nil)
	if err != nil {
		t.Fatal(err)
	}
	UploadID := v.UploadID // 16516751077a264ca9b5c8f560c1d1b5ea5e9d242dee047f62b02b4958491c0d90aa167d51
	fmt.Println(UploadID)
}

// 分片上传
func TestPartUpload(t *testing.T) {
	u, _ := url.Parse("https://getcharzp-1256268070.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.CosSecretID,
			SecretKey: define.CosSecretKey,
		},
	})
	key := "cloud-disk/exampleobject.jpeg"
	UploadID := "16516751077a264ca9b5c8f560c1d1b5ea5e9d242dee047f62b02b4958491c0d90aa167d51"
	f, err := os.ReadFile("0.chunk") // md5 : 108e92d35fe1695fbf29737d0b24561d
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
	PartETag := resp.Header.Get("ETag")
	fmt.Println(PartETag)
}

// 分片上传完成
func TestPartUploadComplete(t *testing.T) {
	u, _ := url.Parse("https://getcharzp-1256268070.cos.ap-chengdu.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.CosSecretID,
			SecretKey: define.CosSecretKey,
		},
	})
	key := "cloud-disk/exampleobject.jpeg"
	UploadID := "16516751077a264ca9b5c8f560c1d1b5ea5e9d242dee047f62b02b4958491c0d90aa167d51"

	opt := &cos.CompleteMultipartUploadOptions{}
	opt.Parts = append(opt.Parts, cos.Object{
		PartNumber: 1, ETag: "108e92d35fe1695fbf29737d0b24561d"},
	)
	_, _, err := client.Object.CompleteMultipartUpload(
		context.Background(), key, UploadID, opt,
	)
	if err != nil {
		t.Fatal(err)
	}
}
