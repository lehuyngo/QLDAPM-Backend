package file

import (
	"io"
	"os"
)

func CopyFile(src, dst string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}

	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	return err
}

func GetFileSize(src string) (int64, error) {
	fileInfo, err := os.Stat(src)

	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}
