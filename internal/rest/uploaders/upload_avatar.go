package uploaders

import (
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	uspb "socio/internal/grpc/user/proto"

	"github.com/google/uuid"
)

const (
	BatchSize = 1 << 23
)

func UploadAvatar(r *http.Request, userClient uspb.UserClient, avatarFH *multipart.FileHeader) (string, error) {
	fileName := uuid.NewString() + filepath.Ext(avatarFH.Filename)
	contentType := avatarFH.Header.Get("Content-Type")

	stream, err := userClient.Upload(r.Context())
	if err != nil {
		return "", err
	}

	file, err := avatarFH.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, BatchSize)
	batchNumber := 1

	for {
		num, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		chunk := buf[:num]

		err = stream.Send(&uspb.UploadRequest{
			FileName:    fileName,
			Chunk:       chunk,
			ContentType: contentType,
		})

		if err != nil {
			return "", err
		}
		batchNumber += 1
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.FileName, nil
}
