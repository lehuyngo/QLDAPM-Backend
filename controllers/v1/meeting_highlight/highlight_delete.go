package meeting_highlight

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteHighlight(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	hightlightUUID := c.Param("highlight_uuid")
	highlight, err := h.Highlight.ReadByUUID(ctx, hightlightUUID)
	if err != nil {
		log.For(c).Error("[delete-highlight] query project info failed", log.Field("user_id", userID), log.Field("uuid", hightlightUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-highlight] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != highlight.GetMeetingNote().GetProject().OrganizationID {
		log.For(c).Error("[delete-highlight] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("meeting_organizationID", highlight.MeetingNote.Project.OrganizationID))
		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.Highlight.Delete(ctx, highlight.ID)
	if err != nil {
		log.For(c).Error("[delete-highlight] delete highlight failed", log.Field("user_id", userID),
			log.Field("highlight_uuid", hightlightUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
