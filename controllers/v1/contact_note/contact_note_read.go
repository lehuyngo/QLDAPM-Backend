package contact_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ReadContactNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	contactUUID := c.Param("uuid")
	noteUUID := c.Param("note_uuid")

	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[read-contact-note] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	note, err := h.ContactNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[read-contact-note] query contact note failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if note.ContactID != contact.ID {
		log.For(c).Error("[read-contact-note] query contact id isnot match", log.Field("user_id", userID), 
			log.Field("contact_uuid", contactUUID), log.Field("note_uuid", noteUUID),
			log.Field("contact_id", contact.ID), log.Field("note_contact_id", note.ContactID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[read-contact-note] query user info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[read-contact-note] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	log.For(c).Info("[read-contact-note] process success", log.Field("user_id", userID), 
			log.Field("contact_uuid", contactUUID), log.Field("note_uuid", noteUUID))
			
	c.JSON(http.StatusOK, &apis.ContactNote{
		UUID : note.GetUUID(),
		Title: note.Title,
		Content: note.Content,
		Color: note.Color,
	})
}
