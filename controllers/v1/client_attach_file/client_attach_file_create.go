package client_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateClientAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[create-client-attach-file] query contact info", log.Field("user_id", userID),
			log.Field("client_id", client.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	file, err := services.UploadFile(c, "file")
	if err != nil {
		log.For(c).Error("[create-client-attach-file] write file failed", log.Field("user_id", userID),
			log.Field("client_id", client.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
	}

	data := &entities.ClientAttachFile{
		OriginalName:   file.OriginalName,
		RelativePath:   file.RelativePathFile,
		Ext:            file.FileExt,
		OrganizationID: user.OrganizationID,
		ClientID:       client.ID,
		CreatedBy:      userID,
		FileSize:       file.FileSize,
	}
	_, err = h.ClientAttachFile.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contact] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
