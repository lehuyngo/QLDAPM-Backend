package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ReadClientLogo(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Client.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[read-client-logo] query client info", log.Field("uuid", uuid), 
			log.Field("user_id", userID), log.Field("client_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[read-client-logo] query user info", log.Field("uuid", uuid), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		log.For(c).Error("[read-client-logo] organization id not match", 
			log.Field("user_id", userID), log.Field("client_uuid", uuid), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("data_organization_id", data.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + data.GetLogo().GetRelativePath()
	fileName := data.GetLogo().GetOriginalName() + data.GetLogo().GetExt()
	log.For(c).Info("[read-client-logo] download", log.Field("uuid", uuid), 
			log.Field("user_id", userID), log.Field("client_uuid", uuid), 
			log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
