package services

import (
	"context"
	"fmt"
	"net/mail"
	"os"

	"github.com/google/uuid"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/models"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/email"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/file"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/queue"
)

var mailQueue *queue.LocalQueue

var batchMailReceiver *models.BatchMailReceiver

func isEmail(email string) bool {
	emailAddress, err := mail.ParseAddress(email)
	return err == nil && emailAddress.Address == email
  }

func SendMailHandleFunc(ctx context.Context, ticket queue.WorkTicket) error {
	mailReceiverID := ticket.Data.(int64)
	if mailReceiverID < 1 {
		log.WithContext(ctx).Error("[send-batch-mail-func] receiver id invalid", log.Field("mail_receiver_id", mailReceiverID))
		return fmt.Errorf("receiver id invalid")
	}

	receiver, err := batchMailReceiver.Read(ctx, mailReceiverID)
	if err != nil {
		log.WithContext(ctx).Error("[send-batch-mail-func] query receiver failed", log.Field("mail_receiver_id", mailReceiverID), log.Err(err))
		return err
	}

	if !isEmail(receiver.Email) {
		log.WithContext(ctx).Error("[send-batch-mail-func] receiver email is invalid", log.Field("mail_receiver_id", mailReceiverID), log.Field("email", receiver.Email))
		receiver.Status = entities.MailFailed
		batchMailReceiver.Update(ctx, receiver)

		return fmt.Errorf("receiver email is invalid")
	}

	if receiver.Mail == nil {
		log.WithContext(ctx).Error("[send-batch-mail-func] mail is nil", log.Field("mail_receiver_id", mailReceiverID))
		receiver.Status = entities.MailFailed
		batchMailReceiver.Update(ctx, receiver)

		return fmt.Errorf("mail is nil")
	}

	userShortenLink := (receiver.Mail.UserShortenLink == 1)

	var codes map[string]entities.EncodedURL
	receiver.Mail.Content, codes, _ = UpgradeURLContent(ctx, receiver.ContactID, receiver.Mail.Content, userShortenLink)
	if err != nil {
		log.WithContext(ctx).Error("[create-mail] upgrade url content failed", log.Field("user_id", receiver.CreatedBy), log.Err(err))
		return err
	}

	msg := &email.Message{
		To: []string{receiver.Email},
		Subject: receiver.Mail.Subject,
		BodyHTML: receiver.Mail.Content,
	}

	for _, val := range receiver.Mail.CarbonCopies {
		if !isEmail(val.MailAddress) {
			log.WithContext(ctx).Error("[send-batch-mail-func] email is invalid", log.Field("mail_receiver_id", mailReceiverID), log.Field("email", val.MailAddress))
			continue
		}

		msg.Cc = append(msg.Cc, val.MailAddress)
	}
	
	for _, val := range receiver.Mail.AttachFiles {
		folder := Config.FileStorage.TempFolder + val.UUID
		if err := os.Mkdir(folder, os.ModePerm); err != nil {
			continue
		}

		filePath := folder + "/" + val.OriginalName + val.Ext
		err = file.CopyFile(Config.FileStorage.Folder + val.RelativePath, filePath)
		if err != nil {
			continue
		}
		
		msg.AttachFiles = append(msg.AttachFiles, filePath)
	}

	err = Mailer.Send(msg)
	if err != nil {
		log.WithContext(ctx).Error("[send-batch-mail-func] send mail failed", log.Field("mail_receiver_id", mailReceiverID))
		receiver.Status = entities.MailFailed
		batchMailReceiver.Update(ctx, receiver)

		for _, val := range msg.AttachFiles {
			os.Remove(val)
		}

		for _, val := range receiver.Mail.AttachFiles {
			os.Remove(Config.FileStorage.TempFolder + val.UUID)
		}

		return err
	}

	log.WithContext(ctx).Info("[send-batch-mail-func] send mail success", log.Field("mail_receiver_id", mailReceiverID))
	receiver.Status = entities.MailSuccess
	batchMailReceiver.Update(ctx, receiver)

	for _, val := range msg.AttachFiles {
		os.Remove(val)
	}

	for _, val := range receiver.Mail.AttachFiles {
		os.Remove(Config.FileStorage.TempFolder + val.UUID)
	}

	trackedURLHandler := NewTrackedURL()
	for code, url := range codes {
		val := &entities.TrackedURL{
			UUID: uuid.NewString(),
			Code: code,
			URL: url.OriginalURL,
			OriginalURL: url.OriginalURL,
			ContactID: receiver.GetContact().GetID(),
			BatchMailID: receiver.MailID,
			CreatedBy: receiver.CreatedBy,
			Status: entities.New,
			OrganizationID: receiver.Mail.OrganizationID,
		}
		if userShortenLink == true {
			val.URL = url.ShortenLink
		}
		trackedURLHandler.Create(ctx, val)
	}

	contactModel := &models.Contact{}
 	contactModel.UpdateLastActiveTime(ctx, receiver.GetContact().GetID())
	
	return nil
}

func StartMailQueue() {
	batchMailReceiver = &models.BatchMailReceiver{}
	mailQueue = &queue.LocalQueue{
		NumberWorker: 5,
		TimeoutPerTicket: 20,
		QueueSize: 300,
		Handler: SendMailHandleFunc,
	}

	mailQueue.StartWorkerDispatcher()
}

func StopMailQueue() {
	mailQueue.StopWorkerDispatcher()
}

func AddMailQueue(receiverID int64) {
	mailQueue.SendTicket(queue.WorkTicket{
		Data: receiverID,
	})
}
