package static

import (
	"io"
	"mime/multipart"
	"os"
)

const DefaultAvatarFileName = "default_avatar.png"
const DefaultGroupAvatarFileName = "default_group_avatar.png"

func SaveFile(fh *multipart.FileHeader, dst string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func RemoveFile(path string) error {
	return os.Remove(path)
}
