package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) DeleteProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Project.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit project of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	// check allow same name
	err = h.Project.Delete(ctx, data.GetID())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}