package meeting_note

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListMeetingNote(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-meeting-note] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[list-meeting-note] query project info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[list-meeting-note] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("meeting_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.MeetingNote.List(ctx, project.ID)
	if err != nil {
		log.For(c).Error("[list-meeting-note] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListMeetingNoteResponse{}
	for _, val := range data {
		meetingNote := &apis.MeetingNote{
			UUID:      val.UUID,
			StartTime: val.StartTime,
			Link:      val.Link,
			Location:  val.Location,
			Note:      val.Note,
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime: val.GetCreatedAt().UnixMilli(),
		}

		for _, contributor := range val.GetContributors() {
			element := &apis.Contributor{
				UUID: contributor.GetUUID(),
			}

			if contributor.GetContact() != nil {
				element.Contact = &apis.Contact{
					UUID:      contributor.GetContact().GetUUID(),
					Email:     contributor.GetContact().GetEmail(),
					FullName:  contributor.GetContact().GetFullName(),
					ShortName: contributor.GetContact().GetShortName(),
				}
			}

			if contributor.GetUser() != nil {
				element.User = &apis.User{
					UUID:        contributor.GetUser().GetUUID(),
					DisplayName: contributor.GetUser().GetDisplayName(),
					Email:       contributor.GetUser().GetEmail(),
				}
			}

			meetingNote.Contributors = append(meetingNote.Contributors, element)
		}

		for _, editor := range val.GetEditors() {
			meetingNote.Editors = append(meetingNote.Editors, &apis.User{
				UUID:        editor.GetEditor().GetUUID(),
				DisplayName: editor.GetEditor().GetDisplayName(),
			})
		}

		resp.Data = append(resp.Data, meetingNote)
	}

	c.JSON(http.StatusOK, resp)
}
