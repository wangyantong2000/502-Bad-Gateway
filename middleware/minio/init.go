package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

//const (
//	endpoint        string = "192.168.100.129:9000"
//	accessKeyID     string = "douyin"
//	secretAccessKey string = "douyin123"
//	useSSL          bool   = false
//	VideoBucketName string = "video"
//)

var (
	minioClient     *minio.Client
	endpoint        string = "192.168.100.129:9000"
	accessKeyID     string = "douyin"
	secretAccessKey string = "douyin123"
	useSSL          bool   = false
	VideoBucketName string = "video"
)

func Init() {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL})
	if err != nil {
		log.Fatalln("minio连接错误: ", err)
	}
	minioClient = client
	if err := CreateBucket(VideoBucketName); err != nil {
		log.Println(err)
	}
	log.Printf("%#v\n", client)
}
