package client_note

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListClientNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	clientUUID := c.Param("uuid")
	client, err := h.Client.ReadByUUID(ctx, clientUUID)
	if err != nil {
		log.For(c).Error("[list-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != client.OrganizationID {
		log.For(c).Error("[list-client-note] query client id isnot match", log.Field("user_id", userID), log.Field("client_uuid", clientUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("client_organization_id", client.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.ClientNote.ListByClientID(ctx, client.ID)
	if err != nil {
		log.For(c).Error("[list-client-note] query client info failed", log.Field("user_id", userID), log.Field("client_uuid", clientUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListClientNoteResponse{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.ClientNote{
			UUID: val.UUID,
			Title: val.Title,
			Content: val.Content,
			Color: val.Color,
			CreateTime: val.CreatedAt.UnixMilli(),
			Creator: &apis.User{
				UUID: val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreateTime < resp.Data[j].CreateTime
	})

	c.JSON(http.StatusOK, resp)
}
