package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"net/url"
	"time"
)

// 创建bucket
func CreateBucket(bucketName string) error {
	ctx := context.Background()
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: "cn-south-1", ObjectLocking: false})
	if err != nil {
		exists, _ := minioClient.BucketExists(ctx, bucketName)
		if exists {
			log.Printf("bucket: %s已经存在", bucketName)
			return nil
		}
		return err
	} else {
		log.Printf("Successfully created bucket:%s\n", bucketName)
	}
	return nil
}

// 上传文件
func FileUploader(bucketName string, objectName string, filePath string, contentType string) (int64, error) {
	object, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("上传失败")
		return -1, err
	}
	return object.Size, nil
}

// put文件
func PutFile(bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := minioClient.PutObject(context.Background(), bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Printf("上传失败")
		return err
	}
	return nil
}

// 获取文件Url
func GetFileUrl(bucketName string, objectName string) (*url.URL, error) {
	ExpireTime := 3600
	expires := time.Second * time.Duration(ExpireTime) //视频临时链接过期秒数
	//reqParams := make(url.Values)
	presignedUrl, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expires, nil)
	if err != nil {
		log.Printf("获取文件Url失败")
		return nil, err
	}
	return presignedUrl, nil
}
