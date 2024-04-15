package client_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ListClientProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, uuid)
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

	if user.GetOrganization().GetID() != client.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	resp := &apis.ListProjectResponse{}

	// Projects
	projects, err := h.Project.ListByClientID(ctx, client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range projects {
		if !val.Available() {
			continue
		}
		
		resp.Data = append(resp.Data, &apis.Project{
			UUID: val.GetUUID(),
			FullName: val.GetFullName(),
			ShortName: val.GetShortName(),
			Code: val.GetCode(),
			CreatedTime: val.CreatedAt.UnixMilli(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
