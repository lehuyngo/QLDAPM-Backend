package contributor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteContributorBatch(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.DeleteContributorBatchRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-contributor-batch] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get note
	noteUUID := c.Param("uuid")
	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[delete-contributor-batch] query note info failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// organization check
	if user.OrganizationID != note.GetProject().GetOrganizationID() {
		log.For(c).Error("[delete-contributor-batch] query organization id is not match", log.Field("user_id", userID), log.Field("note_uuid", noteUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("note_organization_id", note.GetProject().OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	// list contributors
	contributors, err := h.Contributor.ListByUUIDs(ctx, req.UUIDs)
	if err != nil {
		log.For(c).Error("[delete-contributor-batch] query contributors failed", log.Field("user_id", userID), log.Field("contributor_uuids", req.UUIDs), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// check contributors organization
	for _, contributor := range contributors {
		if contributor.MeetingNoteID != note.GetID() {
			log.For(c).Error("[delete-contributor-batch] note id is not match", log.Field("user_id", userID), log.Field("contributor_uuid", contributor.GetUUID()),
				log.Field("note_id", note.GetID()), log.Field("contributor_note_id", contributor.MeetingNoteID))
			c.JSON(http.StatusForbidden, nil)
			return
		}
	}

	// delete contributors
	err = h.Contributor.DeleteBatch(ctx, req.UUIDs)
	if err != nil {
		log.For(c).Error("[delete-contributor-batch] delete contributors failed", log.Field("user_id", userID), log.Field("contributor_uuids", req.UUIDs), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[delete-contributor-batch] process success", log.Field("user_id", userID), log.Field("contributor_uuids", req.UUIDs))
	c.JSON(http.StatusOK, nil)
}
