package draft_contact

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

func (h Handler) CreateDraftContact(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateDraftContactRequest{}

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

	data := &entities.DraftContact{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},

		FullName: req.FullName,
		Phone: req.Phone,
		Email: req.Email,
		ClientName: req.ClientName,
		ClientWebsite: req.ClientWebsite,
		ClientAddress: req.ClientAddress,
		OrganizationID: user.OrganizationID,
	}
	namecard, err := services.UploadImage(c, "name_card", 125)
	if err != nil {
		log.For(c).Error("[create-contact] save name card failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	
	data.NameCardID = namecard.UUID
	data.NameCard = &entities.File{
		UUID:				namecard.UUID,
		OriginalName:		namecard.OriginalName,
		RelativePath:		namecard.RelativePathFile,
		RelativeThumbnail:	namecard.Thumbnail,
		Ext:				namecard.FileExt,
		CreatedBy:			userID,
	}

	companyLogo, err := services.UploadImage(c, "company_logo", 125)
	if err != nil {
		log.For(c).Error("[create-contact] save avatar failed", log.Field("user_id", userID), log.Err(err))
	} else {
		data.CompanyLogoID = companyLogo.UUID
		data.CompanyLogo = &entities.File{
			UUID:				companyLogo.UUID,
			OriginalName:		companyLogo.OriginalName,
			RelativePath:		companyLogo.RelativePathFile,
			RelativeThumbnail:	companyLogo.Thumbnail,
			Ext:				companyLogo.FileExt,
			CreatedBy:			userID,
		}
	}

	_, err = h.DraftContact.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-contact] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
