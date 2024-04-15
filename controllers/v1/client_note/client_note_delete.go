package client_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteClientNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	clientUUID := c.Param("uuid")
	noteUUID := c.Param("note_uuid")

	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	note, err := h.ClientNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query client note failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if note.ClientID != client.ID {
		log.For(c).Error("[delete-client-note] query client id isnot match", log.Field("user_id", userID),
			log.Field("client_uuid", clientUUID), log.Field("note_uuid", noteUUID),
			log.Field("client_id", client.ID), log.Field("note_client_id", note.ClientID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query user info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[delete-client-note] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.ClientNote.Delete(ctx, user, client.ID, note.GetID())
	if err != nil {
		log.For(c).Error("[delete-client-note] udpate database failed", log.Field("user_id", userID),
			log.Field("client_uuid", clientUUID), log.Field("note_uuid", noteUUID),
			log.Field("client_id", client.ID), log.Field("note_id", note.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-client-note] process success", log.Field("user_id", userID),
		log.Field("client_uuid", clientUUID), log.Field("note_uuid", noteUUID))
	c.JSON(http.StatusOK, nil)
}
