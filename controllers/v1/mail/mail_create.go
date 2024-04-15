package mail

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/email"
	file_utils "gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateMail(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateMailRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// hardcode
	req.UseShortenLink = 1

	useShortenLink := (req.UseShortenLink == 1)

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data := &entities.Mail{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		Subject:		req.Subject,
		Content:		req.Content,
		OrganizationID: user.OrganizationID,
	}

	var to []string
	contacts, err := h.Contact.ListByUUIDs(ctx, strings.Split(req.ReceiverContactUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-mail] query list contact failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range contacts {
		data.Receivers = append(data.Receivers, &entities.MailReceiver{
			UUID: uuid.NewString(),
			ContactID: val.GetID(),
			CreatedBy: user.ID,
		})

		if val.Email != "" {
			to = append(to, val.Email)
		}
	}

	users, err := h.User.ListByUUIDs(ctx, strings.Split(req.ReceiverUserUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-mail] query list user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range users {
		data.Receivers = append(data.Receivers, &entities.MailReceiver{
			UUID: uuid.NewString(),
			UserID: val.GetID(),
			CreatedBy: user.ID,
		})

		if val.Email != "" {
			to = append(to, val.Email)
		}
	}

	receiverMailAddresses := strings.Split(req.ReceiverMailAddresses, ",")
	for _, val := range receiverMailAddresses {
		data.Receivers = append(data.Receivers, &entities.MailReceiver{
			UUID: uuid.NewString(),
			MailAddress: val,
			CreatedBy: user.ID,
		})

		if val != "" {
			to = append(to, val)
		}
	}

	var cc []string
	ccContacts, err := h.Contact.ListByUUIDs(ctx, strings.Split(req.CCContactUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-mail] query list cc contact failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range ccContacts {
		data.CarbonCopies = append(data.CarbonCopies, &entities.MailCarbonCopy{
			UUID: uuid.NewString(),
			ContactID: val.GetID(),
			CreatedBy: user.ID,
		})

		if val.Email != "" {
			cc = append(cc, val.Email)
		}
	}

	ccUsers, err := h.User.ListByUUIDs(ctx, strings.Split(req.CCUserUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-mail] query list cc user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range ccUsers {
		data.CarbonCopies = append(data.CarbonCopies, &entities.MailCarbonCopy{
			UUID: uuid.NewString(),
			UserID: val.GetID(),
			CreatedBy: user.ID,
		})

		if val.Email != "" {
			cc = append(cc, val.Email)
		}
	}

	ccMailAddresses := strings.Split(req.CCMailAddresses, ",")
	for _, val := range ccMailAddresses {
		data.CarbonCopies = append(data.CarbonCopies, &entities.MailCarbonCopy{
			UUID: uuid.NewString(),
			MailAddress: val,
			CreatedBy: user.ID,
		})

		if val != "" {
			cc = append(cc, val)
		}
	}

	var attachFiles []string
	for i := 1; i < 5; i++ {
		file, err := services.UploadFile(c, fmt.Sprintf("attach_file_%d", i))
		if err != nil {
			continue
		}

		data.AttachFiles = append(data.AttachFiles, &entities.MailAttachFile{
			UUID:			file.UUID,
			OriginalName:	file.OriginalName,
			RelativePath:	file.RelativePathFile,
			Ext: 			file.FileExt,
			OrganizationID: user.OrganizationID,
			CreatedBy: 		user.ID,
		})

		// attachFiles = append(attachFiles, services.Config.FileStorage.Folder + file.RelativePathFile)
		folder := services.Config.FileStorage.TempFolder + file.UUID
		if err := os.Mkdir(folder, os.ModePerm); err != nil {
			continue
		}

		filePath := folder + "/" + file.OriginalName + file.FileExt
		err = file_utils.CopyFile(services.Config.FileStorage.Folder + file.RelativePathFile, filePath)
		if err != nil {
			continue
		}
		
		attachFiles = append(attachFiles, filePath)
	}

	// replace content by first contact
	var codes map[string]entities.EncodedURL
	var contactID int64
	if len(contacts) > 0 {
		contactID = contacts[0].ID
	}
	req.Content, codes, err = services.UpgradeURLContent(ctx, contactID, req.Content, useShortenLink)
	if err != nil {
		log.For(c).Error("[create-mail] upgrade url content failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	err = services.Mailer.Send(&email.Message{
		To: to,
		Cc: cc,
		Subject: req.Subject,
		BodyHTML: req.Content,
		AttachFiles: attachFiles,
	})
	if err != nil {
		log.For(c).Error("[create-mail] send mail failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// remove temp file
	for _, val := range attachFiles {
		os.Remove(val)
	}

	for _, val := range data.AttachFiles {
		os.Remove(services.Config.FileStorage.TempFolder + val.UUID)
	}

	// update database
	_, err = h.Mail.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-mail] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	for code, url := range codes {
		val := &entities.TrackedURL{
			UUID: uuid.NewString(),
			Code: code,
			URL: url.OriginalURL,
			OriginalURL: url.OriginalURL,
			ContactID: contactID,
			MailID: data.ID,
			CreatedBy: user.ID,
			Status: entities.New,
			OrganizationID: user.OrganizationID,
		}
		if useShortenLink == true {
			val.URL = url.ShortenLink
		}
		h.TrackedURL.Create(ctx, val)
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
