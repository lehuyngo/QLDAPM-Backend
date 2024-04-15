package url_access

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) AccessURL(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.ReadMailRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("code is empty"))
		return
	}

	url, err := h.TrackedURL.ReadByCode(ctx, req.Code)
	if err != nil {
		log.For(c).Error("[read-mail-history] query database info failed", log.Field("code", req.Code), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	url.Status = entities.Read
	url.ReceivedAt = time.Now().UnixMilli()
	err = h.TrackedURL.Update(ctx, url)
	if err != nil {
		log.For(c).Error("[read-mail-history] update database failed", log.Field("code", req.Code), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// add event
	event := &entities.EventClickURL {
		Code: req.Code,
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
		log.For(c).Error("[redirect-shorten-link] create shorten link failed", log.Field("code", req.Code), log.Err(err))
	}

	c.JSON(http.StatusOK, nil)
}
