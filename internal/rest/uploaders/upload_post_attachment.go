package uploaders

import (
	"io"
	"mime/multipart"
	"net/http"
	"socio/errors"
	postpb "socio/internal/grpc/post/proto"
	"socio/pkg/static"
)

func UploadPostAttachment(r *http.Request, postClient postpb.PostClient, fh *multipart.FileHeader) (string, error) {
	if fh == nil {
		return "", errors.ErrInvalidData
	}
	fileName := static.GetUniqueFileName(fh.Filename)
	contentType := fh.Header.Get("Content-Type")

	stream, err := postClient.Upload(r.Context())
	if err != nil {
		return "", err
	}

	file, err := fh.Open()
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

		err = stream.Send(&postpb.UploadRequest{
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
