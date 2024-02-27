package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

func SaveImage(h *multipart.FileHeader) (fileName string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	img, err := h.Open()
	if err != nil {
		return
	}

	fileName = RandStringRunes(32) + path.Ext(h.Filename)
	filePath := path.Join(wd, "static", fileName)

	fmt.Println(fileName, filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	io.Copy(file, img)

	return
}
