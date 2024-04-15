package client_contact

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateClientContact(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientContactRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-contact] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[create-client-contact] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[create-client-contact] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	if req.NewContact != nil {
		err = h.ClientContact.Create(ctx, client.ID, &entities.Contact{
			Base: &entities.Base{
				CreatedBy: userID,
				UpdatedBy: userID,
			},
			FullName: req.NewContact.FullName,
			ShortName: req.NewContact.ShortName,
			JobTitle: req.NewContact.JobTitle,
			Phone: req.NewContact.Phone,
			Email: req.NewContact.Email,
			Gender: entities.Gender(req.NewContact.Gender),
			OrganizationID: user.OrganizationID,
			LastActiveTime: time.Now().UnixMilli(),
			Birthday : req.NewContact.LastActiveTime,
		})
		if err != nil {
			log.For(c).Error("[create-client-contact] create contact failed", log.Field("user_id", userID), 
				log.Field("client_uuid", clientUUID), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	if len(req.UUIDs) > 0 {
		contacts, err := h.Contact.ListByUUIDs(ctx, req.UUIDs)
		if err != nil {
			log.For(c).Error("[create-client-contact] query contact by uuid failed", log.Field("user_id", userID), 
				log.Field("uuids", req.UUIDs), log.Err(err))

			c.JSON(http.StatusInternalServerError, err)
			return
		}

		var data []entities.ClientContact
		for _, val := range contacts {
			if val.OrganizationID != user.OrganizationID {
				continue
			}
			
			data = append(data, entities.ClientContact{
				ClientID: client.ID,
				ContactID: val.ID,
				CreatedBy: userID,
			})
		}

		err = h.ClientContact.AddBatch(ctx, data)
		if err != nil {
			if !define.IsErrDuplicateKey(err) {
				log.For(c).Error("[create-client-contact] create bacth contact failed", log.Field("user_id", userID), 
				log.Field("client_uuid", clientUUID), log.Err(err))

				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
	}
	
	log.For(c).Info("[create-client-contact] process success", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
	c.JSON(http.StatusOK, nil)
}
