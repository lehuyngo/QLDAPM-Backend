package client_tag

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListClientTag(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.ClientTag.List(ctx, user.OrganizationID)
	if err != nil {
		log.For(c).Error("[list-client-note] query client info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListClientTagRequest{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.ClientTag{
			UUID: val.UUID,
			Name: val.Name,
			Color: val.GetColor(),
			CreatedTime: val.GetCreatedAt().UnixMilli(),
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].Name, resp.Data[j].Name)
	})

	c.JSON(http.StatusOK, resp)
}
