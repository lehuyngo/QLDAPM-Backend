package url_access

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) RedirectShortenLink(c *gin.Context) {
	ctx := c.Request.Context()
	
	code := c.Param("token")
	url, err := h.TrackedURL.ReadByCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// mark url accessed
	url.Status = entities.Read
	url.ReceivedAt = time.Now().UnixMilli()
	err = h.TrackedURL.Update(ctx, url)
	if err != nil {
		log.For(c).Error("[redirect-shorten-link] update database failed", log.Field("code", code), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	
	// add event
	event := &entities.EventClickURL {
		Code: code,
		SenderID: url.GetCreator().GetID(),
		ReceiverID: url.GetContact().GetID(),
		SendTime: url.CreatedAt.UnixMilli(),
		ReadTime: time.Now().UnixMilli(),
		URL: url.URL,
		OriginalURL: url.OriginalURL,
		OrganizationID: url.OrganizationID,
	}

	if url.GetMail() != nil {
		event.MailID = url.GetMail().ID
		event.MailUUID = url.GetMail().GetUUID()
		event.MailSubject = url.GetMail().GetSubject()
	}

	if url.GetBatchMail() != nil {
		event.BatchMailID = url.GetBatchMail().ID
		event.BatchMailUUID = url.GetBatchMail().GetUUID()
		event.BatchMailSubject = url.GetBatchMail().GetSubject()
	}

	err = h.EventClickURL.Create(ctx, event)
	if err != nil {
		log.For(c).Error("[redirect-shorten-link] create shorten link failed", log.Field("code", code), log.Err(err))
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}
