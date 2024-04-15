package meeting_highlight

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListHighlight(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-highlight] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[list-highlight] query project info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[list-highlight] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("meeting_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	meetingNotes, err := h.MeetingNote.List(ctx, project.ID)
	if err != nil {
		log.For(c).Error("[list-highlight] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var meetingNoteIDs []int64
	for _, note := range meetingNotes {
		meetingNoteIDs = append(meetingNoteIDs, note.ID)
	}

	var highlights []*entities.MeetingHighlight

	highlights, err = h.Highlight.ListByMeetingNotes(ctx, meetingNoteIDs)
	if err != nil {
		log.For(c).Error("[list-highlight] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListHighlightResponse{}
	for _, val := range highlights {
		resp.Data = append(resp.Data, &apis.Highlight{
			UUID:            val.UUID,
			Title:           val.Title,
			MeetingNoteUUID: val.MeetingNote.UUID,
			CreatorUUID:     val.Creator.UUID,
			CreatedAt:       val.CreatedAt.UnixMilli(),
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreatedAt < resp.Data[j].CreatedAt
	})
	c.JSON(http.StatusOK, resp)
}
