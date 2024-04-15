package meeting_highlight

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateHighlight(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateMeetingHighlightRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-hightlight] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	meetingNoteUUID := c.Param("uuid")
	meetingNote, err := h.MeetingNote.ReadByUUID(ctx, meetingNoteUUID)
	if err != nil {
		log.For(c).Error("[create-hightlight] query project info failed", log.Field("user_id", userID), log.Field("project_uuid", meetingNoteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != meetingNote.GetProject().GetOrganizationID() {
		log.For(c).Error("[create-hightlight] organization id is not match", log.Field("user_id", userID), log.Field("meeting_note_uuid", meetingNoteUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("project_organization_id", meetingNote.Project.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	resp := &apis.CreateMeetingHighlightResponse{}
	data := []*entities.MeetingHighlight{}
	for _, title := range req.Title {
		existedHighlight, err := h.Highlight.First(ctx, title, meetingNote.ID)
		if err == nil {
			log.For(c).Info("[create-hightlight] create highlight is existed", log.Field("user_id", userID), log.Field("uuid", existedHighlight.UUID))
			resp.UUIDs = append(resp.UUIDs, existedHighlight.UUID)
			continue
		}

		data = append(data, &entities.MeetingHighlight{
			UUID:          uuid.NewString(),
			Title:         title,
			MeetingNoteID: meetingNote.ID,
			CreatedBy:     userID,
		})
	}

	if len(data) < 1 {
		log.For(c).Error("[create-hightlight] all highlight is exist", log.Field("user_id", userID),
			log.Field("meeting_note_uuid", meetingNoteUUID))
		c.JSON(http.StatusOK, resp)
		return
	}

	err = h.Highlight.CreateBatch(ctx, data)
	if err != nil {
		log.For(c).Error("[create-hightlight] create highlight failed", log.Field("user_id", userID),
			log.Field("meeting_note_uuid", meetingNoteUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range data {
		resp.UUIDs = append(resp.UUIDs, val.UUID)
	}

	log.For(c).Info("[create-hightlight] process success", log.Field("user_id", userID),
		log.Field("meeting_note_uuid", meetingNoteUUID), log.Field("resp", resp))
	c.JSON(http.StatusOK, resp)
}
