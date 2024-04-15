package meeting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListMeeting(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-meeting] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[list-meeting] query project info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[list-meeting] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("meeting_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.Meeting.List(ctx, project.ID)
	if err != nil {
		log.For(c).Error("[list-meeting] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListMeetingResponse{}
	for _, val := range data {
		meeting := &apis.Meeting{
			UUID:      val.UUID,
			StartTime: val.StartTime,
			Link:      val.Link,
			Location:  val.Location,
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime:    val.GetCreatedAt().UnixMilli(),
			LastActiveTime: val.GetLastActiveTime(),
		}

		for _, attendee := range val.GetAttendees() {
			element := &apis.Attendee{
				UUID: attendee.GetUUID(),
			}

			if attendee.GetContact() != nil {
				element.Contact = &apis.Contact{
					UUID:      attendee.GetContact().GetUUID(),
					Email:     attendee.GetContact().GetEmail(),
					FullName:  attendee.GetContact().GetFullName(),
					ShortName: attendee.GetContact().GetShortName(),
				}
			}

			if attendee.GetUser() != nil {
				element.User = &apis.User{
					UUID:        attendee.GetUser().GetUUID(),
					DisplayName: attendee.GetUser().GetDisplayName(),
					Email:       attendee.GetUser().GetEmail(),
				}
			}

			meeting.Attendees = append(meeting.Attendees, element)
		}

		notes, err := h.MeetingNote.ListByMeetingID(ctx, val.ID)
		if err != nil {
			log.For(c).Error("[list-meeting] query meeting note info failed", log.Field("user_id", userID), log.Field("meeting_id", val.ID), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		for _, note := range notes {
			noteCreator := &apis.NoteCreator{
				NoteUUID: note.UUID,
				Creator: &apis.User{
					UUID:        note.GetCreator().GetUUID(),
					DisplayName: note.GetCreator().GetDisplayName(),
				},
			}

			meeting.NoteCreators = append(meeting.NoteCreators, noteCreator)
		}

		resp.Data = append(resp.Data, meeting)
	}

	c.JSON(http.StatusOK, resp)
}
