package draft_contact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) DownloadDraftContactCompanyLogo(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.DraftContact.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[download-draft-contact-company-logo] query draft contact info", log.Field("uuid", uuid), 
			log.Field("user_id", userID), log.Field("draft_contact_uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[download-draft-contact-company-logo] query user info", log.Field("uuid", uuid), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		log.For(c).Error("[download-draft-contact-company-logo] organization id not match", 
			log.Field("user_id", userID), log.Field("contact_uuid", uuid), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("data_organization_id", data.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + data.GetCompanyLogo().GetRelativePath()
	fileName := data.GetCompanyLogo().GetOriginalName() + data.GetCompanyLogo().GetExt()
	log.For(c).Info("[download-draft-contact-company-logo] download", log.Field("uuid", uuid), 
			log.Field("user_id", userID), log.Field("client_uuid", uuid), 
			log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
