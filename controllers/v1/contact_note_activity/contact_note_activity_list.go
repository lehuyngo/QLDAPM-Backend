package contact_note_activity

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListContactNoteActivity(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-contact-note-activity] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[list-contact-note-activity] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[list-contact-note-activity] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.ContactNoteActivity.List(ctx, contact.ID)
	if err != nil {
		log.For(c).Error("[list-contact-note-activity] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListActivityContactNote{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.ContactNoteActivity{
			Type: apis.ActivityType(val.Type),
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime: val.CreatedAt.UnixMilli(),
			Note: &apis.ContactNote{
				UUID:  val.GetNote().GetUUID(),
				Title: val.GetNote().GetTitle(),
			},
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreatedTime < resp.Data[j].CreatedTime
	})

	c.JSON(http.StatusOK, resp)
}
