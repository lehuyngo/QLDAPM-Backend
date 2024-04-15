package client_contact

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListClientContact(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Client.ReadByUUID(ctx, uuid)
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

	resp := &apis.ListContactResponse{}

	// Contacts
	contacts, err := h.ClientContact.List(ctx, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range contacts {
		if !val.GetContact().Available() {
			continue
		}

		var avatar *apis.File
		if val.GetContact().GetAvatar() != nil {
			avatar = &apis.File{
				UUID: val.GetContact().GetAvatar().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/contacts/%s/avatar/%s", services.Config.Domain.Domain, val.GetContact().GetUUID(), val.GetContact().GetAvatar().GetUUID()),
			}
		}

		resp.Data = append(resp.Data, &apis.Contact{
			UUID: val.GetContact().GetUUID(),
			FullName: val.GetContact().GetFullName(),
			ShortName: val.GetContact().GetShortName(),
			Phone: val.GetContact().GetPhone(),
			Email: val.GetContact().GetEmail(),
			JobTitle: val.GetContact().GetJobTitle(),
			Gender: val.GetContact().GetGender().Value(),
			BirthDay: val.GetContact().GetBitrhDay(),
			CreatedTime: val.GetContact().GetCreatedAt().UnixMilli(),
			Avatar: avatar,
		})
	}

	c.JSON(http.StatusOK, resp)
}
