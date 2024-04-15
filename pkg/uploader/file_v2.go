package uploader

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileUploaderV2 struct {
	folder string
}

func NewFileUploaderV2(folder string) *FileUploaderV2 {
	return &FileUploaderV2{
		folder: folder,
	}
}

func (u *FileUploaderV2) generateName(filename string) string {
	return u.folder + filename
}

func (u *FileUploaderV2) UploadFile(ctx *gin.Context, fieldFile string) (string, string, error) {
	file, header, err := ctx.Request.FormFile(fieldFile)
	if err != nil {
		return "", "", err
	}

	fileExt := filepath.Ext(header.Filename)
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + uuid.NewString() + fileExt
	filePath := u.generateName(filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return "", "", err
	}

	return originalFileName + fileExt, filePath, nil
}
