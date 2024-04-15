package contributor

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

func (h Handler) CreateContributorBatch(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContributorBatchRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contributor-batch] query user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get note
	noteUUID := c.Param("uuid")
	if noteUUID == "" {
		log.For(c).Error("[create-contributor-batch] note uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[create-contributor-batch] query note by uuid failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != note.GetProject().GetOrganizationID() {
		log.For(c).Error("[create-contributor-batch] organization is not match", log.Field("user_id", userID), log.Field("note_uuid", noteUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("note_organization_id", note.GetProject().OrganizationID))
		c.JSON(http.StatusForbidden, nil)
		return
	}

	contactUUIDs := req.ContactUUIDs
	userUUIDs := req.UserUUIDs

	if len(contactUUIDs) < 1 && len(userUUIDs) < 1 {
		log.For(c).Error("[create-contributor-batch] contact uuid and user uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	contacts, err := h.Contact.ListByUUIDs(ctx, contactUUIDs)
	if err != nil {
		log.For(c).Error("[create-contributor-batch] query contacts by uuids failed", log.Field("user_id", userID), log.Field("contact_uuids", contactUUIDs))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	users, err := h.User.ListByUUIDs(ctx, userUUIDs)
	if err != nil {
		log.For(c).Error("[create-contributor-batch] query users by uuids failed", log.Field("user_id", userID), log.Field("user_uuids", userUUIDs))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data := []*entities.Contributor{}

	for _, contributorContact := range contacts {
		if contributorContact.OrganizationID != user.OrganizationID {
			log.For(c).Error("[create-contributor-batch] organization is not match", log.Field("user_id", userID), log.Field("contact_uuid", contributorContact.UUID),
				log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contributorContact.OrganizationID))
			c.JSON(http.StatusForbidden, nil)
			return
		}

		isExist := false
		for _, contributor := range note.GetContributors() {
			if contributor.GetContact() != nil && contributor.GetContact().GetUUID() == contributorContact.UUID {
				isExist = true
				break
			}
		}

		if !isExist {
			data = append(data, &entities.Contributor{
				UUID:          uuid.NewString(),
				MeetingNoteID: note.ID,
				CreatedBy:     user.ID,
				ContactID:     contributorContact.ID,
			})
		}
	}

	for _, contributorUser := range users {
		if contributorUser.OrganizationID != user.OrganizationID {
			log.For(c).Error("[create-contributor-batch] organization is not match", log.Field("user_id", userID), log.Field("user_uuid", contributorUser.UUID),
				log.Field("user_organization_id", user.OrganizationID), log.Field("user_organization_id", contributorUser.OrganizationID))
			c.JSON(http.StatusForbidden, nil)
			return
		}

		isExist := false
		for _, contributor := range note.GetContributors() {
			if contributor.GetUser() != nil && contributor.GetUser().GetUUID() == contributorUser.UUID {
				isExist = true
				break
			}
		}

		if !isExist {
			data = append(data, &entities.Contributor{
				UUID:          uuid.NewString(),
				MeetingNoteID: note.ID,
				CreatedBy:     user.ID,
				UserID:        contributorUser.ID,
			})
		}
	}

	if len(data) < 1 {
		log.For(c).Error("[create-contributor-batch] all contributors exist", log.Field("user_id", userID))
		c.JSON(http.StatusOK, nil)
		return
	}

	// create contributors
	err = h.Contributor.CreateBatch(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contributor-batch] create contributors failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
