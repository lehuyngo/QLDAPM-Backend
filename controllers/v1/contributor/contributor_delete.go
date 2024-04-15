package contributor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteContributor(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-contributor] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get note
	noteUUID := c.Param("uuid")
	note, err := h.MeetingNote.ReadByUUID(ctx, noteUUID)
	if err != nil {
		log.For(c).Error("[delete-contributor] query note info failed", log.Field("user_id", userID), log.Field("note_uuid", noteUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// organization check
	if user.OrganizationID != note.GetProject().GetOrganizationID() {
		log.For(c).Error("[delete-contributor] query organization id is not match", log.Field("user_id", userID), log.Field("note_uuid", noteUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("note_organization_id", note.GetProject().OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	// get contributor
	contributorUUID := c.Param("contributor_uuid")
	contributor, err := h.Contributor.ReadByUUID(ctx, contributorUUID)
	if err != nil {
		log.For(c).Error("[delete-contributor] query contributor info failed", log.Field("user_id", userID), log.Field("contributor_uuid", contributor), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// delete contributor
	err = h.Contributor.Delete(ctx, contributor.ID)
	if err != nil {
		log.For(c).Error("[delete-contributor] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[delete-contributor] process success", log.Field("user_id", userID), log.Field("contributor_uuid", contributorUUID))
	c.JSON(http.StatusOK, nil)
}
