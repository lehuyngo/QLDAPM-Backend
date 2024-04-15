package client_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ReadClientTag(c *gin.Context) {
	ctx := c.Request.Context()

	uuid := c.Param("uuid")

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[read-client-tag] query user info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tag, err := h.ClientTag.ReadByUUID(ctx, user.OrganizationID, uuid)
	if err != nil {
		log.For(c).Error("[read-client-tag] query client note failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := apis.ClientTag{
		UUID: uuid,
		Name: tag.Name,
	}
	for _, val := range tag.Clients {
		resp.Clients = append(resp.Clients, &apis.Client{
			UUID: val.GetClient().GetUUID(),
			FullName: val.GetClient().GetFullName(),
			ShortName: val.GetClient().GetShortName(),
		})
	}

	log.For(c).Info("[read-client-tag] process success", log.Field("user_id", userID), log.Field("uuid", uuid))
	c.JSON(http.StatusOK, resp)
}
