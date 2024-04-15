package client_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) DownloadClientAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[client-attach-file-download] query contact info", log.Field("client_uuid", clientUUID), 
			log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[client-attach-file-download] query user info", log.Field("client_uuid", clientUUID), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[client-attach-file-download] organization id not match", 
			log.Field("user_id", userID), log.Field("client_uuid", clientUUID), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	fileUUID := c.Param("file_uuid")
	file, err := h.ClientAttachFile.Read(ctx, fileUUID)
	if err != nil {
		log.For(c).Error("[client-attach-file-download] query file info failed", log.Field("client_uuid", clientUUID), 
			log.Field("user_id", userID), log.Field("file_uuid", fileUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if file.ClientID != client.ID {
		log.For(c).Error("[client-attach-file-download] client id not match", 
			log.Field("user_id", userID), log.Field("client_id", client.ID), log.Field("file_id", file.ID),
			log.Field("client_id", file.ClientID), log.Field("file_client_id", file.ClientID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + file.GetRelativePath()
	fileName := file.GetOriginalName() + file.GetExt()
	log.For(c).Info("[client-attach-file-download] download",
			log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Field("file_uuid", file.UUID),
			log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
