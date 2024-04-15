package project

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

func (h Handler) CreateProject(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateProjectRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var client *entities.Client
	var contact *entities.Contact

	// verify client and contact
	if req.IsClientValid() {
		if req.Client.UUID != "" {
			client, err = h.Client.ReadByUUID(ctx, req.Client.UUID)
			if err != nil {
				log.For(c).Error("[create-project] query client by uuid failed", log.Field("user_id", userID), log.Field("client_uuid", req.Client.UUID), log.Err(err))
				c.JSON(http.StatusBadRequest, err)
				return
			}
		} else {
			client = &entities.Client{
				Base: entities.Base{
					CreatedBy: userID,
					UpdatedBy: userID,
				},
				FullName: req.Client.FullName,
				ShortName: req.Client.ShortName,
				Code: req.Client.Code,
				Fax: req.Client.Fax,
				Website: req.Client.Website,
				Phone: req.Client.Phone,
				Email: req.Client.Email,
				CompanySize: req.Client.CompanySize,
				Address: req.Client.Address,
				Status: entities.Active,
				OrganizationID: user.OrganizationID,
				LastActiveTime: time.Now().UnixMilli(),
			}
			_, err = h.Client.Create(ctx, client)
			if err != nil {
				log.For(c).Error("[create-project] create client failed", log.Field("user_id", userID), log.Err(err))
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
	}

	if req.IsContactValid() {
		if req.Contact.UUID != "" {
			contact, err = h.Contact.ReadByUUID(ctx, req.Contact.UUID)
			if err != nil {
				log.For(c).Error("[create-project] query contact by uuid failed", log.Field("user_id", userID), log.Field("contact_uuid", req.Contact.UUID), log.Err(err))
				if client != nil {
					h.Client.Delete(ctx, user.ID, client.ID)
				}

				c.JSON(http.StatusBadRequest, err)
				return
			}
		} else {
			contact = &entities.Contact{
				Base: &entities.Base{
					CreatedBy: userID,
					UpdatedBy: userID,
				},
				FullName: req.Contact.FullName,
				ShortName: req.Contact.ShortName,
				JobTitle: req.Contact.JobTitle,
				Phone: req.Contact.Phone,
				Email: req.Contact.Email,
				Gender: entities.Gender(req.Contact.Gender),
				Status: entities.Active,
				OrganizationID: user.OrganizationID,
				LastActiveTime: time.Now().UnixMilli(),
				Birthday: req.Contact.BirthDay,
			}
			_, err = h.Contact.Create(ctx, contact)
			if err != nil {
				log.For(c).Error("[create-project] create contact failed", log.Field("user_id", userID), log.Err(err))
				if client != nil {
					h.Client.Delete(ctx,user.ID, client.ID)
				}

				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
	}
	
	data := &entities.Project{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName:		req.FullName,
		ShortName:		req.ShortName,
		Code:			req.Code,
		OrganizationID: user.OrganizationID,
		ProjectStatus: entities.ProjectStatus(req.ProjectStatus),
	}
	if client.GetID() > 0 {
		data.ClientID = client.GetID()
	}

	projectID, err := h.Project.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-project] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if contact.GetID() < 1 {
		log.For(c).Info("[create-project] process success", log.Field("user_id", userID), log.Field("uuid", data.UUID))
		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: data.UUID,
		})
		return
	}

	err = h.ContactProject.Add(ctx, &entities.ContactProject{
		ContactID: contact.ID,
		ProjectID: data.ID,
		CreatedBy: userID,
	})
	if err != nil {
		if client != nil {
			h.Client.Delete(ctx, user.ID, client.ID)
		}
		if contact != nil {
			h.Contact.Delete(ctx,user.ID, contact.ID)
		}
		h.Project.Delete(ctx, projectID)

		log.For(c).Error("[create-project] add project contact failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
	}

	log.For(c).Info("[create-project] process success", log.Field("user_id", userID), log.Field("uuid", data.UUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
