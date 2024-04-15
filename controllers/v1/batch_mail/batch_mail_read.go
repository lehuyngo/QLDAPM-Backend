package batch_mail

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ReadBatchMail(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.BatchMail.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.GetOrganization().GetID() != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	resp := &apis.BatchMail{
		UUID: data.UUID,
		Sender:	&apis.User{
			UUID: data.GetCreator().GetUUID(),
			DisplayName: data.GetCreator().GetDisplayName(),
		},
		SendTime: data.GetCreatedAt().UnixMilli(),
		Subject : data.Subject,
		Content : data.Content,
		Status: int(data.Status),
	}

	// Receivers
	for _, val := range data.Receivers {
		resp.Receivers = append(resp.Receivers, &apis.BatchMailReceiver{
			UUID: val.UUID,
			Contact: &apis.Contact{
				UUID: val.GetContact().GetUUID(),
				FullName: val.GetContact().GetFullName(),
				Email: val.Email,
			},
			SendTime: val.SendTime,
			Status: int(val.Status),
		})
	}

	c.JSON(http.StatusOK, resp)
}
