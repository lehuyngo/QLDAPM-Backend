package contact_note

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListContactNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[list-contact-note] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[list-client-note] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.ContactNote.ListByContactID(ctx, contact.ID)
	if err != nil {
		log.For(c).Error("[list-contact-note] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListContactNoteResponse{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.ContactNote{
			UUID: val.UUID,
			Title: val.Title,
			Content: val.Content,
			Color: val.Color,
			CreateTime: val.CreatedAt.UnixMilli(),
			Creator: &apis.User{
				UUID: val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreateTime < resp.Data[j].CreateTime
	})

	c.JSON(http.StatusOK, resp)
}
