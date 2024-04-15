package contact_client

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

func (h Handler) CreateContactClient(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactClientRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contact-client] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[create-contact-client] query client info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[create-contact-client] query client id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	if req.NewClient!= nil {
		err = h.ContactClient.Create(ctx, contact.ID, &entities.Client{
			Base: entities.Base{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
			FullName: req.NewClient.FullName,
			ShortName: req.NewClient.ShortName,
			Code: req.NewClient.Code,
			Fax: req.NewClient.Fax,
			Website: req.NewClient.Website,
			Phone: req.NewClient.Phone,
			Email: req.NewClient.Email,
			CompanySize: req.NewClient.CompanySize,
			Address: req.NewClient.Address,
			OrganizationID: user.OrganizationID,
			LastActiveTime: time.Now().UnixMilli(),
		})
		if err != nil {
			log.For(c).Error("[create-contact-client] create contact failed", log.Field("user_id", userID), 
				log.Field("contact_uuid", contactUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	if len(req.UUIDs) > 0 {
		clients, err := h.Client.ListByUUIDs(ctx, req.UUIDs)
		if err != nil {
			log.For(c).Error("[create-contact-client] query contact by uuid failed", log.Field("user_id", userID), 
				log.Field("uuids", req.UUIDs), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}

		var data []entities.ClientContact
		for _, val := range clients {
			if val.OrganizationID != user.OrganizationID {
				continue
			}
			
			data = append(data, entities.ClientContact{
				ClientID: val.ID,
				ContactID: contact.ID,
				CreatedBy: userID,
			})
		}

		err = h.ContactClient.AddBatch(ctx, data)
		if err != nil {
			log.For(c).Error("[create-contact-client] create bacth contact failed", log.Field("user_id", userID), 
				log.Field("contact_uuid", contactUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}
	
	log.For(c).Info("[create-contact-client] process success", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
	c.JSON(http.StatusOK, nil)
}
