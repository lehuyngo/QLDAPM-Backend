package contact_client

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListContactClient(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Contact.ReadByUUID(ctx, uuid)
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

	resp := &apis.ListClientResponse{}

	// Clients
	clients, err := h.ContactClient.List(ctx, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range clients {
		if !val.GetClient().Available() {
			continue
		}

		var logo *apis.File
		if val.GetClient().GetLogo() != nil {
			logo = &apis.File{
				UUID: val.GetClient().GetLogo().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/clients/%s/logo/%s", services.Config.Domain.Domain, val.GetClient().GetUUID(), val.GetClient().GetLogo().GetUUID()),
			}
		}

		resp.Data = append(resp.Data, &apis.Client{
			UUID: val.GetClient().GetUUID(),
			FullName: val.GetClient().GetFullName(),
			ShortName: val.GetClient().GetShortName(),
			Code: val.GetClient().GetCode(),
			Fax: val.GetClient().GetFax(),
			Website: val.GetClient().GetWebsite(),
			Phone: val.GetClient().GetPhone(),
			Email: val.GetClient().GetEmail(),
			Address: val.GetClient().GetAddress(),
			CreatedTime: val.GetClient().GetCreatedAt().UnixMilli(),
			Logo: logo,
		})
	}

	c.JSON(http.StatusOK, resp)
}
