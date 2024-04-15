package client_activity

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListClientActivity(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-client-activity] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[list-client-activity] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[list-client-activity] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.ClientActivity.List(ctx, client.ID)
	if err != nil {
		log.For(c).Error("[list-client-activity] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListActivityClient{}
	for _, val := range data {
		if val.GetClient() == nil {
			log.For(c).Warn("[client-activity] client is nil", log.Field("user_id", userID))
			continue
		}
		resp.Data = append(resp.Data, &apis.ClientActivity{
			Type: apis.ActivityType(val.Type),
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime: val.CreatedAt.UnixMilli(),
			Client: &apis.Client{
				UUID:      val.GetClient().GetUUID(),
				FullName:  val.GetClient().GetFullName(),
				ShortName: val.GetClient().GetShortName(),
			},
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreatedTime < resp.Data[j].CreatedTime
	})

	c.JSON(http.StatusOK, resp)
}
