package client

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

func (h Handler) CreateClient(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateClientRequest{}

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

	// check exist name
	if !req.Force {
		existData, err := h.Client.ReadByName(ctx, user.OrganizationID, req.FullName)
		if err == nil {
			c.JSON(http.StatusConflict, &apis.Error{
				Message: &apis.ErrorMessage{
					VI: "Client name đã tồn tại",
					EN: "Client name is exist",
				},
				Data:  apis.CreateResponse{
					UUID: existData.UUID,
				},
			})
			return
		}
	}

	data := &entities.Client{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName:		req.FullName,
		ShortName:		req.ShortName,
		Code:			req.Code,
		Fax:			req.Fax,
		Website:		req.Website,
		Phone:			req.Phone,
		Email:			req.Email,
		CompanySize:	req.CompanySize,
		OrganizationID: user.OrganizationID,
		Address: 		req.Address,
	}
	image, err := services.UploadImage(c, "image", 125)
	if err != nil {
		log.For(c).Error("[create-client] save image failed", log.Field("user_id", userID), log.Err(err))
	} else {
		data.LogoID = image.UUID
		data.Logo = &entities.File{
			UUID:				image.UUID,
			OriginalName:		image.OriginalName,
			RelativePath:		image.RelativePathFile,
			RelativeThumbnail:	image.Thumbnail,
			Ext:				image.FileExt,
			CreatedBy:			userID,
		}
	}

	_, err = h.Client.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-client] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
