package client_tag

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteClientTag(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.DeleteClientTagRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[delete-client-note] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	tagUUID := c.Param("tag_uuid")
	tag, err := h.ClientTag.ReadByUUID(ctx, user.OrganizationID, tagUUID)
	if err != nil {
		log.For(c).Error("[delete-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if req.IsFullyDeleted == 0 {
		err = h.ClientClientTag.Delete(ctx, client.ID, tag.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, nil)
		return
	}

	err = h.ClientTag.Delete(ctx, user, client.ID, tag.ID)
	if err != nil {
		log.For(c).Error("[delete-client-note] udpate database failed", log.Field("user_id", userID),
			log.Field("client_uuid", clientUUID), log.Field("tag_uuid", tagUUID),
			log.Field("client_id", client.ID), log.Field("tag_id", tag.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-client-note] process success", log.Field("user_id", userID),
		log.Field("client_uuid", clientUUID), log.Field("tag_uuid", tagUUID))
	c.JSON(http.StatusOK, nil)
}
