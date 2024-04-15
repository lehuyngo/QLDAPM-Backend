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

func (h Handler) CreateContact(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateContactRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// create new
	data := &entities.Contact{
		Base: &entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName:		req.FullName,
		ShortName:		req.ShortName,
		JobTitle:		req.JobTitle,
		Phone:			req.Phone,
		Email:			req.Email,
		OrganizationID: user.OrganizationID,
		Gender: 		entities.Gender(req.Gender),
		Birthday: 		req.BirthDay,
	}
	namecard, err := services.UploadImage(c, "name_card", 125)
	if err != nil {
		log.For(c).Error("[create-contact] save name card failed", log.Field("user_id", userID), log.Err(err))
	} else {
		data.NameCardID = namecard.UUID
		data.NameCard = &entities.File{
			UUID:				namecard.UUID,
			OriginalName:		namecard.OriginalName,
			RelativePath:		namecard.RelativePathFile,
			RelativeThumbnail:	namecard.Thumbnail,
			Ext:				namecard.FileExt,
			CreatedBy:			userID,
		}
	}

	avatar, err := services.UploadImage(c, "avatar", 125)
	if err != nil {
		log.For(c).Error("[create-contact] save avatar failed", log.Field("user_id", userID), log.Err(err))
	} else {
		data.AvatarID = avatar.UUID
		data.Avatar = &entities.File{
			UUID:				avatar.UUID,
			OriginalName:		avatar.OriginalName,
			RelativePath:		avatar.RelativePathFile,
			RelativeThumbnail:	avatar.Thumbnail,
			Ext:				avatar.FileExt,
			CreatedBy:			userID,
		}
	}

	// check exist to update
	existContact, err := h.Contact.ReadByEmail(ctx, user.OrganizationID, req.Email)
	if err == nil {
		data.ID = existContact.ID
		data.UUID = existContact.UUID

		err = h.Contact.Update(ctx, data)
		if err != nil {
			log.For(c).Error("[create-contact] update exist contact failed", log.Field("user_id", userID), log.Err(err))
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, &apis.CreateResponse{
			UUID: data.UUID,
		})
		return
	}

	// create
	_, err = h.Contact.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contact] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
