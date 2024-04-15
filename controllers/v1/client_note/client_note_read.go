package client_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ReadClientNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	clientUUID := c.Param("uuid")
	noteUUID := c.Param("note_uuid")

	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[read-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	note, err := h.ClientNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[read-client-note] query client note failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if note.ClientID != client.ID {
		log.For(c).Error("[read-client-note] query client id isnot match", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("note_uuid", noteUUID),
			log.Field("client_id", client.ID), log.Field("note_client_id", note.ClientID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[read-client-note] query user info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[read-client-note] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	log.For(c).Info("[read-client-note] process success", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("note_uuid", noteUUID))
			
	c.JSON(http.StatusOK, &apis.ClientNote{
		UUID : note.GetUUID(),
		Title: note.Title,
		Content: note.Content,
		Color: note.Color,
	})
}
