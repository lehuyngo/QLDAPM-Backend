package meeting_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteMeetingNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[delete-meeting-note] query project info failed", log.Field("user_id", userID), log.Field("uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can delete meeting note of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-meeting-note] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[delete-meeting-note] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("project_organizationID", project.OrganizationID))
		c.JSON(http.StatusForbidden, err)
		return
	}

	noteUUID := c.Param("note_uuid")
	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[delete-meeting-note] query meeting note info failed", log.Field("user_id", userID), log.Field("uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if note.GetProject().GetID() != project.GetID() {
		log.For(c).Error("[delete-meeting-note] project id not match", log.Field("user_id", userID),
			log.Field("meeting_note_project_id", note.GetProject().GetID()), log.Field("project_id", project.GetID()))
		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.MeetingNote.Delete(ctx, note.GetID(), project.GetID(), note.MeetingID)
	if err != nil {
		log.For(c).Error("[delete-meeting-note] delete meeting note failed", log.Field("user_id", userID),
			log.Field("meeting_note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
