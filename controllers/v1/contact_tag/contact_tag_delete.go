package contact_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteContactTag(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-contact-note] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-note] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[delete-contact-note] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	tagUUID := c.Param("tag_uuid")
	tag, err := h.ContactTag.ReadByUUID(ctx, user.OrganizationID, tagUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-note] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = h.ContactTag.Delete(ctx, user, contact.ID, tag.ID)
	if err != nil {
		log.For(c).Error("[delete-contact-note] udpate database failed", log.Field("user_id", userID),
			log.Field("contact_uuid", contactUUID), log.Field("tag_uuid", tagUUID),
			log.Field("contact_id", contact.ID), log.Field("tag_id", tag.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-contact-note] process success", log.Field("user_id", userID),
		log.Field("contact_uuid", contactUUID), log.Field("tag_uuid", tagUUID))
	c.JSON(http.StatusOK, nil)
}
