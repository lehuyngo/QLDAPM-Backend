package client_tag

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateClientTag(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientTagRequest{}

	clientUUID := c.Param("uuid")

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	req.Name = strings.ToLower(strings.TrimSpace(req.Name))
	req.Color = strings.ToLower(strings.TrimSpace(req.Color))

	userID, _, _ := middlewares.ParseToken(c)
	
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[create-client-tag] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-client-tag] query user info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// check permission
	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[create-client-tag] organization isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	tag, err := h.ClientTag.ReadByName(ctx, user.OrganizationID, req.Name)
	if err != nil {
		tag = &entities.ClientTag{
			Name: 			req.Name,
			OrganizationID: user.OrganizationID,
			Color: 			req.Color,
			CreatedBy:		userID,
		}
		tag.Clients = make([]*entities.ClientClientTag, 0)
		tag.Clients = append(tag.Clients, &entities.ClientClientTag{
			ClientID: client.ID,
			Color: req.Color,
			CreatedBy: userID,
		})
		tag.CreatedBy = user.ID

		_, err = h.ClientTag.Create(ctx, tag)
		if err != nil {
			log.For(c).Error("[create-client-tag] insert database failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		log.For(c).Info("[create-client-tag] insert database failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: tag.UUID,
		})
		return
	}
	
	err = h.ClientTag.Add(ctx, &entities.ClientClientTag{
		ClientID: client.ID,
		TagID: tag.ID,
		Color: req.Color,
		CreatedBy: userID,
	})
	if define.IsErrDuplicateKey(err) {
		tag.Color = req.Color
		h.ClientTag.Update(ctx, tag)
		
		log.For(c).Info("[create-client-tag] tag added", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: tag.UUID,
		})
		return
	}
	if err != nil {
		log.For(c).Error("[create-client-tag] insert database failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tag.Color = req.Color
	h.ClientTag.Update(ctx, tag)

	log.For(c).Info("[create-client-tag] add tag success", log.Field("user_id", userID), log.Field("client_uuid", clientUUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: tag.UUID,
	})
}
