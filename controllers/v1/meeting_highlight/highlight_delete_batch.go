package meeting_highlight

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteHighlightBatch(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	req := &apis.DeleteMeetingHighlightRequest{}
	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		log.For(c).Error("[batch-delete-highlight] binding json failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	meetingNotetUUID := c.Param("uuid")
	meetingNote, err := h.MeetingNote.ReadByUUID(ctx, meetingNotetUUID)
	if err != nil {
		log.For(c).Error("[batch-delete-highlight] query meeting_note info failed", log.Field("meeting_note", meetingNotetUUID), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[batch-delete-highlight] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != meetingNote.GetProject().OrganizationID {
		log.For(c).Error("[batch-delete-highlight] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("meeting_note_organizationID", meetingNote.GetProject().OrganizationID))
		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.Highlight.DeleteBatch(ctx, req.UUIDs, meetingNote.ID)
	if err != nil {
		log.For(c).Error("[batch-delete-highlight] batch delete highlight failed", log.Field("user_id", userID),
			log.Field("highlight_uuids", req.UUIDs), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
