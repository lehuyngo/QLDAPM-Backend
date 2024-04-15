package services

import (
	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/uploader"
)

var (
	imageUploader *uploader.ImageUploader
	fileUploader *uploader.FileUploader
)

func initUploaders() {
	imageUploader = uploader.NewImageUploader(Config.FileStorage.Folder)
	fileUploader = uploader.NewFileUploader(Config.FileStorage.Folder)
}

func UploadImage(ctx *gin.Context, fieldImage string, thumbnailWidth int) (*uploader.Image, error) {
	if imageUploader == nil {
		return nil, define.ErrNotInitial
	}
	
	return imageUploader.UploadImage(ctx, fieldImage, thumbnailWidth)
}

func UploadFile(ctx *gin.Context, fieldFile string) (*uploader.File, error) {
	if fileUploader == nil {
		return nil, define.ErrNotInitial
	}
	
	return fileUploader.UploadFile(ctx, fieldFile)
}