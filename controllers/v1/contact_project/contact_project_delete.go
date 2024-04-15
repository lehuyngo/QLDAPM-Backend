package contact_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteContactProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-contact-project] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	contactUUID := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-project] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("project_uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[delete-contact-project] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if (user.OrganizationID != contact.OrganizationID) || (user.OrganizationID != project.OrganizationID) {
		log.For(c).Error("[delete-contact-project] query contact id isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.ContactProject.Delete(ctx, user.ID, contact.ID, project.ID)
	if err != nil {
		log.For(c).Error("[delete-contact-project] update database failed", log.Field("user_id", userID), 
			log.Field("contact_uuid", contactUUID), log.Field("project_uuid", project.UUID),
			log.Field("contact_id", contact.ID), log.Field("project_id", project.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-contact-project] process success", log.Field("user_id", userID), 
			log.Field("contact_uuid", contactUUID), log.Field("project_uuid", project.UUID))
	c.JSON(http.StatusOK, nil)
}
