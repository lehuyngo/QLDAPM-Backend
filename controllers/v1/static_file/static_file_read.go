package static_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) DownloadFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[static-file-download] query user info", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	uuid := c.Param("uuid")
	file, err := h.StaticFile.Read(ctx, uuid)
	if err != nil {
		log.For(c).Error("[static-file-download] query file info failed",
			log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != file.OrganizationID {
		log.For(c).Error("[static-file-download] file id isnot match", log.Field("user_id", userID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("file_organization_id", file.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + file.GetRelativePath()
	fileName := file.GetOriginalName() + file.GetExt()
	log.For(c).Info("[static-file-download] download",
			log.Field("user_id", userID), log.Field("uuid", file.UUID),
			log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
