package client_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteClientProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-client-project] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[delete-client-project] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("project_uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[delete-client-project] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if (user.OrganizationID != client.OrganizationID) || (user.OrganizationID != project.OrganizationID) {
		log.For(c).Error("[delete-client-project] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	project.ClientID = 0
	err = h.Project.Update(ctx, project)
	if err != nil {
		log.For(c).Error("[delete-client-project] udpate database failed", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("project_uuid", project.UUID),
			log.Field("client_id", client.ID), log.Field("project_id", project.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-client-project] process success", log.Field("user_id", userID), 
			log.Field("client_uuid", clientUUID), log.Field("project_uuid", project.UUID))
	c.JSON(http.StatusOK, nil)
}
