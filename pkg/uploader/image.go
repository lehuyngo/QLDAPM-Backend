package uploader

import (
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Image struct {
	UUID string
	OriginalName string
	FileExt string
	RelativePathFile string
	Thumbnail string
}

type ImageUploader struct {
	folder string
}

func NewImageUploader(folder string) *ImageUploader {
	return &ImageUploader{
		folder: folder,
	}
}

func (u *ImageUploader) generateLocalFile(relativePathFile string) string {
	return u.folder + relativePathFile
}

func (u *ImageUploader) UploadImage(ctx *gin.Context, fieldImage string, thumbnailWidth int) (*Image, error) {
	file, header, err := ctx.Request.FormFile(fieldImage)
	if err != nil {
		return nil, err
	}

	result := &Image{
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

	_, err = io.Copy(out, file)
	if err != nil {
		out.Close()
		return nil, err
	}
	out.Close()

	// thumbnail
	input, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	imageFile, _, err := image.Decode(input)
	if err != nil {
		return nil, err
	}
	src := imaging.Resize(imageFile, thumbnailWidth, 0, imaging.Lanczos)

	result.Thumbnail = result.UUID + "-thumbnail" + result.FileExt
	thumbnailPath := u.generateLocalFile(result.Thumbnail)
	err = imaging.Save(src, thumbnailPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}
