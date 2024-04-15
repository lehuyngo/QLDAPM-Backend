package static_file

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/normalize"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	file, err := services.UploadFile(c, "file")
	if err != nil {
		log.For(c).Error("[create-static-file] write file failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
	}

	data := &entities.StaticFile{
		OriginalName:   file.OriginalName,
		RelativePath:   file.RelativePathFile,
		Ext:            file.FileExt,
		OrganizationID: user.OrganizationID,
		CreatedBy:      userID,
		FileSize:       file.FileSize,
	}
	_, err = h.StaticFile.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-static-file] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	fileName := normalize.URLEncode(data.GetOriginalName()) + data.Ext
	url := fmt.Sprintf("%sapi/v1/static-files/%s/%s", services.Config.Domain.Domain, data.GetUUID(), fileName)

	log.For(c).Info("[create-static-file] process success", log.Field("user_id", userID), log.Field("uuid", data.GetUUID()))
	c.JSON(http.StatusOK, &apis.File{
		UUID: data.UUID,
		URL: url,
	})
}
