package mail

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ReadMail(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Mail.ReadByUUID(ctx, uuid)
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

	resp := &apis.Mail{
		UUID: data.UUID,
		Sender:	&apis.User{
			UUID: data.GetCreator().GetUUID(),
			DisplayName: data.GetCreator().GetDisplayName(),
		},
		SendTime: data.GetCreatedAt().UnixMilli(),
		Subject : data.Subject,
		Content : data.Content,
	}

	// Receivers
	for _, val := range data.Receivers {
		if val.Contact != nil {
			resp.ReceiverContacts = append(resp.ReceiverContacts, &apis.Contact{
				UUID: val.Contact.UUID,
				FullName: val.Contact.FullName,
			})
		}

		if val.User != nil {
			resp.ReceiverUsers = append(resp.ReceiverUsers, &apis.User{
				UUID: val.User.UUID,
				DisplayName: val.User.DisplayName,
			})
		}

		if val.MailAddress != "" {
			resp.ReceiverMailAddresses= append(resp.ReceiverMailAddresses, val.MailAddress)
		}
	}

	c.JSON(http.StatusOK, resp)
}
