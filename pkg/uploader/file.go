package uploader

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	file_utils "gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/file"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	UUID             string
	OriginalName     string
	FileExt          string
	RelativePathFile string
	Thumbnail        string
	FileSize         int64
}

type FileUploader struct {
	folder string
}

func NewFileUploader(folder string) *FileUploader {
	return &FileUploader{
		folder: folder,
	}
}

func (u *FileUploader) generateLocalFile(relativePathFile string) string {
	return u.folder + relativePathFile
}

func (u *FileUploader) UploadFile(ctx *gin.Context, fieldFile string) (*File, error) {
	file, header, err := ctx.Request.FormFile(fieldFile)
	if err != nil {
		return nil, err
	}

	result := &File{
		UUID: uuid.NewString(),
	}

	result.FileExt = filepath.Ext(header.Filename)
	result.OriginalName = strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	result.RelativePathFile = result.UUID + result.FileExt
	filePath := u.generateLocalFile(result.RelativePathFile)
	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}

	result.FileSize, err = file_utils.GetFileSize(filePath)
	if err != nil {
		return nil, err
	}

	return result, nil
}
