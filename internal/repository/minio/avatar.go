package minio

import (
	"github.com/minio/minio-go"
)

const (
	AvatarBucket = "user-avatars"
)

type AvatarStorage struct {
	MinioClient *minio.Client
}

func NewAvatarStorage(minioClient *minio.Client) (storage *AvatarStorage, err error) {
	bucketExists, err := minioClient.BucketExists(AvatarBucket)
	if err != nil {
		return
	}

	if !bucketExists {
		err = minioClient.MakeBucket(AvatarBucket, "ru-central1")
		if err != nil {
			return
		}
	}

	storage = &AvatarStorage{
		MinioClient: minioClient,
	}
	return
}

func (a *AvatarStorage) StoreAvatar(fileName string, filePath string) (err error) {
	_, err = a.MinioClient.FPutObject(AvatarBucket, fileName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return
	}

	return
}

func (a *AvatarStorage) DeleteAvatar(fileName string) (err error) {
	err = a.MinioClient.RemoveObject(AvatarBucket, fileName)
	if err != nil {
		return
	}

	return
}
