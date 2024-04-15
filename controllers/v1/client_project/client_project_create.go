package client_project

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateClientProject(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientProjectRequest{}
	
	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		log.For(c).Error("[create-client-project] query user info failed", log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-project] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[create-client-project] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[create-client-project] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	if req.NewProject != nil {
		_, err = h.Project.Create(ctx, &entities.Project{
			Base: entities.Base{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
			ClientID: client.ID,
			FullName: req.NewProject.FullName,
			ShortName: req.NewProject.ShortName,
			Code: req.NewProject.Code,
			ProjectStatus: entities.Prospect,
			OrganizationID: user.OrganizationID,
			LastActiveTime: time.Now().UnixMilli(),
		})
		if err != nil {
			log.For(c).Error("[create-client-project] create project failed", log.Field("user_id", userID), 
				log.Field("client_uuid", clientUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	if len(req.UUIDs) > 0 {
		projects, err := h.Project.ListByUUIDs(ctx, req.UUIDs)
		if err != nil {
			log.For(c).Error("[create-client-project] query project by uuid failed", log.Field("user_id", userID), 
				log.Field("uuids", req.UUIDs), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}

		for _, val := range projects {
			if val.OrganizationID != user.OrganizationID {
				continue
			}
			
			val.ClientID = client.ID
			err = h.Project.Update(ctx, val)
			if err != nil {
				log.For(c).Error("[create-client-project] update client for project failed", log.Field("user_id", userID), 
					log.Field("project_id", val.ID), log.Err(err))
	
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
	}
	
	log.For(c).Info("[create-client-project] process success", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
	c.JSON(http.StatusOK, nil)
}
