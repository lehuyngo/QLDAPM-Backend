package client

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
	"gorm.io/gorm"
)


func (h Handler) UpdateClient(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Client.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	// check allow same name
	if !req.Force {
		existData, err := h.Client.ReadByName(ctx, user.OrganizationID, req.FullName)
		if err == nil {
			if existData.GetID() != data.GetID() {
				c.JSON(http.StatusConflict, &apis.Error{
					Message: &apis.ErrorMessage{
						VI: "Client name đã tồn tại",
						EN: "Client name is exist",
					},
				})
				return
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	image, err := services.UploadImage(c, "image", 125)
	if err != nil {
		log.For(c).Error("[update-client] save image failed", log.Field("user_id", userID), log.Field("client_id", data.GetID()), log.Err(err))
	} else {
		if data.GetLogo().GetUUID() == "" {
			data.Logo = &entities.File{}
		}

		data.LogoID = image.UUID
		data.Logo.UUID = image.UUID
		data.Logo.OriginalName = image.OriginalName
		data.Logo.RelativePath = image.RelativePathFile
		data.Logo.RelativeThumbnail = image.Thumbnail
		data.Logo.Ext = image.FileExt
		data.Logo.CreatedBy = userID
	}

	data.UpdatedBy = userID
	data.FullName =	req.FullName
	data.ShortName = req.ShortName
	data.Code = req.Code
	data.Fax = req.Fax
	data.Website = req.Website
	data.Address = req.Address
	data.Phone = req.Phone
	data.Email = req.Email
	data.CompanySize = req.CompanySize

	err = h.Client.Update(ctx, data)
	if err != nil {
		log.For(c).Error("[update-client] update database failed", log.Field("user_id", userID), log.Field("client_id", data.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
