package minio

import (
	"github.com/minio/minio-go"
)

const (
	UserAvatarsBucket  = "user-avatars"
	AttachmentsBucket  = "post-attachments"
	GroupAvatarsBucket = "group-avatars"
	StickersBucket     = "stickers"
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

		policy := `{
        "Version": "2012-10-17",
        "Statement": [
            {
                "Sid": "PublicRead",
                "Effect": "Allow",
                "Principal": "*",
                "Action": ["s3:GetObject"],
                "Resource": ["arn:aws:s3:::` + bucketName + `/*"]
            }
        ]
    }`

		err = minioClient.SetBucketPolicy(bucketName, policy)
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
