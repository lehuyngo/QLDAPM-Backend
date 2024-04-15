package client_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteClientAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-client-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[delete-client-attach-file] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	fileUUID := c.Param("file_uuid")
	file, err := h.ClientAttachFile.Read(ctx, fileUUID)
	if err != nil {
		log.For(c).Error("[delete-client-attach-file] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != file.OrganizationID {
		log.For(c).Error("[delete-client-attach-file] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.ClientAttachFile.Delete(ctx, client.ID, file.ID)
	if err != nil {
		log.For(c).Error("[delete-client-contact] udpate database failed", log.Field("user_id", userID), 
		log.Field("client_id", client.ID), log.Field("file_uuid", file.UUID),
		log.Field("file_id", file.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-client-contact] process success", log.Field("user_id", userID), 
		log.Field("client_id", client.ID), log.Field("file_uuid", file.UUID),
		log.Field("file_id", file.ID))
		
	c.JSON(http.StatusOK, nil)
}
