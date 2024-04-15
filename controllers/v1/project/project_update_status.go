package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)


func (h Handler) UpdateProjectStatus(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateProjectStatusRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Project.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[update-project] query project info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit project of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[update-project] query user info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		log.For(c).Error("[update-project] organization id not match", log.Field("user_id", userID), 
			log.Field("user_organizationID", user.OrganizationID), log.Field("data_organizationID", data.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data.UpdatedBy = userID
	data.ProjectStatus = entities.ProjectStatus(req.ProjectStatus)
	err = h.Project.Update(ctx, data)
	if err != nil {
		log.For(c).Error("[update-project] update database failed", log.Field("user_id", userID), log.Field("contact_id", data.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
