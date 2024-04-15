package receive_mail_history

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListReceiveMailHistory(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-receive-mail-history] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.TrackedURL.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		log.For(c).Error("[list-receive-mail-history] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListReceiveMailHistoryResponse{}
	for _, val := range data {
		el := &apis.ReceiveMailHistory{
			UUID: val.UUID,
			Code: val.Code,
			URL: val.URL,
			Creator:  	&apis.User{
				UUID: val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime: val.CreatedAt.UnixMilli(),
			ReceivedTime: val.ReceivedAt,
			Status: val.Status.Value(),
		}
		resp.Data = append(resp.Data, el)
	}

	c.JSON(http.StatusOK, resp)
}
