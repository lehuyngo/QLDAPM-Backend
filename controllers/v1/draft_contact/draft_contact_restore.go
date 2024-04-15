package draft_contact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) RestoreDraftContact(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, _ := middlewares.ParseToken(c)

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	err = h.DraftContact.Restore(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
