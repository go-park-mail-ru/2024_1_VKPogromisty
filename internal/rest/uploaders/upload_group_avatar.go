package uploaders

import (
	"io"
	"mime/multipart"
	"net/http"

	pgpb "socio/internal/grpc/public_group/proto"
	"socio/pkg/static"
)

func UploadPublicGroupAvatar(r *http.Request, publicGroupClient pgpb.PublicGroupClient, avatarFH *multipart.FileHeader) (string, error) {
	if avatarFH == nil {
		return "", nil
	}

	fileName := static.GetUniqueFileName(avatarFH.Filename)
	contentType := avatarFH.Header.Get("Content-Type")

	stream, err := publicGroupClient.Upload(r.Context())
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

		err = stream.Send(&pgpb.UploadRequest{
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
