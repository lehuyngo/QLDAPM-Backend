package draft_contact

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ReadDraftContact(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.DraftContact.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.GetOrganization().GetID() != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	// Name card
	var nameCard *apis.File
	if data.NameCard != nil {
		nameCard = &apis.File {
			UUID: data.GetNameCard().GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/draft-contacts/%s/name-card/%s", services.Config.Domain.Domain, uuid, data.GetNameCard().GetUUID()),
		}
	}

	// Name Card
	var companyLogo *apis.File
	if data.CompanyLogo != nil {
		companyLogo = &apis.File {
			UUID: data.GetNameCard().GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/draft-contacts/%s/company-logo/%s", services.Config.Domain.Domain, uuid, data.GetNameCard().GetUUID()),
		}
	}

	resp := &apis.DraftContact{
		UUID: data.UUID,
		FullName: data.FullName,
		Phone: data.Phone,
		Email: data.Email,
		ClientName: data.GetClientName(),
		ClientWebsite: data.GetClientWebsite(),
		ClientAddress: data.GetClientAddress(),
		NameCard: nameCard,
		CompanyLogo: companyLogo,
	}

	c.JSON(http.StatusOK, resp)
}
