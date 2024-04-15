package client

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ReadClient(c *gin.Context) {
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

	// Logo
	var logo *apis.File
	if data.Logo != nil {
		logo = &apis.File {
			UUID: data.GetLogo().GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/clients/%s/logo/%s", services.Config.Domain.Domain, uuid, data.GetLogo().GetUUID()),
		}
	}

	resp := &apis.Client{
		UUID: data.UUID,
		FullName: data.FullName,
		ShortName: data.ShortName,
		Code: data.Code,
		Fax: data.Fax,
		Website: data.Website,
		Phone: data.Phone,
		Email: data.Email,
		CompanySize: data.CompanySize,
		Address: data.Address,
		Logo: logo,
		LastActiveTime: data.LastActiveTime,
		CreatedTime: data.GetCreatedAt().UnixMilli(),
	}

	// Notes
	for _, val := range data.Notes {
		resp.Notes = append(resp.Notes, &apis.ClientNote{
			UUID: val.GetUUID(),
			Title: val.Title,
			Content: val.Content,
			Color: val.Color,
		})
	}

	// Tags
	for _, val := range data.Tags {
		resp.Tags = append(resp.Tags, &apis.ClientTag{
			UUID: val.GetTag().GetUUID(),
			Name: val.GetTag().GetName(),
			CreatedTime: val.GetTag().GetCreatedAt().UnixMilli(),
			Color: val.GetTag().GetColor(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
