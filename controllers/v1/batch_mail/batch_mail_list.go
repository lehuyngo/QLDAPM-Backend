package batch_mail

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListBatchMail(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Debug("[list-batch-mail] start process", log.Field("user_id", userID))

	data, err := h.BatchMail.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListBatchMailResponse{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.BatchMail{
			UUID: val.UUID,
			Sender:	&apis.User{
				UUID: val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			SendTime: val.GetCreatedAt().UnixMilli(),
			Subject : val.Subject,
			Content : val.Content,
			Status: int(val.Status),
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].SendTime < resp.Data[j].SendTime
	})

	log.For(c).Info("[list-batch-mail] process success", log.Field("user_id", userID), log.Field("resp_len", len(resp.Data)))
	c.JSON(http.StatusOK, resp)
}
