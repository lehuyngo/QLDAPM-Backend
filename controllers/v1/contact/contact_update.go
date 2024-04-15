package contact

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) UpdateContact(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Contact.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}
	// check if exist return error
	if data.Email != req.Email {
		_, err := h.Contact.ReadByEmail(ctx, user.OrganizationID, req.Email)
		if err == nil {
			log.For(c).Error("[update-contact] contact's email is exist", log.Field("user_id", userID), log.Err(err))
			c.JSON(http.StatusConflict, nil)
			return
		}
	}
	namecard, err := services.UploadImage(c, "name_card", 125)
	if err != nil {
		log.For(c).Error("[update-contact] save namecard failed", log.Field("user_id", userID), log.Field("contact_id", data.GetID()), log.Err(err))
	} else {
		if data.GetNameCard().GetUUID() == "" {
			data.NameCard = &entities.File{}
		}

		data.NameCardID = namecard.UUID
		data.NameCard.UUID = namecard.UUID
		data.NameCard.OriginalName = namecard.OriginalName
		data.NameCard.RelativePath = namecard.RelativePathFile
		data.NameCard.RelativeThumbnail = namecard.Thumbnail
		data.NameCard.Ext = namecard.FileExt
		data.NameCard.CreatedBy = userID
	}

	avatar, err := services.UploadImage(c, "avatar", 125)
	if err != nil {
		log.For(c).Error("[create-contact] save avatar failed", log.Field("user_id", userID), log.Err(err))
	} else {
		if data.GetAvatar().GetUUID() == "" {
			data.Avatar = &entities.File{}
		}

		data.AvatarID = avatar.UUID
		data.Avatar.UUID = avatar.UUID
		data.Avatar.OriginalName = avatar.OriginalName
		data.Avatar.RelativePath = avatar.RelativePathFile
		data.Avatar.RelativeThumbnail = avatar.Thumbnail
		data.Avatar.Ext = avatar.FileExt
		data.Avatar.CreatedBy = userID
	}

	data.UpdatedBy = userID
	data.FullName = req.FullName
	data.ShortName = req.ShortName
	data.Phone = req.Phone
	data.Email = req.Email
	data.JobTitle = req.JobTitle
	data.Gender = entities.Gender(req.Gender)
	data.Birthday = req.BirthDay

	err = h.Contact.Update(ctx, data)
	if err != nil {
		log.For(c).Error("[update-contact] update database failed", log.Field("user_id", userID), log.Field("contact_id", data.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
