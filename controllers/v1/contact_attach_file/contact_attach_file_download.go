package contact_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) DownloadContactAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[contact-attach-file-download] query contact info", log.Field("contact_uuid", contactUUID), 
			log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[contact-attach-file-download] query user info", log.Field("contact_uuid", contactUUID), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[contact-attach-file-download] organization id not match", 
			log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	fileUUID := c.Param("file_uuid")
	file, err := h.ContactAttachFile.Read(ctx, fileUUID)
	if err != nil {
		log.For(c).Error("[contact-attach-file-download] query file info failed", log.Field("contact_uuid", contactUUID), 
			log.Field("user_id", userID), log.Field("file_uuid", fileUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if file.ContactID != contact.ID {
		log.For(c).Error("[contact-attach-file-download] client id not match", 
			log.Field("user_id", userID), log.Field("client_id", contact.ID), log.Field("file_id", file.ID),
			log.Field("contact_id", file.ContactID), log.Field("file_client_id", file.ContactID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + file.GetRelativePath()
	fileName := file.GetOriginalName() + file.GetExt()
	log.For(c).Info("[client-attach-file-download] download",
			log.Field("user_id", userID), log.Field("client_uuid", contactUUID), log.Field("file_uuid", file.UUID),
			log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
