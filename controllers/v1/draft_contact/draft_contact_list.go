package draft_contact

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListDraftContact(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-draft-contact] query user info failed", log.Field("user_id", userID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[list-draft-contact] start process", log.Field("user_id", userID))

	data, err := h.DraftContact.List(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListDraftContactResponse{}
	for _, val := range data {

		var nameCard *apis.File
		if val.NameCard != nil {
			nameCard = &apis.File {
				UUID: val.GetNameCard().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/draft-contacts/%s/name-card/%s", services.Config.Domain.Domain, val.UUID, val.GetNameCard().GetUUID()),
			}
		}

		var companyLogo *apis.File
		if val.CompanyLogo != nil {
			companyLogo = &apis.File {
				UUID: val.GetCompanyLogo().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/draft-contacts/%s/company-logo/%s", services.Config.Domain.Domain, val.UUID, val.GetCompanyLogo().GetUUID()),
			}
		}

		resp.Data = append(resp.Data, &apis.DraftContact{
			UUID: val.UUID,
			FullName: val.GetFullName(),
			Phone: val.GetPhone(),
			Email: val.GetEmail(),
			ClientName: val.GetClientName(),
			ClientWebsite: val.GetClientWebsite(),
			ClientAddress: val.GetClientAddress(),
			NameCard: nameCard,
			CompanyLogo: companyLogo,
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].FullName, resp.Data[j].FullName)
	})

	log.For(c).Info("[list-draft-contact] process success", log.Field("user_id", userID), log.Field("resp_len", len(resp.Data)))
	c.JSON(http.StatusOK, resp)
}
