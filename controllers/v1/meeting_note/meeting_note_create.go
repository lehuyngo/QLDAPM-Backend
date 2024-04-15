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

func (h Handler) CreateMeetingNote(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateMeetingNoteRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-meeting-note] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[create-meeting-note] query project info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[create-meeting-note] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("project_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	meetingUUID := c.Param("meeting_uuid")
	meeting, err := h.Meeting.ReadByUUID(ctx, meetingUUID)
	if err != nil {
		log.For(c).Error("[create-meeting-note] query meeting info failed", log.Field("user_id", userID), log.Field("meeting_uuid", meetingUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if meeting.ProjectID != project.ID {
		log.For(c).Error("[create-meeting-note] project id is not match", log.Field("user_id", userID),
			log.Field("meeting_project_id", meeting.ProjectID), log.Field("project_id", project.ID))
		c.JSON(http.StatusForbidden, err)
		return
	}

	data := &entities.MeetingNote{
		Base: &entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		ProjectID: project.ID,
		MeetingID: meeting.ID,
		StartTime: req.StartTime,
		Link:      req.Link,
		Location:  req.Location,
		Note:      req.Note,
	}

	userList, err := h.User.ListByUUIDs(ctx, req.UserUUIDs)
	if err != nil {
		log.For(c).Error("[create-meeting-note] list users failed", log.Field("user_id", userID), log.Field("user_uuids", req.UserUUIDs), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range userList {
		data.Contributors = append(data.Contributors, &entities.Contributor{
			UUID:      uuid.NewString(),
			UserID:    val.GetID(),
			CreatedBy: userID,
		})
	}

	contactList, err := h.Contact.ListByUUIDs(ctx, req.ContactUUIDs)
	if err != nil {
		log.For(c).Error("[create-meeting-note] list contacts failed", log.Field("user_id", userID), log.Field("contact_uuids", req.ContactUUIDs), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range contactList {
		data.Contributors = append(data.Contributors, &entities.Contributor{
			UUID:      uuid.NewString(),
			ContactID: val.GetID(),
			CreatedBy: userID,
		})
	}

	data.Editors = append(data.Editors, &entities.MeetingNoteEditor{
		UUID:      uuid.NewString(),
		EditorID:  userID,
		CreatedBy: userID,
	})

	_, err = h.MeetingNote.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-meeting-note] create meeting failed", log.Field("user_id", userID),
			log.Field("project_uuid", projectUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[create-meeting-note] process success", log.Field("user_id", userID), log.Field("meeting_uuid", data.UUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
