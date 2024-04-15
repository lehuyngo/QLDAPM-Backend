package contact_tag

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListContactTag(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.ContactTag.List(ctx, user.OrganizationID)
	if err != nil {
		log.For(c).Error("[list-contact-tag] query contact info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListContactTagRequest{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.ContactTag{
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
