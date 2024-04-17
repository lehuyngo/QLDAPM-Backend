package batch_mail

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateBatchMail(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateBatchMailRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// hardcode
	req.UseShortenLink = 1

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data := &entities.BatchMail{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		Subject:			req.Subject,
		Content:			req.Content,
		OrganizationID: 	user.OrganizationID,
		Status: 			entities.MailProcessing,
		UserShortenLink: 	req.UseShortenLink,
	}

	// Receiver list
	
	if err != nil {
		log.For(c).Error("[create-batch-mail] query list contact failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	

	// CC list
	ccUsers, err := h.User.ListByUUIDs(ctx, strings.Split(req.CCUserUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-batch-mail] query list cc user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range ccUsers {
		data.CarbonCopies = append(data.CarbonCopies, &entities.BatchMailCarbonCopy{
			UUID: uuid.NewString(),
			UserID: val.GetID(),
			MailAddress: val.Email,
			CreatedBy: user.ID,
		})
	}

	ccMailAddresses := strings.Split(req.CCMailAddresses, ",")
	for _, val := range ccMailAddresses {
		data.CarbonCopies = append(data.CarbonCopies, &entities.BatchMailCarbonCopy{
			UUID: uuid.NewString(),
			MailAddress: val,
			CreatedBy: user.ID,
		})
	}

	// Attach files
	for i := 1; i < 5; i++ {
		file, err := services.UploadFile(c, fmt.Sprintf("attach_file_%d", i))
		if err != nil {
			continue
		}

		data.AttachFiles = append(data.AttachFiles, &entities.BatchMailAttachFile{
			UUID:			file.UUID,
			OriginalName:	file.OriginalName,
			RelativePath:	file.RelativePathFile,
			Ext: 			file.FileExt,
			OrganizationID: user.OrganizationID,
			CreatedBy: 		user.ID,
		})
	}

	_, err = h.BatchMail.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-batch-mail] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	for _, val := range data.Receivers {
		if val.ContactID < 1 {
			continue
		}

		
	}

	log.For(c).Info("[create-batch-mail] process success", log.Field("user_id", userID), log.Field("uuid", data.UUID))
	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
