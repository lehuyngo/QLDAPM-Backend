package contact_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ListContactProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Contact.ReadByUUID(ctx, uuid)
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

	resp := &apis.ListProjectResponse{}

	// Projects
	projects, err := h.ContactProject.List(ctx, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range projects {
		if !val.GetProject().Available() {
			continue
		}
		
		resp.Data = append(resp.Data, &apis.Project{
			UUID: val.GetProject().GetUUID(),
			FullName: val.GetProject().GetFullName(),
			ShortName: val.GetProject().GetShortName(),
			Code: val.GetProject().GetCode(),
			CreatedTime: val.CreatedAt.UnixMilli(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
