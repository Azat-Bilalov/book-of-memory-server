package config

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetMinioClient() *minio.Client {
	accessKey := "EDjZ6YLH37fDx4obqdbY"
	secretKey := "wmtSLej3MIF3ttSJSduiaJ23rBKr2qjBLhEuGKus"
	endpoint := "localhost:9000"
	useSsl := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSsl,
	})
	if err != nil {
		panic(err)
	}
	return minioClient
}

func CreateBucket(bucketName string) error {
	minioClient := GetMinioClient()
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			return nil
		}
		return err
	}
	return nil
}

func ReadObject(bucketName string, objectName string) (contentBytes []byte, contentType string, err error) {
	minioClient := GetMinioClient()
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	contentBytes, err = io.ReadAll(object)
	if err != nil {
		return
	}
	stat, err := object.Stat()
	if err != nil {
		return
	}
	contentType = stat.ContentType
	return
}

func UploadObject(bucketName string, objectName string, reader io.Reader, size int64, contentType string) error {
	minioClient := GetMinioClient()
	_, err := minioClient.PutObject(context.Background(), bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

func DeleteObject(bucketName string, objectName string) error {
	minioClient := GetMinioClient()
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
