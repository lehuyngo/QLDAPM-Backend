package contact_client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteContactClient(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-contact-client] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-client] query client info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("client_uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-client] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if (user.OrganizationID != client.OrganizationID) || (user.OrganizationID != contact.OrganizationID) {
		log.For(c).Error("[delete-contact-client] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.ContactClient.Delete(ctx, user.ID, contact.ID, client.ID)
	if err != nil {
		log.For(c).Error("[delete-contact-client] udpate database failed", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("contact_uuid", contact.UUID),
			log.Field("client_id", client.ID), log.Field("contact_id", contact.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[delete-contact-client] process success", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("contact_uuid", contact.UUID))
	c.JSON(http.StatusOK, nil)
}
