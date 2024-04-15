package contact_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ReadContactTag(c *gin.Context) {
	ctx := c.Request.Context()

	uuid := c.Param("uuid")

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[read-contact-tag] query user info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tag, err := h.ContactTag.ReadByUUID(ctx, user.OrganizationID, uuid)
	if err != nil {
		log.For(c).Error("[read-contact-tag] query contact note failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := apis.ContactTag{
		UUID: uuid,
		Name: tag.Name,
	}
	for _, val := range tag.Contacts {
		resp.Contacts = append(resp.Contacts, &apis.Contact{
			UUID: val.GetContact().GetUUID(),
			FullName: val.GetContact().GetFullName(),
			ShortName: val.GetContact().GetShortName(),
		})
	}

	log.For(c).Info("[read-contact-tag] process success", log.Field("user_id", userID), log.Field("uuid", uuid))
	c.JSON(http.StatusOK, resp)
}
