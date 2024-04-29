package minio

import (
	"github.com/minio/minio-go"
)

const (
	AvatarBucket      = "user-avatars"
	AttachmentsBucket = "post-attachments"
)

type StaticStorage struct {
	bucketName  string
	MinioClient *minio.Client
}

func NewStaticStorage(minioClient *minio.Client, bucketName string) (storage *StaticStorage, err error) {
	bucketExists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		return
	}

	if !bucketExists {
		err = minioClient.MakeBucket(bucketName, "ru-central1")
		if err != nil {
			return
		}
	}

	storage = &StaticStorage{
		bucketName:  bucketName,
		MinioClient: minioClient,
	}
	return
}

func (a *StaticStorage) Store(fileName string, filePath string, contentType string) (err error) {
	_, err = a.MinioClient.FPutObject(a.bucketName, fileName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return
	}

	return
}

func (a *StaticStorage) Delete(fileName string) (err error) {
	err = a.MinioClient.RemoveObject(a.bucketName, fileName)
	if err != nil {
		return
	}

	return
}
