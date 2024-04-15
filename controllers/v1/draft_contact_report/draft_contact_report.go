package draft_contact_report

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) GetDraftContactReport(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[draft-contact-report] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.For(c).Debug("[draft-contact-report] start process", log.Field("user_id", userID))

	data, err := h.DraftContactReport.GetReport(ctx, user.OrganizationID)
	if err != nil {
		log.For(c).Error("[draft-contact-report] query report failed", log.Field("user_id", userID), log.Field("organization_id", user.OrganizationID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.DraftContactReport{
		Total: data.Total,
		Resolved: data.Resolved,
		Processing: data.Processing,
		Deleted: data.Deleted,
	})
}
