package static

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"socio/errors"

	"github.com/google/uuid"
)

const DefaultAvatarFileName = "default_avatar.png"

var (
	StaticFilePath  = "../static"
	ImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".svg", ".webp"}
)

func SaveImage(h *multipart.FileHeader) (fileName string, err error) {
	if h == nil {
		return DefaultAvatarFileName, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		err = errors.ErrInternal
		return
	}

	img, err := h.Open()
	if err != nil {
		return
	}

	fileExt := path.Ext(h.Filename)
	isImage := false
	for _, ext := range ImageExtensions {
		if fileExt == ext {
			isImage = true
			break
		}
	}

	if !isImage {
		err = errors.ErrInvalidData
		return
	}

	fileName = uuid.NewString() + fileExt
	filePath := path.Join(wd, StaticFilePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	io.Copy(file, img)

	return
}

func RemoveImage(fileName string) (err error) {
	if fileName == DefaultAvatarFileName {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		err = errors.ErrInternal
		return
	}

	filePath := path.Join(wd, StaticFilePath, fileName)

	if _, err = os.Stat(filePath); err == nil {
		err = os.Remove(filePath)
		if err != nil {
			return
		}
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return
	}

	return
}
