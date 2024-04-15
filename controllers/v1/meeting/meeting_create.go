package meeting

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

func (h Handler) CreateMeeting(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateMeetingRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		log.For(c).Error("[create-meeting] query user info failed", log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit meeting of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-meeting] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[create-meeting] query meeting info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[create-meeting] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("project_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data := &entities.Meeting{
		Base: &entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		ProjectID: project.ID,
		StartTime: req.StartTime,
		Link:      req.Link,
		Location:  req.Location,
	}

	userList, err := h.User.ListByUUIDs(ctx, req.UserUUIDs)
	if err != nil {
		log.For(c).Error("[create-meeting] list of users  failed", log.Field("user_id", userID), log.Field("user_uuids", req.UserUUIDs), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range userList {
		data.Attendees = append(data.Attendees, &entities.Attendee{
			UUID:      uuid.NewString(),
			UserID:    val.GetID(),
			CreatedBy: userID,
		})
	}

	contactList, err := h.Contact.ListByUUIDs(ctx, req.ContactUUIDs)
	if err != nil {
		log.For(c).Error("[create-meeting] list contacts failed", log.Field("user_id", userID), log.Field("contact_uuids", req.ContactUUIDs), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range contactList {
		data.Attendees = append(data.Attendees, &entities.Attendee{
			UUID:      uuid.NewString(),
			ContactID: val.GetID(),
			CreatedBy: userID,
		})
	}

	_, err = h.Meeting.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-meeting] create meeting failed", log.Field("user_id", userID),
			log.Field("project_uuid", projectUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[create-meeting] process success", log.Field("user_id", userID), log.Field("meeting_uuid", data.UUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
