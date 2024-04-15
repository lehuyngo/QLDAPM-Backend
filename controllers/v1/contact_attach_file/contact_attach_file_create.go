package contact_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateContactAttachFile(c *gin.Context) {
	ctx := c.Request.Context()
	
	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contact-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[create-contact-attach-file] query contact info", log.Field("user_id", userID), 
			log.Field("contact_id", contact.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != contact.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	file, err := services.UploadFile(c, "file")
	if err != nil {
		log.For(c).Error("[create-contact-attach-file] write file failed", log.Field("user_id", userID), 
			log.Field("contact_id", contact.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
	}

	data := &entities.ContactAttachFile{
		OriginalName: file.OriginalName,
		RelativePath: file.RelativePathFile,
		Ext: file.FileExt,
		OrganizationID: user.OrganizationID,
		ContactID: contact.ID,
		CreatedBy: userID,
	}
	_, err = h.ContactAttachFile.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contact-attach-file] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
