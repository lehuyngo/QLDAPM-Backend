package contributor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateContributor(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContributorRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contributor] query user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get note
	noteUUID := c.Param("uuid")
	if noteUUID == "" {
		log.For(c).Error("[create-contributor] note uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[create-contributor] query note by uuid failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != note.GetProject().GetOrganizationID() {
		log.For(c).Error("[create-contributor] organization is not match", log.Field("user_id", userID), log.Field("note_uuid", noteUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("note_organization_id", note.GetProject().OrganizationID))
		c.JSON(http.StatusForbidden, nil)
		return
	}

	
	contributorUserUUID := req.UserUUID

	
	

	data := &entities.Contributor{
		MeetingNoteID: note.ID,
		CreatedBy:     user.ID,
	}

	
	if contributorUserUUID != "" {
		contributorUser, err := h.User.ReadByUUID(ctx, contributorUserUUID)
		if err != nil {
			log.For(c).Error("[create-contributor] query user by uuid failed", log.Field("user_id", userID), log.Field("user_uuid", contributorUserUUID))
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		if user.OrganizationID != contributorUser.OrganizationID {
			log.For(c).Error("[create-contributor] organization is not match", log.Field("user_id", userID), log.Field("user_uuid", contributorUserUUID),
				log.Field("user_organization_id", user.OrganizationID), log.Field("query_user_organization_id", contributorUser.OrganizationID))
			c.JSON(http.StatusForbidden, nil)
			return
		}
		data.UserID = contributorUser.ID
	}

	// create contributor
	_, err = h.Contributor.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contributor] create contributor failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
