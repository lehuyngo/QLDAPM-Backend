package meeting_note

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

func (h Handler) UpdateMeetingNote(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateMeetingNoteRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[update-meeting-note] query project info failed", log.Field("user_id", userID), log.Field("uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit meeting note of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[update-meeting-note] query user info failed", log.Field("user_id", userID), log.Field("userID", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[update-meeting-note] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("project_organizationID", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	noteUUID := c.Param("note_uuid")
	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[update-meeting-note] query meeting note info failed", log.Field("user_id", userID), log.Field("uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if note.GetProject().GetID() != project.GetID() {
		log.For(c).Error("[update-meeting-note] project id not match", log.Field("user_id", userID),
			log.Field("meeting_note_project_id", note.GetProject().GetID()), log.Field("project_id", project.GetID()))
		c.JSON(http.StatusForbidden, err)
		return
	}

	// if user is not exist in meeting note editors, then add user to editors
	var isExist bool = false
	for _, editor := range note.GetEditors() {
		if editor.GetEditor().GetID() == userID {
			isExist = true
			break
		}
	}
	if !isExist {
		_, err = h.MeetingNoteEditor.Create(ctx, &entities.MeetingNoteEditor{
			UUID:          uuid.New().String(),
			MeetingNoteID: note.ID,
			EditorID:      userID,
			CreatedBy:     userID,
		})

		if err != nil {
			log.For(c).Error("[update-meeting-note] create editor failed", log.Field("user_id", userID), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	note.UpdatedBy = userID
	note.Note = req.Note
	note.Link = req.Link
	note.StartTime = req.StartTime
	note.Location = req.Location

	err = h.MeetingNote.Update(ctx, note)
	if err != nil {
		log.For(c).Error("[update-meeting-note] update meeting note failed", log.Field("user_id", userID), log.Field("meeting_note_id", note.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
