package contact_project

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

func (h Handler) CreateContactProject(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactProjectRequest{}
	
	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		log.For(c).Error("[create-contact-project] query user info failed", log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contact-project] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[create-contact-project] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[create-contact-project] query contact id isnot match", log.Field("user_id", userID), log.Field("client_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	if req.NewProject != nil {
		err = h.ContactProject.Create(ctx, contact.ID, &entities.Project{
			Base: entities.Base{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
			FullName: req.NewProject.FullName,
			ShortName: req.NewProject.ShortName,
			Code: req.NewProject.Code,
			ProjectStatus: entities.Prospect,
			OrganizationID: user.OrganizationID,
			LastActiveTime: time.Now().UnixMilli(),
		})
		if err != nil {
			log.For(c).Error("[create-contact-project] create project failed", log.Field("user_id", userID), 
				log.Field("contact_uuid", contactUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	if len(req.UUIDs) > 0 {
		projects, err := h.Project.ListByUUIDs(ctx, req.UUIDs)
		if err != nil {
			log.For(c).Error("[create-contact-project] query project by uuid failed", log.Field("user_id", userID), 
				log.Field("uuids", req.UUIDs), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}

		var data []entities.ContactProject
		for _, val := range projects {
			if val.OrganizationID != user.OrganizationID {
				continue
			}
			
			data = append(data, entities.ContactProject{
				ContactID: contact.ID,
				ProjectID: val.ID,
				CreatedBy: userID,
			})
		}

		err = h.ContactProject.AddBatch(ctx, data)
		if err != nil {
			log.For(c).Error("[create-contact-project] create bacth project failed", log.Field("user_id", userID), 
				log.Field("contact_uuid", contactUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}
	
	log.For(c).Info("[create-contact-project] process success", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
	c.JSON(http.StatusOK, nil)
}
