package client_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateClientNote(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientNoteRequest{}

	clientUUID := c.Param("uuid")

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)
	
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[create-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-note] query user info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// check permission
	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[create-client-note] organization isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}
	
	data := &entities.ClientNote{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		ClientID: client.ID,
		Title: req.Title,
		Content: req.Content,
		Color: req.Color,
	}

	_, err = h.ClientNote.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-client-note] insert database failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[create-client-note] insert database failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
