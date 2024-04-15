package contact_tag

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/define"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateContactTag(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactTagRequest{}

	contactUUID := c.Param("uuid")

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	req.Name = strings.ToLower(strings.TrimSpace(req.Name))
	req.Color = strings.ToLower(strings.TrimSpace(req.Color))

	userID, _, _ := middlewares.ParseToken(c)

	contact, err := h.Contact.ReadByUUID(ctx, contactUUID)
	if err != nil {
		log.For(c).Error("[create-contact-tag] query contact info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-contact-tag] query user info failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// check permission
	if user.OrganizationID != contact.OrganizationID {
		log.For(c).Error("[create-contact-tag] organization isnot match", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("contact_organization_id", contact.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	tag, err := h.ContactTag.ReadByName(ctx, user.OrganizationID, req.Name)
	if err != nil {
		tag = &entities.ContactTag{
			Name:           req.Name,
			OrganizationID: user.OrganizationID,
			Color: 			req.Color,
			CreatedBy:      userID,
		}
		tag.Contacts = make([]*entities.ContactContactTag, 0)
		tag.Contacts = append(tag.Contacts, &entities.ContactContactTag{
			ContactID: contact.ID,
			Color:     req.Color,
			CreatedBy: userID,
		})
		tag.CreatedBy = userID

		_, err = h.ContactTag.Create(ctx, tag)
		if err != nil {
			log.For(c).Error("[create-contact-tag] insert database failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		log.For(c).Info("[create-contact-tag] insert database failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: tag.UUID,
		})
		return
	}

	err = h.ContactTag.Add(ctx, &entities.ContactContactTag{
		ContactID: contact.ID,
		TagID:     tag.ID,
		Color:     req.Color,
		CreatedBy: userID,
	})
	if define.IsErrDuplicateKey(err) {
		tag.Color = req.Color
		h.ContactTag.Update(ctx, tag)

		log.For(c).Info("[create-contact-tag] tag added", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: tag.UUID,
		})
		return
	}
	if err != nil {
		log.For(c).Error("[create-contact-tag] insert database failed", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	tag.Color = req.Color
	h.ContactTag.Update(ctx, tag)

	log.For(c).Info("[create-contact-tag] add tag success", log.Field("user_id", userID), log.Field("contact_uuid", contactUUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: tag.UUID,
	})
}
