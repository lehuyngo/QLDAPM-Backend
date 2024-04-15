package user

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListUser(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-user] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.User.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		log.For(c).Error("[list-user] query users by organization failed", log.Field("user_id", userID), log.Field("organization_id", user.OrganizationID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListSelfProfile{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.SelfProfile{
			UUID:        val.UUID,
			DisplayName: val.DisplayName,
			Email:       val.Email,
		})
	}

	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].DisplayName, resp.Data[j].DisplayName)
	})

	c.JSON(http.StatusOK, resp)
}
