package contact_mail_shortclick

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateContactMailShortClick(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactMailShortClickRequest{}

	contactUUID := c.Param("uuid")

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)
	
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[create-contact-mail-shortclick] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contact-mail-shortclick] query user info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// check permission
	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[create-contact-mail-shortclick] organization isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), 
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}
	
	data := &entities.ContactMailShortClick{
		ContactID: contact.ID,
		Type: entities.ActivityCreated,
		CreatedBy: user.ID,
	}

	_, err = h.ContactMailShortClick.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contact-note] insert database failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[create-contact-note] insert database failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
