package client

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

func (h Handler) ListClient(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Debug("[list-client] start process", log.Field("user_id", userID))

	data, err := h.Client.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListClientResponse{}
	for _, val := range data {
		var logo *apis.File
		if val.Logo != nil {
			logo = &apis.File {
				UUID: val.GetLogo().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/clients/%s/logo/%s", services.Config.Domain.Domain, val.UUID, val.GetLogo().GetUUID()),
			}
		}

		client := &apis.Client{
			UUID: val.UUID,
			FullName: val.FullName,
			ShortName: val.ShortName,
			Code: val.Code,
			Fax: val.Fax,
			Website: val.Website,
			Phone: val.Phone,
			Email: val.Email,
			CompanySize: val.CompanySize,
			LastActiveTime: val.LastActiveTime,
			Logo: logo,
		}
		for _, tag := range val.Tags {
			client.Tags = append(client.Tags, &apis.ClientTag{
				UUID: tag.GetTag().GetUUID(),
				Name: tag.GetTag().GetName(),
				Color: tag.GetTag().GetColor(),
			})
		}
		for _, project := range val.Projects {
			client.Projects = append(client.Projects, &apis.Project{
				UUID: project.GetUUID(),
				FullName : project.GetFullName(),
				LastActiveTime: project.GetLastActiveTime(),
			})
		}
		for _, contact := range val.Contacts {
			if contact.GetContact() == nil {
				log.For(c).Warn("[list-client] contact is nil", log.Field("user_id", userID))
				continue
			}

			client.Contacts = append(client.Contacts, &apis.Contact{
				FullName : contact.GetContact().GetFullName(),
				ShortName : contact.GetContact().GetShortName(),
				UUID: contact.GetContact().GetUUID(),
			})
		}
		resp.Data = append(resp.Data, client)
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].FullName, resp.Data[j].FullName)
	})

	log.For(c).Info("[list-client] process success", log.Field("user_id", userID), log.Field("resp_len", len(resp.Data)))
	c.JSON(http.StatusOK, resp)
}
