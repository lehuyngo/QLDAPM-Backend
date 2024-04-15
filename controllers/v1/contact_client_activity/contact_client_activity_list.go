package contact_client_activity

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListContactClientActivity(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-contact-client-activity] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[list-contact-client-activity] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[list-contact-client-activity] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.ContactClientActivity.List(ctx, contact.ID)
	if err != nil {
		log.For(c).Error("[list-contact-client-activity] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListActivityContactClient{}
	for _, val := range data {
		if val.GetClient() == nil {
			log.For(c).Error("[list-contact-client-activity] query client info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
			continue
		}
		resp.Data = append(resp.Data, &apis.ContactClientActivity{
			Type: apis.ActivityType(val.Type),
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			Client: &apis.Client{
				UUID:      val.GetClient().GetUUID(),
				FullName:  val.GetClient().GetFullName(),
				ShortName: val.GetClient().GetShortName(),
			},
			CreatedTime: val.CreatedAt.UnixMilli(),

		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreatedTime < resp.Data[j].CreatedTime
	})

	c.JSON(http.StatusOK, resp)
}
