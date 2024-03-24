package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"socio/errors"

	"github.com/google/uuid"
)

const DefaultAvatarFileName = "default_avatar.png"

var StaticFilePath = "../static"

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

	fileName = uuid.NewString() + path.Ext(h.Filename)
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
